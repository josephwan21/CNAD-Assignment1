package main

import (
	"Assg1/CarSharingUserService/package/db"
	"Assg1/CarSharingUserService/package/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Your secret key for signing JWT
var jwtSecret = []byte("your-secret-key")

// User struct for handling both registration and login data
type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Membership string `json:"membership"`
}

// HandleLogin handles the login process
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	var dbPassword string
	var dbID int
	err = dbConn.QueryRow("SELECT id, password FROM users WHERE email = ?", user.Email).Scan(&dbID, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if !utils.CheckPasswordHash(user.Password, dbPassword) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	tokenString, err := generateJWT(user.Email, dbID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleRegister handles the user registration process
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	var existingEmail string
	err = dbConn.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&existingEmail)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Error checking email: %v", err)
		return
	}

	if err == nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	result, err := dbConn.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		log.Printf("Error hashing")
		return
	}

	// Retrieve the ID of the newly inserted user
	userID, err := result.LastInsertId() // MySQL specific method to get the last inserted ID
	if err != nil {
		http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
		log.Printf("Error retrieving user ID: %v", err)
		return
	}

	sendVerificationEmail(user.Email, int(userID))
	fmt.Printf("User ID: %d", user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
}

// generateJWT generates a JWT token for the user
func generateJWT(email string, id int) (string, error) {
	claims := &jwt.MapClaims{
		"userid": id,
		"sub":    email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // expires in 24 hours
	}

	// Create a new token with the claims and sign it with the secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func sendVerificationEmail(email string, userId int) {
	// Generate verification token
	token, err := generateJWT(email, userId)
	if err != nil {
		log.Fatal("Error generating token", err)
		return
	}

	// Create the verification URL
	verificationURL := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)

	// Send the email (simplified for example)
	auth := smtp.PlainAuth("", "josephbwanzj@gmail.com", "ofpi wrtr jnoa iniy", "smtp.gmail.com")
	to := []string{email}
	subject := "Verify Your Email"
	message := []byte(fmt.Sprintf("Subject: %s\n\nWelcome to our Car Sharing service! To verify your email, we would like you to click on the following link: %s", subject, verificationURL))

	err = smtp.SendMail("smtp.gmail.com:587", auth, "josephbwanzj@gmail.com", to, message)
	if err != nil {
		log.Fatal("Failed to send email:", err)
	}
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Token missing", http.StatusBadRequest)
		return
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Claims: %v", claims)

	// Extract user ID and verify the token's expiration
	userId, ok := claims["userid"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	// Update the user's status to "verified"
	_, err = dbConn.Exec("UPDATE Users SET verified = true WHERE id = ?", int(userId))
	if err != nil {
		http.Error(w, "Error verifying user", http.StatusInternalServerError)
		return
	}

	// Send a response
	fmt.Fprintln(w, "Email verified successfully!")
}

func HandleUpgradeMembership(w http.ResponseWriter, r *http.Request) {
	type UpgradeRequest struct {
		Email      string `json:"email"`
		Membership string `json:"membership"`
	}

	var req UpgradeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	_, err = dbConn.Exec("UPDATE users SET membership = ? WHERE email = ?", req.Membership, req.Email)
	if err != nil {
		http.Error(w, "Error updating membership", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Membership successfully upgraded."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {

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

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbConn := db.InitDB()
	defer dbConn.Close()

	_, err = dbConn.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, userIDInt)
	if err != nil {
		http.Error(w, "Error updating profile", http.StatusInternalServerError)
		return
	}

	log.Printf("Updating user with ID: %v, Name: %v, Email: %v", userID, user.Name, user.Email)

	newToken, err := generateJWT(user.Email, userIDInt) // Generate new token with updated email
	if err != nil {
		http.Error(w, "Failed to generate new token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Profile updated successfully.",
		"token":   newToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization header and extract the token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Strip the "Bearer " prefix from the token
	tokenString = tokenString[7:]

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// If the token is invalid, return an error
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Assuming the email is in the 'sub' claim
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		http.Error(w, "Email not found in token", http.StatusUnauthorized)
		return
	}

	email := claims["sub"].(string) // Extract the email from the 'sub' field

	dbConn := db.InitDB()
	defer dbConn.Close()

	var user User
	err = dbConn.QueryRow("SELECT id, name, email, password, membership FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Membership)
	if err != nil {
		http.Error(w, "Error fetching user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/login", HandleLogin)
	router.HandleFunc("/register", HandleRegister)
	router.HandleFunc("/user-profile", HandleGetProfile)
	router.HandleFunc("/update-profile", HandleUpdateProfile)
	router.HandleFunc("/verify", verifyEmailHandler)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)

	fmt.Println("User Service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}
