package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

func commandRegister(c *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("wrong usage of reg command it takes 1 argument as username")
	}
	username := args[0]
	_,err := c.db.GetUser(username)
	if err == nil{
		return errors.New("username already taken please choose other username")
	}else if err != sql.ErrNoRows{
		return err
	}
	
	fmt.Println("Please provide password for the user:",username)
	reader := bufio.NewScanner(os.Stdin)
	for {	
		fmt.Print("Enter new Password: ")
		reader.Scan()
		input := strings.Fields(reader.Text())
		if len(input) != 1 || input[0] == "" || len(input[0]) < 5{
			fmt.Println("please provide valid password") 
			continue
		}else{
			newUser,err := c.db.AddUser(username,input[0])
			if err != nil{
				return err
			}
			fmt.Printf("New User created with username: %v at %v with 0 tasks\n",newUser.Name,newUser.Created_at)
			c.user = *newUser
			fmt.Println("Current user > ",username)
			fmt.Println("You can use help to know more about available commands")
			break
		}
	}
	
	return nil
}

func commandLogin(c *config, args ...string) error {
	if c.user.Name != ""{
		fmt.Printf("logout current user: %v to login\n",c.user.Name)
		fmt.Println("Use logout command to logout the current user")
		return nil
	}
	if len(args) != 1 {
		return errors.New("wrong usage of login command it takes 1 argument as username")
	}
	username := args[0]
	currentUser,err := c.db.GetUser(username)
	if err == sql.ErrNoRows{
		return errors.New("user doesn't exit use the reg command to register new user")
	}else if err != nil{
		return err
	}
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter password: ")
	reader.Scan()
	password := reader.Text()
	if password != currentUser.Password{
		return fmt.Errorf("incorrect password try again! ") 
	}
	
	c.user = currentUser
	return nil
}