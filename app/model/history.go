package model

type History struct {
	IDHistory string
	Quiz      Quiz
	Item      HistoryItem
}

type HistoryItem struct {
	User  User
	Score int
	Point int
}
