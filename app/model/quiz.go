package model

type Quiz struct {
	IDQuiz   string
	Author   string
	Title    string
	Desc     string
	Category string
	Duration int
	Privacy  string
	Picture  string
	Question []Quest
}

type Quest struct {
	Number   int
	Question string
	Media    string
	Answer   string
	OptionA  string
	OPtionB  string
	OptionC  string
	OptionD  string
}
