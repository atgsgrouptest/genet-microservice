package controllers

import (
	"encoding/base64"
	"io"
	//"mime"
	"io/ioutil"
	"github.com/json-iterator/go"
	"mime/multipart"
    "bytes"
	"net/http"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/models"
	"github.com/gofiber/fiber/v2"
)

func SendRequest(c *fiber.Ctx) error {
	var request models.Request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Message:     "Invalid input",
			Description: "Request body is not valid" + err.Error(),
		})
	}

	if request.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Message:     "Prompt not specified",
			Description: "Prompt is required",
		})
	}


    var response models.Response


	image,err:=c.FormFile("image")

	if err==nil{
    response.Images,err=ProcessImage(image)
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(models.Error{
			Message:     "Invalid image",
			Description: "Image is not valid"+ err.Error(),
		})
     }
    }
    
    response.Prompt=request.Prompt+" give the entire response in json format"  

    reqBody, err := jsoniter.Marshal(response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{
			Message:     "Internal Server Error",
			Description: "Failed to marshal request body: " + err.Error(),
		})
	}

	resp, err := http.Post("<ADD  THE ADDRESS OF THE ADAPTER-SERVICE>", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Error{
			Message:     "Failed to send request to model server",
			Description: "Failed to send request to model server: " + err.Error(),
		})
	}
	defer resp.Body.Close()


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return  c.Status(fiber.StatusInternalServerError).JSON(models.Error{
			Message:     "Internal Server Error",
			Description: "Failed to read response body: " + err.Error(),
		})
	}

	var responsefromadapter string
	if err :=  jsoniter.Unmarshal(body, &responsefromadapter); err != nil {
		return  c.Status(fiber.StatusInternalServerError). JSON(models.Error{
			Message:     "Internal Server Error",
			Description: "Failed to unmarshal response body: " + err.Error(),
		})
	}

return c.JSON(responsefromadapter)
} 

func ProcessImage(image *multipart.FileHeader)(string,error){
   file,err:=image.Open()
   if err!=nil{
	  return "",err
   }
   defer file.Close()

   imageBytes,err:= io.ReadAll(file)
   if err!=nil{
	  return "",err
   }

   base64Str:= base64.StdEncoding.EncodeToString(imageBytes)

   mimeType:=image.Header.Get("Content-Type")
   dataURI:="data:"+mimeType+";base64,"+base64Str

   return dataURI,nil
}

func ProcessFile()