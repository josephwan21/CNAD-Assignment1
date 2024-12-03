// billing.go
package billing

import (
	"database/sql"
	"math"
	"time"
)

// Billing represents a billing record
type Billing struct {
	ID          int       `json:"id"`
	RentalID    int       `json:"rental_id"`
	UserID      int       `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Discount    float64   `json:"discount"`
	CreatedAt   time.Time `json:"created_at"`
}

// CalculateBilling calculates the total billing amount based on rental duration, membership, and any applicable discounts
func CalculateBilling(db *sql.DB, rentalID int, userID int, startTime time.Time, endTime time.Time) (float64, float64, error) {
	// Get membership type of the user
	var membership string
	err := db.QueryRow("SELECT membership FROM users WHERE id = ?", userID).Scan(&membership)
	if err != nil {
		return 0, 0, err
	}

	// Calculate rental duration in hours
	duration := endTime.Sub(startTime).Hours()

	// Pricing model (example: basic rate per hour)
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

	// Discount (example: 10% off for Premium members)
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
func CreateBillingRecord(db *sql.DB, rentalID int, userID int, totalAmount float64, discount float64) (Billing, error) {
	query := "INSERT INTO billing (rental_id, user_id, total_amount, discount) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, rentalID, userID, totalAmount, discount)
	if err != nil {
		return Billing{}, err
	}

	billingID, err := result.LastInsertId()
	if err != nil {
		return Billing{}, err
	}

	return Billing{
		ID:          int(billingID),
		RentalID:    rentalID,
		UserID:      userID,
		TotalAmount: totalAmount,
		Discount:    discount,
		CreatedAt:   time.Now(),
	}, nil
}
