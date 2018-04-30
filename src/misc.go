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

// Returns the number of players in the queue
func getNumberOfPlayersInQueue() string {
	queuedPlayersLock.RLock()
	queued := len(queuedPlayers)
	queuedPlayersLock.RUnlock()

	qString := strconv.Itoa(queued)
	return qString + " players waiting in queue"
}
