package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MatthewKandiah/shogi/constant"
	"github.com/MatthewKandiah/shogi/dao"
	"github.com/MatthewKandiah/shogi/view"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(usersDao dao.UsersDao, passwordsDao dao.PasswordsDao) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle sign up")
		ctx := context.Background()
		if r.Method == http.MethodGet {
			fmt.Println("handle GET")
			err := view.SignUpPage().Render(ctx, w)
			if err != nil {
				log.Fatal(err)
			}
		} else if r.Method == http.MethodPost {
			fmt.Println("handle POST")
			err := r.ParseForm()
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			userName := strings.TrimSpace(r.Form.Get("userName"))
			password := r.Form.Get("password")
			fmt.Printf("un: %s, pw: %s\n", userName, password)
			if userName == "" || password == "" {
				fmt.Println("bad/empty userName/password")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// TODO - wrap in transaction
			_, err = usersDao.GetByUserName(userName)
			if err == nil {
				fmt.Println("Failed to register, username already taken")
				w.WriteHeader(http.StatusConflict)
				return
			}
			if err != sql.ErrNoRows {
				fmt.Println("Unexpected error getting user")
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			newRow := dao.UsersRow{Id: uuid.New().String(), UserName: userName}
			err = usersDao.Insert(newRow)
			if err != nil {
				fmt.Println("Failed to insert user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Println("Successfully inserted user")
			// TODO - check how bcrypt works, is this actually sufficiently secure to deploy?
			encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), constant.BCRYPT_STRENGTH)
			if err != nil {
				fmt.Println("Failed to encrypt password")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Println("Successfully encrypted password")
			passwordsRow := dao.PasswordsRow{UserId: newRow.Id, Password: string(encryptedPassword)}
			err = passwordsDao.Insert(passwordsRow)
			if err != nil {
				fmt.Println("Failed to insert password")
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Println("Successfully inserted password")
			err = view.SignUpSuccessSnippet().Render(ctx, w)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
