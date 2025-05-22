package models

type Response struct{
	modelName string `json:"model_name"`
	Prompt string `json:"prompt"`
	Images string `json:"images"`
}

type Error struct{
	
	Message string `json:"message"`
	Description string `json:"description"`
}

type Request struct{
	Prompt string `json:"prompt"`
}