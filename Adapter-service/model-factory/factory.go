package factory

import (
       "github.com/lokesh2201013/genet-microservice/Adapter-service/models"
       "errors"
	)


//ModelAdapter is an interface that defines the method GenerateResponse
// It takes a models.Request as input and returns a models.Response and models.Error
// It is implemented by the modeltype struct	
type ModelAdapter interface {
	GenerateResponse(request models.Request)(string, models.Error)
}


//GetModelType is a function that takes a model name as input and returns a ModelAdapter
// It checks the model name and returns the corresponding ModelAdapter
// If the model name is not valid, it returns an error
func GetModelType(model string)(ModelAdapter,error){
	switch model{
	/*case "llama3.1:8b":
		return &llama3Adapter{},nil*/
	case "gemma3:4b":
		return &gemma3Adapter{},nil
	case "gpt-4o-mini":
		return &gpt4oMiniAdapter{},nil
	default:
		return nil, errors.New("Invalid model type")
	}
}
