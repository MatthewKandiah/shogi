package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MatthewKandiah/shogi/constant"
	"github.com/MatthewKandiah/shogi/dao"
	"github.com/MatthewKandiah/shogi/view"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignInHandler(usersDao dao.UsersDao, passwordsDao dao.PasswordsDao, sessionsDao dao.SessionsDao) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle sign in")
		ctx := context.Background()
		var viewToSend int
		if r.Method == http.MethodGet {
			fmt.Println("handle GET")
			viewToSend = view.SignInPageView
			goto sendView
		} else if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				viewToSend = view.SignInFormSnippetView
				goto sendView
			}
			userName := r.Form.Get("userName")
			password := r.Form.Get("password")
			if userName == "" || password == "" {
				fmt.Println("Bad request - missing username or password")
				w.WriteHeader(http.StatusBadRequest)
				viewToSend = view.SignInFormSnippetView
				goto sendView
			}
			userRow, err := usersDao.GetByUserName(userName)
			if err != nil {
				fmt.Println("Failed to find user")
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				viewToSend = view.SignInFormSnippetView
				goto sendView
			}
			userId := userRow.Id
			passwordRow, err := passwordsDao.Get(userId)
			if err != nil {
				fmt.Println("Failed to find password")
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				viewToSend = view.SignInFormSnippetView
				goto sendView
			}
			if bcrypt.CompareHashAndPassword([]byte(passwordRow.Password), []byte(password)) != nil {
				fmt.Println("Passwords don't match")
				w.WriteHeader(http.StatusUnauthorized)
				viewToSend = view.SignInFormSnippetView
				goto sendView
			}
			fmt.Println("Passwords matched")
			sessionId := uuid.New()
			sessionDuration := 7 * 24 * 60 * 60 * time.Second // a week
			expiryTime := time.Now().Add(sessionDuration).Format(constant.TIME_FORMAT)
			fmt.Printf("expiry time calculated: %s\n", expiryTime)
			sessionsRow := dao.SessionsRow{UserId: userId, SessionId: sessionId.String(), ExpiryTime: expiryTime}
			err = sessionsDao.Insert(sessionsRow)
			if err != nil {
				fmt.Println("Failed to insert into sessions table")
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				viewToSend = view.SignInFormSnippetView
				goto sendView
			}
			fmt.Println("Succesfully inserted session")
			sessionCookie := http.Cookie{Name: "session", Value: sessionId.String()}
			userIdCookie := http.Cookie{Name: "userId", Value: userId}
			http.SetCookie(w, &sessionCookie)
			http.SetCookie(w, &userIdCookie)
			viewToSend = view.SignInSuccessSnippetView
			goto sendView
		} else {
			log.Fatal("Unexpected http method - " + r.Method)
		}
		log.Fatal("Sign in failed to send a view")
	sendView:
		switch viewToSend {
		case view.SignInPageView:
			err := view.SignInPage().Render(ctx, w)
			if err != nil {
				log.Fatal(err)
			}
			return
		case view.SignInFormSnippetView:
			err := view.SignInFormSnippet().Render(ctx, w)
			if err != nil {
				log.Fatal(err)
			}
			return
		case view.SignInSuccessSnippetView:
			err := view.SignInSuccessSnippet().Render(ctx, w)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		log.Fatal("Sign in failed to select a view")
	}
}
