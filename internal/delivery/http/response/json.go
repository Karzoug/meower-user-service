package response

import (
	"net/http"

	"github.com/go-json-experiment/json"
)

func JSON(w http.ResponseWriter, statusCode int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.MarshalWrite(w, body)
}
