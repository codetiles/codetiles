package main

import (
  "net/http"
  "fmt"
  "strconv"

  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// Create a websocket and send a user information about the game they are
// waititng for
func handleWaitForGame(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Socket opened")
  ws, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    fmt.Println("Error creating websocket in wait.go")
    fmt.Println("Error:", err)
    fmt.Println("Socket closed")
    return
  }

  func () {
    for {
      select {
      case <-searchtick:
        w, err := ws.NextWriter(websocket.TextMessage)
        if err != nil {
          fmt.Println("Socket closed")
          return
        }
        w.Write([]byte(getPlayersInQueue()))
      }
    }
  }()
}

func getPlayersInQueue() string {
  queuedPlayersLock.RLock()
  queued := len(queuedPlayers)
  queuedPlayersLock.RUnlock()

  qString := strconv.Itoa(queued)
  return qString + " players waiting in queue"
}
