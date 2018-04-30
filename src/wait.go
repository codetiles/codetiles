package main

import (
  "net/http"
  "fmt"

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
  defer ws.Close()
  defer fmt.Println("Socket closed")

  if err != nil {
    fmt.Println("Error creating websocket in wait.go")
    fmt.Println("Error:", err)
    fmt.Println("Socket closed")
    return
  }

  // The user must send an id that they want to put in a lobby
  _, message, err := ws.ReadMessage()

  if err != nil {
    return
  }

  var uid [8]byte
  copy(uid[:], []byte(message))


  // If a player is queued or doesn't exist, send an error and close connection
  exists, _, queued, _ := checkUserId(uid)
  if !exists || queued {
    w, err := ws.NextWriter(websocket.TextMessage)

    if err != nil {
      return
    }

    if !exists {
      w.Write([]byte("Error: Player doesn't exist"))
    }
    if queued {
      w.Write([]byte("Error: Player is already in a queue"))
    }

    // We don't have to worry about an error because we are closing the socket
    w.Close()
    return
  }

  // Add a player to the queue
  queuedPlayersLock.Lock()
  queuedPlayers = append(queuedPlayers, uid)
  queuedPlayersLock.Unlock()

  // Initial update
  searchtick <- 1
  ww, err := ws.NextWriter(websocket.TextMessage)
  if err != nil {
    return
  }
  ww.Write([]byte(getNumberOfPlayersInQueue()))
  ww.Close()

  // Send the user the # of users online
  for {
    select {
    case <-searchtick:
      w, err := ws.NextWriter(websocket.TextMessage)
      if err != nil {
        return
      }
      w.Write([]byte(getNumberOfPlayersInQueue()))
      w.Close()
    }
  }

}
