package main

import (
	"fmt"
	"os"
)


func commandExit(c *config,args ...string)error{
	fmt.Println("Exiting app\n Goodbye!")
	os.Exit(0)
	return nil
}