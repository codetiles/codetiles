package main

import (
	"crypto/rand"
	"io"
	"net/http"
	"sync"
)

// A player (deleted after 30 seconds of inactivity)
type user struct {
	mu     sync.Mutex
	id     [6]byte
	name   string
	inGame bool
	exp    int
}

func createUser(un string) [6]byte {
	var newUser user
	newUser.name = un
	newUser.exp = 30
	newUser.inGame = false

	var randomId [6]byte
	rand.Read(randomId[:])
	newUser.id = randomId

	users[newUser.id] = newUser
	return newUser.id
}

// API call for someone joining a lobby (create a user id)
func handleJoiningUser(w http.ResponseWriter, r *http.Request) {
	id := createUser("appins")
	io.WriteString(w, string(id[:]))
}
