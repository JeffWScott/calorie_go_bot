package db

import (
	"bytes"
	"calorie_bot/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *ConvexClient) GetSymanticSearchResults(foodNames []string, id string) ([]types.SymanticSearchResult, error) {
	body := map[string]interface{}{
		"path": "foodsVector:getSimilarFoodNamesExternal",
		"args": map[string]interface{}{
			"foodNames": foodNames,
			"API_KEY":   CONVEX_API_KEY,
		},
		"format": "json",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return []types.SymanticSearchResult{}, err
	}

	req, err := http.NewRequest("POST", ACTIONS_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return []types.SymanticSearchResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []types.SymanticSearchResult{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []types.SymanticSearchResult{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []types.SymanticSearchResult{}, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return []types.SymanticSearchResult{}, err
	}

	if status, ok := result["status"].(string); !ok || status != "success" {
		errorMsg := "unknown error"
		if em, ok := result["errorMessage"].(string); ok {
			errorMsg = em
		}
		return []types.SymanticSearchResult{}, fmt.Errorf("action error: %s", errorMsg)
	}

	valueBytes, err := json.Marshal(result["value"])
	if err != nil {
		return []types.SymanticSearchResult{}, err
	}

	var searchResult []types.SymanticSearchResult
	if err := json.Unmarshal(valueBytes, &searchResult); err != nil {
		return []types.SymanticSearchResult{}, err
	}

	return searchResult, nil

}
