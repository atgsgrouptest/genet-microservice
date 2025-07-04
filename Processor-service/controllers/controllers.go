package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/atgsgrouptest/genet-microservice/Processor-service/Error"
	"github.com/atgsgrouptest/genet-microservice/Processor-service/Logger"
	"github.com/atgsgrouptest/genet-microservice/Processor-service/Prompt"
	"github.com/atgsgrouptest/genet-microservice/Processor-service/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetAPISequence(c *fiber.Ctx) error {
	var incommingRequest models.Request
   var body []byte
   var responseFromLLM models.LLMResponse
	if err := c.BodyParser(&incommingRequest); err != nil {
		logger.Log.Error("Error parsing request body", zap.Error(err), zap.String("ServiceName", "Processor-service Package controllers"))
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("Processor-service Package controllers", err, "Error parsing request body"))
	}
	logger.Log.Debug("Received request", zap.Any("Request", incommingRequest), zap.String("ServiceName", "Processor-service Package controllers"))

	//var allResponses []map[string]interface{}

	for _, swaggerObject := range incommingRequest.SwaggerRequest {
		logger.Log.Debug("Processing Swagger Object", zap.String("HostURL", swaggerObject.HostURL), zap.String("SwaggerJSONLink", swaggerObject.SwaggerJSONLink), zap.String("ServiceName", "Processor-service Package controllers"))

		responseswaggerjson, err := http.Get(swaggerObject.SwaggerJSONLink)
		if err != nil {
			logger.Log.Error("Error fetching Swagger JSON", zap.Error(err), zap.String("ServiceName", "Processor-service Package controllers"))
			return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("Processor-service Package controllers", err, "Error fetching Swagger JSON"))
		}
		defer responseswaggerjson.Body.Close()

		var swaggerObjectData models.OpenAPI
		if err := json.NewDecoder(responseswaggerjson.Body).Decode(&swaggerObjectData); err != nil {
			logger.Log.Error("Error decoding Swagger JSON", zap.Error(err), zap.String("ServiceName", "Processor-service Package controllers"))
			return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("Processor-service Package controllers", err, "Error decoding Swagger JSON"))
		}

		logger.Log.Debug("Decoded Swagger JSON", zap.Any("SwaggerObjectData", swaggerObjectData), zap.String("ServiceName", "Processor-service Package controllers"))
		fmt.Printf("Decoded Swagger JSON: %+v\n", swaggerObjectData)

		// Manually building the OpenAPIRequest here as you do
		OpenAPIRequest := models.OpenAPIRequest{
			Request: []models.HttpRequest{
				{
					URL:         "https://petstore.swagger.io/v2/pet",
					HTTPMethod:  "POST",
					Description: "Pet object that needs to be added to the store",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"id": 0,
						"category": map[string]interface{}{
							"id":   0,
							"name": "string",
						},
						"name":      "doggie",
						"photoUrls": []string{"string"},
						"tags": []map[string]interface{}{
							{
								"id":   0,
								"name": "string",
							},
						},
						"status": "available",
					},
				},
				{
					URL:         "https://petstore.swagger.io/v2/pet/{petId}",
					HTTPMethod:  "GET",
					Description: "ID of pet to return",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"id": 0,
						"category": map[string]interface{}{
							"id":   0,
							"name": "string",
						},
						"name":      "doggie",
						"photoUrls": []string{"string"},
						"tags": []map[string]interface{}{
							{
								"id":   0,
								"name": "string",
							},
						},
						"status": "available",
					},
				},
			},
		}

		jsonBytes, err := json.Marshal(OpenAPIRequest)
		if err != nil {
			logger.Log.Error("Error marshalling OpenAPIRequest to JSON", zap.Error(err), zap.String("ServiceName", "Processor-service Package controllers"))
			return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("Processor-service Package controllers", err, "Error marshalling OpenAPIRequest to JSON"))
		}
		logger.Log.Debug("Marshalled OpenAPIRequest to JSON", zap.String("OpenAPIRequestJSON", string(jsonBytes)), zap.String("ServiceName", "Processor-service Package controllers"))

		APISequencePrompt, err := prompt.GeneratePromptAPISequence(string(jsonBytes))
		fmt.Println(APISequencePrompt)
		if err != nil {
			logger.Log.Error("Error generating API sequence prompt", zap.Error(err), zap.String("ServiceName", "Processor-service Package controllers"))
			return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("Processor-service Package controllers", err, "Error generating API sequence prompt"))
		}

		logger.Log.Debug("Generated API sequence prompt", zap.String("APISequencePrompt", APISequencePrompt), zap.String("ServiceName", "Processor-service Package controllers"))

		requestBody, _ := json.Marshal(map[string]string{
			"prompt": APISequencePrompt,
		})

		resp, err := http.Post("http://127.0.0.1:8002/sendRequest", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			logger.Log.Error("Error calling prompt service", zap.Error(err), zap.String("ServiceName", "Processor-service Package controllers"))
			return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("Processor-service Package controllers", err, "Error calling prompt service"))
		}
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
        fmt.Println(string(body))

		
		if err:=json.Unmarshal(body, &responseFromLLM);err!=nil{
       return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("Processor-service Package controllers", err, "Error paarsing prompt service"))
	}
	}

	// Return the **clean** APIWrapper JSON response (if multiple swagger requests, combine here as you wish)
	// For now, return allResponses array (if you want to flatten or merge, you can do it here)
	return c.Status(fiber.StatusOK).JSON(responseFromLLM)
}
