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
	Username string
	Passowrd string
	DbName   string
}

func (auth *Auth) ConnectionGorm() (*gorm.DB, error) {
	var dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", auth.Username, auth.Passowrd, auth.DbName)
	db, err := gorm.Open("postgres", dbInfo)
	return db, err
}

func (i *Auth) Connection() (*sql.DB, error) {
	var dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", i.Username, i.Passowrd, i.DbName)
	db, err := sql.Open("postgres", dbInfo)
	return db, err
}
