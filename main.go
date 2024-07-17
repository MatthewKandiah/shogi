package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	indexComponent := index()
	signInComponent := signIn()

	http.Handle("/", templ.Handler(indexComponent))
	http.Handle("/sign-in", templ.Handler(signInComponent))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
