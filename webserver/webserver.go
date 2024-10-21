package webserver

import (
	"bytes"
	"dsc/config"
	"dsc/database"
	"dsc/lib"
	"fmt"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io/fs"
	"time"
)

func LogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			log.Info().
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", res.Status).
				Dur("duration", time.Since(start)).
				Msg("Request handled")

			return err
		}
	}
}

func CreateAndListen(appConfig config.AppConfig, db *database.Database, staticFiles fs.FS) {
	e := echo.New()
	e.HideBanner = true
	e.Logger = &lib.EchoLogger{ZeroLog: log.Logger}
	e.Use(LogMiddleware())
	e.StaticFS("/public", echo.MustSubFS(staticFiles, ""))
	e.GET("/", CreateHomepageHandler(db))

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", appConfig.WebServerAddress, appConfig.WebserverPort)))
}

func RenderView(c echo.Context, cmp templ.Component) error {
	buf := new(bytes.Buffer)
	err := cmp.Render(c.Request().Context(), buf)
	if err != nil {
		return err
	}
	return c.HTML(200, buf.String())
}
