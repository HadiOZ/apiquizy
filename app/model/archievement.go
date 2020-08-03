package model

type Archievement struct {
	ArchID    string `json:"archivementID" gorm:"unique; not null"`
	QuizRefer string `json:"quizref" gorm:"not null"`
	Point     uint   `json:"point" gorm:"not null"`
	Rank      uint   `json:"rank" gorm:"not null"`
}
