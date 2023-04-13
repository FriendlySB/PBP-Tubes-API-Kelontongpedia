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
func ReviewItem(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentID := getUserIdFromCookie(r)

	if currentID == -1 {
		sendUnauthorizedResponse(w)
	} else {
		err := r.ParseForm()
		if err != nil {
			sendErrorResponse(w, "Error while parsing form")
			return
		}
		itemId := r.Form.Get("itemId")
		rating := r.Form.Get("rating")
		review := r.Form.Get("review")

		var check bool
		row := db.QueryRow("SELECT COUNT(*) FROM users u INNER JOIN transaction t ON t.userId=u.userid INNER JOIN transaction_detail td ON td.transactionId=t.transactionId INNER JOIN item i ON i.itemId=td.itemId WHERE u.userid=? AND i.itemId=? AND t.progress='Done'", currentID, itemId)
		err = row.Scan(&check)
		if err != nil {
			sendErrorResponse(w, "Error while checking purchase history")
			return
		}

		if !check {
			sendErrorResponse(w, "You haven't bought this item")
			return
		}

		_, errQuery := db.Exec("INSERT INTO review(itemID, userID, rating, review) VALUES(?,?,?,?)", itemId, currentID, rating, review)
		
		if errQuery == nil {
			rows, err := db.Query("SELECT reviewId, userId, review_date, rating, review FROM review WHERE reviewId = LAST_INSERT_ID()")
			if err != nil {
				sendErrorResponse(w, "Error while fetching updated data")
				return
			}
			defer rows.Close()

			var review model.Review
			for rows.Next() {
				if err := rows.Scan(&review.ID, &review.UserId, &review.ReviewDate, &review.Rating, &review.Review); err != nil {
					sendErrorResponse(w, "Error while scanning rows")
					return
				}
			}
			response := model.GenericResponse{Status: 200, Message: "Success", Data: review}
			json.NewEncoder(w).Encode(response)
		} else {
			response := model.GenericResponse{Status: 400, Message: "Error"}
			json.NewEncoder(w).Encode(response)
		}
	}
}