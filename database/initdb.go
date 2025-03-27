package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var SQL *sql.DB

func Database() {
	var err error
	SQL, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Erreur ouverture base :", err)
		return
	}

	if err = SQL.Ping(); err != nil {
		fmt.Println("Erreur connexion base :", err)
		SQL = nil
		return
	}

	fmt.Println("Connexion à SQLite OK")
}
