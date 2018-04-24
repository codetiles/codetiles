package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Handle a person uploading their code
func handleUploadCode(w http.ResponseWriter, r *http.Request) {
	var userJson map[string]string
	err := json.NewDecoder(r.Body).Decode(&userJson)

	if err != nil {
		fmt.Println("Err unmarshalling JSON (client -> server)")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	id, existsId := userJson["id"]
	rawCode, existsCode := userJson["code"]

	if !existsId || !existsCode {
		rj, err := json.Marshal(map[string]bool{
			"Success": false,
		})

		if err != nil {
			fmt.Println("Err marshalling JSON (server -> client)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(rj))
		return
	}

	fmt.Println(id, rawCode)

}
