// vehicle.go
package vehicle

import (
	"database/sql"
	"log"
	"time"
)

// Vehicle represents a car available for reservation
type Vehicle struct {
	ID           int       `json:"id"`
	Make         string    `json:"make"`
	Model        string    `json:"model"`
	LicensePlate string    `json:"license_plate"`
	IsAvailable  bool      `json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetAvailableVehicles fetches all available vehicles
func GetAvailableVehicles(db *sql.DB) ([]Vehicle, error) {
	query := "SELECT id, make, model, license_plate, is_available FROM vehicles WHERE is_available = TRUE"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var vehicle Vehicle
		if err := rows.Scan(&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.LicensePlate, &vehicle.IsAvailable); err != nil {
			log.Println("Error scanning vehicle:", err)
			continue
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}
