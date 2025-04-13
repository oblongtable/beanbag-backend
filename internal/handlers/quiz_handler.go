package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	// Import errors if needed for specific error checking, though not strictly required for the fix
	// "errors"

	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/internal/apimodels"
	"github.com/oblongtable/beanbag-backend/internal/services"
	"github.com/oblongtable/beanbag-backend/middleware" // Import middleware if needed for security checks
	// Import db package only if needed for swagger docs, prefer apimodels
)

type QuizHandler struct {
	quizService *services.QuizService
	// userService *services.UserService // Inject if needed for creator ID validation
}

// Inject userService if needed for creator validation
func NewQuizHandler(quizService *services.QuizService /*, userService *services.UserService */) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
		// userService: userService,
	}
}

// CreateQuiz godoc
// @Summary Create a new basic quiz entry (DEPRECATED? Use POST /quizzes for full creation)
// @Description Create only the quiz entry without questions/answers. Consider using POST /quizzes instead.
// @Tags quizzes
// @Accept json
// @Produce json
// @Param quiz body apimodels.QuizApiModel true "Basic Quiz details (Title, CreatorID)" // <-- FIX: Use qualified name apimodels.QuizApiModel
// @Success 201 {object} db.Quiz "The created basic quiz object"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /quizzes/basic [post] // Suggest different path if keeping basic creation
// @Security BearerAuth
func (h *QuizHandler) CreateQuiz(ctx *gin.Context) {
	// This handler might be redundant if CreateQuizMinimal handles full creation at POST /quizzes
	var req apimodels.QuizApiModel // Request body uses the API model

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// SECURITY: Validate req.CreatorID against authenticated user from JWT
	// jwtSub := ctx.GetString(middleware.GinContextKeyUserSub)
	// loggedInUserID, err := h.userService.GetUserIDBySub(ctx, jwtSub) // Requires userService injection
	// if err != nil || loggedInUserID != req.CreatorID {
	// 	ctx.JSON(http.StatusForbidden, gin.H{"error": "Creator ID mismatch or lookup failed"})
	// 	return
	// }
	log.Printf("Warning: CreatorID %d in CreateQuiz request not verified against JWT user.", req.CreatorID)

	// Call the service that creates only the basic quiz row
	quiz, err := h.quizService.CreateQuiz(ctx.Request.Context(), req.Title, req.CreatorID)
	if err != nil {
		log.Printf("Error calling CreateQuiz service: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
		return
	}

	ctx.JSON(http.StatusCreated, quiz) // Returns the db.Quiz object
}

// GetQuiz godoc
// @Summary Get basic quiz details by ID (DEPRECATED? Use GET /quizzes/{id}/full)
// @Description Get only the quiz details without questions/answers. Consider using GET /quizzes/{id}/full instead.
// @Tags quizzes
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} db.Quiz "Basic quiz object"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Quiz not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /quizzes/{id}/basic [get] // Suggest different path if keeping basic get
// @Security BearerAuth
func (h *QuizHandler) GetQuiz(ctx *gin.Context) {
	// This handler might be redundant if GetFullQuiz provides the main GET functionality
	quizIDStr := ctx.Param("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	quiz, err := h.quizService.GetQuiz(ctx.Request.Context(), int32(quizID))
	if err != nil {
		// Check for specific "not found" error from service
		if err.Error() == fmt.Sprintf("quiz with ID %d not found", quizID) { // Check specific error message or use custom error type
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Error calling GetQuiz service for ID %d: %v", quizID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quiz"})
		return
	}

	ctx.JSON(http.StatusOK, quiz) // Returns the db.Quiz object
}

// CreateQuizMinimal godoc
// @Summary Create a full quiz with questions and answers
// @Description Create a quiz, its questions, and their answers in a single transaction.
// @Tags quizzes
// @Accept json
// @Produce json
// @Param quiz body apimodels.QuizApiModel true "Full quiz details including questions and answers" // <-- FIX: Use qualified name
// @Success 201 {object} apimodels.QuizApiModel "The fully created quiz structure" // <-- FIX: Use qualified name
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden (Creator ID mismatch)"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /quizzes [post] // Assumes this is the main POST endpoint now
// @Security BearerAuth
func (h *QuizHandler) CreateQuizMinimal(ctx *gin.Context) {
	var req apimodels.QuizApiModel

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// SECURITY: Verify CreatorID matches authenticated user
	jwtEmail := ctx.GetString(middleware.GinContextKeyUserEmail) // Assuming middleware sets this
	// You might need userService injected here to map email/sub to your internal user ID
	// loggedInUserID, err := h.userService.GetUserIDByEmail(ctx, jwtEmail)
	// if err != nil {
	//     log.Printf("Error getting user ID for email %s: %v", jwtEmail, err)
	//     ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user identity"})
	//     return
	// }
	// if loggedInUserID != req.CreatorID {
	// 	ctx.JSON(http.StatusForbidden, gin.H{"error": "Creator ID does not match authenticated user"})
	// 	return
	// }
	log.Printf("Warning: CreatorID %d in CreateQuizMinimal request not verified against JWT user %s.", req.CreatorID, jwtEmail)

	// Validate input further? (e.g., must have questions, questions must have answers?)
	if len(req.Questions) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Quiz must contain at least one question"})
		return
	}
	for _, q := range req.Questions {
		if len(q.Answers) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Question '%s' must contain at least one answer", q.Text)})
			return
		}
		// Check for at least one correct answer?
		hasCorrect := false
		for _, a := range q.Answers {
			if a.IsCorrect {
				hasCorrect = true
				break
			}
		}
		if !hasCorrect {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Question '%s' must have at least one correct answer", q.Text)})
			return
		}
	}

	// Call the service method that handles the transaction
	// Ensure service method CreateQuizMinimal returns (apimodels.QuizApiModel, error) or (*apimodels.QuizApiModel, error)
	// Based on your service code, it returns (apimodels.QuizApiModel, error)
	createdQuiz, err := h.quizService.CreateQuizMinimal(ctx.Request.Context(), req)
	if err != nil {
		log.Printf("Error calling CreateQuizMinimal service: %v", err)
		// Check for specific error types if needed
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz"})
		return
	}

	ctx.JSON(http.StatusCreated, createdQuiz) // Return the full structure
}

// GetFullQuiz godoc
// @Summary Get full quiz details by ID
// @Description Get quiz details including all questions and their answers
// @Tags quizzes
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} apimodels.QuizApiModel "Full quiz structure" // <-- FIX: Use qualified name apimodels.QuizApiModel
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Quiz not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /quizzes/{id}/full [get] // Use a specific path for the full structure
// @Security BearerAuth
func (h *QuizHandler) GetFullQuiz(ctx *gin.Context) {
	quizIDStr := ctx.Param("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// Call the service method to get the full structure
	// Ensure service method GetFullQuiz returns (*apimodels.QuizApiModel, error)
	fullQuiz, err := h.quizService.GetFullQuiz(ctx.Request.Context(), int32(quizID))
	if err != nil {
		// Check for specific "not found" error from service
		if err.Error() == fmt.Sprintf("quiz with ID %d not found", quizID) { // Check specific error message
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// Log the unexpected error
		log.Printf("Error calling GetFullQuiz service for ID %d: %v", quizID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve full quiz details"})
		return
	}

	// Return the successfully retrieved structure
	ctx.JSON(http.StatusOK, fullQuiz)
}
