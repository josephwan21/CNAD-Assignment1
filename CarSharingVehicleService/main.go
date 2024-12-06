// main.go
package main

import (
	"Assg1/CarSharingVehicleService/package/db"
	"Assg1/CarSharingVehicleService/package/reservation"
	"Assg1/CarSharingVehicleService/package/vehicle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var jwtSecret = []byte("your-secret-key")

// HandleGetAvailableVehicles handles GET requests to retrieve available vehicles
func HandleGetAvailableVehicles(w http.ResponseWriter, r *http.Request) {
	dbConn := db.InitDB()
	defer dbConn.Close()

	vehicles, err := vehicle.GetAvailableVehicles(dbConn)
	if err != nil {
		http.Error(w, "Error fetching available vehicles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

// HandleCreateReservation handles POST requests for creating a reservation
func HandleCreateReservation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    int       `json:"user_id"`
		VehicleID int       `json:"vehicle_id"`
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

	// Check if vehicle is available
	var isAvailable bool
	err = dbConn.QueryRow("SELECT is_available FROM vehicles WHERE id = ?", req.VehicleID).Scan(&isAvailable)
	if err != nil || !isAvailable {
		http.Error(w, "Vehicle not available", http.StatusConflict)
		return
	}

	// Create reservation
	res, err := reservation.CreateReservation(dbConn, req.UserID, req.VehicleID, req.StartTime, req.EndTime)
	if err != nil {
		http.Error(w, "Error creating reservation", http.StatusInternalServerError)
		return
	}

	// Update vehicle availability
	_, err = dbConn.Exec("UPDATE vehicles SET is_available = FALSE WHERE id = ?", req.VehicleID)
	if err != nil {
		http.Error(w, "Error updating vehicle availability", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Example usage in a handler or service
func HandleGetUserReservations(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Strip the "Bearer " prefix from the token
	tokenString = tokenString[7:]

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract email from the JWT claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userid"] == nil {
		http.Error(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Claims: %v", claims)
	userID := claims["userid"].(float64)
	userIDInt := int(userID)
	fmt.Printf("User ID: %d", userIDInt)

	// Initialize database connection
	dbConn := db.InitDB()
	defer dbConn.Close()

	// Call GetReservationsByUserID
	reservations, err := reservation.GetReservationsByUserID(dbConn, userIDInt)
	if err != nil {
		log.Printf("Error retrieving reservations for user %d: %v", userIDInt, err)
		http.Error(w, "Error retrieving reservations", http.StatusInternalServerError)
		return
	}

	// Return the reservations as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

func UpdateReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract reservation ID from URL
	vars := mux.Vars(r)
	reservationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Decode the request body to get new start and end times
	var requestData struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert start and end times to time.Time
	startTime, err := time.Parse(time.RFC3339, requestData.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, requestData.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	// Call the UpdateReservation function to update the timings
	err = reservation.UpdateReservation(dbConn, reservationID, startTime, endTime)
	if err != nil {
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation updated successfully"})
}

// DeleteReservationHandler handles deleting a reservation
func DeleteReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract reservation ID from URL
	vars := mux.Vars(r)
	reservationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	// Call the DeleteReservation function to delete the reservation
	err = reservation.DeleteReservation(dbConn, reservationID)
	if err != nil {
		http.Error(w, "Failed to delete reservation", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation deleted successfully"})
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/vehicles", HandleGetAvailableVehicles)
	router.HandleFunc("/reserve", HandleCreateReservation)
	router.HandleFunc("/reservations", HandleGetUserReservations)
	router.HandleFunc("/reservations/{id}", UpdateReservationHandler).Methods("PUT")
	router.HandleFunc("/reservations/{id}", DeleteReservationHandler).Methods("DELETE")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)

	log.Println("Vehicle Service is running on port 8082")
	err := http.ListenAndServe(":8082", corsHandler(router))
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
