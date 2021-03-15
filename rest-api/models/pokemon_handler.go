package models

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func SavePokemon(id int, name string) {
	db, err := sql.Open("mysql", "poke:pokepassword@tcp(127.0.0.1:3306)/pokedex")

	if err != nil {
		panic(err.Error())
	}
	id_str := strconv.Itoa(id)
	insertQuery := "INSERT INTO pokemons (id,name) VALUES(" + id_str + ",'" + name + "') ON DUPLICATE KEY UPDATE name='" + name + "' "
	insert, err := db.Query(insertQuery)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

}
