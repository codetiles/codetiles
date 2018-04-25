package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
)

var usersArrayLock sync.RWMutex
var users map[[8]byte]user

var queuedPlayersLock sync.RWMutex
var queuedPlayers [][8]byte

// A player (deleted after 30 seconds of inactivity)
type user struct {
	id     [8]byte
	name   string
	inGame bool
	exp    int
	code   string
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

	// Add to users array
	usersArrayLock.Lock()
	users[newUser.id] = newUser
	usersArrayLock.Unlock()

	return newUser.id
}

// API call for someone creating a user (create a user id)
func handleJoiningUser(w http.ResponseWriter, r *http.Request) {
	// We must have a POST request, otherwise send a 422 status code
	if r.Method == "POST" {
		var usernameMap map[string]string
		err := json.NewDecoder(r.Body).Decode(&usernameMap)

		if handleJsonUnmarshalError(w, r, "code.go - upload", err) {
			return
		}

		id := createUser(usernameMap["DisplayName"])

		idMap := map[string]string{
			"Id": string(id[:]),
		}
		resJ, err := json.Marshal(idMap)

		if handleJsonMarshalError(w, r, "user.go - response", err) {
			return
		}

		io.WriteString(w, string(resJ))
		return
	}

	w.WriteHeader(422)
	io.WriteString(w, "Wrong type of request")

}

// Verify that a user exists and send the name as a result
func handleVerifyUser(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Get everything after the forth backslash
	idstring := strings.Join(strings.Split(path, "/")[4:], "/")

	// Copy the user id and look it up in the users array
	var userid [8]byte
	copy(userid[:], []byte(idstring))
	usersArrayLock.RLock()
	user, exists := users[userid]
	defer usersArrayLock.RUnlock()

	if exists {
		j, err := json.Marshal(map[string]string{
			"DisplayName": user.name,
			"Exists":      "true",
		})

		if handleJsonMarshalError(w, r, "users.go - verify/exists", err) {
			return
		}

		io.WriteString(w, string(j))
		return
	}

	j, err := json.Marshal(map[string]string{
		"DisplayName": "err",
		"Exists":      "false",
	})

	if handleJsonMarshalError(w, r, "users.go - verify/does not exist", err) {
		return
	}

	io.WriteString(w, string(j))

}
