package apimodels

type APIAnswer struct {
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"isCorrect" binding:"required"`
}

type APIQuestion struct {
	Text       string   `json:"text" binding:"required"`
	UseTimer   bool     `json:"useTimer" binding:"required"`
	TimerValue int32    `json:"timerValue" binding:"required"`
	Answers    []APIAnswer `json:"answers" binding:"required"`
}