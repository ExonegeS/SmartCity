package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	//_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type CityGuide struct {
	Name         string  `json:"name"`
	Contact      string  `json:"contact"`
	Price        float64 `json:"price"`
	PersonalData string  `json:"personal_data"`
	// Add other fields as needed
}

type CityLandmark struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

type LostItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Contact  string `json:"contact"`
}

type HistoricPhoto struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Year     int    `json:"year"`
	PhotoURL string `json:"photo_url"`
	History  string `json:"history"`
}

// Define a global variable for the database connection
var db *sql.DB
var cityGuides []CityGuide
var cityLandmarks []CityLandmark
var lostItems []LostItem
var historicPhotos []HistoricPhoto
var secretKey = []byte("12345") // Replace with your secret key

// Middleware for JWT token verification and authentication
func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tokenString := r.Header.Get("Authorization")

		// // Check if the token is missing
		// if tokenString == "" {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// // Verify the token
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 	// Make sure the signing method is the same as what you used for token generation
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("Unexpected signing method")
		// 	}
		// 	return secretKey, nil
		// })

		// if err != nil {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// // Check if the token is valid
		// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 	// Set the user ID in the request context
		// 	userID := int(claims["sub"].(float64))
		// 	ctx := context.WithValue(r.Context(), "userID", userID)
		// 	r = r.WithContext(ctx)
		// } else {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	// Status endpoint to check the backend's connection
	r.HandleFunc("/api/status", GetStatus).Methods("GET")

	// Add the logging middleware to log traffic
	r.Use(LoggingMiddleware)

	// Initialize the database connection
	var err error
	db, err = sql.Open("sqlite", "database/sql/users.db")

	if err != nil {
		log.Fatal(err)
	}
	// Check the database connection
	if err := checkDBConnection(db); err != nil {
		log.Fatal("Database connection error:", err)
	}

	// Ensure the database connection is closed when the application exits
	defer db.Close()

	// Create the "users" table
	createUsersTable(db)

	// Create the "city_guides" table
	createCityGuidesTable(db)

	// Insert the admin account
	//if err := insertAdminAccount(db); err != nil {		log.Fatal("Failed to insert admin account:", err)	}

	// Query to retrieve the list of all tables in the database
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the result set and print table names
	fmt.Println("Tables in the database:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
	}

	fmt.Println("Database created and connected successfully.")

	r.HandleFunc("/api/login", Login).Methods("POST")

	// Wrap the CreateCityGuide handler with AuthenticateMiddleware
	createCityGuideHandler := http.HandlerFunc(CreateCityGuide)

	// Use the wrapped handler for the route
	r.Handle("/api/cityguide", AuthenticateMiddleware(createCityGuideHandler)).Methods("POST")

	r.HandleFunc("/api/cityguide", GetCityGuides).Methods("GET")
	r.HandleFunc("/api/cityguide/{name}", DeleteCityGuideByName).Methods("DELETE")

	// Explore the city landmarks
	r.HandleFunc("/api/citytour", GetCityGuides).Methods("GET")

	// // Search for lost items
	// r.HandleFunc("/api/search/lostitems", SearchLostItems).Methods("GET")

	// // Historic part with photos and history
	// r.HandleFunc("/api/history/historicphotos", GetHistoricPhotos).Methods("GET")

	// // City soundscape
	// r.HandleFunc("/api/soundscape", GetCitySoundscape).Methods("GET")

	// // Photo of the day
	// r.HandleFunc("/api/photos/photooftheday", GetPhotoOfTheDay).Methods("GET")

	// // Challenge to try local dishes
	// r.HandleFunc("/api/challenges/localdishes", GetLocalDishChallenges).Methods("GET")

	// // Video call between campers
	// r.HandleFunc("/api/video/call", VideoCall).Methods("GET")

	// Create a new CORS handler with appropriate options
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Add the origin(s) of your React frontend
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Use the CORS handler with your router
	handler := c.Handler(r)

	http.Handle("/", handler) // Use the router for all routes

	fmt.Println("Server is listening on :3001")
	http.ListenAndServe(":3001", nil)
}
func checkDBConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

func createUsersTable(db *sql.DB) {
	// Create the "users" table if it doesn't exist
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL,
            password_hash TEXT NOT NULL
        )
    `)

	if err != nil {
		log.Fatal("Error creating 'users' table:", err)
	}
}
func insertAdminAccount(db *sql.DB) error {
	// Hash the admin's password
	password := "zxcasdqwe" // Replace with the desired password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert the admin account into the "users" table
	_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", "admin@nu.edu.kz", hashedPassword)
	return err
}

func createCityGuidesTable(db *sql.DB) {
	// Create the "city_guides" table if it doesn't exist
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS city_guides (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            contact TEXT,
            price REAL,
            personal_data TEXT
        )
    `)

	if err != nil {
		log.Fatal("Error creating 'city_guides' table:", err)
	}
}

