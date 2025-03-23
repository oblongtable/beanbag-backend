package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/internal/services"
)

type QuestionHandler struct {
	questionService *services.QuestionService
}

func NewQuestionHandler(questionService *services.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

func (h *QuestionHandler) CreateQuestion(ctx *gin.Context) {
	var req struct {
		QuizID      int32  `json:"quiz_id" binding:"required"`
		Description string `json:"description" binding:"required"`
		TimerOption bool   `json:"timer_option"`
		Timer       int32  `json:"timer"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.questionService.CreateQuestion(ctx, req.QuizID, req.Description, req.TimerOption, req.Timer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, question)
}

func (h *QuestionHandler) GetQuestion(ctx *gin.Context) {
	questionIDStr := ctx.Param("id")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	question, err := h.questionService.GetQuestion(ctx, int32(questionID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, question)
}
