package seed

import (
	"context"
	"database/sql" // Import database/sql
	"fmt"
	"log"
	"time"

	"github.com/oblongtable/beanbag-backend/db"
)

// SeedData holds the queries interface and the DB connection needed for seeding.
type SeedData struct {
	dbConn  *sql.DB // Add the database connection
	queries *db.Queries
}

// NewSeedData creates a new SeedData instance.
func NewSeedData(dbConn *sql.DB, queries *db.Queries) *SeedData {
	return &SeedData{dbConn: dbConn, queries: queries}
}

// SeedDatabaseIfNeeded checks if the database needs seeding (e.g., if users table is empty)
// and performs seeding if necessary.
func (s *SeedData) SeedDatabaseIfNeeded(ctx context.Context) error {
	log.Println("Checking if database seeding is required...")

	// Check if any users exist. If yes, assume DB is already seeded or populated.
	userCount, err := s.queries.CountUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to count users for seeding check: %w", err)
	}

	if userCount > 0 {
		log.Printf("Found %d users. Skipping database seeding.\n", userCount)
		return nil // Seeding not needed
	}

	log.Println("No users found. Proceeding with database seeding...")

	// --- Start Seeding ---
	// Use the dbConn to begin the transaction
	tx, err := s.dbConn.BeginTx(ctx, nil) // Use s.dbConn here
	if err != nil {
		return fmt.Errorf("failed to begin transaction for seeding: %w", err)
	}
	// Ensure rollback happens if any error occurs during seeding
	defer tx.Rollback() // Rollback is a no-op if Commit succeeds

	// Use the transaction-aware queries provided by sqlc
	qtx := s.queries.WithTx(tx)
	// Pass the transaction-aware queries to the helper functions
	// No need for a separate seedRunner struct instance here, just use qtx directly
	// Or, if you prefer the runner pattern, ensure it uses qtx:
	// seedRunner := &SeedData{dbConn: s.dbConn, queries: qtx} // Pass dbConn too if needed by helpers

	// Create Users using the transaction queries (qtx)
	users, err := s.createSeedUsersWithQueries(ctx, qtx) // Pass qtx
	if err != nil {
		return fmt.Errorf("failed during user seeding: %w", err)
	}

	// Create Quizzes using the transaction queries (qtx)
	quizzes, err := s.createSeedQuizzesWithQueries(ctx, qtx, users) // Pass qtx
	if err != nil {
		return fmt.Errorf("failed during quiz seeding: %w", err)
	}

	// Create Questions using the transaction queries (qtx)
	questions, err := s.createSeedQuestionsWithQueries(ctx, qtx, quizzes) // Pass qtx
	if err != nil {
		return fmt.Errorf("failed during question seeding: %w", err)
	}

	// Create Answers using the transaction queries (qtx)
	err = s.createSeedAnswersWithQueries(ctx, qtx, questions) // Pass qtx
	if err != nil {
		return fmt.Errorf("failed during answer seeding: %w", err)
	}

	// If all seeding steps succeeded, commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit seeding transaction: %w", err)
	}

	log.Println("Database seeding completed successfully.")
	return nil
}

// --- Helper Creation Functions ---
// Modify helpers to accept *db.Queries directly, as they will receive the transaction-aware 'qtx'

