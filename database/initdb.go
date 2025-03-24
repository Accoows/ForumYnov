package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Connexion et initialisation DB

var DB *sql.DB // Variable globale accessible partout

func InitTempDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/temp_forum.db")
	if err != nil {
		log.Fatal("Erreur ouverture DB :", err)
	}

	schema, err := os.ReadFile("database/temp_schema.sql")
	if err != nil {
		log.Fatal("Erreur lecture schema SQL :", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal("Erreur exécution schema :", err)
	}
}
