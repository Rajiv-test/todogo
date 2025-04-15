package database

import (
	"log"
	"time"
)

type User struct {
	Id         int
	Name       string
	Created_at time.Time
	Updated_at time.Time
	Tasks      int
	Password   string
	IsAdmin    bool
}

func (c *Client) AddUser(username, password string, created_at, updated_at time.Time, isadmin bool) (*User, error) {
	query := `INSERT INTO users
			(name,created_at,updated_at,password,is_admin)
		values (?,?,?,?,?)Returning *;`
	userRow := c.db.QueryRow(query, username, created_at, updated_at, password, isadmin)
	var user User
	err := userRow.Scan(&user.Id, &user.Name, &user.Created_at, &user.Updated_at, &user.Tasks, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	log.Printf("User %s added successfully with ID %d", user.Name, user.Id)
	return &user, nil
}

func (c *Client) GetUser(username string) (User, error) {
	query := `SELECT * FROM users WHERE name = ?;`
	queryRow := c.db.QueryRow(query, username)
	var user User
	err := queryRow.Scan(&user.Id, &user.Name, &user.Created_at, &user.Updated_at, &user.Tasks, &user.Password, &user.IsAdmin)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (c *Client) DeleteUser(username string) error {
	query := `DELETE FROM users WHERE name = ?;`
	_, err := c.db.Exec(query, username)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetUsers() ([]User, error) {
	query := `SELECT * FROM USERS`

	dbUserRows, err := c.db.Query(query)
	if err != nil {
		return []User{}, err
	}
	var users []User
	// Assuming users is an appropriate slice
	for dbUserRows.Next() {
		var user User // or whatever your user type is
		err := dbUserRows.Scan(&user.Id, &user.Name, &user.Created_at, &user.Updated_at, &user.Tasks, &user.Password, &user.IsAdmin)
		if err != nil {
			return []User{}, err
		}
		// Add user to your slice
		users = append(users, user)
	}
	return users, nil

}
