package controllers

import (
	"fmt"
	//"mime"
	"bytes"
	"io"
	"net/http"
    "go.uber.org/zap"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/Logger"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/Error"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/models"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/json-iterator/go"
)

func SendRequest(c *fiber.Ctx) error {

	// Parse the request body into the Request struct
	// This will automatically bind the JSON body to the struct fields
	var request models.Request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("LLM Client Package Controllers", err, "Failed to parse request body"))
	}
    
	// Check if the Prompt is provided
	if request.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("LLM Client Package Controllers", fmt.Errorf("prompt not specified"), "Prompt not specified in request body"))
	}

    //Variable for giving to the adapter-service
    var response models.Response

    //Get image from the request as multipart/form-data
	image,err:=c.FormFile("image")
     
	//If image if present then process it 
	if err==nil{
	//This takes the images and converts to base64
    //Returns the base64 string or error
    response.Images,err=ProcessImage(image)
	if err!=nil{
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to process image"))
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("LLM Client Package Controllers", err, "Failed to process image"))
     }
    }
    
	//Prompt to be given to the model
    response.Prompt=request.Prompt+" give the entire response in json format"  
    
    //Hardcoded model name
	//This is the model name that is used in the adapter service
	//Stream is set to false
	response.Model="llama3:8b"
	fmt.Println(response.Model) //This is the model name

     logger.Log.Info("LLM Client Package Controllers MODELNAME and PROMPT", zap.String("Model", response.Model), zap.String("Prompt", response.Prompt))
	//Convert to Json format the prompt
    reqBody, err := jsoniter.Marshal(response)
	if err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to marshal request to JSON"))
		return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("LLM Client Package Controllers", err, "Failed to marshal request to JSON"))
	}

	// Encrypt the JSON
    encryptedBody, err := utils.Encrypt(string(reqBody))
    if err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to encrypt request body"))
	    return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("LLM Client Package Controllers", err, "Failed to encrypt request body"))
    }
    
	//Send a POST request to the adapter service
	resp, err := http.Post("http://127.0.0.1:8001/modelRequest", "application/json", bytes.NewBuffer([]byte(encryptedBody)))
	if err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to send request to adapter service"))
		return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("LLM Client Package Controllers", err, "Failed to send request to adapter service"))
	}
	defer resp.Body.Close()

    //REad the response body

	responsefromadapter, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to read response body from adapter service"))
		return  c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("LLM Client Package Controllers", err, "Failed to read response body from adapter service"))
	}
	logger.Log.Info("LLM Client Package Controllers", zap.String("Response from adapter service", string(responsefromadapter)))
    
	fmt.Println("Raw adapter response:", string(responsefromadapter))

	//Take the reposne body and put in reponsefromadapter
	//resp return a string not a json object
	
  return c.JSON(fiber.Map{
	"response": string(responsefromadapter),
})

} 

