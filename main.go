package main

import (
	"calorie_bot/api"
	"calorie_bot/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	log.Print("MAIN")
}

func main() {
	var port string
	var host string

	if os.Getenv("FLY_ALLOC_ID") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file in MAIN: %v", err)
		}

		port = "4444"
		host = ""

	} else {
		port = "8080"
		host = "0.0.0.0"
	}

	api := api.New()

	s := server.New(api)

	if err := s.Run(host, port); err != nil { // Call Run on the instance
		log.Fatalf("Failed to run server: %v", err)
	}

}
