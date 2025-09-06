package main

import (
	"bufio"
	msgparser "chat_server/MsgParser"
	network "chat_server/Network"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	IP_ADDRESS = "127.0.0.1"
	PORT       = "8080"
)

func inputMessage(scanner *bufio.Scanner, conn net.Conn) {
	for {
		fmt.Print("Enter message: ")
		if !scanner.Scan() {
			break
		}
		message := scanner.Text()

		exit := "/exit"
		if message == exit {
			fmt.Println("Good bye!")
			break
		}

		msg, err := msgparser.NewMessage(msgparser.MSG, message)
		if err != nil {
			log.Fatalln(err.Error())
		}
		network.SendMsg(msg, conn)
		fmt.Printf("you: %s\n", msg.Value)
	}
}

func main() {
	conn, err := net.Dial("tcp", IP_ADDRESS+":"+PORT)

	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter name: ")
	scanner.Scan()
	nickname := scanner.Text()
	fmt.Println(nickname)
	idRequest, err := msgparser.NewMessage(msgparser.ID_REQUEST, nickname)
	if err != nil {
		log.Fatalln("failed to create ID_REQUEST:", err)
	}

	network.SendMsg(idRequest, conn)
	idCh := make(chan string)

	fmt.Println("Connecting to server ...")
	go network.ReadFromConnection(conn, idCh)
	id := <-idCh
	log.Println("DEBUG log: ", id)
	idMsg, err := msgparser.ParseMsg([]byte(id))
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println("Connected with ID:", idMsg.Value)
	inputMessage(scanner, conn)
}
