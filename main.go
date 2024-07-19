package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./test.db"

func main() {
	dbFileAlreadyExisted := fileExists(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if !dbFileAlreadyExisted {
		log.Println("Initialising fresh DB")
		err = initialiseDb(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	// rows, err := db.Query("SELECT id, username FROM users WHERE id = ?", 1)
	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var username string
		err := rows.Scan(&id, &username)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, username)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", templ.Handler(index()))
	http.Handle("/sign-in", templ.Handler(signIn()))
	http.Handle("/sign-up", templ.Handler(signUp()))
	http.Handle("/home", templ.Handler(home()))
	http.Handle("/game", templ.Handler(game()))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	return true
}

func initialiseDb(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE users (id text, username text);")
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (id, username) VALUES ('1', 'FirstUser')")
	if err != nil {
		return err
	}
	log.Println("FirstUser added")

	_, err = db.Exec("INSERT INTO users (id, username) VALUES ('2', 'SecondUser')")
	if err != nil {
		return err
	}
	log.Println("SecondUser added")

	return nil
}
