package api

import (
	"calorie_bot/types"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

var MEAL_PARTS_SYSTEM_PROMPT = `
## Task
You are an expert food scientist.
You have been tasked with identifying the main PARTS that make up the main nutritional value of a meal.
You need to identify the PARTS at a HIGH LEVEL of the meal without getting granular into specific cooking ingredients.
The goal isn't to be able to reproduce the meal, it's to get the nutrient value of the major parts.

## Examples: (using structured outputs)
----------
User Prompt: "Chicken Caesar Salad with a bag of Smart Sweets Peach Rings"
BAD ✖️:   { "major_ingredients": ["cooked chicken breast", "caesar dressing", "bacon bits", "croutons", "romaine lettuce",  "gelatin", "sugar free sweetener", "peach flavour"] }
GOOD ✅: { "major_ingredients": ["chicken", "caesar salad", "smart sweets peach rings"] }
* No need to split up the ingredient for caesar salad when it's a common item we can find nutritional value for *
----------
User Prompt: "I ate a plate of spaghetti with meat sauce"
BAD ✖️:   { "major_ingredients": ["spaghetti pasta", "tomatoes", "olive oil", "garlic", "onion", "basil", "salt", "pepper", "ground beef"] }
GOOD ✅: {"major_ingredients": ["spaghetti pasta","tomato sauce","ground beef"]}
* No need to split up the ingredient for tomato sauce when it's a common item we can find nutritional value for *
----------
User Prompt: "2 slices of Pepperoni pizza"
BAD ✖️:   {"major_ingredients": ["2 pepperoni pizza slices"]}
GOOD ✅: {"major_ingredients": ["pepperoni pizza"]}
* The goal is to identify PARTS not serving sizes so "2 slices of Pepperoni pizza" is just "pepperoni pizza" *
----------
User Prompt: "McDonalds Big Mac and a Large Diet Pepsi"
BAD ✖️:   {"major_ingredients": ["hamburger", "diet cola"]}
GOOD ✅: {"major_ingredients": ["McDonald's Big Mac hamburger", "Diet Pepsi"]}
* Be as specific as the user is and don't generalize. It's much better for us to know the exact item than to generalize it *
----------

RULES:
- Do not Generalize away Brands of foods. You can break out that part of the mean but keep the branding.
- Split apart the meals into Parts.  eg. "Spaghetti and meat sauce" would be ["spaghetti pasta","tomato sauce","ground beef"]
- Ignore serving sizes, just provide the PART of the meal
- Provide the major PARTS of the meal, not individual ingredients that make up the Major Parts; eg. ["tomato sauce" ] instead of [ "tomatoes", "olive oil", "garlic", "onion", "basil", "salt", "pepper"]
- Provide the output in the Structured Output JSON format
`

func (c *APIClient) GetMealParts(prompt string, id string) (types.MealParts, error) {
	req := openai.ChatCompletionRequest{
		Model: "google/gemini-2.5-flash-lite",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: MEAL_PARTS_SYSTEM_PROMPT,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "structured_output",
				Schema: JSONSchemaMap(mealPartsSchema),
			},
		},
	}

	resp, err := c.CreateChatCompletion(context.Background(), req)
	if err != nil {
		c.Log(id, fmt.Sprintf("Chat completion error: %v", err))
		return types.MealParts{}, err
	}

	if len(resp.Choices) == 0 {
		return types.MealParts{}, nil
	}

	content := resp.Choices[0].Message.Content
	c.Log(id, fmt.Sprintf("response content: %v", resp.Choices[0].Message.Content))

	var mealParts types.MealParts
	if err := json.Unmarshal([]byte(content), &mealParts); err != nil {
		c.Log(id, fmt.Sprintf("JSON unmarshal error: %v", err))
		c.Log(id, fmt.Sprintf("response content: %v", resp.Choices[0].Message.Content))
		return types.MealParts{}, err
	}

	return mealParts, nil
}
