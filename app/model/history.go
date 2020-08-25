package model

import (
	"database/sql"
	"fmt"
)

type History struct {
	HistoryID string   `json:"historyID" gorm:"unique; not null"`
	Date      string   `json:"date" gorm:"not null"`
	QuizRefer string   `json:"quizref" gorm:"not null"`
	Players   []Player `json:"players"`
}

type Player struct {
	HistoryRefer string `json:"ref" gorm:"not null"`
	UserRefer    string `json:"userref"`
	Guest        string `json:"guest"`
	Score        uint   `json:"score" gorm:"not null"`
	Point        uint   `json:"point" gorm:"not null"`
}

func (h *History) InsertHistory(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public.histories(history_id, quiz_refer, date) VALUES ('%s', '%s', '%s')`, h.HistoryID, h.QuizRefer, h.Date)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}

	for _, item := range h.Players {
		_, err := item.InsertPlayer(db)
		if err != nil {
			return 0, err
		}
	}
	return affect, nil
}

func (p *Player) InsertPlayer(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public.players(history_refer, guest, score, point, user_refer) VALUES ('%s', '%s', %d, %d, '%s')`, p.HistoryRefer, p.Guest, p.Score, p.Point, p.UserRefer)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return affect, nil
}

func SelectPlayerByHistory(idHistory string, db *sql.DB) ([]Player, error) {
	query := fmt.Sprintf(`SELECT score, point, user_refer, guest FROM public.players WHERE history_refer = '%s'`, idHistory)
	row, err := db.Query(query)
	if err != nil {
		return []Player{}, err
	}
	var result []Player
	for row.Next() {
		var item Player
		if err := row.Scan(&item.Score, &item.Point, &item.UserRefer, &item.Guest); err != nil {
			return []Player{}, err
		}
		result = append(result, item)
	}
	return result, nil
}

func SelectPlayerByUserRefer(idUser string, db *sql.DB) ([]Player, error) {
	query := fmt.Sprintf(`SELECT history_refer, score, point, user_refer, guest FROM public.players WHERE user_refer = '%s'`, idUser)
	row, err := db.Query(query)
	if err != nil {
		return []Player{}, err
	}
	var result []Player
	for row.Next() {
		var item Player
		if err := row.Scan(&item.HistoryRefer, &item.Score, &item.Point, &item.UserRefer, &item.Guest); err != nil {
			return []Player{}, err
		}
		result = append(result, item)
	}
	return result, nil
}

func SelectHistoryByID(idHistory string, db *sql.DB) (History, error) {
	var result History
	result.Players, _ = SelectPlayerByHistory(idHistory, db)
	query := fmt.Sprintf(`SELECT history_id, quiz_refer, date FROM public.histories WHERE history_id = '%s'`, idHistory)
	row := db.QueryRow(query)
	if err := row.Scan(&result.HistoryID, &result.QuizRefer, &result.Date); err != nil {
		return History{}, err
	}
	return result, nil
}

func SelectHistoryByQuizID(idQuiz string, db *sql.DB) ([]History, error) {
	var sumresult []History
	query := fmt.Sprintf(`SELECT history_id, quiz_refer, date FROM public.histories WHERE quiz_refer = '%s'`, idQuiz)
	row, err := db.Query(query)
	if err != nil {
		return []History{}, err
	}
	for row.Next() {
		var result History
		if err := row.Scan(&result.HistoryID, &result.QuizRefer, &result.Date); err != nil {
			return []History{}, err
		}
		players, _ := SelectPlayerByHistory(result.HistoryID, db)
		result.Players = players
		sumresult = append(sumresult, result)
	}
	return sumresult, nil
}
