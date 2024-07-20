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
	row := pc.Db.QueryRow("SELECT displayName FROM users WHERE id = ?", id)
	var (
		displayName string
	)
	err := row.Scan(&displayName)
	if err != nil {
		displayName = "Not found"
	}
	return profile.ProfilePage(id, displayName).Render(c.Request().Context(), c.Response())
}
