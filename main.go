package main

import (
	"log"
	"os"

	"github.com/Rajiv-test/todogo/internal/database"
	"github.com/joho/godotenv"
)


func main() {
	godotenv.Load(".env")
	dbPath := os.Getenv("PATH_TO_DB")
	if dbPath == ""{
		log.Fatal("please specify path to database")
	}
	db,err := database.NewClient(dbPath)
	if err != nil{
		log.Fatalf("error creating a database client: %v",err)
	}
	user,err := db.AddUser("Rajiv")
	if err != nil {
		log.Fatal(err)
	}
	err = db.DeleteUser(" Rajiv")
	if err != nil {
		log.Fatal(err)
	}	
	log.Print(*user)
}