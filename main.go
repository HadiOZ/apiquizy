package main

import (
	"apiquizyfull/app"
	"apiquizyfull/app/model/dbcontext"
	"fmt"
	"os"
)

func main() {
	var API app.App
	API.Port = fmt.Sprintf(":%s", os.Getenv("PORT"))
	API.Assets = dbcontext.Assets{
		Profile:  os.Getenv("PROFILE"),
		Quiz:     os.Getenv("QUIZ"),
		Question: os.Getenv("QUESTION"),
	}
	API.DBContext = dbcontext.Auth{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Passowrd: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
	}
	API.CreateDataBase(os.Getenv("INIT_DB"))
	API.Run()
}
