package main

import (
	msgparser "chat_server/MsgParser"
	network "chat_server/Network"
	"log"
	"net"
	"os"
)

var (
	IP_ADDRESS = getEnv("IP_ADDRESS", "0.0.0.0")
	PORT       = getEnv("PORT", "8080")
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

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
				if string(msg.Content) == "" {
					continue
				}
				network.ServerMessageProcessor(&msg, &rooms, conn)
			}
		}(conn)
	}
}
