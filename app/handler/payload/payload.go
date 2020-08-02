package payload

import (
	"apiquizyfull/app/model"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PayloadQuiz struct {
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
	fmt.Println(ID)
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

type PayloadQuestion struct {
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
		Question: p.Question,
		Answer:   p.Answer,
		Options:  options,
	}
}

type PayloadOption struct {
	Symbol  string `json:"symbol" bson:"symbol"`
	Comment string `json:"comment" bson:"comment"`
}

func (p *PayloadOption) Convert() model.Option {
	return model.Option{
		Symbol:  p.Symbol,
		Comment: p.Comment,
	}
}

type PayloadSignIn struct {
	Email    string
	Password string
}

func (p *PayloadSignIn) Convert() model.User {
	return model.User{
		Email:    p.Email,
		Password: p.Password,
	}
}

type PayloadSignUp struct {
	Name      string
	BirthDate string
	Email     string
	Password  string
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
	fmt.Println(ID)
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
