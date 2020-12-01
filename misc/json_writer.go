package misc

import (
	"encoding/json"
	"net/http"
)

// WriteResponse writes the map[string]interface{} which includes boolean error signal and message itself
func WriteResponse(isError bool, message interface{}) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	result["Error"] = isError
	result["Message"] = message
	return result
}

// JSONWrite writes to "response body" info given by "data"
func JSONWrite(w http.ResponseWriter, data interface{}, status int) {
	dataBytes, _ := json.Marshal(data);
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataBytes)
	w.WriteHeader(status)
}