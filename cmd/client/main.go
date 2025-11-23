package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	msgparser "chat_server/MsgParser"
	network "chat_server/Network"
)

const (
	SERVER_ADDRESS = "127.0.0.1:8080"
)

var (
	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	otherStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type incomingMsg msgparser.Message
type errMsg error

type model struct {
	conn      net.Conn
	user      network.User
	roomID    uint32
	viewport  viewport.Model
	textInput textinput.Model
	messages  []string
	err       error
}

func initialModel(conn net.Conn, user network.User, roomID uint32) model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	vp := viewport.New(30, 5)
	vp.SetContent("Welcome to the chat room!")

	return model{
		conn:      conn,
		user:      user,
		roomID:    roomID,
		textInput: ti,
		viewport:  vp,
		messages:  []string{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, waitForMessage(m.conn))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textInput.Width = msg.Width
		m.viewport.Height = msg.Height - 3 // Leave space for input and border
		m.viewport.SetContent(strings.Join(m.messages, "\n"))

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input == "" {
				break
			}

			if input == "/exit" {
				return m, tea.Quit
			}

			if input == "/help" {
				helpMsg := `Available commands:
/join <id> - Join room
/list_rooms - List rooms
/sendfile <path> - Send file
/exit - Exit client`
				m.messages = append(m.messages, lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(helpMsg))
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
				m.textInput.SetValue("")
				m.textInput.Reset()
				break
			}

			// Send message
			sendMessage(input, m.conn, m.user, m.roomID)

			// Optimistically display own message
			if !isCommand(input) {
				displayMsg := senderStyle.Render(fmt.Sprintf("You: %s", input))
				m.messages = append(m.messages, displayMsg)
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			}

			m.textInput.SetValue("")
			m.textInput.Reset()
		}

	case incomingMsg:
		// Handle incoming message
		content := string(msg.Content)
		var displayMsg string

		switch msg.Type {
		case msgparser.MSG:
			if msg.UserId == m.user.Id {
				displayMsg = senderStyle.Render(fmt.Sprintf("You: %s", content))
			} else {
				displayMsg = otherStyle.Render(fmt.Sprintf("User %d: %s", msg.UserId, content))
			}
		case msgparser.ACCEPT:
			if content == "joined_room" {
				m.roomID = msg.RoomId
				displayMsg = fmt.Sprintf("Joined room %d", msg.RoomId)
			} else {
				displayMsg = fmt.Sprintf("Server: %s", content)
			}
		case msgparser.FILE:
			parts := strings.SplitN(content, "|", 2)
			if len(parts) == 2 {
				filename := parts[0]
				fileContent := parts[1]
				err := os.WriteFile("received_"+filename, []byte(fileContent), 0644)
				if err != nil {
					displayMsg = errStyle.Render(fmt.Sprintf("Failed to save file %s: %v", filename, err))
				} else {
					displayMsg = fmt.Sprintf("Received file from %d: %s (saved as received_%s)", msg.UserId, filename, filename)
				}
			} else {
				displayMsg = fmt.Sprintf("Received invalid file format from %d", msg.UserId)
			}
		default:
			displayMsg = fmt.Sprintf("Unknown message type from %d: %s", msg.UserId, content)
		}

		m.messages = append(m.messages, displayMsg)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, waitForMessage(m.conn)

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textInput.View(),
	)
}

func waitForMessage(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		// We need to read one message from the connection
		// Since ReadFromConnection is a loop, we can't use it directly here as is.
		// We need a way to read a single message.
		// Let's implement a single read helper or modify the logic.
		// Actually, ReadFromConnection writes to a channel. We can use that if we adapt it.
		// But for tea.Cmd, it's better to have a blocking read that returns one message.

		// Re-implementing single message read logic here for simplicity and control
		var length uint32
		err := binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			return errMsg(err)
		}

		buf := make([]byte, length)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			return errMsg(err)
		}

		processedMsg, err := msgparser.ParseMsg(buf)
		if err != nil {
			return errMsg(err)
		}

		return incomingMsg(processedMsg)
	}
}

// Imports moved to top

func sendMessage(input string, conn net.Conn, user network.User, roomID uint32) {
	var msgType string
	var content string
	var err error

	if isCommand(input) {
		msgType, content, err = commandProcessor(input, roomID)
		if err != nil {
			// In a real app we might want to show this error in the UI
			return
		}
	} else {
		msgType = msgparser.MSG
		content = input
	}

	msg, err := msgparser.NewMessage(msgType, []byte(content), user.Id, roomID)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		return
	}

	network.SendMsg(msg, conn)
}

func isCommand(input string) bool {
	return len(input) > 0 && input[0] == '/'
}

func commandProcessor(input string, currentRoomID uint32) (string, string, error) {
	args := strings.Split(input, " ")
	switch args[0] {
	case "/list_rooms":
		return msgparser.LIST_ROOMS, "1", nil
	case "/join":
		if len(args) == 2 {
			return msgparser.JOIN, args[1], nil
		}
		return "", "", fmt.Errorf("usage: /join <room_id>")
	case "/room":
		return "", "", fmt.Errorf("room number: %d", currentRoomID)
	case "/sendfile":
		if len(args) == 2 {
			filePath := args[1]
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return "", "", fmt.Errorf("failed to read file: %v", err)
			}
			filename := filePath
			if idx := strings.LastIndex(filePath, "/"); idx != -1 {
				filename = filePath[idx+1:]
			}
			return msgparser.FILE, fmt.Sprintf("%s|%s", filename, string(fileContent)), nil
		}
		return "", "", fmt.Errorf("usage: /sendfile <path>")
	default:
		return "", "", fmt.Errorf("invalid command: %s", args[0])
	}
}

func handleConnection() (net.Conn, network.User, uint32, error) {
	conn, err := net.Dial("tcp", SERVER_ADDRESS)
	if err != nil {
		return nil, network.User{}, 0, fmt.Errorf("connection failed: %v", err)
	}

	// Simple initial handshake outside of Bubble Tea for now
	// Or we could make this part of the initial model state (login screen)
	// For now, let's keep it simple: ask for name in stdio then switch to TUI

	fmt.Print("Enter name: ")
	var nickname string
	fmt.Scanln(&nickname)

	idRequest, err := msgparser.NewMessage(msgparser.ID_REQUEST, []byte(nickname), 0, 0)
	if err != nil {
		conn.Close()
		return nil, network.User{}, 0, fmt.Errorf("failed to create ID request: %v", err)
	}

	network.SendMsg(idRequest, conn)

	// Read ID response
	var length uint32
	binary.Read(conn, binary.BigEndian, &length)
	buf := make([]byte, length)
	io.ReadFull(conn, buf)

	idMsg, _ := msgparser.ParseMsg(buf)
	fmt.Println("Connected with ID:", idMsg.UserId)

	user := network.NewUser(nickname, idMsg.UserId)
	return conn, user, 0, nil
}

func main() {
	conn, user, roomID, err := handleConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	p := tea.NewProgram(initialModel(conn, user, roomID), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
