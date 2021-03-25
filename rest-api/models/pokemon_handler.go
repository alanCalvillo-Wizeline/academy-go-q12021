package models

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func SavePokemon(id int, name string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	db, err := sql.Open("mysql", viper.GetString("app.dbuser")+":"+viper.GetString("app.dbpass")+"@tcp("+viper.GetString("app.dbhost")+":"+viper.GetString("app.dbport")+")/"+viper.GetString("app.dbname"))

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
