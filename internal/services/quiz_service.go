package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/oblongtable/beanbag-backend/db"
	"github.com/oblongtable/beanbag-backend/internal/apimodels"
)

type QuizService struct {
	connPool *sql.DB
	queries  *db.Queries
}

func NewQuizService(connPool *sql.DB, queries *db.Queries) *QuizService {
	return &QuizService{
		connPool: connPool,
		queries:  queries,
	}
}

// GetFullQuiz fetches a quiz, its questions, and their answers, returning the combined structure.
func (s *QuizService) GetFullQuiz(ctx context.Context, quizID int32) (*apimodels.QuizApiModel, error) {
	// 1. Get Quiz
	quiz, err := s.queries.GetQuiz(ctx, quizID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("quiz with ID %d not found", quizID) // Specific not found error
		}
		return nil, fmt.Errorf("failed to get quiz %d: %w", quizID, err)
	}

	// 2. Get Questions
	dbQuestions, err := s.queries.ListQuestionsByQuiz(ctx, sql.NullInt32{Int32: quizID, Valid: true})
	if err != nil {
		// If no questions is okay, check for sql.ErrNoRows and continue, otherwise return error
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to list questions for quiz %d: %w", quizID, err)
		}
		// If ErrNoRows is acceptable, dbQuestions will be an empty slice
		log.Printf("No questions found for quiz %d", quizID)
	}

	// Prepare map to hold answers grouped by question ID
	answersMap := make(map[int32][]apimodels.AnswerApiModel)
	questionIDs := make([]int32, 0, len(dbQuestions))

	if len(dbQuestions) > 0 {
		// Extract question IDs
		for _, q := range dbQuestions {
			questionIDs = append(questionIDs, q.QuesID)
		}

		// 3. Get Answers for all questions in one go
		dbAnswers, err := s.queries.ListAnswersByQuestionIDs(ctx, questionIDs)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("failed to list answers for questions of quiz %d: %w", quizID, err)
			}
			// If ErrNoRows is acceptable, dbAnswers will be an empty slice
			log.Printf("No answers found for questions of quiz %d", quizID)
		}

		// 4. Group Answers by Question ID
		for _, a := range dbAnswers {
			if a.QuesID.Valid { // Check if QuesID is not NULL
				quesID := a.QuesID.Int32
				answersMap[quesID] = append(answersMap[quesID], apimodels.AnswerApiModel{
					Text:      a.Description,
					IsCorrect: a.IsCorrect,
				})
			}
		}
	}

	// 5. Construct the final response object
	apiQuestions := make([]apimodels.QuestionApiModel, 0, len(dbQuestions))
	for _, q := range dbQuestions {
		quesID := q.QuesID
		apiAnswers := answersMap[quesID] // Get answers from map (will be nil if no answers)
		if apiAnswers == nil {
			apiAnswers = []apimodels.AnswerApiModel{} // Ensure it's an empty slice, not nil, for JSON
		}

		apiQuestions = append(apiQuestions, apimodels.QuestionApiModel{
			Text:       q.Description,
			UseTimer:   q.TimerOption, 
			TimerValue: q.Timer,       
			Answers:    apiAnswers,
		})
	}

	creatorID := quiz.CreatorID.Int32

	fullQuiz := &apimodels.QuizApiModel{
		Title:     quiz.QuizTitle,
		CreatorID: creatorID,
		Questions: apiQuestions,
	}

	return fullQuiz, nil
}

func (s *QuizService) CreateQuiz(ctx context.Context, title string, creatorID int32) (*db.Quiz, error) {
	params := db.CreateQuizParams{
		QuizTitle: title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatorID: sql.NullInt32{Int32: creatorID, Valid: true},
		IsPriv:    false,
		Timer:     0,
	}

	quiz, err := s.queries.CreateQuiz(ctx, params)
	if err != nil {
		return nil, err
	}

	return &quiz, nil
}

// CreateQuizMinimal handles the full creation process within a transaction
// Returns a json extracted from the DB with the full quiz returned
func (s *QuizService) CreateQuizMinimal(ctx context.Context, input apimodels.QuizApiModel) (apimodels.QuizApiModel, error) {
	// 1. Start Transaction
	tx, err := s.connPool.BeginTx(ctx, nil) // Use default transaction options
	if err != nil {
		return apimodels.QuizApiModel{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Ensure rollback happens if anything goes wrong before commit
	defer tx.Rollback() // Safe to call even if committed, it becomes a no-op

	// 2. Get SQLC Queries bound to the transaction
	qtx := s.queries.WithTx(tx)

	// 3. Create the Quiz entry

	createdQuiz, err := qtx.CreateQuizMinimal(ctx, db.CreateQuizMinimalParams{
		QuizTitle: input.Title,
		CreatorID: sql.NullInt32{Int32: input.CreatorID, Valid: true},
	})
	if err != nil {
		// No need to rollback here, defer tx.Rollback() handles it
		return apimodels.QuizApiModel{}, fmt.Errorf("failed to create quiz entry: %w", err)
	}

	// 4. Loop and Create Questions and Answers
	for _, createQuestionReq := range input.Questions {
		// Create Question
		// question is of type apimodels.QuestionApiModel
		createdQuestion, err := qtx.CreateQuestionMinimal(ctx, db.CreateQuestionMinimalParams{
			QuizID:      sql.NullInt32{Int32: createdQuiz.QuizID, Valid: true}, // Use the ID returned from CreateQuiz
			Description: createQuestionReq.Text,
			TimerOption: createQuestionReq.UseTimer,
			Timer:       createQuestionReq.TimerValue,
		})
		if err != nil {
			return apimodels.QuizApiModel{}, fmt.Errorf("failed to create question '%s': %w", createQuestionReq.Text, err)
		}

		for _, createAnswerReq := range createQuestionReq.Answers {
			// answer is of type apimodels.AnswerApiModel
			_, err := qtx.CreateAnswerMinimal(ctx, db.CreateAnswerMinimalParams{
				QuesID:      sql.NullInt32{Int32: createdQuestion.QuesID, Valid: true},
				Description: createAnswerReq.Text,
				IsCorrect:   createAnswerReq.IsCorrect,
			})
			if err != nil {
				return apimodels.QuizApiModel{}, fmt.Errorf("failed to create answer '%s' for question '%s': %w", createAnswerReq.Text, createQuestionReq.Text, err)
			}
		}
	}

	// 5. Commit Transaction if all steps succeeded
	if err := tx.Commit(); err != nil {
		return apimodels.QuizApiModel{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 6. construct a response object to return as json
	fullQuizReturn, err := s.GetFullQuiz(ctx, createdQuiz.QuizID)
	if err != nil {
		// This shouldn't ideally happen if commit succeeded, but handle defensively
		return apimodels.QuizApiModel{}, fmt.Errorf("failed to retrieve created quiz after commit: %w", err)
	}

	return *fullQuizReturn, nil // Return the created quiz data
}

// GetQuiz might also need modification if you want it to return questions/answers
func (s *QuizService) GetQuiz(ctx context.Context, id int32) (*db.Quiz, error) {
	// Current implementation likely only gets the quiz row.
	// You would need additional SQLC queries (e.g., GetQuestionsByQuizID, GetAnswersByQuestionID)
	// and orchestrate calling them here if you want to return the full nested structure.
	quiz, err := s.queries.GetQuiz(ctx, id) // Assuming GetQuiz is generated by SQLC
	if err != nil {
		// Handle errors like sql.ErrNoRows specifically if desired
		return nil, fmt.Errorf("failed to get quiz: %w", err)
	}
	return &quiz, nil
}
