package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/MatthewKandiah/shogi/controller"
	_ "github.com/mattn/go-sqlite3"

	"github.com/labstack/echo/v4"
)

const dbPath = "./shogi.db"

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

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	profileController := controller.ProfileController{Db: db}
	e.GET("/profile/:id", profileController.HandleShow)

	signInController := controller.SignInController{Db: db}
	e.GET("sign-in", signInController.HandleSignIn)
	e.POST("sign-in/authenticate", signInController.HandleAuthentication)

	e.Logger.Fatal(e.Start(":3000"))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func initialiseDb(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE users (id TEXT, userName TEXT, password TEXT);")
	if err != nil {
		return err
	}

	// temporary test data
	_, err = db.Exec("INSERT INTO users (id, userName, password) VALUES ('1', 'Matthew', 'password')")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (id, userName, password) VALUES ('2', 'Thomas', 'pw')")
	if err != nil {
		return err
	}

	return nil
}
