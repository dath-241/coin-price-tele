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

// Thêm các hàm mới để xử lý symbol
func StoreUserSymbol(userID int, symbol string) error {
	_, err := db.Exec(`INSERT INTO user_symbols (user_id, symbol) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET symbol = EXCLUDED.symbol`, userID, symbol)
	return err
}

func GetUserSymbol(userID int) (string, error) {
	var symbol string
	err := db.QueryRow(`SELECT symbol FROM user_symbols WHERE user_id = $1`, userID).Scan(&symbol)
	if err != nil {
		return "", err
	}
	return symbol, nil
}
