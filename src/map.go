package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
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
	emptyTile.tileType = "n"
	emptyTile.value = 0
	emptyTile.owner = "/"

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

	// Make all users on the map be registered as being in a game
	usersArrayLock.Lock()
	for _, uid := range players {
		user := users[uid]
		user.inGame = true
		users[uid] = user
	}
	usersArrayLock.Unlock()

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

	// Clear the player queue, put every player's inGame variable to true
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

// Take a user id and get the board for them (includes every tile)
func stringifyBoard(uid [8]byte) string {
	var tilemap [30][30]tile

	// Populate the tilemap with the raw data from the map
	lockGameBoards.RLock()
	for _, board := range games {
		for _, player := range board.players {
			if player == uid {
				tilemap = board.tiles
			}
		}
	}
	lockGameBoards.RUnlock()

	if tilemap == [30][30]tile{} {
		return "User not in game"
	}

	// Create and populate a string with the tile map
	var encodedMap string

	for _, cols := range tilemap {
		for _, tile := range cols {
			encodedMap += tile.owner
			sVal := strconv.Itoa(tile.value)
			if len(sVal) == 1 {
				sVal = "0" + sVal
			}
			encodedMap += sVal
		}
		encodedMap += "\n"
	}

	return encodedMap
}
