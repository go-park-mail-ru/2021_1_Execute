package main

import (
	"2021_1_Execute/src/api"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var allowOrigins = []string{"http://localhost:3000"}

func main() {
	users := make([]api.User, 0)
	sessions := make(api.Sessions, 0)
	e := echo.New()
	//This middleware should be registered before any other middleware.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &api.Database{
				Context:  c,
				Users:    &users,
				Sessions: &sessions,
			}
			return next(cc)
		}
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderAccept, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderCookie, echo.HeaderSetCookie},
		AllowCredentials: true,
	}))

	api.Router(e)
	e.Logger.Fatal(e.Start(":1323"))
}
