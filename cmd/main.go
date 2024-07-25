package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MatthewKandiah/shogi/dao"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const dbPath = "./shogi.db"

const BCRYPT_STRENGTH = 10

func main() {
	dbFileAlreadyExisted := fileExists(dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if !dbFileAlreadyExisted {
		// TODO pass in daos instead of db
		err := initialiseDb(db)
		if err != nil {
			log.Fatal(err)
		}
	}
	usersDao := dao.UsersDao{Db: db}
	passwordsDao := dao.PasswordsDao{Db: db}

	// TODO - write a GET/POST wrapper
	http.HandleFunc("/register", registerUserHandler(usersDao, passwordsDao))
	http.HandleFunc("/sign-in", signInHandler(db, usersDao, passwordsDao))
	http.HandleFunc("/sign-out", signOutHandler(db))

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

// TODO - wrapper to pass a db in and return a HttpHandler function
func registerUserHandler(usersDao dao.UsersDao, passwordsDao dao.PasswordsDao) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle register user")
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userName := r.Form.Get("username")
		password := r.Form.Get("password")
		// TODO - strip trailing whitespace from username to avoid blank username
		if userName == "" || password == "" {
			fmt.Println("Failed to register user because name or password not provided")
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully inserted password")
		w.WriteHeader(http.StatusOK)
	}
}

func signInHandler(db *sql.DB, usersDao dao.UsersDao, passwordsDao dao.PasswordsDao) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle sign in")
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userName := r.Form.Get("username")
		password := r.Form.Get("password")
		userRow, err := usersDao.GetByUserName(userName)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userId := userRow.Id
		passwordRow, err := passwordsDao.Get(userId)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(passwordRow.Password), []byte(password)) != nil {
			fmt.Println("Passwords don't match")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println("Passwords matched")
		sessionId := uuid.New()
		sessionDuration := 7 * 24 * 60 * 60 * time.Second // a week
		expiryTime := time.Now().Add(sessionDuration).Format(time.RFC822)
		fmt.Printf("expiry time calculated: %s\n", expiryTime)
		_, err = db.Exec("INSERT INTO sessions (userId, sessionId, expiryTime) VALUES (?, ?, ?)", userId, sessionId, expiryTime)
		fmt.Println("Succesfully inserted session")
		sessionCookie := http.Cookie{Name: "session", Value: sessionId.String()}
		userIdCookie := http.Cookie{Name: "userId", Value: userId}
		http.SetCookie(w, &sessionCookie)
		http.SetCookie(w, &userIdCookie)
		w.WriteHeader(http.StatusOK)
	}
}

func signOutHandler(db *sql.DB) http.HandlerFunc {
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
		res, err := db.Exec("DELETE FROM sessions WHERE userId = ? AND sessionId = ?", userIdCookie.Value, sessionCookie.Value)
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

func initialiseDb(db *sql.DB) error {
	usersDao := dao.UsersDao{Db: db}
	err := usersDao.Create()
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE sessions (userId TEXT, sessionId TEXT, expiryTime TEXT);")
	if err != nil {
		return err
	}

	passwordsDao := dao.PasswordsDao{Db: db}
	err = passwordsDao.Create()
	if err != nil {
		return err
	}

	return nil
}
