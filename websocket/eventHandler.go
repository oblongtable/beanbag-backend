package websocket

func CreateRoomEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.RegisterRoom <- cliEvt
	return nil
}

func JoinRoomEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.JoinRoom <- cliEvt
	return nil
}

func LeaveRoomEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.LeaveRoom <- cliEvt
	return nil
}

func StartQuizEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.StartQuiz<- cliEvt
	return nil
}

func ForwardQuizEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.ForwardQuiz <- cliEvt
	return nil
}

func SubmitAnswerEventHandler(cliEvt *ClientEvent) error {
	cliEvt.Requester.Wssvr.SubmitAnswer <- cliEvt
	return nil
}
