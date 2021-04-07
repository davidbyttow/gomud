package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

type MessageEvent struct {
	msg string
}

type MoveEvent struct {
	dir string
}

type UserJoinedEvent struct {
}

type ClientInput struct {
	user  *User
	event interface{}
}

type User struct {
	id      string
	name    string
	session *Session
}

type Session struct {
	conn net.Conn
}

func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
}

type World struct {
	users        []*User
	startRoom    string
	roomsById    map[string]*Room
	usersToRooms map[string]string
}

func (w *World) GetUser(id string) *User {
	for _, u := range w.users {
		if u.id == id {
			return u
		}
	}
	return nil
}

func (w *World) GetRoom(id string) *Room {
	return w.roomsById[id]
}

func (w *World) GetUsersInRoom(roomId string) []*User {
	var users []*User
	for u, r := range w.usersToRooms {
		if r == roomId {
			users = append(users, w.GetUser(u))
		}
	}
	return users
}

func (w *World) Broadcast(user *User, msg string) {
	room := w.GetRoom(w.usersToRooms[user.id])
	for _, other := range w.GetUsersInRoom(room.id) {
		if other != user {
			other.session.WriteLine(msg)
		}
	}
}

func (w *World) AddToRoom(user *User, roomId string) {
	prevRoom := w.GetRoom(roomId)
	if prevRoom != nil {

	}

	w.usersToRooms[user.id] = roomId
	room := w.GetRoom(roomId)
	user.session.WriteLine(room.desc)
}

func (w *World) Move(user *User, dir string) {
	room := w.GetRoom(w.usersToRooms[user.id])
	for _, link := range room.links {
		if link.verb == dir {
			w.AddToRoom(user, link.roomId)
		}
	}
}

type RoomLink struct {
	verb   string
	roomId string
}

type Room struct {
	id    string
	desc  string
	links []*RoomLink
}

func generateName() string {
	return fmt.Sprintf("User %d", rand.Intn(100)+1)
}

func handleConnection(conn net.Conn, inputChannel chan ClientInput) error {
	buf := make([]byte, 4096)

	session := &Session{conn}
	user := &User{name: generateName(), session: session}

	inputChannel <- ClientInput{
		user,
		&UserJoinedEvent{},
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 {
			log.Println("Zero bytes, closing connection")
			break
		}
		msg := string(buf[0 : n-2])
		log.Println("Received message:", msg)

		if msg == "east" || msg == "west" || msg == "north" || msg == "south" {
			e := ClientInput{user, &MoveEvent{msg}}
			inputChannel <- e
		} else {
			e := ClientInput{user, &MessageEvent{msg}}
			inputChannel <- e
		}
	}

	return nil
}

func startServer(eventChannel chan ClientInput) error {
	log.Println("Starting server")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		go func() {
			if err := handleConnection(conn, eventChannel); err != nil {
				log.Println("Error handling connection", err)
				return
			}
		}()
	}
}

func startGameLoop(w *World, clientInputChannel <-chan ClientInput) {
	for input := range clientInputChannel {
		switch event := input.event.(type) {
		case *MessageEvent:
			// TODO: error handling
			input.user.session.WriteLine(fmt.Sprintf("You said, \"%s\"", event.msg))
			w.Broadcast(input.user, fmt.Sprintf("%s said, \"%s\"", input.user.name, event.msg))
			// for _, user := range w.users {
			// 	if user != input.user {
			// 		user.session.WriteLine(fmt.Sprintf("%s said, \"%s\"", input.user.name, event.msg))
			// 	}
			// }

		case *MoveEvent:
			w.Move(input.user, event.dir)

		case *UserJoinedEvent:
			input.user.session.WriteLine(fmt.Sprintf("Welcome %s", input.user.name))

			w.users = append(w.users, input.user)
			w.AddToRoom(input.user, w.startRoom)

			// for _, user := range w.users {
			// 	if user != input.user {
			// 		user.session.WriteLine(fmt.Sprintf("%s entered the room.", input.user.name))
			// 	}
			// }
		}
	}
}

func createWorld() *World {
	rooms := []*Room{
		{
			id:   "A",
			desc: "This is a room with a sign that has the letter A written on it.",
			links: []*RoomLink{
				{
					verb:   "east",
					roomId: "B",
				},
			},
		},
		{
			id:   "B",
			desc: "This is a room with a sign that has the letter B written on it.",
			links: []*RoomLink{
				{
					verb:   "west",
					roomId: "A",
				},
			},
		},
	}

	w := &World{
		usersToRooms: map[string]string{},
		roomsById:    map[string]*Room{},
		startRoom:    rooms[0].id,
	}

	for _, room := range rooms {
		w.roomsById[room.id] = room
	}

	return w
}

func main() {
	ch := make(chan ClientInput)

	world := createWorld()

	go startGameLoop(world, ch)

	err := startServer(ch)
	if err != nil {
		log.Fatal(err)
	}
}
