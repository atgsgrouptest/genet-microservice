package controllers

import (
	// "github.com/atgsgrouptest/genet-microservice/RAG-service/Error"
	//"fmt"

	"fmt"

	"github.com/atgsgrouptest/genet-microservice/RAG-service/Logger"
	"github.com/atgsgrouptest/genet-microservice/RAG-service/embedding"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func SendFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		logger.Log.Error("RAG Service Package Controllers", zap.String("Message", "Failed to parse multipart form"), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Service Name": "RAG Service Package Controllers",
			"error":        "Unable to receive files",
			"details":      err.Error(),
		})
	}

	files := form.File["files"]
	//var resp []embedding.EmbeddedDocument
	for _, file := range files {
		corpus, _ := embedding.EmbedFileToCorpus(file)
		//resp = append(resp, corpus...)

		// ✅ Store in Qdrant
		err := embedding.StoreInQdrant(corpus, "rag_corpus")
		if err != nil {
			logger.Log.Error("Failed to store in Qdrant", zap.Error(err))
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Files embedded and stored successfully",
	})
}


func GetPromptWithContext(c *fiber.Ctx) error {
	type reqBody struct {
		Query string `json:"query"`
	}
	var body reqBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	// Search in Qdrant
	chunks, err := embedding.SearchInQdrant(body.Query, "rag_corpus", 5)
	if err != nil {
		logger.Log.Error("Failed to search Qdrant", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate response with LLM
	answer, err := embedding.AskLlama(chunks, body.Query)
	if err != nil {
		logger.Log.Error("Failed to query LLM", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("Generated answer:", answer)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"answer": answer,
	})
}
