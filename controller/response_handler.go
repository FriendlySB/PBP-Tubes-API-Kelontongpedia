package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, message string) {
	var response model.ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendSuccessResponse(w http.ResponseWriter, message string, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var response model.GenericResponse
	response.Status = 200
	response.Message = message
	response.Data = value
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
