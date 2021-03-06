package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// Handle a person uploading their code
func handleUploadCode(w http.ResponseWriter, r *http.Request) {
	// Get user's JSON
	var userJson map[string]string

	err := json.NewDecoder(r.Body).Decode(&userJson)

	if handleJsonUnmarshalError(w, r, "code.go - upload", err) {
		return
	}

	// Process id and code
	id, existsId := userJson["id"]
	rawCode, existsCode := userJson["code"]

	// If a value is missing, the request didn't meet the requirements
	if !existsId || !existsCode {
		rj, err := json.Marshal(map[string]string{
			"Success": "false",
			"Reason":  "JSON is missing essential key value.",
		})

		if handleJsonMarshalError(w, r, "code.go - missing vals", err) {
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, string(rj))
		return
	}

	// Process the user's code for errors
	ver, errstr := verifyCode(rawCode)

	// If an error exists, don't change the code. Send them a reason as well.
	if !ver {
		w.WriteHeader(http.StatusNotAcceptable)

		rj, err := json.Marshal(map[string]string{
			"Success": "false",
			"Reason":  errstr,
		})

		if handleJsonMarshalError(w, r, "code.go - err'd code", err) {
			return
		}

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
	users[uid] = user
	usersArrayLock.Unlock()

	rj, err := json.Marshal(map[string]string{
		"Success": "true",
		"Reason":  "",
	})

	if handleJsonMarshalError(w, r, "code.go - userid", err) {
		return
	}

	io.WriteString(w, string(rj))
}

// Handle a person retrieving their code
func handleDownloadCode(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")
	var uid [8]byte
	copy(uid[:], []byte(id))

	e, _, _, _ := checkUserId(uid)
	if !e {
		io.WriteString(w, "User does not exist")
		return
	}

	usersArrayLock.RLock()
	io.WriteString(w, users[uid].code)
	usersArrayLock.RUnlock()
}

// Fuction for verifing code (confirming validity)
func verifyCode(code string) (bool, string) {
	return true, "None"
}
