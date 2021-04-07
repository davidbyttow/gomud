package main

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
