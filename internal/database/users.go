package database

import (
	"errors"
	"log"
	"strings"
	"time"
)


type User struct{
	Id int
	Name string
	Created_at time.Time
	Updated_at time.Time
	Tasks int
}

func (c * Client) AddUser(username string) (*User,error){
	name := strings.ToLower(strings.TrimSpace(username))
	if len(name) == 0 {
		return nil,errors.New("username cannot be empty")
	}
	query := `INSERT INTO users
			(name)
		values (?)Returning *;`
	userRow := c.db.QueryRow(query,strings.ToLower(username))
	var user User
	err := userRow.Scan(&user.Id,&user.Name,&user.Created_at,&user.Updated_at,&user.Tasks)
	if err != nil {
		return nil,err
	}
	log.Printf("User %s added successfully with ID %d", user.Name, user.Id)
	return &user,nil
}

func (c *Client) DeleteUser(username string) error{
	name := strings.ToLower(strings.TrimSpace(username))
	if len(name) == 0 {
		return errors.New("username cannot be empty")
	}
	query := `DELETE FROM users WHERE name = ?;`
	_,err := c.db.Exec(query,name)
	if err != nil {
		return err
	}
	log.Print("Successfully deleted the user ",name)
	return nil
}