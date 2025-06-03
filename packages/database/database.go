package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

type Record struct {
	ID     int64
	Tokens int32
	Mode   string
}

func Initialize(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	statement := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		tokens INTEGER,
		mode TEXT
	);`

	_, err = DB.Exec(statement)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func UpdateTokens(id int64, tokens int32) error {
	statement := `UPDATE users SET tokens = tokens + ? WHERE id = ?`
	_, err := DB.Exec(statement, tokens, id)
	if err != nil {
		return fmt.Errorf("failed to update tokens: %w", err)
	}
	return nil
}

func CreateUser(id int64) error {
	var tokens int32 = 150
	var mode string = "transcribe"

	statement := `INSERT OR IGNORE INTO users (id, tokens, mode) VALUES (?, ?, ?)`

	_, err := DB.Exec(statement, id, tokens, mode)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	log.Printf("User with ID %d created with initial tokens: %d and mode: %s", id, tokens, mode)
	return nil
}

func ResetTokens() error {
	statement := `UPDATE users SET tokens = 150`
	_, err := DB.Exec(statement)
	if err != nil {
		return fmt.Errorf("failed to update all tokens: %w", err)
	}
	log.Println("All users' tokens updated to 150")
	return nil
}
func RemoveTokens(id int64) error {
	var mode string
	err := DB.QueryRow("SELECT mode FROM users WHERE id = ?", id).Scan(&mode)

	if err != nil {
		return fmt.Errorf("failed to get user mode: %w", err)
	}

	var tokens float32
	switch mode {
	case "transcribe":
		tokens = 10
	case "summarize":
		tokens = 15
	default:
		tokens = 10
	}

	_, err = DB.Exec("UPDATE users SET tokens = tokens - ? WHERE id = ?", tokens, id)
	if err != nil {
		return fmt.Errorf("failed to update tokens: %w", err)
	}

	return nil
}
func GetTokens(id int64) (int32, error) {
	var tokens int32
	err := DB.QueryRow("SELECT tokens FROM users WHERE id = ?", id).Scan(&tokens)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user with ID %d not found", id)
		}
		return 0, fmt.Errorf("failed to get tokens: %w", err)
	}
	return tokens, nil
}

func GetMode(id int64) (string, error) {
	var mode string
	err := DB.QueryRow("SELECT mode FROM users WHERE id = ?", id).Scan(&mode)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user with ID %d not found", id)
		}
		return "", fmt.Errorf("failed to get mode: %w", err)
	}
	return mode, nil
}

func SetMode(id int64, mode string) error {
	statement := `UPDATE users SET mode = ? WHERE id = ?`
	_, err := DB.Exec(statement, mode, id)
	if err != nil {
		return fmt.Errorf("failed to set mode: %w", err)
	}
	log.Printf("User with ID %d set to mode: %s", id, mode)
	return nil
}
