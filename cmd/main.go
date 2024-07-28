package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/MatthewKandiah/shogi/constant"
	"github.com/MatthewKandiah/shogi/dao"
	"github.com/MatthewKandiah/shogi/handler"
	"github.com/MatthewKandiah/shogi/util"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbFileAlreadyExisted := util.FileExists(constant.DB_PATH)
	db, err := sql.Open("sqlite3", constant.DB_PATH)
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
		err := util.InitialiseDb(daos)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler.IndexHandler())
	http.HandleFunc("/home", handler.HomeHandler(usersDao, sessionsDao))
	http.HandleFunc("/sign-up", handler.SignUpHandler(usersDao, passwordsDao))
	http.HandleFunc("/sign-in", handler.SignInHandler(usersDao, passwordsDao, sessionsDao))
	http.HandleFunc("/sign-out", handler.SignOutHandler(sessionsDao))

	// TODO - update to use TLS for https
	// TODO - extract port to env variable
	log.Fatal(http.ListenAndServe(":3000", nil))
}
