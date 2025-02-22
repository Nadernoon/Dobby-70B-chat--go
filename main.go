package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const apiURL = "https://api.fireworks.ai/inference/v1/chat/completions"

// RequestPayload represents the JSON payload sent to the Fireworks API
type RequestPayload struct {
	Model            string `json:"model"`
	MaxTokens        int    `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP            float64 `json:"top_p"`
	TopK            int    `json:"top_k"`
	PresencePenalty float64 `json:"presence_penalty"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	Messages        []Message `json:"messages"`
}

// Message represents a message in the conversation history
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ResponsePayload represents the JSON response from the Fireworks API
type ResponsePayload struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func getAPIKey() string {
	apiKey := os.Getenv("FIREWORKS_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: Missing FIREWORKS_API_KEY. Set it in the environment variables.")
		os.Exit(1)
	}
	return apiKey
}

func sendRequest(messages []Message) string {
	apiKey := getAPIKey()
	
	payload := RequestPayload{
		Model:            "accounts/sentientfoundation/models/dobby-unhinged-llama-3-3-70b-new",
		MaxTokens:        1024,
		Temperature:      0.7,
		TopP:            1.0,
		TopK:            40,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
		Messages:        messages,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return ""
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonPayload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %d %s\n", resp.StatusCode, resp.Status)
		return ""
	}

	var responsePayload ResponsePayload
	json.NewDecoder(resp.Body).Decode(&responsePayload)

	if len(responsePayload.Choices) > 0 {
		return responsePayload.Choices[0].Message.Content
	}

	return "No response from Dobby."
}

func chatLoop() {
	fmt.Println("=== Chat with Dobby-70B via Fireworks (Go) ===")
	fmt.Println("Type 'exit' to end the conversation.")

	messages := []Message{
		{Role: "system", Content: "You are Dobby-70B, an AI assistant with an 'unhinged' personality."},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()

		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Exiting.")
			break
		}

		messages = append(messages, Message{Role: "user", Content: userInput})
		response := sendRequest(messages)
		messages = append(messages, Message{Role: "assistant", Content: response})

		fmt.Println("Dobby:", response)
	}
}

func main() {
	chatLoop()
}
