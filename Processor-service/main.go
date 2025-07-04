package main


import (
	//"log"
	"os"
	"go.uber.org/zap"
	"github.com/atgsgrouptest/genet-microservice/Processor-service/Logger"
	//"github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/joho/godotenv"
	"github.com/atgsgrouptest/genet-microservice/Processor-service/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	 
)


func main() {
	// Load variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("No .env file found, using default values")
	}
    
	//Increase the body limit to 20MB
	//Default is 4MB
    app := fiber.New(fiber.Config{
    BodyLimit: 20 * 1024 * 1024, // 20MB
})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
    
	//app.Use(logger.New())
	logger.InitLogger() // Initialize the logger
	// Register routes
	app.Use(logger.ZapLogger())
 
	routes.UseRoutes(app)

	port := os.Getenv("APP_PORT")
	/*if port == "" {
		fmt.Println("APP_PORT not set, using default port 8001")
		port = "8001"
	}*/

	if err := app.Listen(":" + port); err != nil {
		logger.Log.Fatal("Error starting server", zap.Error(err))
	}
}
