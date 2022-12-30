package main

import (
	"Rest-with-postgres/internal/database"
	"log"
)

func main() {
	//Instantiates the database
	postgres, err := database.NewPostgres('localhost', 'james', 'Neew$Things1')
	if err != nil {
		log.Fatal(err.Error())
	}
}
