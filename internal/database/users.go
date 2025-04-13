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
	Password string
}

func (c * Client) AddUser(username,password string) (*User,error){
	name := strings.ToLower(strings.TrimSpace(username))
	if len(name) == 0 {
		return nil,errors.New("username cannot be empty")
	}
	query := `INSERT INTO users
			(name,password)
		values (?,?)Returning *;`
	userRow := c.db.QueryRow(query,strings.ToLower(username),password)
	var user User
	err := userRow.Scan(&user.Id,&user.Name,&user.Created_at,&user.Updated_at,&user.Tasks,&user.Password)
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


func (c *Client) GetUsers()([]User,error){
	query := `SELECT * FROM USERS`

	dbUserRows,err := c.db.Query(query)
	if err != nil {
		return []User{},err
	}
	var users []User
	// Assuming users is an appropriate slice
for dbUserRows.Next() {
    var user User // or whatever your user type is
    err := dbUserRows.Scan(&user.Id, &user.Name, &user.Created_at,&user.Updated_at,&user.Tasks)
    if err != nil {
        return []User{},err
    }
    // Add user to your slice
    users = append(users, user)
}
	return users,nil

}