package main

import (
	"fmt"
	"log"
)

func commandHelp(c *config, args ...string) error {
	if len(args) != 0 {
		log.Print("help command doesn't take arguments")
	}
	commandList := getcommands()
	for _,comm := range commandList{
		fmt.Println("name: ",comm.name)
		fmt.Println("description: ",comm.description)
		fmt.Println("usage: ",comm.usage)
		fmt.Println("---------------------")
	}
	return nil
}