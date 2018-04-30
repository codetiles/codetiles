package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// Handle a user joining a lobby
func handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")

	var uid [8]byte
	copy(uid[:], []byte(id))

	// NOTE: Consider checking if the user is already in a game and put them back
	usersArrayLock.RLock()
	_, ex := users[uid]
	usersArrayLock.RUnlock()

	if !ex {
		j, err := json.Marshal(map[string]bool{
			"Success":  false,
			"IdExists": false,
		})

		if handleJsonMarshalError(w, r, "lobby.go - joinlobby/user does not exist", err) {
			return
		}

		io.WriteString(w, string(j))
		return
	}

	// Allow use for queued players lock for the rest of the function
	queuedPlayersLock.Lock()
	defer queuedPlayersLock.Unlock()

	for _, i := range queuedPlayers {
		if i == uid {
			// If they are already in the queue, set success to false
			j, err := json.Marshal(map[string]bool{
				"Success":  false,
				"IdExists": true,
			})

			if handleJsonMarshalError(w, r, "lobby.go - joinlobby/already in queue", err) {
				return
			}

			io.WriteString(w, string(j))
			return
		}
	}

	queuedPlayers = append(queuedPlayers, uid)
	searchtick <- 0

	// If they are already in the queue, set success to false
	j, err := json.Marshal(map[string]bool{
		"Success":  true,
		"IdExists": true,
	})

	if handleJsonMarshalError(w, r, "lobby.go - joinlobby/already in queue", err) {
		return
	}

	io.WriteString(w, string(j))
	return

}
