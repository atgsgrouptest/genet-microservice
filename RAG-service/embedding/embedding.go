package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
    "sync"
	//"sort"
	"strings"
	//"math"

	//"github.com/atgsgrouptest/genet-microservice/RAG-service/Logger"
	"github.com/atgsgrouptest/genet-microservice/RAG-service/Logger"
	"github.com/atgsgrouptest/genet-microservice/RAG-service/models"
	"github.com/atgsgrouptest/genet-microservice/RAG-service/utils"

	//"github.com/json-iterator/go"
	"go.uber.org/zap"
)

type EmbeddedDocument struct {
	Content   string
	Embedding []float64
}

type errorHandler struct{}

var Error errorHandler

func (e errorHandler) ReturnError(location string, err error, msg string) models.Error {
	fmt.Printf("[%s] %s: %v\n", location, msg, err)
	return models.Error{Message: msg}
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
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Generated embedding length: %d\n", len(result.Embedding))
	return result.Embedding, nil
}
func EmbedFileToCorpus(file *multipart.FileHeader) ([]EmbeddedDocument, error) {
	openedFile, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer openedFile.Close()

	var content string
	if strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		text, err := utils.ExtractTextFromPDF(openedFile)
		if err != nil {
			return nil, fmt.Errorf("failed to extract text from PDF: %w", err)
		}
		content = text
	} else {
		contentBytes, err := io.ReadAll(openedFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		content = string(contentBytes)
	}

	// Split document into chunks
	chunks := SplitDocument(content, 250, 20)

	var wg sync.WaitGroup
	corpusCh := make(chan EmbeddedDocument, len(chunks))

	// Limit concurrency (e.g., max 5 at a time)
	const maxConcurrent = 5
	sem := make(chan struct{}, maxConcurrent)

	for _, chunk := range chunks {
		wg.Add(1)
		sem <- struct{}{} // acquire slot
		go func(ch string) {
			defer wg.Done()
			defer func() { <-sem }() // release slot

			vec, err := EmbedText(ch)
			if err != nil {
				logger.Log.Warn("Embedding failed for chunk", zap.String("chunk", ch), zap.Error(err))
				return
			}
			corpusCh <- EmbeddedDocument{
				Content:   ch,
				Embedding: vec,
			}
		}(chunk)
	}

	wg.Wait()
	close(corpusCh)

	var corpus []EmbeddedDocument
	for doc := range corpusCh {
		corpus = append(corpus, doc)
	}

	fmt.Println("Embedded corpus:", len(corpus), "chunks")
	return corpus, nil
}
