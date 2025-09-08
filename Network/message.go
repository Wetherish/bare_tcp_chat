package network

import (
	msgparser "chat_server/MsgParser"
	"fmt"
	"log"
	"net"
)

func SendMsg(msg msgparser.Message, conn net.Conn) {
	conn.Write(msg.ToBytes())
}

func ReadFromConnection(conn net.Conn, results chan<- msgparser.Message) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		msg := buf[:n]
		processedMsg, err := msgparser.ParseMsg(msg)
		if err != nil {
			log.Println("Error parsing message:", err)
			continue
		}
		results <- processedMsg
	}
}

var id = 0

func assignID(username string) int {
	log.Println("temporary service generating id for: ", username)
	id = id + 1
	return id
}

func ServerMessageProcessor(msg *msgparser.Message, conn net.Conn) {
	switch msg.Type {
	case msgparser.ID_REQUEST:
		responseMsg, err := msgparser.NewMessage(msgparser.ACCEPT, fmt.Sprint(assignID(msg.Value)))
		if err != nil {
			log.Println(err)
		}
		SendMsg(responseMsg, conn)
	default:
	}
}
