package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/initializers"
)

var (
	server *gin.Engine
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	router := server.Group("/")
	router.GET("/", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	router.GET("/db_health", func(ctx *gin.Context) {
		var tm time.Time

		// Use GORM's Raw method to execute the raw SQL query
		result := initializers.DB.Raw("SELECT NOW()").Scan(&tm)
		if result.Error != nil {
			log.Printf("Query failed: %v\n", result.Error)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database time"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"api": "golang",
			"now_from_postgres": tm,
		})
	})

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	log.Fatal(server.Run(":" + config.ServerPort))
}
