package model

import (
	"database/sql"
	"fmt"
)

type User struct {
	IDUser   string
	Email    string
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
	ProfilePicture string
}

func CreateUser(user User, db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public."User"(id_user, email, password, name, birth_date) VALUES ('%s', '%s', '%s', '%s', '%s')`, user.IDUser, user.Email, user.Password, user.Detail.Name, user.Detail.BirthDate)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

func EditUser() {}

func EditProfilePicture(idUser string, path string, db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public."User" SET profile_picture= '%s' WHERE id_user= '%s'`, path, idUser)
	res, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

func SelectUserProfile() {}

func CheckUSer(email string, db *sql.DB) (User, error) {
	query := fmt.Sprintf(`select id_user, password from public."User" where email = '%s'`, email)
	row, err := db.Query(query)
	if err != nil {
		return User{}, err
	}
	var result User
	for row.Next() {
		if err := row.Scan(&result.IDUser, &result.Password); err != nil {
			return User{}, err
		}
	}
	return result, nil
}
