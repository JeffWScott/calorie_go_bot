package db

import (
	"bytes"
	"calorie_bot/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *ConvexClient) GetNutritionalInfo(servingSizes []types.ServingSizeChoices, id string) ([]types.NutritionalInfo, error) {
	body := map[string]interface{}{
		"path": "foods:getNutritionalInfoExternal",
		"args": map[string]interface{}{
			"servingSizes": servingSizes,
			"API_KEY":      CONVEX_API_KEY,
		},
		"format": "json",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return []types.NutritionalInfo{}, err
	}

	req, err := http.NewRequest("POST", ACTIONS_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return []types.NutritionalInfo{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []types.NutritionalInfo{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []types.NutritionalInfo{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []types.NutritionalInfo{}, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return []types.NutritionalInfo{}, err
	}

	if status, ok := result["status"].(string); !ok || status != "success" {
		errorMsg := "unknown error"
		if em, ok := result["errorMessage"].(string); ok {
			errorMsg = em
		}
		return []types.NutritionalInfo{}, fmt.Errorf("action error: %s", errorMsg)
	}

	valueBytes, err := json.Marshal(result["value"])
	if err != nil {
		return []types.NutritionalInfo{}, err
	}

	var nutrientsResult []types.NutritionalInfo
	if err := json.Unmarshal(valueBytes, &nutrientsResult); err != nil {
		return []types.NutritionalInfo{}, err
	}

	return nutrientsResult, nil

}
