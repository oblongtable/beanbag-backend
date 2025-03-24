/**

This is intended to be ran while the `make compose-up-build` command is running and the databse exists,
so that you can go ahead and see an example of creation and deletion of the different table record/model types
and see the associated FK resources getting deleted too, and referencing correctly.

Note: this requires the 5432 port to be forwarded from the DB docker container to your local machine as it uses
localhost to connect to the db. VSCode will do this for you automatically.

Run this with `go run examples/create_user_quiz_question_answer.go`

**/

package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/oblongtable/beanbag-backend/db"
	"github.com/oblongtable/beanbag-backend/initializers"
)

var (
	DBQueries *db.Queries
)

func init() {
	err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	config := initializers.GetConfig()
	config.DBHost = "localhost"

	dbConn, err := initializers.NewDBConnection(config)
	if err != nil {
		log.Fatal("? Could not connect to the database", err)
	}

	DBQueries = db.New(dbConn)
}

func waitForEnter() {
	fmt.Println("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func main() {
	ctx := context.Background()

	// Create Users
	fmt.Println("\nPress Enter to create users...")
	waitForEnter()
	users := createUsers(ctx)

	// Create Quizzes
	fmt.Println("\nPress Enter to create quizzes...")
	waitForEnter()
	quizzes := createQuizzes(ctx, users)

	// Create Questions
	fmt.Println("\nPress Enter to create questions...")
	waitForEnter()
	questions := createQuestions(ctx, quizzes)

	// Create Answers
	fmt.Println("\nPress Enter to create answers...")
	waitForEnter()
	createAnswers(ctx, questions)

	// Delete All
	fmt.Println("\nPress Enter to delete all data...")
	waitForEnter()
	deleteData(ctx, users, quizzes, questions)
}

func createUsers(ctx context.Context) []db.User {
	var users []db.User
	for i := 0; i < 3; i++ {
		userParams := db.CreateUserParams{
			Name:      fmt.Sprintf("Test User %d", i+1),
			Email:     fmt.Sprintf("testuser%d@example.com", rand.Intn(1000000)), // Unique email
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		user, err := DBQueries.CreateUser(ctx, userParams)
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Println("Created User:", user)
		users = append(users, user)
	}
	return users
}

func createQuizzes(ctx context.Context, users []db.User) []db.Quiz {
	var quizzes []db.Quiz
	for i := 0; i < 3; i++ {
		quizParams := db.CreateQuizParams{
			CreatorID:   sql.NullInt32{Int32: users[i%len(users)].UserID, Valid: true}, // Corrected line
			QuizTitle:   fmt.Sprintf("Test Quiz %d", i+1),
			Description: sql.NullString{String: fmt.Sprintf("This is test quiz %d", i+1), Valid: true},
			IsPriv:      false,
			Timer:       60,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		quiz, err := DBQueries.CreateQuiz(ctx, quizParams)
		if err != nil {
			log.Fatalf("Failed to create quiz: %v", err)
		}
		fmt.Println("Created Quiz:", quiz)
		quizzes = append(quizzes, quiz)
	}
	return quizzes
}

func createQuestions(ctx context.Context, quizzes []db.Quiz) []db.Question {
	var questions []db.Question
	for i := 0; i < 5; i++ {
		questionParams := db.CreateQuestionParams{
			QuizID:      sql.NullInt32{Int32: quizzes[i%len(quizzes)].QuizID, Valid: true}, // Corrected line
			Description: fmt.Sprintf("Question %d for quiz %d?", i+1, (i%len(quizzes))+1),
			TimerOption: true,
			Timer:       45,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		question, err := DBQueries.CreateQuestion(ctx, questionParams)
		if err != nil {
			log.Fatalf("Failed to create question: %v", err)
		}
		fmt.Println("Created Question:", question)
		questions = append(questions, question)
	}
	return questions
}

func createAnswers(ctx context.Context, questions []db.Question) {
	for i := 0; i < 10; i++ {
		answerParams := db.CreateAnswerParams{
			QuesID:      sql.NullInt32{Int32: questions[i%len(questions)].QuesID, Valid: true}, // Corrected line
			Description: fmt.Sprintf("Answer %d for question %d", i+1, (i%len(questions))+1),
			IsCorrect:   i%2 == 0, // Alternate correct/incorrect
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		answer, err := DBQueries.CreateAnswer(ctx, answerParams)
		if err != nil {
			log.Fatalf("Failed to create answer: %v", err)
		}
		fmt.Println("Created Answer:", answer)
	}
}

func deleteData(ctx context.Context, users []db.User, quizzes []db.Quiz, questions []db.Question) {
	// Delete Users, this should delete all data due to cascading deletes
	for _, user := range users {
		fmt.Printf("Deleting User:%s with id %d\n", user.Name, user.UserID)
		err := DBQueries.DeleteUser(ctx, user.UserID)
		if err != nil {
			log.Fatalf("Failed to delete user: %v", err)
		}
	}

	fmt.Println("Please check that all tables have the users above removed, and their quizzes, questions, answers should all be gone due to cascading deletes.")

}
