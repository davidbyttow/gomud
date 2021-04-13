package main

import (
	"fmt"
	"log"
	"net"
)

type Session struct {
	id   string
	conn net.Conn
}

func (s *Session) SessionId() string {
	return s.id
}

func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
}

var nextSessionId = 1

func generateSessionId() string {
	var sid = nextSessionId
	nextSessionId++
	return fmt.Sprintf("%d", sid)
}

func handleConnection(conn net.Conn, inputChannel chan SessionEvent) error {
	buf := make([]byte, 4096)

	session := &Session{generateSessionId(), conn}

	inputChannel <- SessionEvent{session, &SessionCreatedEvent{}}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			inputChannel <- SessionEvent{session, &SessionDisconnectEvent{}}
			return err
		}
		if n == 0 {
			log.Println("Zero bytes, closing connection")
			inputChannel <- SessionEvent{session, &SessionDisconnectEvent{}}
			break
		}

		input := string(buf[0 : n-2])
		log.Println("Received message:", input)

		inputChannel <- SessionEvent{session, &SessionInputEvent{input}}
	}

	return nil
}

func startServer(eventChannel chan SessionEvent) error {
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
