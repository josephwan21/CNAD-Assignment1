// main.go
package main

import (
	"Assg1/CarSharingVehicleService/package/db"
	"Assg1/CarSharingVehicleService/package/reservation"
	"Assg1/CarSharingVehicleService/package/vehicle"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

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

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/vehicles", HandleGetAvailableVehicles)
	router.HandleFunc("/reserve", HandleCreateReservation)

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
