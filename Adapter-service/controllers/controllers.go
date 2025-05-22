package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/models"
"github.com/lokesh2201013/genet-microservice/Adapter-service/model-factory"
)

func ModelRequest(c *fiber.Ctx) error {
	var request models.Request

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON( models.Error{
			Status:      400,
			Message:     "Invalid input",
			Description: "Request body is not valid"+ err.Error(),
		})
	}

	if request.Model == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Status:      400,
			Message:     "Model not specified",
			Description: "Model name is required",
		})
	}
	
	if request.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Status:      400,
			Message:     "Prompt not specified",
			Description: "Prompt is required",
		})
	}
    
	request.Stream = false
	adapter,err:= factory.GetModelType(request.Model)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Status:      400,
			Message:     "Invalid model type",
			Description: "Model type is not supported"+ err.Error(),
		})
	}


    var errResponse models.Error
	response,errResponse:= adapter.GenerateResponse(request)
    
	if errResponse!=( models.Error{}) {
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse)
	}

	return c.Status(fiber.StatusOK).JSON(response)
	}