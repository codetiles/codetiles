package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var users map[[8]byte]user
var queuedPlayers map[[8]byte]bool

// A player (deleted after 30 seconds of inactivity)
type user struct {
	mu     sync.Mutex
	id     [8]byte
	name   string
	inGame bool
	exp    int
}

func createUser(un string) [8]byte {
	var newUser user
	newUser.name = un
	newUser.exp = 30
	newUser.inGame = false

	// Generate a secure random byte array then convert to base64.
	var randomByteArray [6]byte
	rand.Read(randomByteArray[:])
	rb64 := base64.StdEncoding.EncodeToString(randomByteArray[:])
	fmt.Println(rb64)
	var randomId [8]byte
	copy([]byte(rb64), randomId[:])
	newUser.id = randomId

	users[newUser.id] = newUser
	return newUser.id
}

// API call for someone joining a lobby (create a user id)
func handleJoiningUser(w http.ResponseWriter, r *http.Request) {
	id := createUser("appins")
	io.WriteString(w, string(id[:]))
}
