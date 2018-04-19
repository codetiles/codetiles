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
	var randomId [8]byte
	copy(randomId[:], []byte(rb64))
	newUser.id = randomId

	users[newUser.id] = newUser
	fmt.Println(newUser)
	return newUser.id
}

type userCreateRequest struct {
	id string
}

// API call for someone creating a user (create a user id)
func handleJoiningUser(w http.ResponseWriter, r *http.Request) {
	// We must have a POST request, otherwise send a 422 status code
	if r.Method == "POST" {
		var usernameMap map[string]string
		err := json.NewDecoder(r.Body).Decode(&usernameMap)
		// User error encoding json
		if err != nil {
			fmt.Println("Err unmarshalling json!")
			w.WriteHeader(400)
			io.WriteString(w, "Unable to unmarshal data (user -> server)")
			return
		}

		id := createUser(usernameMap["DisplayName"])

		idMap := map[string]string{"id": string(id[:])}
		resJ, err := json.Marshal(idMap)
		// Server error encoding id
		if err != nil {
			fmt.Println("Err marshalling json!")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Unable to create marshal json (server -> user)")
			return
		}

		io.WriteString(w, string(resJ))
		return
	}

	io.WriteString(w, "")

}

func handleCheckUser(w http.ResponseWriter, r *http.Request) {

}
