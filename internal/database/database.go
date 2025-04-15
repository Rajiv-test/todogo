package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	db *sql.DB
}

func NewClient(pathToDB string) (Client, error) {
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		return Client{}, err
	}

	// Enable foreign keys for this connection
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return Client{}, err
	}

	c := Client{db}
	err = c.automigrate()
	if err != nil {
		return Client{}, err
	}
	return c, nil
}

func (c *Client) automigrate() error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY NOT NULL,
		name TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		tasks INTEGER DEFAULT 0,
		password TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT 0
	);`
	_, err := c.db.Exec(userTable)
	if err != nil {
		return err
	}

	tasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY NOT NULL,
		username INTEGER NOT NULL REFERENCES users(name) ON DELETE CASCADE,
		taskname TEXT NOT NULL,
		description TEXT DEFAULT 'No Description',
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,	
		deadline TIMESTAMP DEFAULT NULL,
		completed BOOLEAN DEFAULT 0,
		over_deadline BOOLEAN DEFAULT 0,
		CHECK (deadline IS NULL OR deadline > created_at)
	);`
	_, err = c.db.Exec(tasksTable)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Reset() error {
	_, err := c.db.Exec("DELETE FROM users")
	if err != nil {
		return fmt.Errorf("error reseting the users table %w", err)
	}
	_, err = c.db.Exec("DELETE FROM tasks")
	if err != nil {
		return fmt.Errorf("error reseting the tasks table %w", err)
	}
	return nil
}
