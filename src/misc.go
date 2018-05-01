package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Handle an error marshalling (or encoding) json
func handleJsonMarshalError(w http.ResponseWriter, r *http.Request, he string, err error) bool {
	if err != nil {
		fmt.Println("Err marshalling JSON (server -> client)")
		fmt.Println("\t --", he)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Unable to marshal JSON")
		return true
	}

	return false
}

// A much less serious error: handle json that is malformed
func handleJsonUnmarshalError(w http.ResponseWriter, r *http.Request, he string, err error) bool {
	if err != nil {
		fmt.Println("Err unmarshalling JSON (client -> server)")
		fmt.Println("\t --", he)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Malformed JSON")
		return true
	}

	return false
}

// Returns the number of players in the queue in a string format
func getNumberOfPlayersInQueue() string {
	queuedPlayersLock.RLock()
	queued := len(queuedPlayers)
	queuedPlayersLock.RUnlock()

	qString := strconv.Itoa(queued)
	if queued <= 1 {
		return qString + " player waiting in queue..."
	}
	return qString + " players wating in queue..."

}

// Removes a player from the queue
func removePlayerFromQueue(id [8]byte) {
	var qpa [][8]byte

	queuedPlayersLock.RLock()
	for _, j := range queuedPlayers {
		if j != id {
			qpa = append(qpa, j)
		}
	}
	queuedPlayersLock.RUnlock()
	queuedPlayersLock.Lock()
	queuedPlayers = qpa
	queuedPlayersLock.Unlock()
}
