package network

import (
	msgparser "chat_server/MsgParser"
	"encoding/binary"
	"io"
	"log"
	"net"
	"strconv"
)

func SendMsg(msg msgparser.Message, conn net.Conn) {
	bytesMsg := msg.ToBytes()
	length := uint32(len(bytesMsg))
	err := binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		log.Println("Error sending message length:", err)
		return
	}
	conn.Write(bytesMsg)
}

func ReadFromConnection(conn net.Conn, results chan<- msgparser.Message) {
	for {
		var length uint32
		err := binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading message length:", err)
			}
			break
		}

		buf := make([]byte, length)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			log.Println("Error reading message body:", err)
			break
		}

		processedMsg, err := msgparser.ParseMsg(buf)
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
		id := assignID(string(msg.Content))
		idBytes := []byte(strconv.FormatUint(uint64(id), 10))
		responseMsg, err := msgparser.NewMessage(msgparser.ACCEPT, idBytes, id, 0)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(responseMsg.ToBytes()))
		SendMsg(responseMsg, conn)

	case msgparser.JOIN:
		roomID, err := strconv.Atoi(string(msg.Content))
		if err != nil {
			log.Println("Error converting room ID:", err)
			return
		}
		rooms.JoinRoom(uint32(roomID), msg.UserId, conn)
		responseMsg, err := msgparser.NewMessage(msgparser.ACCEPT, []byte("joined_room"), 0, uint32(roomID))
		if err != nil {
			log.Println(err)
		}
		SendMsg(responseMsg, conn)

	case msgparser.LIST_ROOMS:
		responseMsg, err := msgparser.NewMessage(msgparser.ACCEPT, rooms.ListRooms(), 0, 0)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(rooms.ListRooms())
		SendMsg(responseMsg, conn)
	case msgparser.MSG:
		if msg.RoomId == 0 {
			log.Println("MSG Sended into nowhere: ", msg.Content)
		} else {
			rooms.Rooms[msg.RoomId].BroadcastMessage(msgparser.MSG, msg.Content, msg.UserId)
		}
	case msgparser.FILE:
		if msg.RoomId == 0 {
			log.Println("FILE Sended into nowhere: ", msg.Content)
		} else {
			rooms.Rooms[msg.RoomId].BroadcastMessage(msgparser.FILE, msg.Content, msg.UserId)
		}
	default:
	}
}
