package handler

import (
	"fmt"
	"net/http"

	"github.com/MatthewKandiah/shogi/dao"
)

func SignOutHandler(sessionsDao dao.SessionsDao) http.HandlerFunc {
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
		w.Header().Set("Location", "/sign-in")
		w.WriteHeader(http.StatusSeeOther)
	}
}
