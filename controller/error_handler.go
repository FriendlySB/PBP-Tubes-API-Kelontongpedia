package controller

import (
	"encoding/json"
	"net/http"
	"week6/model"
)

func sendErrorResponse(w http.ResponseWriter, message string) {
	var response model.ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
