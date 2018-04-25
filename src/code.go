package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Handle a person uploading their code
func handleUploadCode(w http.ResponseWriter, r *http.Request) {
	// Get user's JSON
	var userJson map[string]string

	err := json.NewDecoder(r.Body).Decode(&userJson)

	if err != nil {
		fmt.Println("Err unmarshalling JSON (client -> server)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Process id and code
	id, existsId := userJson["id"]
	rawCode, existsCode := userJson["code"]

	// If a value is missing, the request didn't meet the requirements
	if !existsId || !existsCode {
		rj, _ := json.Marshal(map[string]string{
			"Success": "false",
			"Reason":  "JSON is missing essential key value.",
		})

		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, string(rj))
		return
	}

	// Process the user's code for errors
	ver, errstr := verifyCode(rawCode)

	// If an error exists, don't change the code. Send them a reason as well.
	if !ver {
		w.WriteHeader(http.StatusNotAcceptable)

		rj, _ := json.Marshal(map[string]string{
			"Success": "false",
			"Reason":  errstr,
		})

		io.WriteString(w, string(rj))
		return
	}

	// Write to user's code
	var uid [8]byte
	copy(uid[:], []byte(id))
	usersArrayLock.RLock()
	_, ex := users[uid]
	usersArrayLock.RUnlock()

	// Check if user exists
	if !ex {
		rj, _ := json.Marshal(map[string]string{
			"Success": "false",
			"Reason":  "Id not found",
		})

		io.WriteString(w, string(rj))
		return
	}

	usersArrayLock.Lock()
	user := users[uid]
	user.code = rawCode
	usersArrayLock.Unlock()

	rj, _ := json.Marshal(map[string]string{
		"Success": "true",
		"Reason":  "",
	})

	io.WriteString(w, string(rj))
}

// Fuction for verifing code (confirming validity)
func verifyCode(code string) (bool, string) {
	return true, "None"
}
