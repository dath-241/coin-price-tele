package services

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	var err error

	// Get the DATABASE_URL environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Open a connection to the PostgreSQL database
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Create the table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_tokens (
        user_id BIGINT PRIMARY KEY,
		is_muted BOOLEAN DEFAULT FALSE,
        token TEXT
    )`)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_symbols (
        user_id BIGINT PRIMARY KEY,
        symbol TEXT
    )`)
	if err != nil {
		log.Fatal("Error creating user_symbols table:", err)
	}
}

// Store the token in the database
func StoreUserToken(userID int, token string) error {
	_, err := db.Exec(`INSERT INTO user_tokens (user_id, token) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET token = EXCLUDED.token`, userID, token)
	return err
}

// Retrieve the token from the database
func GetUserToken(userID int) (string, error) {
	var token string
	err := db.QueryRow(`SELECT token FROM user_tokens WHERE user_id = $1`, userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetMute(userID int) (bool, error) {
	var isMuted bool
	err := db.QueryRow(`SELECT is_muted FROM user_tokens WHERE user_id = $1`, userID).Scan(&isMuted)
	if err != nil {
		return false, err
	}
	return isMuted, nil
}

func SetMute(userID int, isMuted bool) error {
	_, err := db.Exec(`INSERT INTO user_tokens (user_id, is_muted) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET is_muted = EXCLUDED.is_muted`, userID, isMuted)
	return err
}
