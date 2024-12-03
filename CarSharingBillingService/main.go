// main.go
package main

import (
	billing "Assg1/CarSharingBillingService/models"
	"Assg1/CarSharingBillingService/package/db"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// HandleCalculateBilling handles POST requests for calculating billing
func HandleCalculateBilling(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RentalID  int       `json:"rental_id"`
		UserID    int       `json:"user_id"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	// Calculate billing
	totalAmount, discount, err := billing.CalculateBilling(dbConn, req.RentalID, req.UserID, req.StartTime, req.EndTime)
	if err != nil {
		http.Error(w, "Error calculating billing", http.StatusInternalServerError)
		return
	}

	// Create billing record
	bill, err := billing.CreateBillingRecord(dbConn, req.RentalID, req.UserID, totalAmount, discount)
	if err != nil {
		http.Error(w, "Error creating billing record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bill)
}

func main() {
	http.HandleFunc("/billing", HandleCalculateBilling)

	log.Println("Billing Service is running on port 8083...")
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
