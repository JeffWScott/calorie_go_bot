package main

import (
	"calorie_bot/api"
	"calorie_bot/server"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	log.Print("MAIN")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	api := api.New()

	s := server.New(api)

	if err := s.Run(); err != nil { // Call Run on the instance
		log.Fatalf("Failed to run server: %v", err)
	}

}
