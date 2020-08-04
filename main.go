package main

import (
	"apiquizyfull/app"
	"apiquizyfull/app/model/dbcontext"
)

func main() {
	var API app.App
	API.Port = ":8000"
	API.Assets = dbcontext.Assets{
		Profile:  "/home/hadiese/Documents/quizymedia/profile-picture",
		Quiz:     "/home/hadiese/Documents/quizymedia/quiz-picture",
		Question: "/home/hadiese/Documents/quizymedia/question-media",
	}
	API.DBContext = dbcontext.Auth{
		Username: "postgres",
		Passowrd: "root",
		DbName:   "db_quizy",
	}
	API.CreateDataBase(false)
	API.Run()
}
