package controllers

import (
	"fmt"
	"os"
	//"mime"
	"bytes"
	"io"
	"net/http"

	"github.com/atgsgrouptest/genet-microservice/LLM-client/Error"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/Logger"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/models"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/json-iterator/go"
	"go.uber.org/zap"
)

func SendRequest(c *fiber.Ctx) error {

	// Parse the request body into the Request struct
	// This will automatically bind the JSON body to the struct fields
	var request models.Request
	if err := c.BodyParser(&request); err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to parse request body"))
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("LLM Client Package Controllers", err, "Failed to parse request body"))
	}
    
	// Check if the Prompt is provided
	if request.Prompt == "" {
		logger.Log.Error("LLM Client Package Controllers", zap.String("Message", "Prompt not specified in request body"))
		return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("LLM Client Package Controllers", fmt.Errorf("prompt not specified"), "Prompt not specified in request body"))
	}

    //Variable for giving to the adapter-service
    var response models.Response
    
	//Prompt to be given to the model
    response.Prompt=request.Prompt
	//Hardcoded model name
	//This is the model name that is used in the adapter service
	//Stream is set to false
	response.Model="llama3.1:8b" //This is the model name that is used in the adapter service
	logger.Log.Debug("LLM Client Package Controllers",zap.String("Name of model = ",response.Model)) //This is the model name
	 responsefromadapter,error:= ReturnResponse(response)
	  if error!= (models.Error{}) {
		logger.Log.Error("LLM Client Package Controllers", zap.String("Message", "Failed to get response from adapter service"))
		return c.Status(fiber.StatusInternalServerError).JSON(error.ServiceName,error.Message,error.Description)
	  }
	

	return c.JSON(fiber.Map{
	"response": string(responsefromadapter),
}) 

} 

func SendRequestWithImages(c *fiber.Ctx) error{
	//Variable for giving to the adapter-service
    var response models.Response
    
	
	response.Prompt = c.FormValue("prompt") //Get the prompt from the form data
	logger.Log.Debug("LLM Client Package Controllers", zap.String("Prompt", response.Prompt))
    fmt.Println("Prompt:", response.Prompt)
	image, err := c.FormFile("image")

if err != nil && err != http.ErrMissingFile {
    logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to retrieve image"))
    return c.Status(fiber.StatusBadRequest).JSON(Error.ReturnError("LLM Client Package Controllers", err, "Failed to retrieve image"))
}

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

	response.Images = "data:image/png;base64," + response.Images
  fmt.Println("Image Base64 String:", response.Images)
	 //Hardcoded model name
	//This is the model name that is used in the adapter service
	//Stream is set to false
	response.Model="gpt-4o-mini" //This is the model name that is used in the adapter service

	  responsefromadapter,error:= ReturnResponse(response)
	  if error!= (models.Error{}) {
		logger.Log.Error("LLM Client Package Controllers", zap.String("Message", "Failed to get response from adapter service"))
		return c.Status(fiber.StatusInternalServerError).JSON(error.ServiceName,error.Message,error.Description)
	  }
	

	logger.Log.Info("LLM Client Package Controllers", zap.String("Response from adapter service", string(responsefromadapter)))
	return c.JSON(fiber.Map{
	"response": string(responsefromadapter),
}) 
	

}

func ReturnResponse(response models.Response) (string , models.Error) {
	logger.Log.Debug("LLM Client Package Controllers MODELNAME and PROMPT", zap.String("Model", response.Model), zap.String("Prompt", response.Prompt))
	//Convert to Json format the prompt
    reqBody, err := jsoniter.Marshal(response)
	if err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to marshal request to JSON"))
		return "",Error.ReturnError("LLM Client Package Controllers",err , "Failed to marshal request to JSON")
	}

	// Encrypt the JSON
    encryptedBody, err := utils.Encrypt(string(reqBody))
    if err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to encrypt request body"))
	    return "",Error.ReturnError("LLM Client Package Controllers", err, "Failed to encrypt request body")
    }
    
	//Send a POST request to the adapter service
	resp, err := http.Post(os.Getenv("ADAPTER_SERVICE_HOST")+"/modelRequest", "application/json", bytes.NewBuffer([]byte(encryptedBody)))
	if err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to send request to adapter service"))
		return "",Error.ReturnError("LLM Client Package Controllers",err , "Failed to send request to adapter service")
	}
	defer resp.Body.Close()

    //REad the response body

	responsefromadapter, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("LLM Client Package Controllers", zap.Error(err), zap.String("Message", "Failed to read response body from adapter service"))
		return  "",Error.ReturnError("LLM Client Package Controllers", err, "Failed to read response body from adapter service")
	}
	logger.Log.Debug("LLM Client Package Controllers", zap.String("Response from adapter service", string(responsefromadapter)))

	//Take the reposne body and put in reponsefromadapter
	//resp return a string not a json object
	fmt.Println("Response from adapter service:", string(responsefromadapter))
	return string(responsefromadapter),models.Error{}
}

