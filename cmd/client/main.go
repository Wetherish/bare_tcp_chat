package main

import (
	"fmt"
	"log"
	"net"
)

const (
	IP_ADDRESS = "127.0.0.1"
	PORT       = "8080"
)

func main() {
	conn, err := net.Dial("tcp", IP_ADDRESS+":"+PORT)
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	msg := ""
	for {
		fmt.Scanln(&msg)
		conn.Write([]byte(msg))
		fmt.Printf("you: %s\n", msg)
	}
	log.Println("Client connected to the server: ")
}
