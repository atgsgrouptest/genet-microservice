package routes

import(
	"github.com/gofiber/fiber/v2"
	"github.com/atgsgrouptest/genet-microservice/AI-testing/controllers"
)

func UseRoutes(app *fiber.App) {
app.Post("/makePrompt", controllers.MakePrompt)

}