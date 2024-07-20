package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	signin "github.com/MatthewKandiah/shogi/view/signIn"
	"github.com/labstack/echo/v4"
)

type SignInController struct {
	Db *sql.DB
}

func (sc SignInController) HandleSignIn(c echo.Context) error {
	return signin.ProfilePage().Render(c.Request().Context(), c.Response())
}

func (sc SignInController) HandleAuthentication(c echo.Context) error {
	userName := c.FormValue("userName")
	password := c.FormValue("password")
	fmt.Printf("user trying to sign in = %s\n", userName)
	fmt.Printf("password they are using = %s\n", password)

	row := sc.Db.QueryRow("SELECT password FROM users WHERE userName = ?", userName)
	var (
		dbPassword string
	)
	err := row.Scan(&dbPassword)
	if err != nil {
		return err
	}

	if password == dbPassword {
		fmt.Printf("Successful authentication UN=%s PW=%s DPW=%s\n", userName, password, dbPassword)
	} else {
		fmt.Printf("Failed authentication UN=%s PW=%s DPW=%s\n", userName, password, dbPassword)
	}
	return c.HTML(http.StatusNoContent, "")
}
