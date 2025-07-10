package controllers
import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/atgsgrouptest/genet-microservice/AI-testing/Logger"
	"github.com/atgsgrouptest/genet-microservice/AI-testing/database"
	"github.com/atgsgrouptest/genet-microservice/AI-testing/models"
)

func ReturnResult(c *fiber.Ctx) error {
	requestId := c.Query("requestId")
	if requestId == "" {
		logger.Log.Error("Missing requestId", zap.String("endpoint", "ReturnResult"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "RequestId not specified",
			"details": "requestId not found in query",
		})
	}

	// Convert to ObjectID for the main requests collection
	objID, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		logger.Log.Error("Invalid ObjectID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid requestId format",
			"details": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the request document
	requestCollection := database.MongoDB.Collection("requests")
	var request models.Request

	err = requestCollection.FindOne(ctx, bson.M{"requestId": objID}).Decode(&request)
	if err != nil {
		logger.Log.Error("Completed request not found", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Request not completed",
			"details": err.Error(),
		})
	}

	// Use string version of ObjectID for matching in responses collections
	requestIdStr := objID.Hex()

	// Fetch positive responses
	var posResponses []models.PositiveResponse
	posCollection := database.MongoDB.Collection("positive_responses")

	posCursor, err := posCollection.Find(ctx, bson.M{"requestId": requestIdStr})
	if err != nil {
		logger.Log.Error("Failed to fetch positive responses", zap.Error(err))
	} else {
		if err := posCursor.All(ctx, &posResponses); err != nil {
			logger.Log.Error("Error decoding positive responses", zap.Error(err))
		}
		posCursor.Close(ctx)
	}

	// Fetch negative responses
	var negResponses []models.NegativeResponse
	negCollection := database.MongoDB.Collection("negative_responses")

	negCursor, err := negCollection.Find(ctx, bson.M{"requestId": requestIdStr})
	if err != nil {
		logger.Log.Error("Failed to fetch negative responses", zap.Error(err))
	} else {
		if err := negCursor.All(ctx, &negResponses); err != nil {
			logger.Log.Error("Error decoding negative responses", zap.Error(err))
		}
		negCursor.Close(ctx)
	}

	fmt.Println("Positive Responses:", posResponses)
	fmt.Println("Negative Responses:", negResponses)

	completed := models.CompletedRequest{
		RequestId:         request.RequestID,
		CompanyID:         request.CompanyID,
		RequestMaterial: request.RequestMaterial,
		PositiveResponses: posResponses,
		NegativeResponses: negResponses,
	}

	return c.Status(fiber.StatusOK).JSON(completed)
} 