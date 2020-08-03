package model

import (
	"database/sql"
	"fmt"
	"log"
)

type Quiz struct {
	QuizID    string         `json:"quizID" gorm:"unique; not null"`
	UserRefer string         `json:"author" gorm:"not null"`
	Title     string         `json:"title" gorm:"not null"`
	Desc      sql.NullString `json:"description"`
	Category  sql.NullString `json:"category"`
	Duration  uint           `json:"duraion" gorm:"not null"`
	Privacy   string         `json:"privacy" gorm:"not null"`
	Picture   sql.NullString `json:"picture" gorm:"unique"`
	Questions []Quest        `json:"questions"`
}

type Quest struct {
	QuestID   string         `json:"questID" gorm:"unique;not null"`
	QuizRefer string         `json:"quizref" gorm:"not null"`
	Question  string         `json:"question" gorm:"not null"`
	Media     sql.NullString `json:"media" gorm:"unique"`
	Answer    string         `json:"answer" gorm:"not null"`
	Options   []Option       `json:"options"`
}

type Option struct {
	QuizRefer  string `json:"quizref" gorm:"not null"`
	QuestRefer string `json:"questref" gorm:"not null"`
	Symbol     string `json:"symbol" gorm:"not null"`
	Comment    string `json:"comment" gorm:"not null"`
}

//checked
func (q *Quiz) CreateQuiz(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public."quizzes"(quiz_id, user_refer, title, "desc", category, duration, privacy) VALUES ('%s', '%s', '%s', '%s', '%s', %d, '%s')`, q.QuizID, q.UserRefer, q.Title, q.Desc.String, q.Category.String, q.Duration, q.Privacy)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	for _, item := range q.Questions {
		affect, err := item.AddQuestion(db)
		if affect <= 0 || err != nil {
			return affect, err
		}
	}
	return affect, nil
}

func (q *Quiz) EditQuiz(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public.quizzes SET title= '%s', "desc"= '%s', category= '%s', duration= %d, privacy= '%s' WHERE quiz_id= '%s'`, q.Title, q.Desc.String, q.Category.String, q.Duration, q.Privacy, q.QuizID)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	for _, item := range q.Questions {
		affect, err := item.EditQuestion(db)
		if err != nil {
			return affect, err
		}
	}
	return affect, nil
}

func DeleteQuiz(quizid string, db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM public.quizzes WHERE quiz_id= '%s'`, quizid)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	query = fmt.Sprintf(`DELETE FROM public.quests WHERE quiz_refer= '%s'`, quizid)
	res, err = db.Exec(query)
	if err != nil {
		return 0, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}

	query = fmt.Sprintf(`DELETE FROM public.options WHERE quiz_refer= '%s'`, quizid)
	res, err = db.Exec(query)
	if err != nil {
		return 0, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

//cheked
func (q *Quiz) UploadPicture(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public.quizzes SET picture= '%s' WHERE  quiz_id= '%s'`, q.Picture.String, q.QuizID)
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
func (q *Quest) AddQuestion(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public.quests(quest_id, quiz_refer, question, answer) VALUES ('%s', '%s', '%s', '%s')`, q.QuestID, q.QuizRefer, q.Question, q.Answer)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	for _, item := range q.Options {
		affect, err := item.AddOption(db)
		if affect <= 0 || err != nil {
			return affect, err
		}
	}
	return affect, nil
}

//checked
func (q *Quest) EditQuestion(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public.quests SET question= '%s', answer= '%s' WHERE quiz_refer= '%s' AND quest_id= '%s'`, q.Question, q.Answer, q.QuizRefer, q.QuestID)
	log.Println(query)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	for _, item := range q.Options {
		affect, err := item.UpdateOption(db)
		if err != nil {
			return affect, err
		}
	}
	return affect, nil
}

