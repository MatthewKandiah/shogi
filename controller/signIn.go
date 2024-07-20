package controller

import (
	"database/sql"
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

	row := sc.Db.QueryRow("SELECT password FROM users WHERE userName = ?", userName)
	var (
		dbPassword string
	)
	err := row.Scan(&dbPassword)
	if err != nil {
		return c.HTML(http.StatusOK, "<h1>Failed! User doesn't exist</h1><br><a href=\"/sign-in\">Try again</a>")
	}

	if password == dbPassword {
		return c.HTML(http.StatusOK, "<h1>Success!</h1>")
	} else {
		return c.HTML(http.StatusOK, "<h1>Failed! Incorrect Password</h1><br><a href=\"/sign-in\">Try again</a>")
	}
}
