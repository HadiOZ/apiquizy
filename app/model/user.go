package model

type User struct {
	IDUser   string
	Username string
	Password string
	Detail   Profile
}

type Profile struct {
	Name           string
	Gender         string
	BirthDate      string
	Job            string
	Institution    string
	Phone          string
	Email          string
	ProfilePicture string
}
