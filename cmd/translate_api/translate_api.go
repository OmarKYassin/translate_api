package main

import (
	"fmt"
	"os"
	"time"

	"github.com/OmarKYassin/translate_api/pkg/logging"
	"github.com/OmarKYassin/translate_api/pkg/translate_api/routes"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize Zap logger
	logger := logging.Logger()
	logging.SyncLogger()

	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found, loading environment variables from system.")
	}

	// Initialize a Gin router
	router := gin.New()

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	routes.RegisterRoutes(router)

	// Start the server on PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(fmt.Sprintf(":%s", port))
}
