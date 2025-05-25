package factory

import (
	"bytes"
	"fmt"

	"github.com/json-iterator/go"

	"io"
	"net/http"

	"github.com/lokesh2201013/genet-microservice/Adapter-service/models"
		"github.com/lokesh2201013/genet-microservice/Adapter-service/Error"
)

//Each model type implements the ModelAdapter interface
//ModelAdapter is an interface that defines the method GenerateResponse
// It takes a models.Request as input and returns a string response and models.Error
// It is implemented by the individual response struct of each model
type llama3Adapter struct{}

//This is the member function of the llama3Adapter struct
func (l *llama3Adapter) GenerateResponse(request models.Request) (string, models.Error) {
//This is reponse from the model
//If done is true we are a success
    type Response struct {
	Model                string    `json:"model"`
	CreatedAt            string    `json:"created_at"`
	Response             string    `json:"response"`
	Done                 bool      `json:"done"`
	DoneReason           string    `json:"done_reason"`
	Context              []int     `json:"context"`
	TotalDuration        int64     `json:"total_duration"`
	LoadDuration         int64     `json:"load_duration"`
	PromptEvalCount      int       `json:"prompt_eval_count"`
	PromptEvalDuration   int64     `json:"prompt_eval_duration"`
	EvalCount            int       `json:"eval_count"`
	EvalDuration         int64     `json:"eval_duration"`
}


    //Converts the request to JSON format
	reqBody, err := jsoniter.Marshal(request)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory",err,"Failed to marshal request to JSON")
	}
     
	fmt.Println(request)
	//Send a POST request to the model server that is currently running locally 
	// The URL is hardcoded to http://localhost:11434/api/generate
	// The content type is set to application/json
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "models.Response{}", Error.ReturnError("Adapter Service Package factory", err, "Failed to send request to model server")
	}
	defer resp.Body.Close()
    
	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to read response body")
	}
    
	//Unmarshal the response body from []bytes to json to the Response struct
	var response Response
	if err :=  jsoniter.Unmarshal(body, &response); err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to unmarshal response JSON")
	}
    fmt.Println("Raw adapter response:", string(body))
	return response.Response, models.Error{}
}

type gemma3Adapter struct{}

func (g *gemma3Adapter) GenerateResponse(request models.Request) (string, models.Error) {
	 type Response struct {
	Model                string    `json:"model"`
	CreatedAt            string    `json:"created_at"`
	Response             string    `json:"response"`
	Done                 bool      `json:"done"`
	DoneReason           string    `json:"done_reason"`
	Context              []int     `json:"context"`
	TotalDuration        int64     `json:"total_duration"`
	LoadDuration         int64     `json:"load_duration"`
	PromptEvalCount      int       `json:"prompt_eval_count"`
	PromptEvalDuration   int64     `json:"prompt_eval_duration"`
	EvalCount            int       `json:"eval_count"`
	EvalDuration         int64     `json:"eval_duration"`
}
  
reqBody, err := jsoniter.Marshal(request)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory",err,"Failed to marshal request to JSON")
	}
     
	fmt.Println(request)
	//Send a POST request to the model server that is currently running locally 
	// The URL is hardcoded to http://localhost:11434/api/generate
	// The content type is set to application/json
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "models.Response{}", Error.ReturnError("Adapter Service Package factory", err, "Failed to send request to model server")
	}
	defer resp.Body.Close()
    
	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to read response body")
	}
    
	//Unmarshal the response body from []bytes to json to the Response struct
	var response Response
	if err :=  jsoniter.Unmarshal(body, &response); err != nil {
		return "", Error.ReturnError("Adapter Service Package factory", err, "Failed to unmarshal response JSON")
	}
    fmt.Println("Raw adapter response:", string(body))
	return response.Response, models.Error{}
}