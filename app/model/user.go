package model

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	UserID         string         `json:"userID"`
	Email          string         `json:"email"`
	Password       string         `json:"password"`
	Name           string         `json:"name"`
	Gender         sql.NullString `json:"gender"`
	Country        sql.NullString `json:"country"`
	BirthDate      sql.NullString `json:"birthdate"`
	Job            sql.NullString `json:"job"`
	Institution    sql.NullString `json:"institution"`
	Phone          sql.NullString `json:"phone"`
	ProfilePicture sql.NullString `json:"profile"`
}

//function active
func (u *User) CreateUser(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO public.users(user_id, email, password, name, birth_date) VALUES ('%s', '%s', '%s', '%s', '%s');`, u.UserID, u.Email, u.Password, u.Name, u.BirthDate.String)
	res, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return affect, nil
}

//function active
func (u *User) CheckUSer(db *sql.DB) (result User, err error) {
	query := fmt.Sprintf(`select user_id, password from public.users where email = '%s'`, u.Email)
	row := db.QueryRow(query)
	if err := row.Scan(&result.UserID, &result.Password); err != nil {
		return User{}, err
	}
	return result, nil
}

//function clear
func (u *User) EditProfile(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public.users SET name= '%s', gender= '%s', country= '%s', job= '%s', phone= '%s', institution= '%s', birth_date= '%s' WHERE user_id= '%s'`, u.Name, u.Gender.String, u.Country.String, u.Job.String, u.Phone.String, u.Institution.String, u.BirthDate.String, u.UserID)
	log.Println(query)
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

//function active
func (u *User) EditProfilePicture(db *sql.DB) (int64, error) {
	query := fmt.Sprintf(`UPDATE public.users SET profile_picture= '%s' WHERE user_id= '%s'`, u.ProfilePicture.String, u.UserID)
	res, err := db.Exec(query)
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

//function active
func SelectUserProfile(iduser string, db *sql.DB) (result User, err error) {
	query := fmt.Sprintf(`SELECT name, gender, country, job, institution, phone, birth_date, profile_picture FROM public.users WHERE user_id= '%s'`, iduser)
	row := db.QueryRow(query)
	if err := row.Scan(&result.Name, &result.Gender, &result.Country, &result.Job, &result.Institution, &result.Phone, &result.BirthDate, &result.ProfilePicture); err != nil {
		return User{}, err
	}
	return result, nil
}
