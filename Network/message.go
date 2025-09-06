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

func ReadFromConnection(conn net.Conn, results chan<- string) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		msg := buf[:n]
		processedMsg, err := msgparser.ParseMsg(buf)
		if err != nil {
			log.Println("Error parsing message:", err)
			continue
		}
		messageProcessor(&processedMsg, conn)
		results <- string(msg)
	}
}

var id = 0

func assignID(username string) int {
	log.Println("temporary service generating id for: ", username)
	id = id + 1
	return id
}

func messageProcessor(msg *msgparser.Message, conn net.Conn) {
	switch msg.Type {
	case msgparser.ACCEPT:
	case msgparser.ID_REQUEST:
		responseMsg, err := msgparser.NewMessage(msgparser.ACCEPT, fmt.Sprintf("%d", assignID(msg.Value)))
		if err != nil {
			log.Println(err)
		}
		SendMsg(responseMsg, conn)
	case msgparser.RECONNECT:
	case msgparser.NEW_CONNECTION:
	case msgparser.MSG:
	case msgparser.DISCONNECT:
	}
}
