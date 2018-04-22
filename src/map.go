package main

import (
	"sync"
)

var lockGameBoards sync.RWMutex
var games []gameBoard

type tile struct {
	tileType string
	value    int
	owner    string
	x        int
	y        int
}

type gameBoard struct {
	tiles   [30][30]tile
	players [][8]byte
	id      [12]byte
}

func formMap(players [][8]byte) gameBoard {
	// Create a sample tile
	var emptyTile tile
	emptyTile.tileType = "empty"
	emptyTile.value = 0
	emptyTile.owner = "none"

	// Populate a map of tiles
	var tiles [30][30]tile
	for i := 0; i < 30; i++ {
		for j := 0; j < 30; j++ {
			emptyTile.x = i
			emptyTile.y = j
			tiles[i][j] = emptyTile
		}
	}

	var newMap gameBoard
	newMap.tiles = tiles
	newMap.players = players
	return newMap
}
