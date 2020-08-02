package model

type Archievement struct {
	ArchID    string `json:"archivementID"`
	QuizRefer string `json:"quizref"`
	Point     uint   `json:"point"`
	Rank      uint   `json:"rank"`
}
