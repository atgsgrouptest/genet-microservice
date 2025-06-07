package factory

import (
		"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/lokesh2201013/genet-microservice/Adapter-service/Error"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/models"
)

//Each model type implements the ModelAdapter interface
//ModelAdapter is an interface that defines the method GenerateResponse
// It takes a models.Request as input and returns a string response and models.Error
// It is implemented by the individual response struct of each model
type llama3Adapter struct{}

//This is the member function of the llama3Adapter struct
type EmbeddedDocument struct {
	Content   string
	Embedding []float64
}


type errorHandler struct{}

func (e errorHandler) ReturnError(location string, err error, msg string) models.Error {
	fmt.Printf("[%s] %s: %v\n", location, msg, err)
	return models.Error{Message: msg}
}

type modelsRequest struct {
	Model  string
	Prompt string
}

type modelsResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type modelsError struct {
	Message string
}

func FetchOpenAPIDocument(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func SplitDocument(content string, chunkSize, overlap int) []string {
	var chunks []string
	for i := 0; i < len(content); i += chunkSize - overlap {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}
		chunks = append(chunks, content[i:end])
	}
	return chunks
}

func EmbedText(chunk string) ([]float64, error) {
	payload := map[string]string{
		"model":  "nomic-embed-text",
		"prompt": chunk,
	}
	data, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Embedding []float64 `json:"embedding"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Embedding, nil
}

func CosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func RetrieveTopK(question string, corpus []EmbeddedDocument, k int) []string {
	qVec, _ := EmbedText(question)
	sort.Slice(corpus, func(i, j int) bool {
		return CosineSimilarity(qVec, corpus[i].Embedding) > CosineSimilarity(qVec, corpus[j].Embedding)
	})
	var topK []string
	for i := 0; i < k && i < len(corpus); i++ {
		topK = append(topK, corpus[i].Content)
	}
	return topK
}

func (l *llama3Adapter) GenerateResponse(request models.Request) (string, models.Error) {
	docText, err := FetchOpenAPIDocument("https://petstore3.swagger.io/api/v3/openapi.json")
	if err != nil {
		return "", Error.ReturnError("Adapter Service", err, "Failed to fetch document")
	}
	chunks := SplitDocument(docText, 250, 20)

	var corpus []EmbeddedDocument
	for _, chunk := range chunks {
		vec, err := EmbedText(chunk)
		if err != nil {
			continue
		}
		corpus = append(corpus, EmbeddedDocument{Content: chunk, Embedding: vec})
	}

	relevant := RetrieveTopK(request.Prompt, corpus, 4)
	context := strings.Join(relevant, "\n")

	ragPrompt := fmt.Sprintf(`
Question: %s
Documents: %s
Answer:`, request.Prompt, context)

	bodyily := map[string]interface{}{
		"model": request.Model,
		"stream": false,
		"messages": []map[string]string{
			{"role": "user", "content": ragPrompt},
		},
		"temperature": 0,
	}

	reqBody, err := jsoniter.Marshal(bodyily)
	if err != nil {
		return "", Error.ReturnError("Adapter Service", err, "Failed to marshal request to JSON")
	}

	resp, err := http.Post("http://"+os.Getenv("LLM_VM_HOST")+":11434/api/chat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", Error.ReturnError("Adapter Service", err, "Failed to send request to model server")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", Error.ReturnError("Adapter Service", err, "Failed to read response body")
	}

	var response modelsResponse
	if err := jsoniter.Unmarshal(body, &response); err != nil {
		return "", Error.ReturnError("Adapter Service", err, "Failed to unmarshal response JSON")
	}

	fmt.Println("Raw adapter response:", string(body))
	return response.Message.Content, models.Error{}
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
	resp, err := http.Post("http://"+os.Getenv("LLM_VM_HOST")+":11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
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