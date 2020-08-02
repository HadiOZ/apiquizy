package model

import (
	"database/sql"
	"fmt"
	"log"
)

type Quiz struct {
	QuizID    string         `json:"quizID"`
	UserRefer string         `json:"author"`
	Title     string         `json:"title"`
	Desc      sql.NullString `json:"description"`
	Category  sql.NullString `json:"category"`
	Duration  uint           `json:"duraion"`
	Privacy   string         `json:"privacy"`
	Picture   sql.NullString `json:"picture"`
	Questions []Quest        `json:"questions"`
}

type Quest struct {
	QuestID   string         `json:"questID"`
	QuizRefer string         `json:"quizref"`
	Question  string         `json:"question"`
	Media     sql.NullString `json:"media"`
	Answer    string         `json:"answer"`
	Options   []Option       `json:"options"`
}

type Option struct {
	QuizRefer  string `json:"quizref"`
	QuestRefer string `json:"questref"`
	Symbol     string `json:"symbol"`
	Comment    string `json:"comment"`
}

//checked
func (q *Quiz) CreateQuiz(db *sql.DB) (int64, error) {
	log.Println(q)
	query := fmt.Sprintf(`INSERT INTO public."quizzes"(quiz_id, user_refer, title, "desc", category, duration, privacy) VALUES ('%s', '%s', '%s', '%s', '%s', %d, '%s')`, q.QuizID, q.UserRefer, q.Title, q.Desc.String, q.Category.String, q.Duration, q.Privacy)
	log.Println(query)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	for _, item := range q.Questions {
		affect, err := item.AddQuestion(q.QuizID, db)
		if affect <= 0 || err != nil {
			return affect, err
		}
	}
	return affect, nil
}

func (q *Quiz) EditQuiz(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public."Quiz" SET title= '%s', description= '%s', category= '%s', duration= %d, privacy= '%s' WHERE id_quiz= '%s'`, q.Title, q.Desc.String, q.Category.String, q.Duration, q.Privacy, q.QuizID)
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

func (q *Quiz) DeleteQuiz(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM public."Quiz" WHERE id_quiz= '%s'`, q.QuizID)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	query = fmt.Sprintf(`DELETE FROM public."Questions" WHERE id_quiz= '%s'`, q.QuizID)
	res, err = db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

func (q *Quiz) UploadPicture(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public."Quiz" SET picture= '%s' WHERE id_quiz= '%s'`, q.Picture.String, q.QuizID)
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

// //checked
func (q *Quest) AddQuestion(idQuiz string, db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public.quests(quest_id, quiz_refer, question, answer) VALUES ('%s', '%s', '%s', '%s')`, q.QuestID, idQuiz, q.Question, q.Answer)
	log.Println(query)
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

func (q *Quest) EditQuestion(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public."Questions" SET question= '%s', answer= '%s' WHERE id_quiz= '%s' AND id_question= '%s'`, q.Question, q.Answer, q.QuizRefer, q.QuestID)
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

func (q *Quest) DeleteQuestion(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM public."Questions" WHERE id_quiz= '%s' AND id_question= '%s'`, q.QuizRefer, q.QuestID)
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

func (q *Quest) UploadMedia(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public."Questions" SET media= '%s' WHERE id_quiz= '%s' AND id_question= '%s'`, q.Media.String, q.QuizRefer, q.QuestID)
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

func (q *Quiz) SelectQuizByAuthor(idAuthor string, db *sql.DB) ([]Quiz, error) {
	query := fmt.Sprintf(`SELECT id_quiz, author, title, description, category, duration, privacy, picture FROM public."Quiz" WHERE author = '%s'`, q.UserRefer)
	row, err := db.Query(query)
	if err != nil {
		return []Quiz{}, err
	}
	var result []Quiz
	for row.Next() {
		item := Quiz{}
		if err := row.Scan(&item.QuizID, &item.UserRefer, &item.Title, &item.Desc, &item.Category, &item.Duration, &item.Privacy, &item.Picture); err != nil {
			return []Quiz{}, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (q *Quiz) SelectQuizByID(db *sql.DB) (Quiz, error) {
	var item Quiz
	questions, err := SelectQuestions(q.QuizID, db)
	if err != nil {
		return Quiz{}, err
	}
	item.Questions = questions
	query := fmt.Sprintf(`SELECT id_quiz, author, title, description, category, duration, privacy, picture FROM public."Quiz" WHERE iq_quiz = '%s'`, q.QuizID)
	row, err := db.Query(query)
	if err != nil {
		return Quiz{}, err
	}
	for row.Next() {
		if err := row.Scan(&item.QuizID, &item.UserRefer, &item.Title, &item.Desc, &item.Category, &item.Duration, &item.Privacy, &item.Picture); err != nil {
			return Quiz{}, err
		}
	}
	return item, nil
}

func SelectQuestions(idQuiz string, db *sql.DB) ([]Quest, error) {
	query := fmt.Sprintf(`SELECT id_question, question, media, answer, option FROM public."Questions" WHERE id_quiz= '%s'`, idQuiz)
	row, err := db.Query(query)
	if err != nil {
		return []Quest{}, err
	}
	var result []Quest
	for row.Next() {
		item := Quest{}
		if err := row.Scan(&item.QuestID, &item.Question, &item.Media, &item.Answer, &item.Options); err != nil {
			return []Quest{}, err
		}
		result = append(result, item)
	}
	return result, nil
}
