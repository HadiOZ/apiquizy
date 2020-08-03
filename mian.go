package main

import (
	"apiquizyfull/app"
	"apiquizyfull/app/model/dbcontext"
	"log"
)

func main() {
	var API app.App
	API.Port = ":8080"
	API.Assets = dbcontext.Assets{
		Profile:  "/home/hadiese/Documents/quizymedia/profile-picture",
		Quiz:     "/home/hadiese/Documents/quizymedia/quiz-picture",
		Question: "/home/hadiese/Documents/quizymedia/question-media",
	}
	database := dbcontext.Auth{
		Username: "postgres",
		Passowrd: "root",
		DbName:   "db_quizy",
	}
	var err error
	// db, err := database.ConnectionGorm()
	// db.DropTableIfExists(&model.User{}, &model.Quiz{}, &model.Quest{}, &model.Player{}, &model.Option{}, &model.History{}, &model.Archievement{})
	// db.Debug().CreateTable(&model.User{}, &model.Quiz{}, &model.Quest{}, &model.Player{}, &model.Option{}, &model.History{}, &model.Archievement{})
	// db.Close()
	API.DB, err = database.Connection()
	if err != nil {
		log.Panic(err)
	}

	API.Run()
}
