package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Rajiv-test/todogo/internal/database"
	"github.com/joho/godotenv"
)
type config struct{
	db database.Client
	dbPath string
	user string
}

func main() {
	godotenv.Load(".env")
	databasePath := os.Getenv("PATH_TO_DB")
	if databasePath == ""{
		log.Fatal("please specify path to database")
	}
	dbClient,err := database.NewClient(databasePath)
	if err != nil{
		log.Fatalf("error creating a database client: %v",err)
	}
	cfg := &config{
		db: dbClient,
		dbPath: databasePath,
		user: "",
	}
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("todo > ")
		reader.Scan()
		input := cleanInput(reader.Text())
		if len(input) == 0{
			continue
		}
		commandName := input[0]
		args := []string{}
		if len(input) >1 {
			args = input[1:]
		}
		command,exists := getcommands()[commandName]
		if exists{
			err = command.callback(cfg,args...)
			if err != nil{
				fmt.Println(err)
			}
			continue
		}else{
			fmt.Println("Unknown command use help to find existing commands")
			continue
	}
}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}