func DeleteQuestion(quizid string, questid string, db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM public.quests WHERE quiz_refer= '%s' AND quest_id= '%s'`, quizid, questid)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	affect, err = DeleteOption(quizid, questid, db)
	return affect, nil
}

//checked
func (q *Quest) UploadMedia(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public.quests SET media= '%s' WHERE quiz_refer = '%s' AND quest_id = '%s'`, q.Media.String, q.QuizRefer, q.QuestID)
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

//checked
func SelectQuizDetail(idquiz string, idAuthor string, db *sql.DB) ([]Quiz, error) {
	var sql string
	if idquiz != "" && idAuthor != "" {
		sql = `SELECT quiz_id, user_refer, title, "desc", category, duration, privacy, picture FROM public.quizzes WHERE user_refer = '%s' AND quiz_id = '%s'`
	} else {
		sql = `SELECT quiz_id, user_refer, title, "desc", category, duration, privacy, picture FROM public.quizzes WHERE user_refer = '%s' OR quiz_id = '%s'`
	}
	query := fmt.Sprintf(sql, idAuthor, idquiz)
	log.Println(query)
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
		item.Questions, err = SelectQuestions(item.QuizID, db)
		if err != nil {
			log.Println(err)
		}
		result = append(result, item)
	}
	return result, nil
}

//checked
func SelectQuiz(idquiz string, idAuthor string, db *sql.DB) ([]Quiz, error) {
	var sql string
	if idquiz != "" && idAuthor != "" {
		sql = `SELECT quiz_id, user_refer, title, "desc", category, duration, privacy, picture FROM public.quizzes WHERE user_refer = '%s' AND quiz_id = '%s'`
	} else {
		sql = `SELECT quiz_id, user_refer, title, "desc", category, duration, privacy, picture FROM public.quizzes WHERE user_refer = '%s' OR quiz_id = '%s'`
	}
	query := fmt.Sprintf(sql, idAuthor, idquiz)
	log.Println(query)
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

//function checked
func SelectQuestions(idQuiz string, db *sql.DB) ([]Quest, error) {
	query := fmt.Sprintf(`SELECT quest_id, question, media, answer FROM public.quests WHERE quiz_refer = '%s'`, idQuiz)
	row, err := db.Query(query)
	if err != nil {
		return []Quest{}, err
	}
	var result []Quest
	for row.Next() {
		item := Quest{}
		if err := row.Scan(&item.QuestID, &item.Question, &item.Media, &item.Answer); err != nil {
			return []Quest{}, err
		}
		item.Options, err = SelectOptions(idQuiz, item.QuestID, db)
		if err != nil {
			log.Println(err)
		}
		result = append(result, item)
	}
	return result, nil
}

//checked
func (o *Option) AddOption(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public.options(quest_refer, symbol, comment, quiz_refer) VALUES ('%s', '%s', '%s', '%s')`, o.QuestRefer, o.Symbol, o.Comment, o.QuizRefer)
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

//fumction checked
func SelectOptions(idQuiz string, idQuest string, db *sql.DB) ([]Option, error) {
	query := fmt.Sprintf(`SELECT symbol, comment FROM public.options WHERE quiz_refer = '%s' AND quest_refer = '%s'`, idQuiz, idQuest)
	row, err := db.Query(query)
	if err != nil {
		return []Option{}, err
	}
	var result []Option
	for row.Next() {
		item := Option{}
		if err := row.Scan(&item.Symbol, &item.Comment); err != nil {
			return []Option{}, err
		}
		result = append(result, item)
	}
	return result, nil
}

//make options update
func (o *Option) UpdateOption(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public.options SET symbol= '%s', comment= '%s' WHERE quiz_refer= '%s' AND quest_refer= '%s'`, o.Symbol, o.Comment, o.QuizRefer, o.QuestRefer)
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

func DeleteOption(quizid string, questid string, db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM public.options WHERE quiz_refer= '%s' AND quest_refer= '%s'`, quizid, questid)
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
