package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrResponseJSON struct {
	ErrMsg string `json:"error"`
}

type MsgResponseJSON struct {
	Msg string `json:"message"`
}

func RespondWithMessage(w http.ResponseWriter, receivedMsg string, status int) {
	msg := MsgResponseJSON{Msg: receivedMsg}
	if err := WriteJSON(w, msg, status); err != nil {
		log.Println("error writing message:", err)
	}
}

func RespondWithError(w http.ResponseWriter, errorMessage string, status int) {
	errMsg := ErrResponseJSON{ErrMsg: errorMessage}
	if err := WriteJSON(w, errMsg, status); err != nil {
		log.Println("error writing error message: ", err)
	}
}

func WriteJSON(w http.ResponseWriter, data interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}
