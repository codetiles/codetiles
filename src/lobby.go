package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Create a websocket and send a user information about the game they are
// waititng for
func handleWaitForGame(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Socket opened")
	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.WriteMessage(websocket.CloseMessage, []byte{})
	defer ws.Close()
	defer fmt.Println("Socket closed")

	if err != nil {
		fmt.Println("Error creating websocket in wait.go")
		fmt.Println("Error:", err)
		fmt.Println("Socket closed")
		return
	}

	go tickUser()

	// The user must send an id that they want to put in a lobby
	_, message, err := ws.ReadMessage()

	if err != nil {
		return
	}

	var uid [8]byte
	copy(uid[:], []byte(message))

	// If a player is queued or doesn't exist, send an error and close connection
	exists, _, queued, _ := checkUserId(uid)
	if !exists {
		ws.WriteMessage(websocket.TextMessage, []byte("User id does not exist"))
		return
	}

	if !queued {
		queuedPlayersLock.Lock()
		queuedPlayers = append(queuedPlayers, uid)
		queuedPlayersLock.Unlock()
	}

	defer removePlayerFromQueue(uid)
	defer func() { go tickUser() }()

	// Currently unimplented
	quit := make(chan int)

	err = ws.WriteMessage(websocket.TextMessage, []byte(getNumberOfPlayersInQueue()))
	if err != nil {
		return
	}

	// Send the user the # of users online
	for {
		select {
		case <-searchtick:
			err := ws.WriteMessage(websocket.TextMessage, []byte(getNumberOfPlayersInQueue()))
			if err != nil {
				return
			}
		case <-gametick:
			err := ws.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		case <-quit:
			return
		}
	}

}
