package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	config "github.com/ridhamz/AI-cloud-cost-optimizer-agent/configs"
)

type ClaudeRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

func CallClaude(prompt string) (string, error) {
	apiKey := config.AppConfig.AI.ClaudeAPIKey
	if apiKey == "" {
		return "", fmt.Errorf("Claude API key not set in config")
	}

	reqBody := ClaudeRequest{
		Model:     "claude-v1",
		Prompt:    prompt,
		MaxTokens: 300,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/complete", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Completion string `json:"completion"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Completion, nil
}
