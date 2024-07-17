package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	http.Handle("/", templ.Handler(index()))
	http.Handle("/sign-in", templ.Handler(signIn()))
	http.Handle("/sign-up", templ.Handler(signUp()))
	http.Handle("/home", templ.Handler(home()))
	http.Handle("game", templ.Handler(game()))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
