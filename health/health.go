package health

import (
	"encoding/json"
	"net/http"
)

// Health runs the health action.
func Health(rw http.ResponseWriter, req *http.Request) {
	response, _ := json.Marshal("Ok")

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}
