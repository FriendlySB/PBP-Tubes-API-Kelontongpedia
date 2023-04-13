package controller

import (
	"PBP-Tubes-API-Tokopedia/model"
	"encoding/json"
	"log"
	"net/http"
)

func GetItemReview(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentID := getUserIdFromCookie(r)

	if currentID == -1 {
		sendUnauthorizedResponse(w)
	} else {
		itemId := r.URL.Query()["itemId"]
		query := "SELECT r.reviewId, u.userId, r.review_date, r.rating, r.review "
		query += "FROM item i "
		query += "INNER JOIN review r ON r.itemId=i.itemId "
		query += "INNER JOIN users u ON u.userid=r.userId "
		query += "WHERE i.itemId=" + itemId[0]

		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
			sendErrorResponse(w, "Error")
			return
		}
		var review model.Review
		var reviews []model.Review
		for rows.Next() {
			if err := rows.Scan(&review.ID, &review.UserId, &review.ReviewDate, &review.Rating, &review.Review); err != nil {
				sendErrorResponse(w, "Error while scanning rows")
				return
			} else {
				reviews = append(reviews, review)
			}
		}
		var response model.GenericResponse
		w.Header().Set("Content=Type", "application/json")
		if err == nil {
			response.Status = 200
			response.Message = "Success"
			response.Data = reviews
			json.NewEncoder(w).Encode(response)
		} else {
			response.Status = 400
			response.Message = "Error"
			json.NewEncoder(w).Encode(response)
		}
	}
}
