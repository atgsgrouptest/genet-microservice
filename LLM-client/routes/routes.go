package routes


import(
	"github.com/gofiber/fiber/v2"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/controllers"
)


func UseRoutes(app *fiber.App) {
	//SendRequest to the controller
  app.Post("/sendRequest",controllers.SendRequest)
	//SendRequestWithImages to the controller it includes images
	//Use multipart/form-data to send images with key "images"->file and key "prompt"->text
  app.Post("/sendRequestImages",controllers.SendRequestWithImages)
}