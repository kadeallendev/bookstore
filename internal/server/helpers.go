package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// Error response
type ErrResponseJSON struct {
	ErrMsg string `json:"error"`
}

// Message response
type MsgResponseJSON struct {
	Msg string `json:"message"`
}

// Respond with a message encoded as JSON
func RespondWithMessage(w http.ResponseWriter, receivedMsg string, status int) {
	msg := MsgResponseJSON{Msg: receivedMsg}
	if err := WriteJSON(w, msg, status); err != nil {
		log.Println("error writing message:", err)
	}
}

// Respond with an error message encoded as json
func RespondWithError(w http.ResponseWriter, errorMessage string, status int) {
	errMsg := ErrResponseJSON{ErrMsg: errorMessage}
	if err := WriteJSON(w, errMsg, status); err != nil {
		log.Println("error writing error message: ", err)
	}
}

// Respond with an error message with the http status being InternalServerError
func RespondWithServerError(w http.ResponseWriter) {
	RespondWithError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Encode the given data as JSON and write it to the response
func WriteJSON(w http.ResponseWriter, data interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}
