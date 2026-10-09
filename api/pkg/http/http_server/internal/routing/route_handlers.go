package routing

import (
	"encoding/json"
	"log"
	"net/http"

	"frisboo-bank/openapi-generator-service/pkg/http/http_server/constants"
)

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding the JSON response: %v", err)
	}
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSON)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}
