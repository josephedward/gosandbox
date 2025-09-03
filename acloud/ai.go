package acloud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func GetLoginSelectorsFromAI(pageContent string) (usernameSelector, passwordSelector, loginButtonSelector string, err error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", "", "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(`Here is the HTML content of a web page. Please provide CSS selectors for the username input field, the password input field, and the login button. Return the selectors as a JSON object with keys "username", "password", and "loginButton". Example: {"username": "#email", "password": "#password", "loginButton": "#submit"}.
HTML:
%s`, pageContent),
				},
			},
		},
	)

	if err != nil {
		return "", "", "", err
	}

	var selectors struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		LoginButton string `json:"loginButton"`
	}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &selectors)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	return selectors.Username, selectors.Password, selectors.LoginButton, nil
}

