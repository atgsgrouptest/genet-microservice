package controllers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"encoding/json"
	//"strings"
    "context"
   // "log"
   // "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
    //"go.mongodb.org/mongo-driver/mongo"
    //"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/atgsgrouptest/genet-microservice/AI-testing/database"

	"github.com/atgsgrouptest/genet-microservice/AI-testing/Logger"
	"github.com/atgsgrouptest/genet-microservice/AI-testing/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"github.com/atgsgrouptest/genet-microservice/AI-testing/models"
)

func MakePrompt(c *fiber.Ctx) error {
	//language := c.FormValue("language")
	image, err := c.FormFile("image")
	if err != nil {
		logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to parse multipart form"), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Service Name": "AI Testing Package Controllers",
			"error":        "Unable to receive files",
			"details":      err.Error(),
		})
	}

	openimage, err := image.Open()
	if err != nil {
		logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to open image"), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Service Name": "AI Testing Package Controllers",
			"error":        "Unable to open image",
			"details":      err.Error(),
		})
	}
 
	openimagebytes, err := io.ReadAll(openimage)
	if err != nil {
		logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to read image data"), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Service Name": "AI Testing Package Controllers",
			"error":        "Unable to read image data",
			"details":      err.Error(),
		})
	}

	fmt.Println(openimagebytes)

	defer openimage.Close()

	prompt := `Please analyze this flowchart diagram and convert it into a JSON format representing a tree structure.
	Respond with pure JSON. Do not wrap it in triple backticks or code blocks. without any additional text or explanation.
Each node in the tree should have the following fields:

1. name: Title or label of the step/box in the chart.
2. description: Short summary of what that step does.
3. requires_input: Boolean indicating whether user input is needed at this step.
4. children: An array of child steps that follow from this one.

If the flow has terminal points (like 'End'), make those leaf nodes with empty children arrays.
Return the structure as valid JSON that can be traversed using a depth-first search algorithm.`

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add prompt field
	if err := writer.WriteField("prompt", prompt); err != nil {
		logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to write prompt field"), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Service Name": "AI Testing Package Controllers",
			"error":        "Unable to add prompt field",
			"details":      err.Error(),
		})
	}

	// Add image file
	part, _ := writer.CreateFormFile("image", image.Filename)
	_, _ = part.Write(openimagebytes)

	writer.Close()

	var outer models.Outer

	// Send POST request
	resp, err := http.Post("http://localhost:8002/sendRequestImages", writer.FormDataContentType(), &b)
	if err != nil {
		logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to send request"), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Service Name": "AI Testing Package Controllers",
			"error":        "Unable to send request",
			"details":      err.Error(),
		})
	}
	
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	fmt.Println("Response Body:", string(responseBody))


