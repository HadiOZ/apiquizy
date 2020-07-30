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
}

func (app *App) setMiddlewere() {
	app.Router.Use(middleware.Loging)
}

func (app *App) setRoutes() {
	app.Router.HandleFunc("/", app.handelRequest(handler.TestAPI))
	app.Router.HandleFunc("/signup", app.handelRequest(handler.SignUp))
	app.Router.HandleFunc("/signin", app.handelRequest(handler.SignIn))
	app.Router.HandleFunc("/upload", app.handelRequest(handler.UploadProfile))
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
	app.setRoutes()
	app.setMiddlewere()
	fmt.Println("strat server at :8080")
	if err := http.ListenAndServe(":8080", app.Router); err != nil {
		log.Fatal(err)
	}
}

func (app *App) handelRequest(handler func(w http.ResponseWriter, r *http.Request, db *sql.DB)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, app.DB)
	}
}
