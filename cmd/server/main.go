package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const (
	IP_ADDRESS = "127.0.0.1"
	PORT       = "8080"
)

func handleConnection(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 1024)
	for {
		num, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection lost")
			} else{
				log.Fatalln(err)
			}
			return
		}
		data := buf[:num]
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
