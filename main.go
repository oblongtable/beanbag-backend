package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/db"
	"github.com/oblongtable/beanbag-backend/initializers"
	"github.com/oblongtable/beanbag-backend/internal/handlers"
	"github.com/oblongtable/beanbag-backend/internal/services"
	"github.com/oblongtable/beanbag-backend/middleware"
	"github.com/oblongtable/beanbag-backend/websocket"

	adaptor "github.com/gwatts/gin-adapter"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/oblongtable/beanbag-backend/docs" // docs is generated by Swag CLI, you have to import it.
)

// @title Beanbag Backend API
// @version 1.0
// @description This is the API for the Beanbag Backend quiz application.

// @host localhost:8080
// @BasePath /
// @schemes http
// @query.collection.format multi
var (
	server    *gin.Engine
	DBQueries *db.Queries
	db_conn   *sql.DB
)

func init() {
	err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	// Connect to the database
	dbConn, err := initializers.NewDBConnection(initializers.GetConfig())
	if err != nil {
		log.Fatal("? Could not connect to the database", err)
	}

	// Run migrations
	initializers.MigrateDB(dbConn)

	// Create a new Queries instance
	queries := db.New(dbConn)

	// Make the connection object available globally
	db_conn = dbConn

	// Make the queries available globally
	DBQueries = queries

	server = gin.Default()
}

func main() {
	config := initializers.GetConfig()
	wssvr := websocket.NewWebSockServer()

	// Initialize services
	quizService := services.NewQuizService(DBQueries)
	userService := services.NewUserService(DBQueries)
	questionService := services.NewQuestionService(DBQueries)
	answerService := services.NewAnswerService(DBQueries)

	// Initialize handlers
	quizHandler := handlers.NewQuizHandler(quizService)
	userHandler := handlers.NewUserHandler(userService)
	questionHandler := handlers.NewQuestionHandler(questionService)
	answerHandler := handlers.NewAnswerHandler(answerService)

	router := server.Group("/")
	router.GET("/", func(ctx *gin.Context) {
		message := "Welcome to Golang api with SQLC, Goose and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	router.GET("/db_health", func(ctx *gin.Context) {
		var tm time.Time

		// Use DBQueries.db.QueryRowContext() to execute the raw SQL query
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := db_conn.QueryRowContext(ctxTimeout, "SELECT NOW()").Scan(&tm)
		if err != nil {
			log.Printf("Query failed: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database time"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"api":               "golang",
			"now_from_postgres": tm,
		})
	})

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	// API routes
	api := router.Group("/api")
	{
		api.Use(adaptor.Wrap(middleware.VerifyToken()))

		// User routes
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)

		// Quiz routes
		api.POST("/quizzes", quizHandler.CreateQuiz)
		api.GET("/quizzes/:id", quizHandler.GetQuiz)

		// Question routes
		api.POST("/questions", questionHandler.CreateQuestion)
		api.GET("/questions/:id", questionHandler.GetQuestion)

		// Answer routes
		api.POST("/answers", answerHandler.CreateAnswer)
		api.GET("/answers/:id", answerHandler.GetAnswer)

		api.GET("/ws", wssvr.ServeWs)
	}

	// Swagger route
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	log.Fatal(server.Run(":" + config.ServerPort))
}
