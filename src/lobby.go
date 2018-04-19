package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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
	newUser.exp = 0
	newUser.inGame = false

	// Generate a secure random byte array then convert to base64.
	var randomByteArray [6]byte
	rand.Read(randomByteArray[:])
	rb64 := base64.StdEncoding.EncodeToString(randomByteArray[:])
	fmt.Println(rb64)
	var randomId [8]byte
	copy(randomId[:], []byte(rb64))
	newUser.id = randomId

	users[newUser.id] = newUser
	return newUser.id
}

type userCreateRequest struct {
	id string
}

// API call for someone creating a user (create a user id)
func handleJoiningUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := createUser("")

		idMap := map[string]string{"id": string(id[:])}
		resJ, err := json.Marshal(idMap)
		if err != nil {
			fmt.Println("Error marshalling json!")
		}

		io.WriteString(w, string(resJ))
	}

}

func handleUserCheck(w http.ResponseWriter, r *http.Request) {

}
