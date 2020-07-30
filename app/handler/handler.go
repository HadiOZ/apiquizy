package handler

import (
	"apiquizyfull/app/model"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func TestAPI(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	resposeJSON(w, http.StatusOK, "halo dari server")
}

//checked
func SignUp(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}

	decoder := json.NewDecoder(r.Body)

	payload := struct {
		IDUser    string
		Name      string
		BirthDate string
		Email     string
		Password  string
	}{}

	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
	}
	user := model.User{
		IDUser:   payload.IDUser,
		Email:    payload.Email,
		Password: payload.Password,
	}
	user.Detail.BirthDate = payload.BirthDate
	user.Detail.Name = payload.Name

	if affect, err := model.CreateUser(user, db); err != nil || affect <= 0 {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resposeJSON(w, http.StatusOK, "data seved")
}

//checked
func SignIn(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
	if err := r.ParseForm(); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Parse form error")
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	auth, err := model.CheckUSer(email, db)
	if err != nil {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}
	if password != auth.Password {
		resposeErrorJSON(w, http.StatusBadRequest, "password wrong")
		log.Println(err)
		return
	}

	resposeJSON(w, http.StatusOK, auth.IDUser)
}

func CreateQuiz(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		resposeErrorJSON(w, http.StatusBadRequest, "Just Allow Method POST")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var payload model.Quiz
	if err := decoder.Decode(&payload); err != nil {
		resposeErrorJSON(w, http.StatusBadRequest, "Data Structure wrong")
	}

	if affect, err := model.CreateQuiz(payload, db); err != nil || affect <= 0 {
		resposeErrorJSON(w, http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	resposeJSON(w, http.StatusOK, "data seved")
}

//checked
func UploadProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	reader, err := r.MultipartReader()
	var errid error
	dir, _ := os.Getwd()

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
		encrype := base64.StdEncoding.EncodeToString(bit)

		ext := filepath.Ext(path.FileName())
		filename := strings.Join([]string{encrype, ext}, "")
		userID := path.FormName()

		if affect, _ := model.EditProfilePicture(userID, filename, db); affect <= 0 {
			errid = errors.New("IDUser not found")
			break
		}

		filelocation := filepath.Join(dir, filename)
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
