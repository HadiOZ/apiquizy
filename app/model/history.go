package model

type History struct {
	HistoryID string   `json:"historyID"`
	Quiz      Quiz     `json:"quiz"`
	QuizRefer string   `json:"quizref"`
	Players   []Player `json:"players"`
}

type Player struct {
	HistoryRefer string `json:"ref"`
	User         User   `json:"user"`
	UserRefer    uint   `json:"userref"`
	Score        uint   `json:"score"`
	Point        uint   `json:"point"`
}
