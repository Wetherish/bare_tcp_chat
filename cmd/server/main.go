package main

import (
	msgparser "chat_server/MsgParser"
	network "chat_server/Network"
	"log"
	"net"
)

const (
	IP_ADDRESS = "127.0.0.1"
	PORT       = "8080"
)

func main() {
	ln, err := net.Listen("tcp", IP_ADDRESS+":"+PORT)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	log.Printf("server is listening on %s", ln.Addr().String())
	rooms := network.NewRooms()
	rooms.CreateRoom(1, "Default")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go func(c net.Conn) {
			defer func() {
				c.Close()
			}()
			result := make(chan msgparser.Message)
			defer close(result)

			go network.ReadFromConnection(c, result)

			for msg := range result {
				if msg.Value == "" {
					continue
				}
				network.ServerMessageProcessor(&msg, &rooms, conn)
			}
		}(conn)
	}
}
