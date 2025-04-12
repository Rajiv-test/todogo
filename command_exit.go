package main

import (
	"fmt"
	"os"
)


func commandExit(c *config,args ...string)error{
	fmt.Println("Exiting currect user's session")
	os.Exit(0)
	return nil
}