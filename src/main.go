package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var lockGameBoards sync.Mutex
var games []gameBoard
var users map[[6]byte]user

func main() {
	PORT := "8080"
	fmt.Println("Starting codetiles server on port " + PORT)

	users = make(map[[6]byte]user)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/api/joinlobby", handleJoiningUser)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Get path, if it ends in a '/', then we should append "index.html" to it.
	fileName := r.URL.Path
	if fileName[len(fileName)-1] == '/' {
		fileName += "index.html"
	}

	http.ServeFile(w, r, "public/"+fileName)
	return
}
