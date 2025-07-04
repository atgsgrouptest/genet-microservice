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
	//"github.com/atgsgrouptest/genet-microservice/Processor-service/httpclient"
	"github.com/atgsgrouptest/genet-microservice/Processor-service/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	)

func TestAPI(c *fiber.Ctx) error{
    var SequencedResponse models.APIWrapper
	if err:= c.BodyParser(&SequencedResponse);err!=nil{
		return err
	}
	
	PositivePrompt,_:=prompt.PositiveCasePrompt(SequencedResponse)
	PositiveReq,_:= json.Marshal(map[string]string{
			"prompt": PositivePrompt,
		})
	PositiveTestCase,_:=http.Post("http://127.0.0.1:8002/sendRequest", "application/json", bytes.NewBuffer(PositiveReq))
    PositiveRespBody, err := io.ReadAll(PositiveTestCase.Body)
if err != nil {
	logger.Log.Error("Error reading LLM positive test response body", zap.Error(err))
	return c.Status(fiber.StatusInternalServerError).JSON(Error.ReturnError("Processor-service", err, "Error reading response body"))
}
fmt.Println(string(PositiveRespBody))
 var reponseFromLLM models.LLMResponse
if err:=json.Unmarshal(PositiveRespBody, &reponseFromLLM);err!=nil{
	return err
}

	return c.Status(fiber.StatusOK).JSON(reponseFromLLM)
}