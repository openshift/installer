package server

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// writeJSON writes the given bytes with a JSON Content-Type.
func writeJSON(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(data)
	if err != nil {
		log.Errorf("failed to write JSON data: %v", err)
	}
}

// writeJSONData marshals the given data into json and writes it with the JSON
// Content-Type
func writeJSONData(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("failed to unmarshal JSON data: %v", err)
		return
	}
	writeJSON(w, data)
}
