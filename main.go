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
)

var (
	server    *gin.Engine
	DBQueries *db.Queries
	db_conn   *sql.DB
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	// Connect to the database
	dbConn, err := initializers.NewDBConnection(&config)
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
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

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

	// Quiz routes
	router.POST("/quizzes", quizHandler.CreateQuiz)
	router.GET("/quizzes/:id", quizHandler.GetQuiz)

	// User routes
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.GetUser)

	// Question routes
	router.POST("/questions", questionHandler.CreateQuestion)
	router.GET("/questions/:id", questionHandler.GetQuestion)

	// Answer routes
	router.POST("/answers", answerHandler.CreateAnswer)
	router.GET("/answers/:id", answerHandler.GetAnswer)

	log.Fatal(server.Run(":" + config.ServerPort))
}
