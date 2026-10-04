package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/database"
	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	// This is the main function

	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMovies!")
	})

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	config := cors.Config{}
	if allowedOrigins == "*" {
		config.AllowOriginFunc = func(origin string) bool {
			return true
		}
		log.Println("Allowed Origin: * (dynamic reflection)")
	} else {
		var origins []string
		if allowedOrigins != "" {
			origins = strings.Split(allowedOrigins, ",")
			for i := range origins {
				origins[i] = strings.TrimSpace(origins[i])
				log.Println("Allowed Origin:", origins[i])
			}
		} else {
			origins = []string{"http://localhost:5173"}
			log.Println("Allowed Origin: http://localhost:5173")
		}
		config.AllowOrigins = origins
	}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	//config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Logger())

	var client *mongo.Client
	var pingErr error
	for attempt := 1; attempt <= 10; attempt++ {
		client = database.Connect()
		if client != nil {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
			pingErr = client.Ping(pingCtx, nil)
			pingCancel()
			if pingErr == nil {
				log.Println("Connected to MongoDB successfully!")
				break
			}
		}
		log.Printf("MongoDB connection attempt %d/10 failed (%v), retrying in 2 seconds...\n", attempt, pingErr)
		time.Sleep(2 * time.Second)
	}

	if pingErr != nil || client == nil {
		log.Fatalf("Failed to reach MongoDB after 10 retries: %v", pingErr)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}

	}()

	routes.SetupUnProtectedRoutes(router, client)
	routes.SetupProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}
