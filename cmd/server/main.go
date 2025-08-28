package main

import (
	msgparser "chat_server/MsgParser"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	IP_ADDRESS = "127.0.0.1"
	PORT       = "8080"
)

func sendMsg(msg *msgparser.Message, conn net.Conn) {
    conn.Write(msg.ToBytes())
}

func messageProcessor(msg *msgparser.Message, c *net.Conn) {
	switch msg.Type {
	case msgparser.ACCEPT:
	case msgparser.RECONNECT:
	case msgparser.NEW_CONNECTION:
	case msgparser.MSG:
	case msgparser.DISCONNECT:
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 1024)
	for {
		num, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("connection lost")
			} else {
				log.Fatalln(err)
			}
			return
		}
		data := buf[:num]
		message, err := msgparser.ParseMsg(data)
		println(message.Value)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Printf("received: %s\n", string(data))
	}
}

func main() {
	ln, err := net.Listen("tcp", IP_ADDRESS+":"+PORT)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	log.Printf("server is listening on %s", ln.Addr().String())
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handleConnection(conn)
	}
}
