package routes

import (
	//"github.com/gofiber/fiber/middleware"
	  "github.com/gofiber/fiber/v2"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/controllers"
)

func UseRoutes(app *fiber.App) {
	//This is the endpoint to send the requst to the model
	app.Post("/modelRequest",controllers.ModelRequest)
}