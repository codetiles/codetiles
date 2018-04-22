package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type returnPlayers struct {
	InGame     bool
	UserExists bool
	Players    []string
}

// Return the players that are in the game to the client (see struct above)
func handleRetrievePlayers(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")
	var userId [8]byte
	copy(userId[:], []byte(id))

	usersArrayLock.RLock()
	_, exists := users[userId]
	usersArrayLock.RUnlock()

	if !exists {
		j, err := json.Marshal(returnPlayers{
			InGame:     false,
			UserExists: false,
			Players:    []string{},
		})

		if err != nil {
			fmt.Println("Error marshalling player list (server -> client)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(j))
		return
	}

	if exists {
		var returnPlayers returnPlayers
		returnPlayers.UserExists = true

		usersArrayLock.RLock()
		returnPlayers.InGame = users[userId].inGame
		usersArrayLock.RUnlock()

		if returnPlayers.InGame == true {
			// TODO: Read from maps
			returnPlayers.Players = []string{"appins", "zane", "yev"}
			return
		}

		j, err := json.Marshal(returnPlayers)

		if err != nil {
			fmt.Println("Error marshalling returnPlayers  (server -> client)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(j))
	}
}
