package api

import (
	"bytes"
	"calorie_bot/types"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

var MEAL_MACROS_SYSTEM_PROMPT = `
	You are a professional  calorie and macro nutrients counter. 
	A user will provide you a description of a meal and you will provide them back two things:
	1. an array consisting of each part of the mean with it's macros.
	
	Food Grouping example:
	USER - 'I had two baked potatoes with sour cream and three buttered rolls for lunch'

	✅ CORRECT grouping response:
	response [{..., food_name: "Baked Potatoes"}, {..., food_name: "Buttered Rolls"}]

	✖️ BAD grouping response:
	response [{..., food_name: "potato_1"}, {..., food_name: "sour_cream_on_potato_1"}, {..., food_name: "potato_2"}, {..., food_name: "sour_cream_on_potato_2"}, {..., food_name: "roll_1"}, {..., food_name: "butter_on_roll_1"}, {..., food_name: "roll_2"},  {..., food_name: "butter_on_roll_2"}]

	2. a SHORT (6 word MAX) description of the meal they said they had, being as concise as possible. We don't need to specify the meal (lunch) or that it was a meal. Just focue on describing what the macros describe.
	✅ CORRECT meal_description response:
	response {..., meal_description: "Baked Potatoes and Rolls"}

	✖️ BAD meal_description response:
	response {..., meal_description: "Two Baked Potatoes with Sour Cream with Two Buttered Rolls for Lunch."}
	
	Report all macros in the numerical representation of the units below.
	Include Estimated Serving Size of the FoodGrouping in grams (g).
	Report Total Fat, Saturated Fat, Trans Fat, Total Carbohydrates, Dietary Fiber, Total Sugars, Added Sugars, Sugar Alcohol, and Protein in grams (g).
	Report Cholesterol and Sodium in milligrams (mg).
	For Vitamins and Minerals, report Vitamin D in micrograms (mcg), and Calcium, Iron, and Potassium in milligrams (mg).
	Report Calories as a numeric value (no unit).

	Example Value response:
	✅ CORRECT numerical representation:
		{..., "Cholestero_mg": 10}  // numerical
	✖️ BAD numerical representation:
		{..., "Cholestero_mg": "10mg"}  // Not numerical
`

var MEAL_MACROS_PREAMBLE = `
Use this nutritional info to help you estimate the macros for the user's meal. 
This is the only info I have so you will have to estimate anythign that isn't provied:
`

func (c *APIClient) GetMealMacros(prompt string, nutritionalInfo []types.NutritionalInfo, model string, id string) (string, error) {
	jsonBytes, err := json.MarshalIndent(nutritionalInfo, "", "  ")
	if err != nil {
		return "", err
	}
	userContent := "MY MEAL: " + prompt + "\n\n" + MEAL_MACROS_PREAMBLE + "\n" + string(jsonBytes)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: MEAL_MACROS_SYSTEM_PROMPT,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userContent,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "structured_output",
				Schema: JSONSchemaMap(mealMacroSchema),
			},
		},
	}

	resp, err := c.CreateChatCompletion(context.Background(), req)
	if err != nil {
		c.Log(id, fmt.Sprintf("Chat completion error: %v", err))
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", nil
	}

	content := resp.Choices[0].Message.Content
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(content)); err != nil {
		return "", err
	}
	return buf.String(), nil
}
