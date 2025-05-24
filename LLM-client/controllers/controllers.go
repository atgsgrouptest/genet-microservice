package controllers

import (
	"encoding/base64"
	"fmt"
	"io"

	//"mime"
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"

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
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			ServiceName: "LLM Client",
			Message:     "Invalid input",
			Description: "Request body is not valid" + err.Error(),
		})
	}
    
	// Check if the Prompt is provided
	if request.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			ServiceName: "LLM Client",
			Message:     "Prompt not specified",
			Description: "Prompt is required",
		})
	}

    //Variable for giving to the adapter-service
    var response models.Response

    //Get image from the request as multipart/form-data
	image,err:=c.FormFile("image")
     
	//If image if present then process it 
	if err==nil{
	//This gives the image in base64 format
    response.Images,err=ProcessImage(image)
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			ServiceName: "LLM Client",
			Message:     "Invalid image",
			Description: "Image is not valid"+ err.Error(),
		})
     }
    }
    
	//Prompt to be given to the model
    response.Prompt=request.Prompt+" give the entire response in json format"  
    

	response.Model="llama3:8b"
	fmt.Println(response.Model) //This is the model name
	//Convert to Json format the prompt
    reqBody, err := jsoniter.Marshal(response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{
			ServiceName: "LLM Client",
			Message:     "Internal Server Error",
			Description: "Failed to marshal request body: " + err.Error(),
		})
	}

	// Encrypt the JSON
    encryptedBody, err := utils.Encrypt(string(reqBody))
    if err != nil {
	    return c.Status(fiber.StatusInternalServerError).JSON(models.Error{
		Message:     "Encryption Error",
		Description: err.Error(),
	   })
    }
    
	//Send a POST request to the adapter service
	resp, err := http.Post("http://127.0.0.1:8001/modelRequest", "application/json", bytes.NewBuffer([]byte(encryptedBody)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{
			ServiceName: "LLM Client",
			Message:     "Failed to send request to model server",
			Description: "Failed to send request to model server: " + err.Error(),
		})
	}
	defer resp.Body.Close()

    //REad the response body
	responsefromadapter, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return  c.Status(fiber.StatusInternalServerError).JSON(models.Error{
			ServiceName: "LLM Client",
			Message:     "Internal Server Error",
			Description: "Failed to read response body: " + err.Error(),
		})
	}
    
	fmt.Println("Raw adapter response:", string(responsefromadapter))

	//Take the reposne body and put in reponsefromadapter
	//resp return a string not a json object
	
  return c.JSON(fiber.Map{
	"response": string(responsefromadapter),
})

} 

//This takes the images and converts to base64
//Returns the base64 string or error
func ProcessImage(image *multipart.FileHeader)(string,error){
	//open the file
   file,err:=image.Open()
   if err!=nil{
	  return "",err
   }
   defer file.Close()
  
   //Read the file 
   imageBytes,err:= io.ReadAll(file)
   if err!=nil{
	  return "",err
   }
   
   //Convert to base64
   base64Str:= base64.StdEncoding.EncodeToString(imageBytes)
   
   //Get the mime type example image/jpeg
   mimeType:=image.Header.Get("Content-Type")

   //final string given example data:image/jpeg;base64,/9j/4AAQSkZJRgABAQ...
   dataURI:="data:"+mimeType+";base64,"+base64Str

   return dataURI,nil
}

//Potentially this function can be used to process the file
//func ProcessFile()