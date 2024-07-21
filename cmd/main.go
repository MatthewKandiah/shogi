package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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
		err = initialiseDb(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	// TODO - write a POST wrapper
	http.HandleFunc("/register", registerUserHandler(db))

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
func registerUserHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			return
		}

		userName := r.Form.Get("username")
		password := r.Form.Get("password")
		// TODO - strip trailing whitespace from username to avoid blank username
		if userName == "" || password == "" {
			fmt.Println("Failed to register user because name or password not provided")
			return
		}
		// TODO - wrap in transaction
		row := db.QueryRow("SELECT username, id FROM users WHERE username = ?", userName)
		var existingUserId string
		var existingUsername string
		err = row.Scan(&existingUserId, &existingUsername)
		if err == nil {
			fmt.Println("Failed to register, username already taken")
			w.WriteHeader(http.StatusConflict)
			return
		}
		if err != sql.ErrNoRows{
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userId := uuid.New()
		fmt.Printf("new userId generated: %s\n", userId)
		_, err = db.Exec("INSERT INTO users (id, userName) VALUES (?, ?)", userId, userName)
		if err != nil {
			fmt.Println("Failed to insert user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully inserted user")
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_STRENGTH)
		if err != nil {
			fmt.Println("Failed to encrypt password")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully encrypted password")
		_, err = db.Exec("INSERT INTO passwords (userId, password) VALUES (?, ?)", userId, encryptedPassword)
		if err != nil {
			fmt.Println("Failed to insert password")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully inserted password")
		w.WriteHeader(http.StatusOK)
	}
}

func initialiseDb(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE users (id TEXT, userName TEXT);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE sessions (userId TEXT, sessionId TEXT, expiryTime TEXT);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE passwords (userId TEXT, password TEXT);")
	if err != nil {
		return err
	}

	return nil
}
