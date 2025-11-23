package network

import (
	"bytes"
	msgparser "chat_server/MsgParser"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type Room struct {
	Id          uint32
	Name        string
	Connections map[uint32]net.Conn
	mu          sync.Mutex
}

func NewRoom(id uint32, name string) *Room {
	return &Room{
		Id:          id,
		Name:        name,
		Connections: make(map[uint32]net.Conn),
	}
}

func (r *Room) AddConnection(id uint32, conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Connections[id] = conn
}

func (r *Room) RemoveConnection(id uint32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Connections, id)
}

func (r *Room) BroadcastMessage(msgType string, message []byte, senderId uint32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, conn := range r.Connections {
		if id != senderId {
			sendMsg, err := msgparser.NewMessage(msgType, message, senderId, r.Id)
			if err != nil {
				log.Println(err)
				continue
			}
			SendMsg(sendMsg, conn)
		}
	}
}

type Rooms struct {
	Rooms map[uint32]*Room
	mu    sync.Mutex
}

func NewRooms() Rooms {
	return Rooms{
		Rooms: make(map[uint32]*Room),
		mu:    sync.Mutex{},
	}
}

func (rooms *Rooms) JoinRoom(roomId, userId uint32, conn net.Conn) error {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()
	room, ok := rooms.Rooms[roomId]
	if !ok {
		return errors.New("Room does not exist")
	}

	room.AddConnection(userId, conn)
	return nil
}

func (rooms *Rooms) ListRooms() []byte {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	var buf bytes.Buffer
	for _, room := range rooms.Rooms {
		buf.WriteString(room.Name + ".........." + fmt.Sprint(room.Id) + "\n") //TODO replace with prettier version
	}
	return buf.Bytes()
}

func (rooms *Rooms) LeaveRoom(roomId, userId uint32) error {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	room, ok := rooms.Rooms[roomId]
	if !ok {
		return errors.New("Room does not exists")
	}
	room.RemoveConnection(userId)
	return nil
}

func (rooms *Rooms) ListUsers(roomId uint32) ([]byte, error) {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	room, ok := rooms.Rooms[roomId]
	if !ok {
		return nil, errors.New("Room does not exists")
	}

	var buf bytes.Buffer
	for id, conn := range room.Connections {
		buf.WriteString(fmt.Sprint(id) + ": " + conn.LocalAddr().String())
	}

	return buf.Bytes(), nil
}

func (rooms *Rooms) CreateRoom(id uint32, name string) {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	rooms.Rooms[id] = NewRoom(id, name)
}
