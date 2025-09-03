package acloud

import (
	
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestGetSelectorFromAI(t *testing.T) {
	// Create a mock server to simulate the OpenAI API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "gpt-3.5-turbo-0613",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": "button:contains('Start Sandbox')"
					},
					"finish_reason": "stop"
				}
			]
		}`)
	}))
	defer mockServer.Close()

	// Set the OPENAI_API_KEY environment variable to a dummy value
	os.Setenv("OPENAI_API_KEY", "test-key")

	// Create a new OpenAI client with the mock server URL
	config := openai.DefaultConfig("test-key")
	config.BaseURL = mockServer.URL
	client := openai.NewClientWithConfig(config)

	// Sample HTML content
	htmlContent := `
		<html>
			<body>
				<button>Start Sandbox</button>
			</body>
		</html>
	`

	// Call the getSelectorFromAI function
	selector, err := getSelectorFromAIWithClient(htmlContent, client)

	// Assert that there is no error and the selector is correct
	assert.NoError(t, err)
	assert.Equal(t, "button:contains('Start Sandbox')", selector)
}


