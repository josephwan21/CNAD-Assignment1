// main.go
package main

import (
	billing "Assg1/CarSharingBillingService/models"
	"Assg1/CarSharingBillingService/package/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// HandleCalculateBilling handles POST requests for calculating billing
func HandleCalculateBilling(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    int       `json:"user_id"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConn := db.GetDBConn()

	// Calculate billing
	totalAmount, discount, err := billing.CalculateBilling(dbConn, req.UserID, req.StartTime, req.EndTime)
	if err != nil {
		log.Printf("Error: %s", err)
		http.Error(w, "Error calculating billing", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		TotalAmount float64 `json:"total_amount"`
		Discount    float64 `json:"discount"`
	}{
		TotalAmount: totalAmount,
		Discount:    discount,
	}

	json.NewEncoder(w).Encode(response)
}

func CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {

	// Define a struct to bind the request body to
	var req struct {
		UserID        int       `json:"user_id"`
		VehicleID     int       `json:"vehicle_id"`
		ReservationID int       `json:"reservation_id"`
		StartTime     time.Time `json:"start_time"`
		EndTime       time.Time `json:"end_time"`
	}

	// Decode the incoming request body into the struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the database connection
	dbConn := db.GetDBConn()

	// Calculate the billing (total amount and discount) based on the reservation times
	totalAmount, discount, err := billing.CalculateBilling(dbConn, req.UserID, req.StartTime, req.EndTime)
	if err != nil {
		http.Error(w, "Error calculating billing", http.StatusInternalServerError)
		return
	}

	// Create the billing record in the database
	billingRecord, err := billing.CreateBillingRecord(dbConn, req.ReservationID, req.UserID, req.VehicleID, totalAmount, discount)
	if err != nil {
		fmt.Printf("Error creating invoice for UserID: %d, Total Amount: %v, Discount: %v, Error: %v\n",
			req.UserID, totalAmount, discount, err)
		http.Error(w, "Error creating billing record", http.StatusInternalServerError)
		return
	}

	// Set the response header to indicate JSON format
	w.Header().Set("Content-Type", "application/json")

	// Prepare the response struct with the billing details
	response := struct {
		ReservationID int `json:"reservation_id"`
		UserID        int `json:"user_id"`
		VehicleID     int `json:"vehicle_id"`

		TotalAmount float64 `json:"total_amount"`
		Discount    float64 `json:"discount"`
	}{
		ReservationID: billingRecord.ReservationID,
		UserID:        billingRecord.UserID,
		TotalAmount:   billingRecord.TotalAmount,
		Discount:      billingRecord.Discount,
	}

	// Send the response back to the client
	json.NewEncoder(w).Encode(response)
}

// HandleGetInvoicesByUser handles the API request to fetch all invoices for a user
func GetInvoicesByUserHandler(w http.ResponseWriter, r *http.Request) {
	// Decode user_id from query parameters
	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	dbConn := db.GetDBConn()
	invoices, err := billing.GetInvoicesByUser(dbConn, userID)
	if err != nil {
		fmt.Printf("Error: %s", err)
		http.Error(w, "Error fetching invoices", http.StatusInternalServerError)
		return
	}

	// Respond with the invoices as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

// DeleteInvoiceHandler handles the DELETE request for deleting an invoice by reservation ID
func DeleteInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Get the reservation ID from the URL
	reservationIDStr := r.URL.Query().Get("reservation_id")
	if reservationIDStr == "" {
		http.Error(w, "Reservation ID is required", http.StatusBadRequest)
		return
	}

	// Convert reservation ID from string to int
	reservationID, err := strconv.Atoi(reservationIDStr)
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Get the database connection (assuming you have a function to connect to the database)
	dbConn := db.GetDBConn()

	// Call the function to delete the invoice by reservation ID
	err = billing.DeleteInvoiceByReservationID(dbConn, reservationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting invoice: %v", err), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invoice deleted successfully"})
}

// Handler for updating the invoice
func updateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract reservation ID from URL or request body
	reservationIDStr := r.URL.Query().Get("reservation_id")
	if reservationIDStr == "" {
		http.Error(w, "Reservation ID is required", http.StatusBadRequest)
		return
	}
	reservationID, err := strconv.Atoi(reservationIDStr)
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	var invoiceData struct {
		UserID    int       `json:"user_id"`
		VehicleID int       `json:"vehicle_id"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}
	err = json.NewDecoder(r.Body).Decode(&invoiceData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConn := db.GetDBConn()

	// Calculate the billing (total amount and discount) based on the reservation times
	totalAmount, discount, err := billing.CalculateBilling(dbConn, invoiceData.UserID, invoiceData.StartTime, invoiceData.EndTime)
	if err != nil {
		http.Error(w, "Error calculating billing", http.StatusInternalServerError)
		return
	}

	// Update the invoice in the database
	err = billing.UpdateInvoice(dbConn, reservationID, invoiceData.UserID, invoiceData.VehicleID, totalAmount, discount)
	if err != nil {
		http.Error(w, "Failed to update invoice", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invoice updated successfully"})
}

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialise database: %v", err)
	}
	defer db.CloseDB()

	router := mux.NewRouter()
	router.HandleFunc("/calculatebilling", HandleCalculateBilling)
	router.HandleFunc("/createinvoice", CreateInvoiceHandler)
	router.HandleFunc("/invoices", GetInvoicesByUserHandler).Methods("GET")
	router.HandleFunc("/updateinvoice", updateInvoiceHandler).Methods("PUT")
	router.HandleFunc("/invoices", DeleteInvoiceHandler).Methods("DELETE")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)

	log.Println("Billing Service is running on port 8083")
	if err := http.ListenAndServe(":8083", corsHandler(router)); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
