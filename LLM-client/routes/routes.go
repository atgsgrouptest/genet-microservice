package routes


import(
	"github.com/gofiber/fiber/v2"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/controllers"
)


func UseRoutes(app *fiber.App) {
	//SendRequest to the controller
  app.Post("/sendRequest",controllers.SendRequest)
  app.Post("/sendRequestImages",controllers.SendRequestWithImages)
}