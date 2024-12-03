// reservation.go
package reservation

import (
	"database/sql"
	"time"
)

// Reservation represents a user's vehicle reservation
type Reservation struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VehicleID int       `json:"vehicle_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateReservation creates a new reservation
func CreateReservation(db *sql.DB, userID, vehicleID int, startTime, endTime time.Time) (Reservation, error) {
	query := "INSERT INTO reservations (user_id, vehicle_id, start_time, end_time, status) VALUES (?, ?, ?, ?, 'Active')"
	result, err := db.Exec(query, userID, vehicleID, startTime, endTime)
	if err != nil {
		return Reservation{}, err
	}

	reservationID, err := result.LastInsertId()
	if err != nil {
		return Reservation{}, err
	}

	return Reservation{
		ID:        int(reservationID),
		UserID:    userID,
		VehicleID: vehicleID,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    "Active",
		CreatedAt: time.Now(),
	}, nil
}
