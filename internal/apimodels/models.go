package apimodels

type AnswerApiModel struct {
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"isCorrect" binding:"required"`
}

type QuestionApiModel struct {
	Text       string           `json:"text" binding:"required"`
	UseTimer   bool             `json:"useTimer" binding:"required"`
	TimerValue int32            `json:"timerValue" binding:"required"`
	Answers    []AnswerApiModel `json:"answers"`
}

type QuizApiModel struct {
	Title     string             `json:"title" binding:"required"`
	QuizID    int32              `json:"quiz_id"`
	CreatorID int32              `json:"creator_id" binding:"required"`
	Questions []QuestionApiModel `json:"questions"`
}
