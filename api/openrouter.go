package api

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

//go:embed meal_macros_schema.json
var mealMacroSchemaJSON embed.FS
var mealMacroSchema map[string]interface{}

//go:embed meal_parts_schema.json
var mealPartsSchemaJSON embed.FS
var mealPartsSchema map[string]interface{}

//go:embed serving_sizes_schema.json
var servingSizesSchemaJSON embed.FS
var servingSizesSchema map[string]interface{}

func init() {
	data, err := mealMacroSchemaJSON.ReadFile("meal_macros_schema.json")
	if err != nil {
		log.Fatalf("Failed to read schema.json: %v", err)
	}
	if err := json.Unmarshal(data, &mealMacroSchema); err != nil {
		log.Fatalf("Failed to unmarshal schema.json: %v", err)
	}

	data, err = mealPartsSchemaJSON.ReadFile("meal_parts_schema.json")
	if err != nil {
		log.Fatalf("Failed to read schema.json: %v", err)
	}
	if err := json.Unmarshal(data, &mealPartsSchema); err != nil {
		log.Fatalf("Failed to unmarshal schema.json: %v", err)
	}

	data, err = servingSizesSchemaJSON.ReadFile("serving_sizes_schema.json")
	if err != nil {
		log.Fatalf("Failed to read schema.json: %v", err)
	}
	if err := json.Unmarshal(data, &servingSizesSchema); err != nil {
		log.Fatalf("Failed to unmarshal schema.json: %v", err)
	}
}

type JSONSchemaMap map[string]interface{}

func (j JSONSchemaMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(j))
}

type OpenRouter struct {
	apiKey  string
	baseURL string
}

type APIClient struct {
	*openai.Client
}

func New() *OpenRouter {
	OPENROUTER_API_KEY, exists := os.LookupEnv("OPENROUTER_API_KEY")
	if !exists {
		log.Fatalf("OPENROUTER_API_KEY not set")
	}

	return &OpenRouter{
		apiKey:  OPENROUTER_API_KEY,
		baseURL: "https://openrouter.ai/api/v1",
	}
}

func (o *OpenRouter) NewClient() *APIClient {
	config := openai.DefaultConfig(o.apiKey)
	config.BaseURL = o.baseURL
	client := openai.NewClientWithConfig(config)

	return &APIClient{Client: client}
}

func (c *APIClient) Log(id string, message string) {
	fmt.Printf("[OpenRouter] [%s] %s\n", id, message)
}
