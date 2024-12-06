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
	Make      string    `json:"make"`
	Model     string    `json:"model"`
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

// GetReservationsByUserID retrieves all reservations made by a specific user
func GetReservationsByUserID(db *sql.DB, userID int) ([]Reservation, error) {
	// Prepare the query to retrieve reservations for the given user

	query := "SELECT r.id, r.user_id, r.vehicle_id, v.make, v.model, r.start_time, r.end_time, r.status, r.created_at FROM reservations r INNER JOIN vehicles v ON r.vehicle_id = v.id WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation

	// Iterate over the rows returned by the query
	for rows.Next() {
		var reservation Reservation
		// Scan each row into the reservation struct
		err := rows.Scan(&reservation.ID, &reservation.UserID, &reservation.VehicleID, &reservation.Make, &reservation.Model, &reservation.StartTime, &reservation.EndTime, &reservation.Status, &reservation.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Append the reservation to the slice
		reservations = append(reservations, reservation)
	}

	// Check for errors after iterating through the rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the slice of reservations
	return reservations, nil
}
