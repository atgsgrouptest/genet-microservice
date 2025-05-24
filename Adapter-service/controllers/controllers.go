package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/json-iterator/go"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/model-factory"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/models"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/utils"
)

func ModelRequest(c *fiber.Ctx) error {
    encryptedBody:=c.Body()

	decryptedJson, err := utils.Decrypt(string(encryptedBody))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			ServiceName: "Adapter Service",
			Message:     "Decryption Failed",
			Description: err.Error(),
		})
	}

	var request models.Request
	if err := jsoniter.Unmarshal([]byte(decryptedJson), &request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			ServiceName: "Adapter Service",
			Message:     "Invalid input",
			Description: "Failed to parse decrypted JSON: " + err.Error(),
		})
	}
    fmt.Println("Decrypted JSON:", decryptedJson)

	if request.Model == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			ServiceName: "Adapter Service",
			Message:     "Model not specified in adapter service",
			Description: "Model name is required",
		})
	}
	
	if request.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			ServiceName: "Adapter Service",
			Message:     "Prompt not specified",
			Description: "Prompt is required",
		})
	}
    
	request.Stream = false
	adapter,err:= factory.GetModelType(request.Model)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			ServiceName: "Adapter Service",
			Message:     "Invalid model type",
			Description: "Model type is not supported"+ err.Error(),
		})
	}


    var errResponse models.Error
	response,errResponse:= adapter.GenerateResponse(request)
    
	if errResponse!=( models.Error{}) {
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse)
	}

	fmt.Println("Response from model:", response)
	return c.Status(fiber.StatusOK).JSON(response)
	}