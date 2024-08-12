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
	gamesDao := dao.GamesDao{Db: db}
	if !dbFileAlreadyExisted {
		daos := []dao.Dao{
			usersDao,
			passwordsDao,
			sessionsDao,
			gamesDao,
		}
		err := util.InitialiseDb(daos)
		if err != nil {
			log.Fatal(err)
		}
	}

	// TODO - put build artefacts in the static dir to simplify this
	//		not sure yet if I just want to copy-paste the js files into there, or if I want a "build" script that copies the current version into there and auto-increment a version
	//		moving zig build artefacts there should be as easy as setting a build destination somewhere in the build config!
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/zig-out/", http.StripPrefix("/zig-out/", http.FileServer(http.Dir("zig-out"))))

	http.HandleFunc("/", handler.IndexHandler())
	http.HandleFunc("/home", handler.HomeHandler(usersDao, sessionsDao))
	http.HandleFunc("/sign-up", handler.SignUpHandler(usersDao, passwordsDao))
	http.HandleFunc("/sign-in", handler.SignInHandler(usersDao, passwordsDao, sessionsDao))
	http.HandleFunc("/sign-out", handler.SignOutHandler(sessionsDao))

	// TODO - update to use TLS for https
	// TODO - extract port to env variable
	log.Fatal(http.ListenAndServe(":3000", nil))
}
