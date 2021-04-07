package main

import (
	"log"
	"net"
)

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
