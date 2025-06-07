package httpclient

import (
	"encoding/json"
	//"fmt"
	//"strconv"
	"fmt"
	//"strings"
	"io"
	"log"
	"net/http"
	"bytes"
	"github.com/atgsgrouptest/genet-microservice/Processor-service/models"
	//"regexp"
)

// Define the structs to unmarshal into

// HTTPRequest parses the nested LLM response JSON and returns the extracted APIWrapper

// Struct for each API request
type APIRequest struct {
	SequenceNumber       int               `json:"SequenceNumber"`
	Description          string            `json:"Description"`
	URL                  string            `json:"URL"`
	Path                 string            `json:"Path"`
	HTTPMethod           string            `json:"HTTPMethod"`
	ContentType          string            `json:"ContentType"`
	Headers              map[string]string `json:"Headers"`
	RequestBody          map[string]any    `json:"RequestBody"`
	ExpectedResponseCode string            `json:"ExpectedResponseCode"`
}

// Wrapper to hold API requests
type APIWrapper struct {
	APIs []APIRequest `json:"apis"`
}

func ParseLLMResponse(rawInput string) (*APIWrapper, error) {
	// Step 1: Parse the outer JSON
	var outer struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal([]byte(rawInput), &outer); err != nil {
		return nil, fmt.Errorf("error unmarshalling outer JSON: %w", err)
	}

	// Step 2: Unescape the JSON string (remove outer quotes)
	var innerJSON string
	if err := json.Unmarshal([]byte(outer.Response), &innerJSON); err != nil {
		return nil, fmt.Errorf("error unmarshalling quoted JSON string: %w", err)
	}

	// Step 3: Unmarshal the unescaped string into []APIRequest
	var apiRequests []APIRequest
	if err := json.Unmarshal([]byte(innerJSON), &apiRequests); err != nil {
		return nil, fmt.Errorf("error unmarshalling inner JSON array into []APIRequest: %w", err)
	}

	// Step 4: Wrap and return
	return &APIWrapper{APIs: apiRequests}, nil
}


func HTTPRequest(FormatedData models.APIWrapper){
	
	
	results := []map[string]interface{}{}

	for _, api := range FormatedData.APIs {
		// Marshal the body
		bodyBytes, err := json.Marshal(api.RequestBody)
		if err != nil {
			log.Printf("Error marshaling body for API #%d: %v", api.SequenceNumber, err)
			continue
		}

		// Build the HTTP request
		req, err := http.NewRequest(api.HTTPMethod, api.URL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			log.Printf("Error building request #%d: %v", api.SequenceNumber, err)
			continue
		}

		// Set headers
		for key, val := range api.Headers {
			req.Header.Set(key, val)
		}
		req.Header.Set("Content-Type", api.ContentType)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Request #%d failed: %v", api.SequenceNumber, err)
			continue
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		results = append(results, map[string]interface{}{
			"sequenceNumber": api.SequenceNumber,
			"url":            api.URL,
			"statusCode":     resp.StatusCode,
			"responseBody":   string(respBody),
		})
	}

}