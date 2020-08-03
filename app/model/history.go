package model

type History struct {
	HistoryID string   `json:"historyID" gorm:"unique; not null"`
	Quiz      Quiz     `json:"quiz"`
	QuizRefer string   `json:"quizref" gorm:"not null"`
	Players   []Player `json:"players"`
}

type Player struct {
	HistoryRefer string `json:"ref" gorm:"not null"`
	User         User   `json:"user"`
	UserRefer    uint   `json:"userref" gorm:"not null"`
	Score        uint   `json:"score" gorm:"not null"`
	Point        uint   `json:"point" gorm:"not null"`
}
