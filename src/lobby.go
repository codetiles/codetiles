package main

import (
	"net/http"
	"encoding/json"
	"io"
	"fmt"
)

// Handle a user joining a lobby
func handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")
	gt := r.Header.Get("gametype")

	gametype := "testing"

	switch gt {
	case "testing":
		gametype = "testing"
	case "1v1":
		gametype = "1v1"
	}

	var uid [8]byte
	copy(uid[:], []byte(id))

	// NOTE: Consider checking if the user is already in a game
	usersArrayLock.RLock()
	_, ex := users[uid]
	usersArrayLock.RUnlock()

	if !ex {
		j, err := json.Marshal(map[string]bool{
			"Success" : false,
			"IdExists" : false,
		})

		if handleJsonMarshalError(w, r, "lobby.go - joinlobby/user does not exist", err) {
			return
		}

		io.WriteString(w, string(j))
		return
	}

	fmt.Println(gametype)
}
