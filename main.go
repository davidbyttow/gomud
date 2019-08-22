package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	request := make([]byte, 4096)
	defer conn.Close()
	for {
		len, err := conn.Read(request)
		if err != nil {
			panic(err)
		}
		if len == 0 {
			break
		}
		input := string(request[:len])
		input = strings.TrimSpace(input)
		echo := fmt.Sprintf("You said, \"%s\"", input)
		if _, err := conn.Write([]byte(echo + "\n")); err != nil {
			panic(err)
		}
	}
}

func startServer() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		println("New connection accepted")
		if err != nil {
			return err
		}
		go handleConnection(conn)
	}
}

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}
