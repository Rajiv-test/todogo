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
	c := Client{db}
	_ = c.automigrate()
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
	log.Print("successfully created a database client")
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

