package main

import (
	"apiquizyfull/app"
)

func main() {
	var Api app.App
	Api.Port = ":8080"
	Api.SetAssetPath("profile", "hadiese/home/Documents/quizymedia/profile-picture")
	Api.SetAssetPath("picture", "hadiese/home/Documents/quizymedia/quiz-picture")
	Api.SetAssetPath("media", "hadiese/home/Documents/quizymedia/question-media")
	Api.Run()
}
