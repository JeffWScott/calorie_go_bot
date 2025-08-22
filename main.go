package main

import (
	"calorie_bot/api"
	"calorie_bot/db"
	"calorie_bot/pipeline"
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

	db := db.New()
	api := api.New()
	pipeline := pipeline.New(api, db)

	s := server.New(pipeline)

	if err := s.Run(host, port); err != nil { // Call Run on the instance
		log.Fatalf("Failed to run server: %v", err)
	}

}

/* Currently accepts a prompt and feeds it to OpenRouter and gives back the response. */

/* We need to change this to an agentic model

1. Take the prompt and feed it to a fast model that can split out the main ingredients.
2. Take that ingredients list and call the symantic search on Convex to get DB possibilities
3. Food those back into a model and ask it to pick the relevant ones
4. Take the response and get the calorie info
5. Feed that calorie info back into a model and have it esimate the calcories for the meal and provide the structured output
6. Return that to Convex

*/
