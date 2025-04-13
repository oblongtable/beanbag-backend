package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/internal/services"
)

type QuizHandler struct {
	quizService *services.QuizService
}

func NewQuizHandler(quizService *services.QuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}

// CreateQuiz godoc
// @Summary Create a new quiz
// @Description Create a new quiz with the given details
// @Tags quizzes
// @Accept json
// @Produce json
// @Param quiz body CreateQuizRequest true "Quiz details"
// @Success 201 {object} db.Quiz
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quizzes [post]
func (h *QuizHandler) CreateQuiz(ctx *gin.Context) {

	type answer struct {
		Text      string `json:"text" binding:"required"`
		IsCorrect bool   `json:"isCorrect" binding:"required"`
	}

	type question struct {
		Text       string   `json:"text" binding:"required"`
		UseTimer   bool     `json:"useTimer" binding:"required"`
		TimerValue int32    `json:"timerValue" binding:"required"`
		Answers    []answer `json:"answers" binding:"required"`
	}

	var req struct {
		Title     string     `json:"title" binding:"required"`
		CreatorID int32      `json:"creatorId" binding:"required"`
		Questions []question `json:"questions" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz, err := h.quizService.CreateQuiz(ctx, req.Title, req.CreatorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, quiz)
}

// GetQuiz godoc
// @Summary Get a quiz by ID
// @Description Get a quiz by its ID
// @Tags quizzes
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} db.Quiz
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /quizzes/{id} [get]
func (h *QuizHandler) GetQuiz(ctx *gin.Context) {
	quizIDStr := ctx.Param("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	quiz, err := h.quizService.GetQuiz(ctx, int32(quizID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, quiz)
}

// CreateQuizRequest represents the request body for creating a quiz.
// @Description Quiz details
type CreateQuizRequest struct {
	Title     string `json:"title" binding:"required"`
	CreatorID int32  `json:"creator_id" binding:"required"`
}
