package payload

import (
	"apiquizyfull/app/model"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

type PayloadQuiz struct {
	QuizID    string            `json:"id" bson:"id"`
	Picture   string            `json:"picture" bson:"picture"`
	Author    string            `json:"author" bson:"author"`
	Title     string            `json:"title" bson:"title"`
	Desc      string            `json:"description" bson:"description"`
	Category  string            `json:"category" bson:"category"`
	Duration  uint              `json:"duration" bson:"duration"`
	Privacy   string            `json:"privacy" bson:"privacy"`
	Questions []PayloadQuestion `json:"questions" bson:"questions"`
}

func (p *PayloadQuiz) GetID() string {
	time := time.Now()
	year := strconv.Itoa(time.Year())
	month := strings.ToUpper(time.Month().String())
	nano := strconv.Itoa(time.Nanosecond())

	var UserID []string
	UserID = append(UserID, "QUIZ")
	UserID = append(UserID, year)
	UserID = append(UserID, month[0:3])
	UserID = append(UserID, nano)

	ID := strings.Join(UserID, "-")
	encode := base64.StdEncoding.EncodeToString([]byte(ID))
	return encode
}

func (p *PayloadQuiz) Convert() model.Quiz {
	ID := p.GetID()
	var questions []model.Quest
	for index, item := range p.Questions {
		conv := item.Convert(ID, strconv.Itoa(index+1))
		conv.QuizRefer = ID
		conv.QuestID = strconv.Itoa(index + 1)
		questions = append(questions, conv)
	}

	result := model.Quiz{
		QuizID:    ID,
		UserRefer: p.Author,
		Title:     p.Title,
		Duration:  p.Duration,
		Privacy:   p.Privacy,
		Questions: questions,
	}
	result.Desc.String = p.Desc
	result.Category.String = p.Category
	return result
}

func QuizToPayload(quiz model.Quiz) PayloadQuiz {
	var questions []PayloadQuestion
	for _, item := range quiz.Questions {
		quest := questToPayload(item)
		questions = append(questions, quest)
	}
	return PayloadQuiz{
		QuizID:    quiz.QuizID,
		Picture:   quiz.Picture.String,
		Author:    quiz.UserRefer,
		Title:     quiz.Title,
		Desc:      quiz.Desc.String,
		Category:  quiz.Category.String,
		Duration:  quiz.Duration,
		Privacy:   quiz.Privacy,
		Questions: questions,
	}
}

type PayloadQuestion struct {
	QuestID  string          `json:"id" bson:"id"`
	Media    string          `json:"media" bson:"media"`
	Question string          `json:"question" bson:"question"`
	Answer   string          `json:"answer" bson:"answer"`
	Options  []PayloadOption `json:"options" bson:"options"`
}

func (p *PayloadQuestion) Convert(quizid string, questid string) model.Quest {
	var options []model.Option
	for _, item := range p.Options {
		conv := item.Convert()
		conv.QuestRefer = questid
		conv.QuizRefer = quizid
		options = append(options, conv)
	}

	return model.Quest{
		QuizRefer: quizid,
		QuestID:   questid,
		Question:  p.Question,
		Answer:    p.Answer,
		Options:   options,
	}
}

func questToPayload(quest model.Quest) PayloadQuestion {
	var options []PayloadOption
	for _, item := range quest.Options {
		option := optionToPayload(item)
		options = append(options, option)
	}

	return PayloadQuestion{
		QuestID:  quest.QuestID,
		Media:    quest.Media.String,
		Question: quest.Question,
		Answer:   quest.Answer,
		Options:  options,
	}
}

type PayloadOption struct {
	Symbol  string `json:"symbol" bson:"symbol"`
	Comment string `json:"comment" bson:"comment"`
}

func optionToPayload(option model.Option) PayloadOption {
	return PayloadOption{
		Symbol:  option.Symbol,
		Comment: option.Comment,
	}
}

func (p *PayloadOption) Convert() model.Option {
	return model.Option{
		Symbol:  p.Symbol,
		Comment: p.Comment,
	}
}

type PayloadSignIn struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (p *PayloadSignIn) Convert() model.User {
	return model.User{
		Email:    p.Email,
		Password: p.Password,
	}
}

type PayloadSignUp struct {
	Name      string `json:"name" bson:"name"`
	BirthDate string `json:"birthdate" bson:"birthdate"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

func (p *PayloadSignUp) GetID() string {
	time := time.Now()
	year := strconv.Itoa(time.Year())
	month := strings.ToUpper(time.Month().String())
	second := strconv.Itoa(time.Second())
	nano := strconv.Itoa(time.Nanosecond())

	var UserID []string
	UserID = append(UserID, year)
	UserID = append(UserID, month[0:3])
	UserID = append(UserID, second)
	UserID = append(UserID, nano)
	ID := strings.Join(UserID, "Q")
	encode := base64.StdEncoding.EncodeToString([]byte(ID))
	return encode
}

func (p *PayloadSignUp) Convert() model.User {
	result := model.User{
		UserID:   p.GetID(),
		Email:    p.Email,
		Password: p.Password,
		Name:     p.Name,
	}
	result.BirthDate.String = p.BirthDate
	return result
}

type PayloadProfile struct {
	Name        string `json:"name" bson:"name"`
	BirthDate   string `json:"birthdate" bson:"birthdate"`
	Gender      string `json:"gender" bson:"gender"`
	Country     string `json:"country" bson:"country"`
	Job         string `json:"job" bson:"job"`
	Institution string `json:"institution" bson:"institution"`
	Phone       string `json:"phone" bson:"phone"`
	Profile     string `json:"picture" bson:"picture"`
}

func (p *PayloadProfile) Convert() model.User {
	var result model.User
	result.Name = p.Name
	result.BirthDate.Scan(p.BirthDate)
	result.Gender.Scan(p.Gender)
	result.Country.Scan(p.Country)
	result.Job.Scan(p.Job)
	result.Institution.Scan(p.Institution)
	result.Phone.Scan(p.Phone)
	result.ProfilePicture.Scan(p.Profile)
	return result
}

type PayloadUpload struct {
	ID       string
	Filename string
}

func (p *PayloadUpload) ConvertToProfile() model.User {
	var result model.User
	result.ProfilePicture.String = p.Filename
	result.UserID = p.ID
	return result
}

func (p *PayloadUpload) ConvertToQuiz() model.Quiz {
	var result model.Quiz
	result.Picture.String = p.Filename
	result.QuizID = p.ID
	return result
}

func (p *PayloadUpload) ConvertToQuestion() model.Quest {
	var result model.Quest
	result.Media.String = p.Filename
	result.QuestID = p.ID
	return result
}

type PayloadHistory struct {
	HistoryID string          `json:"historyID" bson:"historyID"`
	Date      string          `json:"date" bson:"date"`
	QuizRefer string          `json:"quizrefer" bson:"quizrefer"`
	Players   []PayloadPlayer `json:"players" bson:"players"`
}

func (p *PayloadHistory) GetID() string {
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := strings.ToUpper(now.Month().String())
	nano := strconv.Itoa(now.Nanosecond())

	var UserID []string
	UserID = append(UserID, "HIZ")
	UserID = append(UserID, year)
	UserID = append(UserID, month[0:3])
	UserID = append(UserID, nano)

	ID := strings.Join(UserID, "")
	encode := base64.StdEncoding.EncodeToString([]byte(ID))
	return encode
}

func (p *PayloadHistory) Convert() model.History {
	var players []model.Player
	ID := p.GetID()
	for _, item := range p.Players {
		player := item.Convert(ID)
		players = append(players, player)
	}

	t := time.Now()
	var dates []string
	dates = append(dates, strconv.Itoa(t.Day()))
	dates = append(dates, t.Month().String())
	dates = append(dates, strconv.Itoa(t.Year()))
	date := strings.Join(dates, "-")

	result := model.History{
		HistoryID: ID,
		Date:      date,
		QuizRefer: p.QuizRefer,
		Players:   players,
	}
	return result
}

type PayloadPlayer struct {
	UserRefer string `json:"userref" bson:"userref"`
	Nickname  string `json:"nickname" bson:"nickname"`
	Score     uint   `json:"score" bson:"score"`
	Point     uint   `json:"point" bson:"point"`
}

func (p *PayloadPlayer) Convert(idhistory string) model.Player {
	return model.Player{
		UserRefer:    p.UserRefer,
		Guest:        p.Nickname,
		HistoryRefer: idhistory,
		Score:        p.Score,
		Point:        p.Point,
	}
}

func ConvertToPayloadHistory(history model.History) PayloadHistory {
	var players []PayloadPlayer
	for _, item := range history.Players {
		player := ConvertToPayloadPlayer(item)
		players = append(players, player)
	}
	return PayloadHistory{
		HistoryID: history.HistoryID,
		Date:      history.Date,
		QuizRefer: history.QuizRefer,
		Players:   players,
	}
}

func ConvertToPayloadPlayer(player model.Player) PayloadPlayer {
	return PayloadPlayer{
		UserRefer: player.UserRefer,
		Nickname:  player.Guest,
		Score:     player.Score,
		Point:     player.Point,
	}
}
