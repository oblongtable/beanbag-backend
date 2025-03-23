package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/internal/services"
)

type AnswerHandler struct {
	answerService *services.AnswerService
}

func NewAnswerHandler(answerService *services.AnswerService) *AnswerHandler {
	return &AnswerHandler{answerService: answerService}
}

// CreateAnswer godoc
// @Summary Create a new answer
// @Description Create a new answer with the given details
// @Tags answers
// @Accept json
// @Produce json
// @Param answer body CreateAnswerRequest true "Answer details"
// @Success 201 {object} db.Answer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /answers [post]
func (h *AnswerHandler) CreateAnswer(ctx *gin.Context) {
	var req struct {
		QuestionID  int32  `json:"question_id" binding:"required"`
		Description string `json:"description" binding:"required"`
		IsCorrect   bool   `json:"is_correct"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer, err := h.answerService.CreateAnswer(ctx, req.QuestionID, req.Description, req.IsCorrect)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, answer)
}

// GetAnswer godoc
// @Summary Get an answer by ID
// @Description Get an answer by its ID
// @Tags answers
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {object} db.Answer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /answers/{id} [get]
func (h *AnswerHandler) GetAnswer(ctx *gin.Context) {
	answerIDStr := ctx.Param("id")
	answerID, err := strconv.Atoi(answerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	answer, err := h.answerService.GetAnswer(ctx, int32(answerID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, answer)
}

// CreateAnswerRequest represents the request body for creating an answer.
// @Description Answer details
type CreateAnswerRequest struct {
	QuestionID  int32  `json:"question_id" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsCorrect   bool   `json:"is_correct"`
}
