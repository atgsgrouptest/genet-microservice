package main

import (
	//"fmt"
	"log"
	"os"
    "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/lokesh2201013/genet-microservice/Adapter-service/routes"
)

func main() {
	// Load variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
   
	 app.Use(logger.New()) 

	routes.UseRoutes(app)

	port := os.Getenv("APP_PORT")
	/*if port == "" {
		fmt.Println("APP_PORT not set, using default port 8001")
		port = "8001"
	}*/

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
