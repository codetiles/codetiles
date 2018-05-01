package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var openws []*websocket.Conn
var wwslock []*sync.Mutex
var wslock sync.RWMutex

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Create a websocket and send a user information about the game they are
// waititng for
func handleWaitForGame(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.WriteMessage(websocket.CloseMessage, []byte{})
	defer ws.Close()

	// WARNING: If this is not unlocked before return, the lobby sys will hang.
	// This is pointed to in wwslock
	mut := new(sync.Mutex)

	// Add both the write lock and the websocket writter
	wslock.Lock()
	openws = append(openws, ws)
	wwslock = append(wwslock, mut)
	wslock.Unlock()

	// Remove the ws and lock when returning
	defer func() {
		wslock.Lock()
		var arr []*websocket.Conn
		for _, j := range openws {
			if j != ws {
				arr = append(arr, j)
			}
		}

		var arr2 []*sync.Mutex
		for _, j := range wwslock {
			if j != mut {
				arr2 = append(arr2, j)
			}
		}

		openws = arr
		wwslock = arr2
		wslock.Unlock()
	}()
	// defer fmt.Println("Socket closed")

	if err != nil {
		fmt.Println("Error creating websocket in wait.go")
		fmt.Println("Error:", err)
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

	// A player can join a lobby if they are already queued, so no double queueing
	if !queued {
		queuedPlayersLock.Lock()
		queuedPlayers = append(queuedPlayers, uid)
		queuedPlayersLock.Unlock()
	}

	// Update other sockets when this user joins
	go tickUser()

	defer removePlayerFromQueue(uid)
	defer func() {
		go tickUser()
	}()

	mut.Lock()
	err = ws.WriteMessage(websocket.TextMessage, []byte(getNumberOfPlayersInQueue()))
	if err != nil {
		mut.Unlock()
		return
	}
	mut.Unlock()

	// To prevent connections being closed, recieve a message from the client every now and then
	go func() {
		for {
			ws.SetReadDeadline(time.Now().Add(time.Duration(time.Second * 5)))
			_, _, err := ws.ReadMessage()
			if err != nil {
				return
			}
			mut.Lock()
			ws.SetWriteDeadline(time.Now().Add(time.Duration(time.Second * 1)))
			err = ws.WriteMessage(websocket.TextMessage, []byte(getNumberOfPlayersInQueue()))
			if err != nil {
				mut.Unlock()
				return
			}
			mut.Unlock()
		}

	}()

	// Send the user the # of users online when one joins or when an autosearchtick
	// Elapses for every user
	for {
		select {
		case <-searchtick:
			mut.Lock()
			ws.SetWriteDeadline(time.Now().Add(time.Duration(time.Second * 1)))
			err := ws.WriteMessage(websocket.TextMessage, []byte(getNumberOfPlayersInQueue()))
			if err != nil {
				mut.Unlock()
				return
			}
			mut.Unlock()
		case <-autosearchtick:
			mut.Lock()
			ws.SetWriteDeadline(time.Now().Add(time.Duration(time.Second * 1)))
			err := ws.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				mut.Unlock()
				return
			}
			mut.Unlock()
		}
	}

}
