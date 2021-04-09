package main

import (
	//FilesHttpDelivery "2021_1_Execute/internal/files/delivery/http"
	//"2021_1_Execute/internal/files"
	//SessionsHttpDelivery "2021_1_Execute/internal/session/delivery"

	//UserHttpDelivery "2021_1_Execute/internal/users/delivery/http"
	"flag"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var allowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000", "http://localhost:1323", "http://89.208.199.114:3000"}

func main() {
	clientPort := flag.Int("client-port", 3000, "")
	serverPort := flag.Int("server-port", 1323, "")
	flag.Parse()
	allowOrigins = append(allowOrigins, fmt.Sprint("http://89.208.199.114:", *clientPort))
	fmt.Println(allowOrigins)

	//todo настроить подключение к бд dbConn := ...

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	Router(e)

	//todo инициализировать userRepo и sessionsRepo

	//userUC := usecase.NewUserUsecase()
	//fileUtil := files.NewFileUtil()

	//UserHttpDelivery.NewUserHandler(e,)
	//SessionsHttpDelivery.NewSessionHandler(e,)
	//FilesHttpDelivery.NewFilesHandler(e,)
	e.Logger.Fatal(e.Start(fmt.Sprint(":", *serverPort)))
}

func Router(e *echo.Echo) {
	e.File("/api/", "docs/index.html")
}
