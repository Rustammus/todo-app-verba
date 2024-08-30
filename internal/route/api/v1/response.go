package v1

import (
	"encoding/json"
	"net/http"
)

type errorJSON struct {
	Err string `json:"error"`
}

func writeResponse(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data == nil {
		return
	}

	respJSON, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(respJSON)
	//TODO handle err
}

func writeResponseErr(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err == nil {
		return
	}

	respJSON, err := json.Marshal(errorJSON{err.Error()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(respJSON)
	//TODO handle err
}
