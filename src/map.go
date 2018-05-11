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

// Stores the last tick of the board for comparison
var lockLastTickBoards sync.RWMutex
var lastTickBoards [][30][30]tile

type tile struct {
	tileType string
	value    int
	owner    string
}

type gameBoard struct {
	tiles   [30][30]tile
	players [][8]byte
	colors  []string
	id      [8]byte
	// Boards need a tick offset to figure out when tiles should grow
	tOffset int
}

// Possible constant colors
var colors []string = []string{"b", "r", "g", "y"}

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

	var emptyMap [30][30]tile
	copy(emptyMap[:], tiles[:])

	// Create the map
	var newMap gameBoard

	newMap.players = players

	// Assign each player to a color
	for i := range players {
		// Limit it to 2 players
		if i > 2 {
			break
		}
		newMap.colors = append(newMap.colors, colors[i])
	}

	// Place a tile randomly on the map for that player's color
	// (verify another player main tile isn't there too)
	for _, j := range newMap.colors {
		for {
			var b [2]byte
			rand.Read(b[:])
			rTileX := b[0] % 30
			rTileY := b[1] % 30
			if tiles[rTileX][rTileY].owner == "/" {
				tiles[rTileX][rTileY].owner = j
				tiles[rTileX][rTileY].value = 2
				break
			}
		}
	}

	newMap.tiles = tiles

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
	lockLastTickBoards.Lock()

	games = append(games, newMap)
	lastTickBoards = append(lastTickBoards, emptyMap)

	lockLastTickBoards.Unlock()
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
	defer lockGameBoards.RUnlock()
	for _, board := range games {
		for _, player := range board.players {
			if player == uid {
				tilemap = board.tiles
			}
		}
	}

	if tilemap == [30][30]tile{} {
		return "User not in game"
	}

	// Create and populate a string with the tile map
	var encodedMap string

	for _, cols := range tilemap {
		for _, tile := range cols {
			encodedMap += tile.owner
			encodedMap += twoDigitITOA(tile.value)
		}
		encodedMap += "\n"
	}

	return encodedMap
}

func twoDigitITOA(val int) string {
	ret := strconv.Itoa(val)
	if len(ret) == 1 {
		ret = "0" + ret
	}

	return ret
}

func stringifyBoardDifference(uid [8]byte) string {
	var tilemap [30][30]tile
	var lTilemap [30][30]tile

	// Populate the tilemap with the raw data from the map
	lockGameBoards.RLock()
	lockLastTickBoards.RLock()
	defer lockLastTickBoards.RUnlock()
	defer lockGameBoards.RUnlock()
	for i, board := range games {
		for _, player := range board.players {
			if player == uid {
				tilemap = board.tiles
				lTilemap = lastTickBoards[i]
			}
		}
	}

	if tilemap == [30][30]tile{} {
		return "User not in game"
	}

	var returnString string

	for i := range tilemap {
		for j := range tilemap[i] {
			if tilemap[i][j] != lTilemap[i][j] {
				tile := tilemap[i][j]
				returnString += twoDigitITOA(i)
				returnString += twoDigitITOA(j)
				returnString += tile.owner
				returnString += twoDigitITOA(tile.value)
			}

		}
	}

	return returnString
}
