package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *ConvexClient) ScheduleGetNutrition(foodNames []string, id string) error {
	body := map[string]interface{}{
		"path": "schedules:scheduleFoodNutritionExtenal",
		"args": map[string]interface{}{
			"foodNames": foodNames,
			"API_KEY":   CONVEX_API_KEY,
		},
		"format": "json",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", ACTIONS_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	if status, ok := result["status"].(string); !ok || status != "success" {
		errorMsg := "unknown error"
		if em, ok := result["errorMessage"].(string); ok {
			errorMsg = em
		}
		return fmt.Errorf("action error: %s", errorMsg)
	}

	return nil
}
