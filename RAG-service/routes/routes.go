package routes

import(
	"github.com/gofiber/fiber/v2"
	"github.com/atgsgrouptest/genet-microservice/RAG-service/controllers"
)

func UseRoutes(app *fiber.App) {
 app.Post("/sendFiles", controllers.SendFiles)
 app.Get("/getPromptWithContext", controllers.GetPromptWithContext)
}