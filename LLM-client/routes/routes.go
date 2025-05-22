package routes


import(
	"github.com/gofiber/fiber/v2"
	"github.com/atgsgrouptest/genet-microservice/LLM-client/controllers"
)


func UseRoutes(app *fiber.App) {
  app.Post("/sendRequest",controllers.SendRequest)
}