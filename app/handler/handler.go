package handler

import (
	"apiquizyfull/app/handler/payload"
	"apiquizyfull/app/model"
	"apiquizyfull/app/model/dbcontext"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// func TestAPI(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	resposeJSON(w, http.StatusOK, "halo dari server")
// }

//SignUpFunc checked
func SignUpFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}

	decoder := json.NewDecoder(r.Body)

	var payload payload.PayloadSignUp

	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
		return
	}
	user := payload.Convert()
	if _, err := user.CreateUser(db); err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	resposeJSON(w, http.StatusOK, "Data Recorded")
}

//SignInFunc checked
func SignInFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var payload payload.PayloadSignIn

	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
		return
	}
	user := payload.Convert()
	auth, err := user.CheckUSer(db)
	if err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}
	if user.Password != auth.Password {
		resposeErrorJSON(w, http.StatusBadRequest, "password wrong")
		log.Println(err)
		return
	}

	resposeJSON(w, http.StatusOK, auth.UserID)
}

//UploadProfilePictureFunc checked
func UploadProfilePictureFunc(w http.ResponseWriter, r *http.Request, db *sql.DB, root dbcontext.Assets) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
	UserID := r.URL.Query().Get("id")
	reader, err := r.MultipartReader()
	var errid error
	//dir, _ := os.Getwd()

	if err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	for {
		path, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		name := time.Now()
		bit, _ := name.MarshalText()
		filename := strings.Join([]string{path.FormName(), string(bit)}, "%")
		encrype := base64.StdEncoding.EncodeToString([]byte(filename))

		ext := filepath.Ext(path.FileName())
		payload := payload.PayloadUpload{
			ID:       UserID,
			Filename: strings.Join([]string{encrype, ext}, ""),
		}
		item := payload.ConvertToProfile()
		if affect, _ := item.EditProfilePicture(db); affect <= 0 {
			errid = errors.New("IDUser " + path.FormName() + " not found")
			break
		}

		filelocation := filepath.Join(root.Profile, payload.Filename)
		file, err := os.OpenFile(filelocation, os.O_WRONLY|os.O_CREATE, 0666)
		defer file.Close()

		if err != nil {
			log.Panic(err)
			return

		}

		if _, err := io.Copy(file, path); err != nil {
			resposeErrorJSON(w, http.StatusInternalServerError, "can't save file")
			return
		}
	}

	if errid != nil {
		resposeErrorJSON(w, http.StatusBadRequest, errid.Error())
		return
	}

	resposeJSON(w, http.StatusOK, "data seved")
}

func EditUserFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}

	UserID := r.URL.Query().Get("id")
	decoder := json.NewDecoder(r.Body)
	var payload payload.PayloadProfile

	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
		return
	}
	profile := payload.Convert()
	profile.UserID = UserID
	if _, err := profile.EditProfile(db); err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resposeJSON(w, http.StatusOK, "data seved")
}

func SelectUserProfileFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "GET" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method GET")
		return
	}
	idUser := r.URL.Query().Get("id")
	profile, err := model.SelectUserProfile(idUser, db)

	payload := payload.PayloadProfile{
		Name:        profile.Name,
		BirthDate:   profile.BirthDate.String,
		Gender:      profile.Gender.String,
		Country:     profile.Country.String,
		Job:         profile.Job.String,
		Institution: profile.Institution.String,
		Phone:       profile.Phone.String,
		Profile:     profile.ProfilePicture.String,
	}

	if err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	resposeJSON(w, http.StatusOK, payload)

}

//CreateQuizFunc checked
func CreateQuizFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
	decoder := json.NewDecoder(r.Body)

	var payload payload.PayloadQuiz
	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
		return
	}

	quiz := payload.Convert()
	if affect, err := quiz.CreateQuiz(db); err != nil || affect <= 0 {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resposeJSON(w, http.StatusOK, "data was seved")
}

func EditQuizFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	QuizID := r.URL.Query().Get("id")
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var payload payload.PayloadQuiz

	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
		return
	}
	quiz := payload.Convert()
	quiz.QuizID = QuizID
	if affect, err := quiz.EditQuiz(db); err != nil || affect <= 0 {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resposeJSON(w, http.StatusOK, "data seved")
}

func EditQuestionFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
	QuizID := r.URL.Query().Get("quizid")
	QuestID := r.URL.Query().Get("id")
	decoder := json.NewDecoder(r.Body)
	var payload payload.PayloadQuestion

	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
		return
	}
	question := payload.Convert(QuizID, QuestID)
	log.Println(question)
	if affect, err := question.EditQuestion(db); affect <= 0 {
		resposeErrorJSON(w, http.StatusInternalServerError, "Data not change please check again your data that will be change")
		log.Println(err)
		return
	}

	resposeJSON(w, http.StatusOK, "data seved")
}

