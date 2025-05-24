package models


//This model is the Request which we will give to the adapter-service
//what to do with the model name we have to see
//Images are in base64encoded format
type Response struct{
	Model string `json:"model"`
	Prompt string `json:"prompt"`
	Images string `json:"images"`
}


//This is the error reponse if prompt is not valid or Reponse has a problem

type Error struct{
	ServiceName string `json:"service_name"`
	Message string `json:"message"`
	Description string `json:"description"`
}

//this is the requst we will take from the user
//The images will be coming in form of Multipart form-data whihch we will convert to base64
//The prompt will be in the form of string
type Request struct{
	Prompt string `json:"prompt"`
}