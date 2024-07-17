package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.Handle("/", templ.Handler(index()))
	http.Handle("/sign-in", templ.Handler(signIn()))
	http.Handle("/sign-up", templ.Handler(signUp()))
	http.Handle("/home", templ.Handler(home()))
	http.Handle("/game", templ.Handler(game()))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
