package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func DatabaseCreate() error {
	// Create Data directory if it doesn't already exist
	if _, err := os.Stat("./Data"); os.IsNotExist(err) {
		err = os.Mkdir("./Data", 0755)
		if err != nil {
			return fmt.Errorf("error creating Data directory: %v", err)
		}
	}
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Create the table if it doesn't already exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (user_id TEXT PRIMARY KEY, token TEXT UNIQUE)`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

func DatabaseInsert(user_id string, token string) error {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Insert the token into the database
	_, err = db.Exec(`INSERT INTO tokens (user_id, token) VALUES (?, ?)`, user_id, token)
	if err != nil {
		return fmt.Errorf("error inserting token: %v", err)
	}

	return nil
}

func DatabaseDelete(user_id string) error {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Delete the token from the database
	_, err = db.Exec(`DELETE FROM tokens WHERE user_id = ?`, user_id)
	if err != nil {
		return fmt.Errorf("error deleting token: %v", err)
	}

	return nil
}

type User struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func DatabaseSelectAll() ([]User, error) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Select the token from the database
	rows, err := db.Query(`SELECT * FROM tokens`)
	if err != nil {
		return nil, fmt.Errorf("error selecting token: %v", err)
	}
	defer rows.Close()

	// user map ({"users": ["user_id": "...", "token": "..."], ...})
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserID, &user.Token)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}
		users = append(users, user)
	}

	return users, nil

}

// Verify admin key
func VerifyAdminKey(key string) bool {
	return key == os.Args[2]
}

func VerifyToken(token string) (bool, error) {
	// Check if token is admin key
	if VerifyAdminKey(token) {
		return true, nil
	}
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./Data/auth.db")
	if err != nil {
		return false, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Select the token from the database
	rows, err := db.Query(`SELECT * FROM tokens WHERE token = ?`, token)
	if err != nil {
		return false, fmt.Errorf("error selecting token: %v", err)
	}
	defer rows.Close()

	// Check if the token exists
	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}
