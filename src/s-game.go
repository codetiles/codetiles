package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Pointers to websockets and also write-locks
var pGameWS []*websocket.Conn
var pGameLocks []*sync.Mutex

// Array that associates those with id's (multiple games)
var pGameWSid [][8]byte

// Track when the user last requested their board
var pGameTickN []*int

// Lock for above arrays
var gameLock sync.RWMutex

// WSHandleGameBoard is a Websocket handler to send the game board to a user.
// It requires heavy use of tick's for updating game board correctly.
func WSHandleGameBoard(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	defer ws.Close()
	defer ws.WriteMessage(websocket.CloseMessage, []byte{})

	if err != nil {
		fmt.Println("Error upgrading to websocket (s-game.go)")
		return
	}

	// Once the websocket is open, get the user's id before doing anything else.
	_, message, err := ws.ReadMessage()

	if err != nil {
		return
	}

	var uid [8]byte
	copy(uid[:], []byte(message))

	exists, inGame, _, _ := checkUserId(uid)
	if !exists {
		ws.WriteMessage(websocket.TextMessage, []byte("User does not exist"))
		return
	}

	if !inGame {
		ws.WriteMessage(websocket.TextMessage, []byte("User is not in a game"))
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte(stringifyBoard(uid)))

	wrL := new(sync.Mutex)
	tickN := new(int)

	ws.WriteMessage(websocket.TextMessage, []byte(stringifyBoard(uid)))

	// When function returns remove both lock and websocket pointers
	defer func() {
		gameLock.Lock()

		var wsArr []*websocket.Conn
		for _, j := range pGameWS {
			if j != ws {
				wsArr = append(wsArr, j)
			}
		}
		pGameWS = wsArr

		var lArr []*sync.Mutex
		for _, j := range pGameLocks {
			if j != wrL {
				lArr = append(lArr, j)
			}
		}
		pGameLocks = lArr

		var iArr []*int
		for _, j := range pGameTickN {
			if j != tickN {
				iArr = append(iArr, j)
			}
		}
		pGameTickN = iArr

		gameLock.Unlock()
	}()

	// Add both lock and websocket to the game-pointer arrays
	gameLock.Lock()
	pGameWS = append(pGameWS, ws)
	pGameLocks = append(pGameLocks, wrL)
	pGameTickN = append(pGameTickN, tickN)
	pGameWSid = append(pGameWSid, uid)
	fmt.Println(len(pGameWS), "game socket(s) open.")
	gameLock.Unlock()

	go func() {
		for {
			ws.SetWriteDeadline(time.Now().Add(time.Duration(time.Second * 2)))
			_, m, err := ws.ReadMessage()
			if err != nil {
				return
			}
			wrL.Lock()
			if string(m) != getTickStr() {
				ws.WriteMessage(websocket.TextMessage, []byte(stringifyBoard(uid)))
			} else {
				ws.WriteMessage(websocket.TextMessage, []byte("--pong--"))
			}

			wrL.Unlock()

		}
	}()

	time.Sleep(time.Second * 10)

}
