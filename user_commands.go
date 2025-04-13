package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func commandRegister(c *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("wrong usage of login command it takes 1 argument as username")
	}
	username := args[0]
	fmt.Println("Please provide password for the user:",username)
	reader := bufio.NewScanner(os.Stdin)
	for {	
		fmt.Print("Enter new Password: ")
		reader.Scan()
		input := strings.Fields(reader.Text())
		if len(input) != 1 || input[0] == ""{
			fmt.Println("please provide valid password") 
			continue
		}else{
			user,err := c.db.AddUser(username,input[0])
			if err != nil {
				return err
			}
			fmt.Printf("New User created with username: %v at %v with 0 tasks\n",user.Name,user.Created_at)
			c.user = username
			fmt.Println("Current user > ",username)
			fmt.Println("You can use help to know more about available commands")
			break
		}
	}
	
	
	return nil
}