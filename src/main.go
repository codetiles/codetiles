package main

import (
  "fmt"
  "io"
  "net/http"
  "log"
  "os"
)

func main() {
  PORT := "8080"

  fmt.Println("Starting codetiles server on port " + PORT)

  http.HandleFunc("/", handleRoot)
  log.Fatal(http.ListenAndServe(":" + PORT, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
  // Get path, if it ends in a '/', then we should append "index.html" to it.
  fileName := r.URL.Path
  if fileName[len(fileName) - 1] == '/' {
    fileName += "index.html"
  }

  http.ServeFile(w, r, "public/" + fileName)
  return
}
