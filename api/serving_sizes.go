package api

import (
	"calorie_bot/types"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

var SERVING_SIZE_SYSTEM_PROMPT = `
## Task
This is a conversation where you had previously requested the nutritional information about parts of a meal.
The User has taken your meal_parts list, looked in their semantic database and provided you a list of foods they have nutritional information for. 
Look at the list and request any items you think that will help you later on to determine the nutritional information of the user's meal.

## Rules
- Don't request things that don't make sense for the meal or the meal_parts
- Only request items that are relevant to the meal; eg.  If the user says they had Carrot Cake for dessert don't just select a random cake if Carrot Cake isn't listed. Not selecting something is okay. 
-  Avoid redundant choices; eg. selecting both "Ground Beef" and "Hamburger Meat" to satisfy a mean_part of "ground beef"
- It's okay if there isn't any semantic matches for all meal_parts, we'll deal with those later on.
- Request only ONE serving size per foodName and provide the servingID for it.
- Pick the serving size that will help you the best later to estimate the nutrition
- Make sure the match the servingID with the correct serving, this is very important because the user will have use that servingID to lookup the nutritional information in their database.

## Examples
{
  "meal_parts": [
    "spaghetti pasta",
    "tomato sauce",
    "ground beef",
    "smart sweets peach rings"
  ]
}

User provides:
[
  {
    "foodName": "Pizza",
    "servings": [
      {
		"servingId": "abc",
        "servingSize": "1 piece",
        "servingWeight": "119g"
      },
      {
	  	"servingId": "def",
        "servingSize": "1 piece",
        "servingWeight": "119g"
      }
},
  {
    "foodName": "Ground Beef",
    "servings": [
      {
		"servingId": "ghi",
        "servingSize": "1 oz, cooked",
        "servingWeight": "28.35g"
      },
      {
	  	"servingId": "jkl",
        "servingSize": "1 oz, raw (yield after cooking)",
        "servingWeight": "20g"
      },
},
  {
    "foodName": "Hamburger Meat",
    "servings": [
      {
		"servingId": "mno",
        "servingSize": "1 oz, cooked",
        "servingWeight": "28.35g"
      },
      {
	  	"servingId": "pqr",
        "servingSize": "1 oz, raw (yield after cooking)",
        "servingWeight": "20g"
      },
}
]

You would choose: "Ground Beef", "1 oz, cooked" as it is the only option that help us. Pizza is irrelevant and Hamburger Meat is redundant to Ground Beef and also a worse choice.
`

var SERVING_SIZE_PREAMPLE = `
I have done a semantic search in my nutritional database  and this is the list of ingredients that might match the meal_parts you have identified.
Return to me the your choices from this list that you think will help you with estimating the nutrients of my meal.
`

func (c *APIClient) GetServingSizes(
	semanticSearchResults []types.SymanticSearchResult,
	id string,
) ([]types.ServingSizeChoices, error) {
	jsonBytes, err := json.MarshalIndent(semanticSearchResults, "", "  ")
	if err != nil {
		return []types.ServingSizeChoices{}, err
	}
	userContent := SERVING_SIZE_PREAMPLE + "\n" + string(jsonBytes)

	req := openai.ChatCompletionRequest{
		Model: "google/gemini-2.5-flash-lite",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: SERVING_SIZE_SYSTEM_PROMPT,
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
				Schema: JSONSchemaMap(servingSizesSchema),
			},
		},
	}

	resp, err := c.CreateChatCompletion(context.Background(), req)
	if err != nil {
		c.Log(id, fmt.Sprintf("Chat completion error: %v", err))
		return []types.ServingSizeChoices{}, err
	}

	if len(resp.Choices) == 0 {
		return []types.ServingSizeChoices{}, nil
	}

	content := resp.Choices[0].Message.Content
	//c.Log(id, fmt.Sprintf("response content: %v", resp.Choices[0].Message.Content))

	var servingSizeChoices []types.ServingSizeChoices
	if err := json.Unmarshal([]byte(content), &servingSizeChoices); err != nil {
		c.Log(id, fmt.Sprintf("JSON unmarshal error: %v", err))
		c.Log(id, fmt.Sprintf("response content: %v", resp.Choices[0].Message.Content))
		return []types.ServingSizeChoices{}, err
	}

	return servingSizeChoices, nil
}
