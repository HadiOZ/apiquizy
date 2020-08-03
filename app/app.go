package app

import (
	"apiquizyfull/app/handler"
	"apiquizyfull/app/middleware"
	"apiquizyfull/app/model/dbcontext"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
	Assets dbcontext.Assets
	Port   string
}

func (app *App) setMiddlewere() {
	app.Router.Use(middleware.Loging)
}

func (app *App) setRoutes() {
	// 	app.Router.HandleFunc("/", app.handelRequest(handler.TestAPI))
	app.Router.HandleFunc("/signup", app.handelRequest(handler.SignUpFunc))
	app.Router.HandleFunc("/signin", app.handelRequest(handler.SignInFunc))
	app.Router.HandleFunc("/uploadpp", app.handelRequestUpload(handler.UploadProfilePictureFunc))
	app.Router.HandleFunc("/uploadqp", app.handelRequestUpload(handler.UploadQuizPictureFunc))
	app.Router.HandleFunc("/uploadqm", app.handelRequestUpload(handler.UploadQuestionMediaFunc))
	app.Router.HandleFunc("/editprofile", app.handelRequest(handler.EditUserFunc))
	app.Router.HandleFunc("/createquiz", app.handelRequest(handler.CreateQuizFunc))
	app.Router.HandleFunc("/editquiz", app.handelRequest(handler.EditQuizFunc))
	app.Router.HandleFunc("/deletequiz", app.handelRequest(handler.DeleteQuiz))
	app.Router.HandleFunc("/addquestion", app.handelRequest(handler.AddQuestion))
	app.Router.HandleFunc("/editquestion", app.handelRequest(handler.EditQuestionFunc))
	app.Router.HandleFunc("/deletequestion", app.handelRequest(handler.DeleteQuestion))
	app.Router.HandleFunc("/user", app.handelRequest(handler.SelectUserProfileFunc))
	app.Router.HandleFunc("/quizdetail", app.handelRequest(handler.SelectQuizDetailFunc))
	app.Router.HandleFunc("/quiz", app.handelRequest(handler.SelectQuizFunc))
}

func (app *App) Run() {
	app.Router = mux.NewRouter()
	app.setRoutes()
	app.setMiddlewere()
	fmt.Println("strat server at" + app.Port)
	if err := http.ListenAndServe(app.Port, app.Router); err != nil {
		log.Fatal(err)
	}
}

func (app *App) handelRequest(handler func(w http.ResponseWriter, r *http.Request, db *sql.DB)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, app.DB)
	}
}

func (app *App) handelRequestUpload(handler func(w http.ResponseWriter, r *http.Request, db *sql.DB, path dbcontext.Assets)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, app.DB, app.Assets)
	}
}
