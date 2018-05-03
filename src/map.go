package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
)

var lockGameBoards sync.RWMutex
var games []gameBoard

type tile struct {
	tileType string
	value    int
	owner    string
}

type gameBoard struct {
	tiles   [30][30]tile
	players [][8]byte
	id      [8]byte
	// Boards need a tick offset to figure out when tiles should grow
	tOffset int
}

func formMap(players [][8]byte, tOffset int) [8]byte {
	// Create a sample tile
	var emptyTile tile
	emptyTile.tileType = "/"
	emptyTile.value = 0
	emptyTile.owner = "none"

	// Populate a map of tiles
	var tiles [30][30]tile
	for i := 0; i < 30; i++ {
		for j := 0; j < 30; j++ {
			tiles[i][j] = emptyTile
		}
	}

	// Create the map
	var newMap gameBoard

	newMap.tiles = tiles
	newMap.players = players

	// Store a random base64 array as the id
	var randomByteArray [6]byte
	rand.Read(randomByteArray[:])
	rb64 := base64.StdEncoding.EncodeToString(randomByteArray[:])
	var randomId [8]byte
	copy(randomId[:], []byte(rb64))
	newMap.id = randomId
	newMap.tOffset = tOffset

	lockGameBoards.Lock()
	games = append(games, newMap)
	lockGameBoards.Unlock()

	return newMap.id
}

// Create a map and clear the queued players
func startGame() {

	// Create a map with a proper tick offset
	gTickLock.RLock()
	formMap(queuedPlayers, gTick)
	gTickLock.RUnlock()

	// Clear the player queue
	queuedPlayersLock.Lock()
	queuedPlayers = [][8]byte{}
	queuedPlayersLock.Unlock()

	// List # active games
	lockGameBoards.RLock()
	fmt.Println("Starting game...")
	fmt.Println(len(games), "active game(s)")
	fmt.Println()
	lockGameBoards.RUnlock()

}
