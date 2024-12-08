// billing.go
package billing

import (
	"database/sql"
	"math"
	"time"
)

// Billing represents a billing record
type Billing struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	VehicleID     int       `json:"vehicle_id"`
	ReservationID int       `json:"reservation_id"`
	Make          string    `json:"make"`
	Model         string    `json:"model"`
	TotalAmount   float64   `json:"total_amount"`
	Discount      float64   `json:"discount"`
	CreatedAt     time.Time `json:"created_at"`
}

// CalculateBilling calculates the total billing amount based on rental duration, membership, and any applicable discounts
func CalculateBilling(db *sql.DB, userID int, startTime time.Time, endTime time.Time) (float64, float64, error) {
	// Get membership type of the user
	var membership string
	err := db.QueryRow("SELECT membership FROM carsharinguserservice.users WHERE id = ?", userID).Scan(&membership)
	if err != nil {
		return 0, 0, err
	}

	// Calculate rental duration based on hours
	duration := endTime.Sub(startTime).Hours()

	// Pricing model, in this case, rate per hour
	var ratePerHour float64
	if membership == "Premium" {
		ratePerHour = 10.00 // Premium rate
	} else if membership == "VIP" {
		ratePerHour = 8.00 // VIP rate
	} else {
		ratePerHour = 12.00 // Basic rate
	}

	// Calculate the total cost based on the duration and rate
	totalAmount := ratePerHour * duration

	// Discount (10% off for Premium members)
	var discount float64
	if membership == "Premium" {
		discount = totalAmount * 0.10 // 10% discount for Premium
	}

	// Apply discount to total amount
	finalAmount := totalAmount - discount

	// Round to 2 decimal places
	finalAmount = math.Round(finalAmount*100) / 100
	discount = math.Round(discount*100) / 100

	return finalAmount, discount, nil
}

// CreateBillingRecord creates a new billing record in the database
func CreateBillingRecord(db *sql.DB, reservationID int, userID int, vehicleID int, totalAmount float64, discount float64) (Billing, error) {
	query := "INSERT INTO billing (reservation_id, user_id, vehicle_id, total_amount, discount) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Exec(query, reservationID, userID, vehicleID, totalAmount, discount)
	if err != nil {
		return Billing{}, err
	}

	billingID, err := result.LastInsertId()
	if err != nil {
		return Billing{}, err
	}

	return Billing{
		ID:            int(billingID),
		ReservationID: reservationID,
		UserID:        userID,
		VehicleID:     vehicleID,
		TotalAmount:   totalAmount,
		Discount:      discount,
		CreatedAt:     time.Now(),
	}, nil
}

// GetInvoicesByUser fetches all invoices for a given user from the database
func GetInvoicesByUser(db *sql.DB, userID int) ([]Billing, error) {
	query := "SELECT b.id, b.reservation_id, b.user_id, b.vehicle_id, v.make, v.model, b.total_amount, b.discount, b.created_at FROM billing b INNER JOIN CarSharingVehicleService.Vehicles v ON b.vehicle_id = v.id WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Billing
	for rows.Next() {
		var invoice Billing
		err := rows.Scan(&invoice.ID, &invoice.ReservationID, &invoice.UserID, &invoice.VehicleID, &invoice.Make, &invoice.Model, &invoice.TotalAmount, &invoice.Discount, &invoice.CreatedAt)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

// DeleteInvoiceByReservationID deletes an invoice by its associated reservation ID
func DeleteInvoiceByReservationID(db *sql.DB, reservationID int) error {
	query := "DELETE FROM billing WHERE reservation_id = ?"
	_, err := db.Exec(query, reservationID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateInvoice(db *sql.DB, reservationID int, userID int, vehicleID int, totalAmount float64, discount float64) error {
	query := `
        UPDATE billing
        SET total_amount = ?, discount = ?
        WHERE reservation_id = ?
    `
	_, err := db.Exec(query, totalAmount, discount, reservationID)
	if err != nil {
		return err
	}
	return nil
}
