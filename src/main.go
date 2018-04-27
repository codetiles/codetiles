package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	PORT := "8080"
	fmt.Println("Starting codetiles server on port " + PORT)

	users = make(map[[8]byte]user)

	http.HandleFunc("/", handleRoot)

	// API endpoints:
	http.HandleFunc("/api", handleGetVersion)
	http.HandleFunc("/api/v1", handleGetVersion)
	http.HandleFunc("/api/v1/createuser", handleJoiningUser)       // users.go
	http.HandleFunc("/api/v1/verifyuser/", handleVerifyUser)       // users.go
	http.HandleFunc("/api/v1/uploadcode", handleUploadCode)        // code.go
	http.HandleFunc("/api/v1/findgame", handleJoinLobby)           // lobby.go
	http.HandleFunc("/api/v1/game/players", handleRetrievePlayers) // game.go

	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Get path, make changes if needed
	filename := r.URL.Path

	if filename[len(filename)-1] == '/' {
		filename += "index.html"
	}
	// remove need to append .html to URL
	if filename == "/game" || filename == "/unsupported" {
		filename += ".html"
	}

	http.ServeFile(w, r, "public/"+filename)
}

func handleGetVersion(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "API Version 0.1 Pre-alpha")
}
