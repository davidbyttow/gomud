package main

import (
	"fmt"
	"math/rand"
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
