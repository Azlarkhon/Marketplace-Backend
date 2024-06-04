package main

import (
	"log"
	"os"

	"example.com/m/database"
	"example.com/m/middleware"
	"example.com/m/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.DBinstance()
	defer database.CloseDB()

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(middleware.CorsMiddleware())

	routes.UserRoutes(router)
	routes.ProductRoutes(router)
	routes.OrderRoutes(router)
	routes.CartRoutes(router)

	router.Run(":" + port)
}
