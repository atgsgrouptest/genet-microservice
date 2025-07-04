package routes


import(
	"github.com/gofiber/fiber/v2"
	"github.com/atgsgrouptest/genet-microservice/Processor-service/controllers"
)


func UseRoutes(app *fiber.App) {
	app.Post("/getAPIsequence", controllers.GetAPISequence)
	app.Post("/testAPI", controllers.TestAPI)
}