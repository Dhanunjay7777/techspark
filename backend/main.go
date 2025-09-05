package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
    "github.com/redis/go-redis/v9"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    client *mongo.Client
    ctx    = context.Background()
)

func main() {
    // Load .env
    _ = godotenv.Load()

    mongoURI := os.Getenv("MONGO_URI")
    redisURL := os.Getenv("REDIS_URL")
    port := os.Getenv("PORT")
	frontendURL := os.Getenv("FRONTEND_URL")


    if port == "" {
        port = "4000"
    }

    // MongoDB - store client for later use
    var err error
    client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("MongoDB error:", err)
    }
    defer client.Disconnect(ctx)
    
    // Test MongoDB connection
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("MongoDB ping error:", err)
    }
    fmt.Println("✅ Connected to MongoDB")

    // Redis (Fixed for Upstash TLS)
    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        log.Fatal("Redis URL parse error:", err)
    }
    rdb := redis.NewClient(opt)
    _, err = rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Redis error:", err)
    }
    fmt.Println("✅ Connected to Redis")

    // Fiber App
    app := fiber.New()
    
    app.Use(cors.New(cors.Config{
        AllowOrigins: frontendURL,
        AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
    }))

    // Routes
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "OK",
            "message": "Backend is running!",
        })
    })

    app.Get("/api/test", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Hello from Go backend!",
            "mongo": "Connected",
            "redis": "Connected",
        })
    })

    fmt.Println("🚀 Backend running on port", port)
    log.Fatal(app.Listen(":" + port))
}
