package main

import (
    "github.com/atgsgrouptest/genet-microservice/AI-testing/Logger"
	"github.com/atgsgrouptest/genet-microservice/AI-testing/routes"
	"os"
	"go.uber.org/zap"
     "fmt"
	"context"
	"github.com/joho/godotenv"
     "github.com/atgsgrouptest/genet-microservice/AI-testing/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func main() {
	// Load variables from .env file
	logger.InitLogger() // Initialize the logger
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
	
	// Register routes
	app.Use(logger.ZapLogger())
    database.ConnectMongo()
	routes.UseRoutes(app)


	port := os.Getenv("APP_PORT")
	if port == "" {
		fmt.Println("APP_PORT not set, using default port 8004")
		port = "8004"
	}

	if err := app.Listen(":" + port); err != nil {
		logger.Log.Fatal("Error starting server", zap.Error(err))
	}
	defer func() {
    if err := database.MongoClient.Disconnect(context.TODO()); err != nil {
        logger.Log.Error("Error disconnecting MongoDB", zap.Error(err))
    } else {
        logger.Log.Info("MongoDB disconnected successfully")
    }
}()

}
