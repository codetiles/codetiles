package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// use <-searchtick in select to check when a user joins the queue
var searchtick chan int
var autosearchtick chan int

var gTick int
var gTickLock sync.RWMutex

func tickUser() {
	searchtick <- 0
}

func performSearchTick() {
	time.Sleep(250 * time.Millisecond)
	autosearchtick <- 0
	performSearchTick()
}

// Check if the countdown has started every 100 milliseconds
func checkCountdown() {
	defer func() {
		time.Sleep(time.Millisecond * 100)
		checkCountdown()
	}()

	queuedPlayersLock.RLock()
	lenQ := len(queuedPlayers)
	queuedPlayersLock.RUnlock()

	// If there are >= 2 players, start the countdown
	if lenQ >= 2 {
		wslock.RLock()
		defer wslock.RUnlock()

		// Lock every sockets writer
		// WARNING: If a socket's locker is stuck in the locked state, it will hang
		for _, j := range wwslock {
			j.Lock()
		}

		for i := 5; i >= 1; i-- {

			ss := "Game starts in " + strconv.Itoa(i)
			if i > 1 {
				ss += " seconds"
			} else {
				ss += " second"
			}

			// Write the timestamp to every user
			for _, j := range openws {
				j.SetWriteDeadline(time.Now().Add(time.Duration(time.Millisecond * 100)))
				err := j.WriteMessage(websocket.TextMessage, []byte(ss))
				if err != nil {
					fmt.Println(err)
					j.Close()
				}
			}

			fmt.Println(ss)
			time.Sleep(time.Second)
		}

		for _, j := range openws {
			j.SetWriteDeadline(time.Now().Add(time.Duration(time.Millisecond * 100)))
			j.WriteMessage(websocket.TextMessage, []byte("..."))
			j.Close()
		}

		// Unlock every websocket writter
		for _, j := range wwslock {
			j.Unlock()
		}

		startGame()

	}
}

func gameTick() {
	gTickLock.Lock()
	gTick++
	gTickLock.Unlock()

	gameLock.RLock()
	for i, j := range(pGameWS) {
		pGameLocks[i].Lock()
		j.SetWriteDeadline(time.Now().Add(time.Duration(time.Millisecond * 200)))
		j.WriteMessage(websocket.TextMessage, []byte(stringifyBoard(pGameWSid[i])))
		pGameLocks[i].Unlock()
	}
	gameLock.RUnlock()
	time.Sleep(time.Second)
	gameTick()
}
