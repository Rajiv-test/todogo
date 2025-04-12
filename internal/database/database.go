package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct{
	db *sql.DB
}

func NewClient(pathToDB string) (Client,error){
	db,err := sql.Open("sqlite3",pathToDB)
	if err != nil {
		return Client{},err
	}
	
	// Enable foreign keys for this connection
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return Client{},err
	}

	c := Client{db}
	err = c.automigrate()
	if err != nil {
		return Client{},err
	}
	return c,nil
}


func (c *Client) automigrate() error{
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		tasks INTEGER DEFAULT 0
	);`
	_,err := c.db.Exec(userTable)
	if err != nil {
		return err
	}
	log.Print("successfully created a users table")

	tasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY NOT NULL,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		completed BOOLEAN DEFAULT 0
	);`
	_,err = c.db.Exec(tasksTable)
	if err != nil {
		return err
	}
	log.Print("successfully created a tasks table")
	return nil
}

func (c * Client) Reset() error{
	_,err := c.db.Exec("DELETE FROM users")
	if err != nil {
		return fmt.Errorf("error reseting the users table %w",err)
	}
	log.Print("reset completed")
	return nil
}

