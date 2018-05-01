package main

import (
	"time"
)

// This file contains 2 types of functions:
// Functions that call themselves and update a channel
// Functions that are called as goroutines to change a channel

// use <-searchtick in select to check when a user joins the queue
var searchtick chan int
var gametick chan int
var countdown chan string

func tickUser() {
	searchtick <- 0
}

func performGameTick() {
	time.Sleep(250 * time.Millisecond)
	gametick <- 0
	performGameTick()
}

func checkCountdown() {
	queuedPlayersLock.RLock()
	lqp := len(queuedPlayers)
	queuedPlayersLock.RUnlock()

	// Once the # of queued players reaches 2, we start a countdown for a game
	if lqp < 2 {
		return
	}

	for {

	}
}
