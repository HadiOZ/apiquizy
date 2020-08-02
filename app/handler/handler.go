package handler

import (
	"apiquizyfull/app/handler/payload"
	"apiquizyfull/app/model"
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
	}

	s := payload.GetID()
	fmt.Println(s)

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
func UploadProfilePictureFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
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
			ID:       path.FormName(),
			Filename: strings.Join([]string{encrype, ext}, ""),
		}
		item := payload.ConvertToProfile()
		if affect, _ := item.EditProfilePicture(db); affect <= 0 {
			errid = errors.New("IDUser " + path.FormName() + " not found")
			break
		}

		filelocation := filepath.Join("/home/hadiese/Documents/quizymedia/profile-picture", payload.Filename)
		log.Println(filelocation)
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

	UserID := r.URL.Query().Get("userid")
	log.Println(UserID)
	decoder := json.NewDecoder(r.Body)
	var payload payload.PayloadProfile

	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
		return
	}
	profile := payload.Convert()
	profile.UserID = UserID
	log.Println(profile)
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
	idUser := r.URL.Query().Get("iduser")
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

	// var ap payload.PayloadQuiz
	// ap.Author = "hadi"
	// ap.Category = "infotec"
	// ap.Desc = "salah"
	// ap.Duration = 86
	// var quests []payload.PayloadQuestion
	// var quest payload.PayloadQuestion
	// quest.Answer = "kita kula"
	// var options []payload.PayloadOption
	// var opsi payload.PayloadOption
	// opsi.Comment = "adab"
	// options = append(options, opsi)
	// quest.Options = options
	// quests = append(quests, quest)
	// ap.Questions = quests

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

// func EditQuizFunc(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
// 	if r.Method != "POST" {
// 		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
// 		return
// 	}
// 	decoder := json.NewDecoder(r.Body)
// 	var payload payload.PayloadQuiz

// 	if err := decoder.Decode(&payload); err != nil {
// 		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
// 		return
// 	}
// 	quiz := payload.Convert()
// 	if affect, err := quiz.EditQuiz(db); err != nil || affect <= 0 {
// 		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
// 		log.Println(err)
// 		return
// 	}

// 	resposeJSON(w, http.StatusOK, "data seved")
// }

// func EditQuestionFunc(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
// 	if r.Method != "POST" {
// 		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
// 		return
// 	}
// 	decoder := json.NewDecoder(r.Body)
// 	payload := struct {
// 		idQuiz   string
// 		question payload.PayloadQuestion
// 	}{}

// 	if err := decoder.Decode(&payload); err != nil {
// 		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
// 		return
// 	}
// 	question := payload.question.Convert()
// 	if affect, err := question.EditQuestion(payload.idQuiz, db); err != nil || affect <= 0 {
// 		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
// 		log.Println(err)
// 		return
// 	}

// 	resposeJSON(w, http.StatusOK, "data seved")
// }

// func UploadQuizPictureFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	if r.Method != "POST" {
// 		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
// 		return
// 	}
// 	reader, err := r.MultipartReader()
// 	var errid error
// 	dir, _ := os.Getwd()

// 	if err != nil {
// 		resposeErrorJSON(w, http.StatusInternalServerError, "internal server error")
// 		return
// 	}

// 	for {
// 		path, err := reader.NextPart()
// 		if err == io.EOF {
// 			break
// 		}

// 		timer := time.Now()
// 		bit, _ := timer.MarshalText()
// 		filename := strings.Join([]string{path.FormName(), string(bit)}, "%")
// 		encrype := base64.StdEncoding.EncodeToString([]byte(filename))

// 		ext := filepath.Ext(path.FileName())
// 		payload := payload.PayloadUpload{
// 			ID:       path.FormName(),
// 			Filename: strings.Join([]string{encrype, ext}, ""),
// 		}
// 		item := payload.ConvertToQuiz()
// 		if affect, _ := item.UploadPicture(db); affect <= 0 {
// 			errid = errors.New("IDQuiz " + path.FormName() + " not found")
// 			break
// 		}

// 		filelocation := filepath.Join(dir, payload.Filename)
// 		file, err := os.OpenFile(filelocation, os.O_WRONLY|os.O_CREATE, 0666)
// 		defer file.Close()

// 		if err != nil {
// 			log.Panic(err)
// 			return

// 		}

// 		if _, err := io.Copy(file, path); err != nil {
// 			resposeErrorJSON(w, http.StatusInternalServerError, "can't save file")
// 			return
// 		}
// 	}

// 	if errid != nil {
// 		resposeErrorJSON(w, http.StatusBadRequest, errid.Error())
// 		return
// 	}

// 	resposeJSON(w, http.StatusOK, "data seved")
// }

// func UploadQuestionMediaFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	if r.Method != "POST" {
// 		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
// 		return
// 	}
// 	reader, err := r.MultipartReader()
// 	var errid error
// 	dir, _ := os.Getwd()

// 	if err != nil {
// 		resposeErrorJSON(w, http.StatusInternalServerError, "internal server error")
// 		return
// 	}

// 	for {
// 		path, err := reader.NextPart()
// 		if err == io.EOF {
// 			break
// 		}

// 		timer := time.Now()
// 		bit, _ := timer.MarshalText()
// 		id := strings.Split(path.FormName(), "-")
// 		idQuiz := id[0]
// 		idQuestion := id[1]
// 		filename := strings.Join([]string{idQuestion, string(bit)}, "%")
// 		encrype := base64.StdEncoding.EncodeToString([]byte(filename))

// 		ext := filepath.Ext(path.FileName())
// 		payload := payload.PayloadUpload{
// 			ID:       idQuestion,
// 			Filename: strings.Join([]string{encrype, ext}, ""),
// 		}
// 		item := payload.ConvertToQuestion()
// 		if affect, _ := item.UploadMedia(idQuiz, db); affect <= 0 {
// 			errid = errors.New("IDQuiz " + path.FormName() + " not found")
// 			break
// 		}

// 		filelocation := filepath.Join(dir, payload.Filename)
// 		file, err := os.OpenFile(filelocation, os.O_WRONLY|os.O_CREATE, 0666)
// 		defer file.Close()

// 		if err != nil {
// 			log.Panic(err)
// 			return

// 		}

// 		if _, err := io.Copy(file, path); err != nil {
// 			resposeErrorJSON(w, http.StatusInternalServerError, "can't save file")
// 			return
// 		}
// 	}

// 	if errid != nil {
// 		resposeErrorJSON(w, http.StatusBadRequest, errid.Error())
// 		return
// 	}

// 	resposeJSON(w, http.StatusOK, "data seved")
// }

// func SelectQuizByAuthorFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	author := r.URL.Query().Get("author")
// 	quizs, err := model.SelectQuizByAuthor(author, db)
// 	if err != nil {
// 		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	resposeJSON(w, http.StatusOK, quizs)

// }

// func SelectQuizByIDFunc(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	id := r.URL.Query().Get("id")
// 	quiz, err := model.SelectQuizByID(id, db)
// 	if err != nil {
// 		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	resposeJSON(w, http.StatusOK, quiz)

// }
