package main

import (
	"2021_1_Execute/src/api"
	"flag"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var allowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000", "http://localhost:1323", "http://89.208.199.114:3000"}

func main() {
	clintPort := flag.Int("-clint-port", 3000, "")
	serverPort := flag.Int("-server-port", 1323, "")
	flag.Parse()
	allowOrigins = append(allowOrigins, fmt.Sprint("http://89.208.199.114:", *clintPort))
	fmt.Println(allowOrigins)

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
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	api.Router(e)
	e.Logger.Fatal(e.Start(fmt.Sprint(":", *serverPort)))
}
