package common

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data       interface{}
	Status     string
	StatusCode int
	Error      string
}

const SecretKey = "secret"

func ReturnResponse(w http.ResponseWriter, status string, statusCode int, error string, data interface{}) {
	var result Response
	result.Status = status
	result.StatusCode = statusCode
	result.Error = error
	result.Data = data
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(result)
	return
}
