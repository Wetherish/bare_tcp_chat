package main

import (
	"bufio"
	msgparser "chat_server/MsgParser"
	network "chat_server/Network"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	SERVER_ADDRESS = "127.0.0.1:8080"
)

var roomID uint32 = 0

func isCommand(input string) bool {
	return input[0] == '/'
}

func commandProcessor(input string) (string, string, error) {
	var msgType string
	var content string
	args := strings.Split(input, " ")
	switch args[0] {
	case "/list_rooms":
		msgType = msgparser.LIST_ROOMS
		content = "1"
	case "/join":
		msgType = msgparser.JOIN
		if len(args) == 2 {
			content = args[1]
		} else {
			return "", "", fmt.Errorf("invalid amount of arguments (%v) for this command", len(args))
		}
	case "/room":
		return "", "", fmt.Errorf("room number: %d", roomID)
	default:
		return "", "", fmt.Errorf("invalid command: %s", args[0])
	}
	return msgType, content, nil
}

func sendMessage(input string, conn net.Conn, user network.User) bool {
	if input == "/exit" {
		return true
	}

	var msgType string
	var content string
	var err error

	if isCommand(input) {
		msgType, content, err = commandProcessor(input)
		if err != nil {
			fmt.Println(err)
			return false
		}
	} else {
		msgType = msgparser.MSG
		content = input
	}

	msg, err := msgparser.NewMessage(msgType, []byte(content), user.Id, roomID)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		return false
	}

	network.SendMsg(msg, conn)
	return false
}

func handleConnection() (net.Conn, network.User, error) {
	conn, err := net.Dial("tcp", SERVER_ADDRESS)
	if err != nil {
		return nil, network.User{}, fmt.Errorf("connection failed: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter name: ")
	scanner.Scan()
	nickname := scanner.Text()

	idRequest, err := msgparser.NewMessage(msgparser.ID_REQUEST, []byte(nickname), 0, 0)
	if err != nil {
		conn.Close()
		return nil, network.User{}, fmt.Errorf("failed to create ID request: %v", err)
	}

	network.SendMsg(idRequest, conn)
	fmt.Println("Connecting to server...")

	msgCh := make(chan msgparser.Message, 10)
	go network.ReadFromConnection(conn, msgCh)

	idMsg := <-msgCh
	fmt.Println("Connected with ID:", idMsg.UserId)

	user := network.NewUser(nickname, idMsg.UserId)

	go func() {
		for msg := range msgCh {
			if msg.Type == msgparser.ACCEPT && string(msg.Content) == "joined_room" {
				roomID = msg.RoomId
				fmt.Println("You joined room nr: ", msg.RoomId)
			} else {
				fmt.Println(msg.UserId, ": ", string(msg.Content))
			}
		}
	}()

	return conn, user, nil
}

func main() {
	fmt.Println("Connected to chat server. Type your messages:")

	conn, user, err := handleConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if sendMessage(scanner.Text(), conn, user) {
			break
		}
	}
}
