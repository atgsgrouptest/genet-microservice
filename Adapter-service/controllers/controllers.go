package controllers

import (
	"fmt"
	"go.uber.org/zap"
	"github.com/gofiber/fiber/v2"
	"github.com/json-iterator/go"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/Error"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/Logger"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/model-factory"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/models"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/utils"
)

// ModelRequest handles the incoming request for generating a response from the model
// It decrypts the request body, unmarshals it into a Request struct, and then uses the appropriate model adapter to generate a response.
//Uses AES encryption for the request body.
// It returns the response in JSON format or an error if something goes wrong.
func ModelRequest(c *fiber.Ctx) error {

    encryptedBody:=c.Body()
  

	// Decrypt the request body using AES encryption
	decryptedJson, err := utils.Decrypt(string(encryptedBody))
	if err != nil {
		logger.Log.Error("Adapter Service Package controllers", zap.Error(err), zap.String("Message", "Failed to decrypt request body"))
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("Adapter Service Package controllers",err,"Failed to decrypt request body",))
	}
    
	//Convert the decrypted JSON to []bytes and then to  a Request struct
	var request models.Request
	if err := jsoniter.Unmarshal([]byte(decryptedJson), &request); err != nil {
		logger.Log.Error("Adapter Service Package controllers", zap.Error(err), zap.String("Message", "Failed to unmarshal request body to models.Request"))  
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("Adapter Service Package controllers",err,"Failed to unmarshal request body to models.Request",
		))
	}
    
    
     logger.Log.Debug("Adapter Service Package controllers", zap.String("Decrypted JSON", decryptedJson))
	// Check if the model and prompt are specified in the request and model is multimodal
	if request.Model == "" {
		logger.Log.Error("Adapter Service Package controllers", zap.Error(fmt.Errorf("model not specified")), zap.String("Message", "Model not specified in request body"))
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("Adapter Service Package controllers",fmt.Errorf("model not specified"),"Model not specified in request body",
		))
	}

	if request.Prompt == "" {
		logger.Log.Error("Adapter Service Package controllers", zap.Error(fmt.Errorf("prompt not specified")), zap.String("Message", "Prompt not specified in request body"))
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("Adapter Service Package controllers",fmt.Errorf("prompt not specified"),"Prompt not specified in request body",))
	}

    if request.Model=="llama3:8b"&& request.Images!=""{
		logger.Log.Error("Adapter Service Package controllers", zap.Error(fmt.Errorf("Images not supported for llama3:8b model")), zap.String("Message", "Images not supported for llama3:8b model"))						
       return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("Adapter Service Package controllers",fmt.Errorf("Images not supported for llama3:8b model"),"Images not supported for llama3:8b model",))
	}
    
	// Set the stream to false as per the request
	request.Stream = false
	
    logger.Log.Debug("Adapter Service Package controllers", zap.String("Model", request.Model), zap.String("Prompt", request.Prompt), zap.Any("Images", request.Images))
	// Get the appropriate model adapter based on the model specified in the request
	adapter,err:= factory.GetModelType(request.Model)

	if err != nil {
		logger.Log.Error("Adapter Service Package controllers", zap.Error(err), zap.String("Message", "Invalid model type specified in request body"))
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("Adapter Service Package controllers",err,"Invalid model type specified in request body",))
	}

	// Generate the response using the model adapter
    var errResponse models.Error
	response,errResponse:= adapter.GenerateResponse(request)
    
	if errResponse!=( models.Error{}) {
		logger.Log.Error("Adapter Service Package controllers", zap.Error(fmt.Errorf(errResponse.Message)),zap.String("Message", "Error generating response from model"))
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse)
	}

	logger.Log.Debug("Adapter Service Package controllers", zap.String("Model", request.Model), zap.String("Prompt", request.Prompt), zap.String("Response", response))
	return c.Status(fiber.StatusOK).JSON(response)
	}