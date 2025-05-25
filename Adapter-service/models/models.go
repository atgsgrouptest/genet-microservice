package models


//Request given to ollama model images must be in base 64 encoding
//Model string must contain available and valid model name
type Request struct{
    Model  string   `json:"model"`
    Prompt string   `json:"prompt"`
   Images []string `json:"images"`
	Stream bool	    `json:"stream"`
}
//This is the error reponse if prompt is not valid or Reponse has a problem
type Error struct{
	ServiceName string `json:"service_name"`
	Message string `json:"error"`
	Description string `json:"description"`
}