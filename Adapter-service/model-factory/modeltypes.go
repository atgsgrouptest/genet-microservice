package factory

import (
		"bytes"
	
	"fmt"
	"io"
	
	"net/http"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/lokesh2201013/genet-microservice/Adapter-service/Error"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/models"

)

//Each model type implements the ModelAdapter interface
//ModelAdapter is an interface that defines the method GenerateResponse
// It takes a models.Request as input and returns a string response and models.Error
// It is implemented by the individual response struct of each model

type gemma3Adapter struct{}

func (g *gemma3Adapter) GenerateResponse(request models.Request) (string, models.Error) {
  
reqBody, err := jsoniter.Marshal(request)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory",err,"Failed to marshal request to JSON")
	}
     
	fmt.Println(request)
	//Send a POST request to the model server that is currently running locally 
	// The URL is hardcoded to http://localhost:11434/api/generate
	// The content type is set to application/json
	resp, err := http.Post("http://"+os.Getenv("LLM_VM_HOST")+":11434/api/chat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "models.Response{}", Error.ReturnError("Adapter Service Package factory", err, "Failed to send request to model server")
	}
	defer resp.Body.Close()
    
	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to read response body")
	}
    fmt.Println("Response BODY:", string(body))
	//Unmarshal the response body from []bytes to json to the Response struct
	var response models.Response
	if err :=  jsoniter.Unmarshal(body, &response); err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to unmarshal response JSON")
	}
    fmt.Println("Raw adapter response:", string(body))
	return response.Message.Content, models.Error{}
}

type gpt4oMiniAdapter struct{}

func (g *gpt4oMiniAdapter) GenerateResponse(request models.Request) (string, models.Error) {
	type ImageURL struct {
		URL string `json:"url"`
	}

	type ChatContent struct {
		Type     string    `json:"type"`               // "text" or "image_url"
		Text     string    `json:"text,omitempty"`     // required if type == "text"
		ImageURL *ImageURL `json:"image_url,omitempty"` // required if type == "image_url"
	}

	type ChatMessage struct {
		Role    string        `json:"role"`
		Content []ChatContent `json:"content"`
	}

	type ChatRequest struct {
		Model       string        `json:"model"`
		Messages    []ChatMessage `json:"messages"`
		MaxTokens   int           `json:"max_tokens,omitempty"`
		Temperature float32       `json:"temperature,omitempty"`
	}

	type ChatResponse struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Index        int `json:"index"`
			Message      struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage interface{} `json:"usage"`
	}

	// Build message content
	content := []ChatContent{
		{
			Type: "text",
			Text: request.Prompt,
		},
	}

	// Add image_url only if valid
	if strings.TrimSpace(request.Images) != "" {
		content = append(content, ChatContent{
			Type: "image_url",
			ImageURL: &ImageURL{
				URL: request.Images,
			},
		})
	}

	// Build the full chat request
	chatRequest := ChatRequest{
		Model:       "gpt-4.1-nano",
		Messages:    []ChatMessage{{Role: "user", Content: content}},
		MaxTokens:   1024,
		Temperature: 0,
	}

	// Marshal the request
	reqBody, err := jsoniter.Marshal(chatRequest)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to marshal ChatRequest")
	}

	fmt.Println("Final request body:", string(reqBody))

	// Create HTTP request to OpenAI
	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to create HTTP request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("OPEN_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to send request to GPT-4o-mini server")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to read response body")
	}

	fmt.Println("Response body:", string(body))

	// Unmarshal the response
	var chatResponse ChatResponse
	if err := jsoniter.Unmarshal(body, &chatResponse); err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to unmarshal GPT-4o-mini response")
	}

	if len(chatResponse.Choices) == 0 {
		return "", Error.ReturnError("Adapter Service Package factory", fmt.Errorf("no choices returned"), "Empty response from GPT-4o-mini")
	}

	return chatResponse.Choices[0].Message.Content, models.Error{}
}

type llama3Adapter struct{}
func (l *llama3Adapter) GenerateResponse(request models.Request) (string, models.Error) {
	if request.Prompt == "" {
		return "", models.Error{
			ServiceName: "llama3Adapter",
			Message:     "Invalid prompt",
			Description: "Prompt cannot be empty",
		}
	}

	// Chat format payload
	payload := map[string]interface{}{
		"model": request.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": request.Prompt,
			},
		},
		"stream": request.Stream,
	}
	if request.Images != "" {
		payload["images"] = []string{request.Images}
	}

	bodyBytes, err := jsoniter.Marshal(payload)
	if err != nil {
		return "", models.Error{
			ServiceName: "llama3Adapter",
			Message:     "JSON encoding error",
			Description: err.Error(),
		}
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", models.Error{
			ServiceName: "llama3Adapter",
			Message:     "Chat API request failed",
			Description: err.Error(),
		}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", models.Error{
			ServiceName: "llama3Adapter",
			Message:     "Failed to read response",
			Description: err.Error(),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return "", models.Error{
			ServiceName: "llama3Adapter",
			Message:     "Non-200 response from Ollama chat API",
			Description: string(respBody),
		}
	}

	var apiResp models.Response
	if err := jsoniter.Unmarshal(respBody, &apiResp); err != nil {
		return "", models.Error{
			ServiceName: "llama3Adapter",
			Message:     "JSON decode error",
			Description: err.Error(),
		}
	}

	return apiResp.Message.Content, models.Error{}
}