func UploadQuizPictureFunc(w http.ResponseWriter, r *http.Request, db *sql.DB, root dbcontext.Assets) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}

	QuizID := r.URL.Query().Get("id")
	reader, err := r.MultipartReader()
	var errid error
	//dir, _ := os.Getwd()

	if err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	for {
		path, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		timer := time.Now()
		bit, _ := timer.MarshalText()
		filename := strings.Join([]string{QuizID, string(bit)}, "%")
		encrype := base64.StdEncoding.EncodeToString([]byte(filename))

		ext := filepath.Ext(path.FileName())
		payload := payload.PayloadUpload{
			ID:       QuizID,
			Filename: strings.Join([]string{encrype, ext}, ""),
		}
		item := payload.ConvertToQuiz()
		if affect, _ := item.UploadPicture(db); affect <= 0 {
			errid = errors.New("IDQuiz " + path.FormName() + " not found")
			break
		}

		filelocation := filepath.Join(root.Quiz, payload.Filename)
		file, err := os.OpenFile(filelocation, os.O_WRONLY|os.O_CREATE, 0666)
		defer file.Close()

		if err != nil {
			log.Panic(err)
			return

		}

		if _, err := io.Copy(file, path); err != nil {
			resposeErrorJSON(w, http.StatusInternalServerError, "can't save file")
			return
		}
	}

	if errid != nil {
		resposeErrorJSON(w, http.StatusBadRequest, errid.Error())
		return
	}

	resposeJSON(w, http.StatusOK, "data seved")
}

func UploadQuestionMediaFunc(w http.ResponseWriter, r *http.Request, db *sql.DB, root dbcontext.Assets) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}

	QuizID := r.URL.Query().Get("quizid")

	reader, err := r.MultipartReader()
	var missing []string
	//dir, _ := os.Getwd()

	if err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	for {
		path, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		timer := time.Now()
		bit, _ := timer.MarshalText()

		filename := strings.Join([]string{QuizID, path.FormName(), string(bit)}, "%")
		encrype := base64.StdEncoding.EncodeToString([]byte(filename))

		ext := filepath.Ext(path.FileName())
		payload := payload.PayloadUpload{
			ID:       path.FormName(),
			Filename: strings.Join([]string{encrype, ext}, ""),
		}
		item := payload.ConvertToQuestion()
		item.QuizRefer = QuizID
		if affect, _ := item.UploadMedia(db); affect <= 0 {
			missing = append(missing, path.FormName())
			continue
		}

		filelocation := filepath.Join(root.Question, payload.Filename)
		file, err := os.OpenFile(filelocation, os.O_WRONLY|os.O_CREATE, 0666)
		defer file.Close()

		if err != nil {
			log.Panic(err)
			return

		}

		if _, err := io.Copy(file, path); err != nil {
			resposeErrorJSON(w, http.StatusInternalServerError, "can't save file")
			return
		}
	}

	if missing != nil {
		message := fmt.Sprint("ID ", missing, " not Found")
		resposeErrorJSON(w, http.StatusBadRequest, message)
		return
	}

	resposeJSON(w, http.StatusOK, "data seved")
}

func SelectQuizDetailFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "GET" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method GET")
		return
	}
	author := r.URL.Query().Get("author")
	quizID := r.URL.Query().Get("id")
	quizs, err := model.SelectQuizDetail(quizID, author, db)
	if err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	var pld []payload.PayloadQuiz
	for _, item := range quizs {
		q := payload.QuizToPayload(item)
		pld = append(pld, q)
	}
	resposeJSON(w, http.StatusOK, pld)

}

func SelectQuizFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "GET" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method GET")
		return
	}
	id := r.URL.Query().Get("id")
	author := r.URL.Query().Get("author")
	quizs, err := model.SelectQuiz(id, author, db)
	if err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	var pld []payload.PayloadQuiz
	for _, item := range quizs {
		q := payload.QuizToPayload(item)
		pld = append(pld, q)
	}

	resposeJSON(w, http.StatusOK, pld)

}

func DeleteQuestion(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "DELETE" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method DELETE")
		return
	}
	QuestID := r.URL.Query().Get("id")
	QuizID := r.URL.Query().Get("quizid")
	affect, err := model.DeleteQuestion(QuizID, QuestID, db)
	if affect <= 0 || err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Question ID or Quiz ID not Found")
		return
	}

	resposeJSON(w, http.StatusOK, "Data deleted")
}

func DeleteQuiz(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "DELETE" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method DELETE")
		return
	}
	QuizID := r.URL.Query().Get("id")
	affect, err := model.DeleteQuiz(QuizID, db)
	log.Println(affect, err)
	if affect <= 0 || err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Quiz ID not Found")
		return
	}

	resposeJSON(w, http.StatusOK, "Data deleted")
}

func AddQuestion(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
	decoder := json.NewDecoder(r.Body)
	QuizID := r.URL.Query().Get("quizid")
	QuestID := r.URL.Query().Get("id")
	var payload payload.PayloadQuestion
	log.Println(QuestID, QuizID)

	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
		return
	}

	question := payload.Convert(QuizID, QuestID)
	log.Println(question)
	if _, err := question.AddQuestion(db); err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	resposeJSON(w, http.StatusOK, "Data Recorded")
}
