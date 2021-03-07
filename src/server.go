package main

import (
	"net/http"
	"server/api"

	"github.com/labstack/echo"
)

func main() {
	users := make([]api.User, 0)
	sessions := make(api.SessionsMap, 0)
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

	api.Router(e)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

/*
curl -X POST http://localhost:1323/api/users/ \
-H 'Content-Type: application/json' \
-d '{"email": "user@example.com","password": "stringc","username": "string"}'
*/