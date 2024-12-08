package rental_history

import (
	"database/sql"
	"time"
)

// RentalHistoryEntry represents a single rental history record
type RentalHistoryEntry struct {
	ID            int       `json:"id"`
	ReservationID int       `json:"reservation_id"`
	UserID        int       `json:"user_id"`
	VehicleID     int       `json:"vehicle_id"`
	Make          string    `json:"make"`
	Model         string    `json:"model"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	RentalStatus  string    `json:"rental_status"`
	TotalAmount   float64   `json:"total_amount"`
}

// AddRentalHistoryEntry inserts a new record into the rental history table
func AddRentalHistoryEntry(db *sql.DB, history RentalHistoryEntry) error {
	query := `INSERT INTO Rentals (reservation_id, user_id, vehicle_id, start_time, end_time, rental_status) 
              VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, history.ReservationID, history.UserID, history.VehicleID, history.StartTime, history.EndTime, history.RentalStatus)
	return err
}

// GetRentalHistory fetches all rental history records for a specific user
func GetRentalHistory(db *sql.DB, userID int) ([]RentalHistoryEntry, error) {
	query := `SELECT r.id, r.reservation_id, r.user_id, r.vehicle_id, v.make, v.model, r.start_time, r.end_time, r.rental_status, b.total_amount
              FROM Rentals r INNER JOIN CarSharingVehicleService.Vehicles v ON r.vehicle_id = v.id INNER JOIN CarSharingBillingService.Billing b ON r.reservation_id = b.reservation_id WHERE r.user_id = ?`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []RentalHistoryEntry
	for rows.Next() {
		var entry RentalHistoryEntry
		err := rows.Scan(&entry.ID, &entry.ReservationID, &entry.UserID, &entry.VehicleID, &entry.Make, &entry.Model, &entry.StartTime, &entry.EndTime, &entry.RentalStatus, &entry.TotalAmount)
		if err != nil {
			return nil, err
		}
		history = append(history, entry)
	}

	return history, nil
}
