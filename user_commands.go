package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

func commandRegister(c *config, args ...string) error {
	if len(args) != 1 && (len(args) == 2 && args[0] != "-a") {
		return errors.New("wrong usage of reg command use help command to learn more")
	}
	if c.user != nil {
		return fmt.Errorf("logout to register new user. use reg command to register new user")
	}

	username := args[0]
	isadmin := false
	if len(args) == 2 {
		username = args[1]
		isadmin = true
	}
	_, err := c.db.GetUser(username)
	if err == nil {
		return errors.New("username already taken please choose other username")
	} else if err != sql.ErrNoRows {
		return err
	}

	fmt.Println("Please provide password for the user:", username)
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter new Password: ")
		reader.Scan()
		input := strings.Fields(reader.Text())
		if len(input) != 1 || input[0] == "" || len(input[0]) < 5 {
			fmt.Println("please provide valid password")
			continue
		} else {
			newUser, err := c.db.AddUser(username, input[0], time.Now().Local().UTC(), time.Now().Local().UTC(), isadmin)
			if err != nil {
				return err
			}
			fmt.Printf("New User created with username: %v at %v with 0 tasks\n", newUser.Name, newUser.Created_at)
			c.user = newUser
			fmt.Println("Current user > ", username)
			fmt.Println("You can use help to know more about available commands")
			break
		}
	}

	return nil
}

func commandLogin(c *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("wrong usage of login command use help command to learn more")
	}
	username := args[0]
	if c.user != nil && c.user.Name == username {
		return fmt.Errorf("user already logged in")
	}
	if c.user != nil {
		fmt.Printf("logout current user: %v to login\n", c.user.Name)
		fmt.Println("Use logout command to logout the current user")
		return nil
	}

	currentUser, err := c.db.GetUser(username)
	if err == sql.ErrNoRows {
		return errors.New("user doesn't exit use the reg command to register new user")
	} else if err != nil {
		return err
	}
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter password: ")
	reader.Scan()
	password := reader.Text()
	if password != currentUser.Password {
		return fmt.Errorf("incorrect password try again! ")
	}

	c.user = &currentUser
	return nil
}

func commandLogout(c *config, args ...string) error {
	if len(args) > 0 {
		return fmt.Errorf("logout command takes no arguments use help to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("no user to logout")
	}
	c.user = nil
	return nil
}

func commandDelete(c *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("wrong usage of delete use help command to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("login in to delete users")
	}
	if c.user.Name != args[0] && !c.user.IsAdmin {
		return fmt.Errorf("cannot delete other users unless you are admin")
	}
	err := c.db.DeleteUser(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("%v successfully deleted the user: %v\n", c.user.Name, args[0])
	if c.user.Name == args[0] {
		c.user = nil
	}
	return nil
}

func commandListUsers(c *config, args ...string) error {
	if len(args) != 0 {
		return errors.New("wrong usage of login command use help command to learn more")
	}
	if c.user == nil {
		return fmt.Errorf("login in as admin to view all registered users")
	}
	if !c.user.IsAdmin {
		return fmt.Errorf("not authorized to list users")
	}
	users, err := c.db.GetUsers()
	if err != nil {
		return err
	}

	fmt.Printf("%-5s %-8s %-14s %-12s %-5s %-5s\n", "No", "Name", "CreatedAt", "UpdatedAt", "Tasks", "admin")
	fmt.Println(strings.Repeat("-", 60)) // Add a separator line

	// Print each user with aligned columns
	for i, user := range users {
		fmt.Printf("%-4d %-9s %-14s %-14s %-5d %-6v\n",
			i,
			user.Name,
			user.Created_at.Format(time.DateOnly),
			user.Updated_at.Format(time.DateOnly),
			user.Tasks,
			user.IsAdmin)
	}
	return nil
}
