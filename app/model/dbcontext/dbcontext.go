package dbcontext

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Assets struct {
	Profile  string
	Quiz     string
	Question string
}

type Auth struct {
	Host     string
	Port     string
	Username string
	Passowrd string
	DbName   string
}

func (auth *Auth) ConnectionGorm() (*gorm.DB, error) {
	var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", auth.Host, auth.Port, auth.Username, auth.Passowrd, auth.DbName)
	db, err := gorm.Open("postgres", dbInfo)
	return db, err
}

func (auth *Auth) Connection() (*sql.DB, error) {
	var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", auth.Host, auth.Port, auth.Username, auth.Passowrd, auth.DbName)
	db, err := sql.Open("postgres", dbInfo)
	return db, err
}
