package main

import (
	"log"
	"os"

 "github.com/lokesh2201013/genet-microservice/Adapter-service/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func main(){
	app:=fiber.New()

	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
      
   routes.UseRoutes(app)
 
	if err:=app.Listen(os.Getenv(":APP_PORT"));err!=nil{
		log.Fatalf("Error starting server: %v", err)
		return
	}
}