func (s *SeedData) createSeedUsersWithQueries(ctx context.Context, q *db.Queries) ([]db.User, error) {
	log.Println("Seeding users...")
	var users []db.User
	userEmails := []string{"seed.user1@example.com", "seed.user2@example.com"}

	for i, email := range userEmails {
		userParams := db.CreateUserParams{
			Name:      fmt.Sprintf("Seed User %d", i+1),
			Email:     email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// Use the passed-in queries (q), which might be transaction-aware
		user, err := q.CreateUser(ctx, userParams)
		if err != nil {
			return nil, fmt.Errorf("failed to create user %s: %w", email, err)
		}
		log.Printf("Created User: ID=%d, Name=%s, Email=%s\n", user.UserID, user.Name, user.Email)
		users = append(users, user)
	}
	return users, nil
}

func (s *SeedData) createSeedQuizzesWithQueries(ctx context.Context, q *db.Queries, users []db.User) ([]db.Quiz, error) {
	log.Println("Seeding quizzes...")
	var quizzes []db.Quiz
	quizTitles := []string{"Sample Geography Quiz", "Basic Science Quiz"}

	if len(users) == 0 {
		return nil, fmt.Errorf("cannot create quizzes without users")
	}

	for i, title := range quizTitles {
		creator := users[i%len(users)]
		quizParams := db.CreateQuizParams{
			CreatorID:   sql.NullInt32{Int32: creator.UserID, Valid: true},
			QuizTitle:   title,
			Description: sql.NullString{String: fmt.Sprintf("A sample quiz about %s.", title), Valid: true},
			IsPriv:      false,
			Timer:       0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		// Use the passed-in queries (q)
		quiz, err := q.CreateQuiz(ctx, quizParams)
		if err != nil {
			return nil, fmt.Errorf("failed to create quiz '%s': %w", title, err)
		}
		log.Printf("Created Quiz: ID=%d, Title=%s, CreatorID=%d\n", quiz.QuizID, quiz.QuizTitle, creator.UserID)
		quizzes = append(quizzes, quiz)
	}
	return quizzes, nil
}

func (s *SeedData) createSeedQuestionsWithQueries(ctx context.Context, q *db.Queries, quizzes []db.Quiz) ([]db.Question, error) {
	log.Println("Seeding questions...")
	var questions []db.Question
	questionData := map[int][]string{
		0: {"What is the capital of France?", "Which is the largest ocean?"},
		1: {"What is H2O?", "What force pulls objects towards Earth?"},
	}

	if len(quizzes) == 0 {
		return nil, fmt.Errorf("cannot create questions without quizzes")
	}

	qCount := 0
	for quizIndex, quiz := range quizzes {
		if texts, ok := questionData[quizIndex]; ok {
			for _, text := range texts {
				questionParams := db.CreateQuestionParams{
					QuizID:      sql.NullInt32{Int32: quiz.QuizID, Valid: true},
					Description: text,
					TimerOption: false,
					Timer:       0,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				// Use the passed-in queries (q)
				question, err := q.CreateQuestion(ctx, questionParams)
				if err != nil {
					return nil, fmt.Errorf("failed to create question '%s' for quiz %d: %w", text, quiz.QuizID, err)
				}
				log.Printf("Created Question: ID=%d, Text=%s, QuizID=%d\n", question.QuesID, question.Description, quiz.QuizID)
				questions = append(questions, question)
				qCount++
			}
		}
	}
	if qCount == 0 {
		log.Println("Warning: No questions were seeded.")
	}
	return questions, nil
}

func (s *SeedData) createSeedAnswersWithQueries(ctx context.Context, q *db.Queries, questions []db.Question) error {
	log.Println("Seeding answers...")
	answerData := map[string][]struct {
		Text      string
		IsCorrect bool
	}{
		"What is the capital of France?":          {{"Paris", true}, {"London", false}, {"Berlin", false}, {"Madrid", false}},
		"Which is the largest ocean?":             {{"Atlantic", false}, {"Indian", false}, {"Arctic", false}, {"Pacific", true}},
		"What is H2O?":                            {{"Salt", false}, {"Water", true}, {"Sugar", false}, {"Oxygen", false}},
		"What force pulls objects towards Earth?": {{"Magnetism", false}, {"Friction", false}, {"Gravity", true}, {"Tension", false}},
	}

	if len(questions) == 0 {
		return fmt.Errorf("cannot create answers without questions")
	}

	aCount := 0
	for _, question := range questions {
		if answers, ok := answerData[question.Description]; ok {
			for _, ans := range answers {
				answerParams := db.CreateAnswerParams{
					QuesID:      sql.NullInt32{Int32: question.QuesID, Valid: true},
					Description: ans.Text,
					IsCorrect:   ans.IsCorrect,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				// Use the passed-in queries (q)
				_, err := q.CreateAnswer(ctx, answerParams)
				if err != nil {
					log.Printf("Error creating answer '%s' for question %d: %v\n", ans.Text, question.QuesID, err)
				} else {
					aCount++
				}
			}
		} else {
			log.Printf("Warning: No answer data defined for question ID %d ('%s')\n", question.QuesID, question.Description)
		}
	}
	log.Printf("Attempted to seed %d answers.\n", aCount)
	return nil
}
