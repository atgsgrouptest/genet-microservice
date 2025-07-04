package database

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongo() {
    mongoURI := os.Getenv("MONGO_URI")
    dbName := os.Getenv("MONGO_DB")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatalf("MongoDB Connection Failed: %v", err)
    }

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("MongoDB Ping Failed: %v", err)
    }

    fmt.Println("✅ Connected to MongoDB")

    MongoClient = client
    MongoDB = client.Database(dbName)
}
