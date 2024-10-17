package webserver

import (
	"dsc/database"
	"github.com/labstack/echo/v4"
)

func CreateHomepageHandler(db *database.Database) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.HTML(200, "<b>docker stats collector</b>")
	}
}
