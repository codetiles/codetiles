package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// use <-searchtick in select to check when a user joins the queue
var searchtick chan int
var autosearchtick chan int

func tickUser() {
	searchtick <- 0
}

func performSearchTick() {
	time.Sleep(250 * time.Millisecond)
	autosearchtick <- 0
	performSearchTick()
}

func checkCountdown() {
	defer func() {
		time.Sleep(time.Millisecond * 100)
		checkCountdown()
	}()

	queuedPlayersLock.RLock()
	lenQ := len(queuedPlayers)
	queuedPlayersLock.RUnlock()

	if lenQ > 1 {
		wslock.RLock()
		defer wslock.RUnlock()

		for _, j := range wwslock {
			fmt.Println(j)
			j.Lock()
		}

		for i := 5; i >= 1; i-- {

			ss := "Game starts in " + strconv.Itoa(i)
			if i > 1 {
				ss += " seconds"
			} else {
				ss += " second"
			}

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

		for _, j := range wwslock {
			fmt.Println(j)
			j.Unlock()
		}

		fmt.Println(queuedPlayers)

		queuedPlayersLock.Lock()
		queuedPlayers = [][8]byte{}
		queuedPlayersLock.Unlock()
	}
}
