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
	Assets []map[string]string
	Port   string
}

func (app *App) setMiddlewere() {
	app.Router.Use(middleware.Loging)
}

func (app *App) SetAssetPath(key string, value string) {
	switch key {
	case "profile":
		path := map[string]string{"profile": value}
		app.Assets = append(app.Assets, path)
	case "media":
		path := map[string]string{"media": value}
		app.Assets = append(app.Assets, path)
	case "picture":
		path := map[string]string{"picture": value}
		app.Assets = append(app.Assets, path)
	}
}

func (app *App) setRoutes() {
	// 	app.Router.HandleFunc("/", app.handelRequest(handler.TestAPI))
	app.Router.HandleFunc("/signup", app.handelRequest(handler.SignUpFunc))
	app.Router.HandleFunc("/signin", app.handelRequest(handler.SignInFunc))
	app.Router.HandleFunc("/uploadpp", app.handelRequest(handler.UploadProfilePictureFunc))
	// 	app.Router.HandleFunc("/uploadqp", app.handelRequest(handler.UploadQuizPictureFunc))
	// 	app.Router.HandleFunc("/uploadqm", app.handelRequest(handler.UploadQuestionMediaFunc))
	app.Router.HandleFunc("/editprofile", app.handelRequest(handler.EditUserFunc))
	app.Router.HandleFunc("/createquiz", app.handelRequest(handler.CreateQuizFunc))
	// 	app.Router.HandleFunc("/editquiz", app.handelRequest(handler.EditQuizFunc))
	// 	app.Router.HandleFunc("/editquestion", app.handelRequest(handler.EditQuestionFunc))
	app.Router.HandleFunc("/user", app.handelRequest(handler.SelectUserProfileFunc))
	// 	app.Router.HandleFunc("/quizs", app.handelRequest(handler.SelectQuizByAuthorFunc))
	// 	app.Router.HandleFunc("/quiz", app.handelRequest(handler.SelectQuizByIDFunc))
}

func (app *App) Run() {
	app.Router = mux.NewRouter()
	testAuth := dbcontext.Auth{
		Username: "postgres",
		Passowrd: "root",
		DbName:   "db_quizy",
	}
	var err error
	app.DB, err = testAuth.Connection()
	if err != nil {
		log.Panic(err)
	}

	//app.DB.DropTableIfExists(&model.User{}, &model.Quiz{}, &model.Quest{}, &model.Player{}, &model.Option{}, &model.History{}, &model.Archievement{})
	//app.DB.Debug().CreateTable(&model.User{}, &model.Quiz{}, &model.Quest{}, &model.Player{}, &model.Option{}, &model.History{}, &model.Archievement{})
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
