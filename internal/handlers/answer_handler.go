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
