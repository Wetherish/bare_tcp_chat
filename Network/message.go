package network

import (
	msgparser "chat_server/MsgParser"
	"fmt"
	"log"
	"net"
	"strconv"
)

func SendMsg(msg msgparser.Message, conn net.Conn) {
	conn.Write(msg.ToBytes())
}

func ReadFromConnection(conn net.Conn, results chan<- msgparser.Message) {
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

var id uint32 = 0

func assignID(username string) uint32 {
	log.Println("temporary service generating id for: ", username)
	id = id + 1
	return id
}

func ServerMessageProcessor(msg *msgparser.Message, rooms *Rooms, conn net.Conn) {
	switch msg.Type {
	case msgparser.ID_REQUEST:
		id := assignID(msg.Value)
		responseMsg, err := msgparser.NewMessage(msgparser.ACCEPT, fmt.Sprint(id), id, 0)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(responseMsg.ToBytes()))
		SendMsg(responseMsg, conn)

	case msgparser.JOIN:
		roomID, err := strconv.Atoi(msg.Value)
		if err != nil {
			log.Println("Error converting room ID:", err)
			return
		}
		rooms.JoinRoom(uint32(roomID), msg.UserId, conn)
		responseMsg, err := msgparser.NewMessage(msgparser.ACCEPT, "joined_room", 0, uint32(roomID))
		if err != nil {
			log.Println(err)
		}
		SendMsg(responseMsg, conn)

	case msgparser.LIST_ROOMS:
		responseMsg, err := msgparser.NewMessage(msgparser.ACCEPT, string(rooms.ListRooms()), 0, 0)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(rooms.ListRooms())
		SendMsg(responseMsg, conn)
	case msgparser.MSG:
		if msg.RoomId == 0 {
			log.Println("MSG Sended into nowhere: ", msg.Value)
		} else {
			rooms.Rooms[msg.RoomId].BroadcastMessage(msg.Value, msg.UserId)
		}
	default:
	}
}
