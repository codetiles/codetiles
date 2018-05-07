package main

import (
  "net/http"
  "sync"
  "fmt"

  "github.com/gorilla/websocket"
)

// Pointers to websockets and also write-locks
var pGameWS    []*websocket.Conn
var pGameLocks []*sync.Mutex
// Array that associates those with id's (multiple games)
var pGameWSid  [][8]byte
// Track when the user last requested their board
var pGameTickN []*int
// Lock for above arrays
var gameLock   sync.RWMutex

// Websocket handler to send the game board to a user. Requires heavy use of
// tick's for updating game board correctly.
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
  }

  if !inGame {
    ws.WriteMessage(websocket.TextMessage, []byte("User is not in a game"))
  }

  wrL := new(sync.Mutex)
  tickN := new(int)

  // Add both lock and websocket to the game-pointer arrays
  gameLock.Lock()
  pGameWS = append(pGameWS, ws)
  pGameLocks = append(pGameLocks, wrL)
  pGameTickN = append(pGameTickN, tickN)
  fmt.Println(pGameWS, pGameLocks, tickN)
  gameLock.Unlock()

  // When function returns remove both lock and websocket pointers
  defer func(){
    gameLock.Lock()

    var wsArr []*websocket.Conn
    for _, j := range(pGameWS) {
      if j != ws {
        wsArr = append(wsArr, ws)
      }
    }
    pGameWS = wsArr

    var lArr []*sync.Mutex
    for _, j := range(pGameLocks) {
      if j != wrL {
        lArr = append(lArr, wrL)
      }
    }
    pGameLocks = lArr

    var iArr []*int
    for _, j := range(pGameTickN) {
      if j != tickN {
        iArr = append(iArr, j)
      }
    }
    pGameTickN = iArr

    gameLock.Unlock()
  }()



}
