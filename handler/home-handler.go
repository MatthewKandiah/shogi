package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MatthewKandiah/shogi/dao"
	"github.com/MatthewKandiah/shogi/util"
	"github.com/MatthewKandiah/shogi/view"
)

func HomeHandler(usersDao dao.UsersDao, sessionsDao dao.SessionsDao) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle home")

		// TODO - pull out a requiresValidSession helper
		cookieSessionId := util.ValueFromCookie("session", r)
		cookieUserId := util.ValueFromCookie("userId", r)
		sessionsRow, err := sessionsDao.Get(cookieSessionId)
		if err != nil {
			fmt.Println("Couldn't find session in DB")
			w.Header().Set("Location", "/sign-in")
			w.WriteHeader(http.StatusSeeOther)
			return
		}
		if sessionsRow.UserId != cookieUserId {
			fmt.Println("User id in cookie and DB do not match")
			w.Header().Set("Location", "/sign-in")
			w.WriteHeader(http.StatusSeeOther)
			return
		}
		validSession, err := util.HasValidSession(cookieUserId, cookieSessionId, sessionsRow.ExpiryTime, sessionsDao)
		if err != nil || !validSession {
			fmt.Println("Failed to validate session")
			w.Header().Set("Location", "/sign-in")
			w.WriteHeader(http.StatusSeeOther)
			return
		}
		fmt.Println("Valid session confirmed")

		ctx := context.Background()
		userIdCookie, err := r.Cookie("userId")
		var userNameString string
		if err != nil {
			// TODO - these should probably be switched to asserts
			userNameString = "ERROR - you aren't logged in?!"
		} else {
			usersRow, err := usersDao.Get(userIdCookie.Value)
			if err != nil {
				// TODO - these should probably be switched to asserts
				userNameString = "ERROR = you aren't in the users table?!"
			} else {
				userNameString = usersRow.UserName
			}
		}
		err = view.HomeView(userNameString).Render(ctx, w)
		if err != nil {
			log.Fatal("Error serving home page")
		}
	}
}