// Logger
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Capture the response status code
		rw := NewResponseLogger(w)
		next.ServeHTTP(rw, r)

		// Log response details
		log.Printf("Response: %d", rw.Status())
	})
}

type ResponseLogger struct {
	http.ResponseWriter
	status int
}

func NewResponseLogger(w http.ResponseWriter) *ResponseLogger {
	return &ResponseLogger{w, http.StatusOK}
}
func (rw *ResponseLogger) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
func (rw *ResponseLogger) Status() int {
	return rw.status
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request data")
		return
	}

	// Authenticate the user based on the provided email and password
	user, err := authenticateUser(loginData.Email, loginData.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Authentication failed")
		return
	}

	// Generate a JWT token for the authenticated user
	token, err := generateJWTToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to generate token")
		return
	}
	fmt.Fprintln(w, "No error")
	// Return the token as the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
func authenticateUser(email, password string) (*User, error) {
	// Query the database to find the user by email
	user, err := getUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
func getUserByEmail(email string) (*User, error) {
	var err error
	// Query the database to find the user by email
	var user User
	err = db.QueryRow("SELECT id, email, password_hash FROM users WHERE email=?", email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	//err = db.QueryRow("SELECT * FROM users LIMIT 1;").Scan(user.ID)
	//log.Printf("User: %d %s %s", user.ID, user.Email, user.PasswordHash)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err // Handle the error appropriately
	}
	return &user, nil
}
func generateJWTToken(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iss": "your-issuer",                         // Replace with your issuer
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expiration (e.g., 24 hours)
	})

	return token.SignedString(secretKey)
}

// Rewrite the CreateCityGuide function to use SQLite
func CreateCityGuide(w http.ResponseWriter, r *http.Request) {
	var newGuide CityGuide
	// Parse the request body into a CityGuide struct
	if err := json.NewDecoder(r.Body).Decode(&newGuide); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid request data")
		log.Printf("newGuide2: %s", err)
		return
	}

	// Insert the new city guide entry into the database
	stmt, err := db.Prepare("INSERT INTO city_guides (name, contact, price, personal_data) VALUES (?, ?, ?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to prepare SQL statement")
		log.Printf("newGuide3: %s", err)
		return
	}
	defer stmt.Close() // Defer closing the prepared statement

	_, err = stmt.Exec(newGuide.Name, newGuide.Contact, newGuide.Price, newGuide.PersonalData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to insert data into the database")
		log.Printf("newGuide4: %s", err)
		return
	}

	// Return the inserted city guide entry as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newGuide)
}

// Handler to retrieve all city guide entries
func GetCityGuides(w http.ResponseWriter, r *http.Request) {
	// Query the database to retrieve all city guides
	rows, err := db.Query("SELECT name, contact, price, personal_data FROM city_guides")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to fetch city guides")
		log.Printf("Error fetching city guides: %v", err)
		return
	}
	defer rows.Close()

	// Create a slice to hold the city guides
	var cityGuides []CityGuide

	// Iterate through the rows and populate the cityGuides slice
	for rows.Next() {
		var guide CityGuide
		if err := rows.Scan(&guide.Name, &guide.Contact, &guide.Price, &guide.PersonalData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Failed to fetch city guides")
			log.Printf("Error scanning city guide row: %v", err)
			return
		}
		cityGuides = append(cityGuides, guide)
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the city guides as a JSON response
	if err := json.NewEncoder(w).Encode(cityGuides); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to encode city guides as JSON")
		log.Printf("Error encoding city guides: %v", err)
		return
	}
}

// Handler to delete a city guide entry by name
func DeleteCityGuideByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Execute a SQL DELETE statement to delete the city guide entry by name.
	result, err := db.Exec("DELETE FROM city_guides WHERE name = ?", name)
	if err != nil {
		// Handle the error, e.g., log it or return an appropriate response.
		fmt.Println("Error deleting city guide:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error deleting city guide")
		return
	}

	// Check if any rows were affected to determine if the entry was found.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Handle the error.
		fmt.Println("Error getting rows affected:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error deleting city guide")
		return
	}

	if rowsAffected == 0 {
		// No rows were affected, indicating that the city guide entry was not found.
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "City guide not found")
		return
	}

	// City guide entry deleted successfully.
	w.WriteHeader(http.StatusNoContent)
}

// Implement the handler functions for each endpoint
// (ExploreCityLandmarks, SearchLostItems, GetHistoricPhotos, GetCitySoundscape,
// GetPhotoOfTheDay, GetLocalDishChallenges, and VideoCall)
// ...

func GetStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{"status": "OK"}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the status as a JSON response
	jsonResponse(w, status)
}

// Implement the handler functions for each endpoint with appropriate logic
// ...

func jsonResponse(w http.ResponseWriter, data interface{}) {
	// Encode data as JSON and send as the response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
