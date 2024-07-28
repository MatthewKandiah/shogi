package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/MatthewKandiah/shogi/dao"
	"github.com/MatthewKandiah/shogi/view"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const DB_PATH = "./shogi.db"

const BCRYPT_STRENGTH = 10

const TIME_FORMAT = time.RFC822

func main() {
	dbFileAlreadyExisted := fileExists(DB_PATH)
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Fatal(err)
	}

	usersDao := dao.UsersDao{Db: db}
	passwordsDao := dao.PasswordsDao{Db: db}
	sessionsDao := dao.SessionsDao{Db: db}
	if !dbFileAlreadyExisted {
		daos := []dao.Dao{
			usersDao,
			passwordsDao,
			sessionsDao,
		}
		err := initialiseDb(daos)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler())
	http.HandleFunc("/home", homeHandler(usersDao))
	http.HandleFunc("/sign-up", signUpHandler(usersDao, passwordsDao))
	http.HandleFunc("/sign-in", signInHandler(usersDao, passwordsDao, sessionsDao))
	http.HandleFunc("/sign-out", signOutHandler(sessionsDao))

	// TODO - update to use TLS for https
	// TODO - extract port to env variable
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// TODO - if cookies for valid session exist, redirect to home
func indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle index")
		ctx := context.Background()
		err := view.IndexView().Render(ctx, w)
		if err != nil {
			log.Fatal("Error serving index page")
		}
	}
}

// TODO - redirect to sign in if you are not logged in
func homeHandler(usersDao dao.UsersDao, sessionsDao dao.SessionsDao) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle home")

		// TODO - pull out a requiresValidSession helper
		cookieSessionId := valueFromCookie("sessionId", r)
		cookieUserId := valueFromCookie("userId", r)
		sessionsRow, err := sessionsDao.Get(cookieSessionId)
		if err != nil {
			fmt.Println("Couldn't find session in DB")
			return
		}
		if sessionsRow.UserId != cookieUserId {
			fmt.Println("User id in cookie and DB do not match")
			return
		}
		validSession, err := hasValidSession(cookieUserId, cookieSessionId, sessionsRow.ExpiryTime, sessionsDao)
		if err != nil || !validSession {
			fmt.Println("Failed to validate session")
			return
		}
		fmt.Println("Valid session confirmed")

		ctx := context.Background()
		userIdCookie, err := r.Cookie("userId")
		var userNameString string
		if err != nil {
			userNameString = "ERROR - you aren't logged in?!"
		} else {
			usersRow, err := usersDao.Get(userIdCookie.Value)
			if err != nil {
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

func signUpHandler(usersDao dao.UsersDao, passwordsDao dao.PasswordsDao) http.HandlerFunc {
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
			encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_STRENGTH)
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

func signInHandler(usersDao dao.UsersDao, passwordsDao dao.PasswordsDao, sessionsDao dao.SessionsDao) http.HandlerFunc {
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
			expiryTime := time.Now().Add(sessionDuration).Format(TIME_FORMAT)
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

func signOutHandler(sessionsDao dao.SessionsDao) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handle sign out")
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			fmt.Println("Did not find a session cookie")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userIdCookie, err := r.Cookie("userId")
		if err != nil {
			fmt.Println("Did not find a userId")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := sessionsDao.Delete(userIdCookie.Value, sessionCookie.Value)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("Deleted %d rows\n", rowsAffected)
		updatedSessionCookie := http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		}
		updatedUserIdCookie := http.Cookie{
			Name:   "userId",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, &updatedSessionCookie)
		http.SetCookie(w, &updatedUserIdCookie)
		if rowsAffected > 1 {
			fmt.Printf("Unexpectedly deleted %d rows, expected 1\n", rowsAffected)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func initialiseDb(daos []dao.Dao) error {
	for _, dao := range daos {
		err := dao.Create()
		if err != nil {
			return err
		}
	}
	return nil
}

func hasValidSession(userId string, sessionId string, expiryTime string, sessionsDao dao.SessionsDao) (bool, error) {
	dbSessionIds, err := sessionsDao.GetAll(userId)
	if err != nil {
		return false, err
	}
	matchingIndex := slices.IndexFunc(dbSessionIds, func(sr dao.SessionsRow) bool {
		return sr.SessionId == sessionId
	})
	if matchingIndex == -1 {
		return false, nil
	}
	parsedExpiryTime, err := time.Parse(TIME_FORMAT, expiryTime)
	if err != nil {
		return false, nil
	}
	isValid := time.Now().Compare(parsedExpiryTime) == -1
	return isValid, nil
}

func valueFromCookie(key string, r *http.Request) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}
