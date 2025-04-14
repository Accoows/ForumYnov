package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "database/forum.db")
	if err != nil {
		log.Fatal("Erreur ouverture base :", err)
	}
	defer db.Close()

	printTableNames(db)
	printTableContent(db, "Users")
	printTableContent(db, "Categories")
	printTableContent(db, "Posts")
	printTableContent(db, "Comments")
	printTableContent(db, "Likes_Dislikes")
}

func printTableNames(db *sql.DB) {
	fmt.Println("\n📋 Tables présentes dans forum.db :")
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Println("Erreur lecture tables :", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println("→", name)
	}
}

func printTableContent(db *sql.DB, table string) {
	fmt.Printf("\n📦 Contenu de %s :\n", table)
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		log.Printf("Impossible de lire la table %s : %v\n", table, err)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println("Erreur lecture colonnes :", err)
		return
	}

	values := make([]interface{}, len(cols))
	for i := range values {
		var val interface{}
		values[i] = &val
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			log.Println("Erreur lecture ligne :", err)
			continue
		}
		for i, col := range cols {
			fmt.Printf("%s = %v | ", col, *(values[i].(*interface{})))
		}
		fmt.Println()
		count++
	}

	if count == 0 {
		fmt.Println("Aucune ligne.")
	}
}
