package main

import (
	"fmt"
	"log"
	"net"
)

type Session struct {
	conn net.Conn
}

func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
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

func main() {
	ch := make(chan ClientInput)

	world := createWorld()

	go startGameLoop(world, ch)

	err := startServer(ch)
	if err != nil {
		log.Fatal(err)
	}
}