err = json.Unmarshal(responseBody, &outer);if err != nil {
	logger.Log.Error("Failed to unmarshal outer JSON", zap.Error(err))
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"Service Name": "AI Testing Package Controllers",
		"error":        "Invalid JSON format from LLM",
		"details":      err.Error(),
	})
}

	
 //outer.Response ="\"{\\n  \\\"name\\\": \\\"Start\\\",\\n  \\\"description\\\": \\\"Beginning of the process.\\\",\\n  \\\"requires_input\\\": false,\\n  \\\"children\\\": [\\n    {\\n      \\\"name\\\": \\\"Select an option\\\",\\n      \\\"description\\\": \\\"User selects an option to proceed.\\\",\\n      \\\"requires_input\\\": true,\\n      \\\"children\\\": [\\n        {\\n          \\\"name\\\": \\\"Get Sales\\\",\\n          \\\"description\\\": \\\"Proceed to sales information.\\\",\\n          \\\"requires_input\\\": false,\\n          \\\"children\\\": [\\n            {\\n              \\\"name\\\": \\\"Get account details\\\",\\n              \\\"description\\\": \\\"Retrieve account information.\\\",\\n              \\\"requires_input\\\": false,\\n              \\\"children\\\": [\\n                {\\n                  \\\"name\\\": \\\"Input your data\\\",\\n                  \\\"description\\\": \\\"User inputs necessary data.\\\",\\n                  \\\"requires_input\\\": true,\\n                  \\\"children\\\": [\\n                    {\\n                      \\\"name\\\": \\\"Get pricing\\\",\\n                      \\\"description\\\": \\\"Obtain pricing information.\\\",\\n                      \\\"requires_input\\\": false,\\n                      \\\"children\\\": [\\n                        {\\n                          \\\"name\\\": \\\"Contract Support\\\",\\n              \\\"description\\\": \\\"Contact support for assistance.\\\",\\n                          \\\"requires_input\\\": false,\\n                         \\\"children\\\": [\\n                            {\\n                              \\\"name\\\": \\\"Decision\\\",\\n                              \\\"description\\\": \\\"Decision point for further action.\\\",\\n                              \\\"requires_input\\\": false,\\n                              \\\"children\\\": [\\n                                {\\n                                  \\\"name\\\": \\\"Premium pricing\\\",\\n                                  \\\"description\\\": \\\"Select premium pricing option.\\\",\\n                                  \\\"requires_input\\\": false,\\n                                  \\\"children\\\": []\\n                                },\\n                                {\\n                                  \\\"name\\\": \\\"Standard\\\",\\n                                  \\\"description\\\": \\\"Select standard pricing option.\\\",\\n                                  \\\"requires_input\\\": false,\\n                                  \\\"children\\\": []\\n                                },\\n                                {\\n                                   \\\"name\\\": \\\"Enterprise Pricing\\\",\\n                                  \\\"description\\\": \\\"Select enterprise pricing option.\\\",\\n                                  \\\"requires_input\\\": false,\\n                                  \\\"children\\\": []\\n                                },\\n                                {\\n                                  \\\"name\\\": \\\"Indian Support\\\",\\n                                  \\\"description\\\": \\\"Select Indian support option.\\\",\\n                                  \\\"requires_input\\\": false,\\n                                  \\\"children\\\": []\\n                                }\\n                              ]\\n                            }\\n                          ]\\n                    }\\n                      ]\\n                    }\\n                  ]\\n                }\\n              ]\\n            }\\n          ]\\n        }\\n      ]\\n    }\\n  ]\\n}\""

	DFSresult, err := utils.CleanJSON(outer.Response)
	
	if err != nil {
		logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to clean JSON"), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Service Name": "AI Testing Package Controllers",
			"error":        "Unable to clean JSON",
			"details":      err.Error(),
		})
	}

	var allNegativeCases [][]string

    for _, flow := range DFSresult {
		//fmt.Printf("Flow %d: %s\n", i+1, flow)
		type negativeCase struct{
			Prompt string `json:"prompt"`
		}
		var negativeCaseData negativeCase

		negativeCaseData.Prompt = fmt.Sprintf(
    "I will give you a flow for chatbot. You have to return []string with negative case. "+
    "It should include tests for SQL injection and other security vulnerabilities in the flow. "+
	"The response flow should match the input flow with  ->"+
    "The response should be in a JSON format with a single key \"negative_cases\" and the value should be an array of strings. "+
    "No additional text or explanation is needed. Do not wrap it in triple backticks or code blocks.\n\n%s", flow)


        flowEncodedData, err := json.Marshal(negativeCaseData)
		if err != nil {
			logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to marshal negative case data"), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Service Name": "AI Testing Package Controllers",
				"error":        "Unable to marshal negative case data",
				"details":      err.Error(),
			})
		}
		//fmt.Println("Negative Flow Response:", string(flowEncodedData))

		resp, err := http.Post("http://localhost:8002/sendRequest", "application/json", bytes.NewBuffer(flowEncodedData))
		if err != nil {
			logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to send negative case request"), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Service Name": "AI Testing Package Controllers",
				"error":        "Unable to send negative case request",
				"details":      err.Error(),
			})
		}
		defer resp.Body.Close()
		negativeCaseResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to read negative case response"), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Service Name": "AI Testing Package Controllers",
				"error":        "Unable to read negative case response",
				"details":      err.Error(),
			})
		}

		fmt.Println("Negative Case Response:", string(negativeCaseResponse))

		err = json.Unmarshal(negativeCaseResponse, &outer); if err != nil {
			logger.Log.Error("Failed to unmarshal outer JSON", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Service Name": "AI Testing Package Controllers",
				"error":        "Invalid JSON format from LLM",
		"details":      err.Error(),
	})}

	 

		 negativeCaseResult,err := utils.ParseEscapedJSON[models.NegativeCaseResult](outer.Response)

		if err != nil {
			logger.Log.Error("AI Testing Package Controllers", zap.String("Message", "Failed to parse negative case result"), zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Service Name": "AI Testing Package Controllers",
				"error":        "Unable to parse negative case result",
				"details":      err.Error(),
			})
		}

		fmt.Println("Negative Case Result:", negativeCaseResult)
        allNegativeCases = append(allNegativeCases, negativeCaseResult.NegativeCases)
	}
     
	
	
	request := models.Request{
        RequestID:       primitive.NewObjectID(),
        CompanyID:       c.FormValue("company_id"), // Replace with actual User ID
        RequestMaterial: openimagebytes,
        PositiveCases:   DFSresult,
        NegativeCases:  allNegativeCases,
    }
collection := database.MongoDB.Collection("requests")
	ctx := context.TODO()
	_, err = collection.InsertOne(ctx, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert request",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"negative_cases": allNegativeCases,
	})

}