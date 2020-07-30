package model

import (
	"database/sql"
	"fmt"
)

type Quiz struct {
	IDQuiz    string
	Author    string
	Title     string
	Desc      string
	Category  []string
	Duration  int
	Privacy   string
	Picture   string
	Questions []Quest
}

type Quest struct {
	IDQuest  int
	Question string
	Media    string
	Answer   string
	Options  []string
}

func CreateQuiz(quiz Quiz, db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public."Quiz"(id_quiz, author, title, description, category, duration, privacy) VALUES (?, ?, ?, ?, ?, ?, ?)`, quiz.IDQuiz, quiz.Author, quiz.Title, quiz.Desc, quiz.Category, quiz.Duration, quiz.Privacy)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	for _, item := range quiz.Questions {
		affect, err := AddQuestion(quiz.IDQuiz, item, db)
		if affect <= 0 || err != nil {
			return affect, err
		}
	}
	return affect, nil
}

func EditQuiz()   {}
func DeleteQuiz() {}

func AddQuestion(idQuiz string, question Quest, db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public."Questions"(id_quiz, id_question, question, answer, option) VALUES (?, ?, ?, ?, ?)`, idQuiz, question.IDQuest, question.Question, question.Answer, question.Options)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}
func EditQuestion()   {}
func DeleteQuestion() {}
