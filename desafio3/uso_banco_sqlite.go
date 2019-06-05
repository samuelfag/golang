package main 

import (
	"fmt"
	"log"
	//"encoding/json"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Dados struct {
	Actor string `json:"actor"`
	Quote string `json:"quote"`
}

func main() {
	retorno := conectaDB("SELECT actor, detail as quote FROM scripts WHERE actor IS NOT NULL AND detail IS NOT NULL LIMIT 10")

	fmt.Printf("%T", retorno)
}

func checkError(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func conectaDB(query string) *sql.Rows {
	var db, err = sql.Open("sqlite3", "database.sqlite")
	rows, err := db.Query(query)

	checkError("Erro: ", err)
	return rows
}


