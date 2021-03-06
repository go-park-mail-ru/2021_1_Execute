package main

import (
	"net/http"
	"server/api"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	api.Router(e)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
