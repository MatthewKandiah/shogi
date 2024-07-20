package controller

import (
	"database/sql"

	"github.com/MatthewKandiah/shogi/view/profile"
	"github.com/labstack/echo/v4"
)

type ProfileController struct {
	Db *sql.DB
}

func(pc ProfileController) HandleShow(c echo.Context) error {
	id := c.Param("id")
	row := pc.Db.QueryRow("SELECT userName FROM users WHERE id = ?", id)
	var (
		userName string
	)
	err := row.Scan(&userName)
	if err != nil {
		userName = "Not found"
	}
	return profile.ProfilePage(id, userName).Render(c.Request().Context(), c.Response())
}
