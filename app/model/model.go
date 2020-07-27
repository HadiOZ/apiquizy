package model 

import (
	"database/sql"
)

type Calon struct {
	IDCalon     string `json:"id"`
	Nama        string `json:"name"`
	Partai      string `json:"party"`
	NoUrut      int    `json:"number"`
	JenisPemilu string `json:"candidate"`
	Region      string `json:"region"`
	FotoAddress string `json:"photo"`
}

func SelectAllCalon(db *sql.DB) []Calon {
	defer db.Close()
	query := "select id_calon, nama, partai, no_urut, region, jenis_pemilu, foto_address from calon"
	row, _ := db.Query(query)
	var result []Calon
	for row.Next() {
		each := Calon{}
		row.Scan(&each.IDCalon, &each.Nama, &each.Partai, &each.NoUrut, &each.Region, &each.JenisPemilu, &each.FotoAddress)
		result = append(result, each)
	}
	return result
}
