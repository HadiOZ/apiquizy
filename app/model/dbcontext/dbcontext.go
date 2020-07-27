package dbcontext

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Auth struct {
	Username string
	Passowrd string
	DbName   string
}

func (auth *Auth) Connection() (*sql.DB, error) {
	var dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", auth.Username, auth.Passowrd, auth.DbName)
	db, err := sql.Open("postgres", dbInfo)
	return db, err
}
