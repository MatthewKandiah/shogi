package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MatthewKandiah/shogi/view"
)

// TODO - if cookies for valid session exist, redirect to home
func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle index")
		ctx := context.Background()
		err := view.IndexView().Render(ctx, w)
		if err != nil {
			log.Fatal("Error serving index page")
		}
	}
}
