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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, existsId := userJson["id"]
	rawCode, existsCode := userJson["code"]

	if !existsId || !existsCode {
		rj, err := json.Marshal(map[string]string{
			"Success": "false",
			"Reason" : "JSON is missing essential key value.",
		})

		if err != nil {
			fmt.Println("Err marshalling JSON (server -> client)")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(rj))
		return
	}

	ver, errstr := verifyCode(rawCode)

	if !ver {
		w.WriteHeader(http.StatusNotAcceptable)
		rj, err := json.Marshal(map[string]string{
			"Success" : "false",
			"Reason" : errstr,
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


func verifyCode(code string) (bool, string) {
	return true, "None"
}
