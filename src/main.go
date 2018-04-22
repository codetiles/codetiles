package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	PORT := "8080"
	fmt.Println("Starting codetiles server on port " + PORT)

	users = make(map[[8]byte]user)
	queuedPlayers = make(map[[8]byte]bool)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/api/v1/createuser", handleJoiningUser)
	http.HandleFunc("/api/v1/verifyuser/", handleVerifyUser)
	http.HandleFunc("/api/v1/game/players", handleRetrievePlayers)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Get path, if it ends in a '/', then we should append "index.html" to it.
	filename := r.URL.Path
	if filename[len(filename)-1] == '/' {
		filename += "index.html"
	}
	// remove need to append .html to URL
	if filename == "/game" {
		filename += ".html"
	}
	if filename == "/unsupported" {
		filename += ".html"
	}

	http.ServeFile(w, r, "public/"+filename)
	return
}
