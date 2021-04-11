package main

import (
	"2021_1_Execute/database"
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/files"
	FilesHttpDelivery "2021_1_Execute/internal/files/delivery/http"
	SessionsDelivery "2021_1_Execute/internal/session/delivery"
	"2021_1_Execute/internal/session/repository/postgreSessionsRepository"
	"log"

	UserHttpDelivery "2021_1_Execute/internal/users/delivery/http"
	"2021_1_Execute/internal/users/repository/postgreUserRepository"
	"2021_1_Execute/internal/users/usecase"
	"flag"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var allowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000", "http://localhost:1323", "http://89.208.199.114:3000"}

func main() {
	clientPort := flag.Int("client-port", 3000, "")
	serverPort := flag.Int("server-port", 1323, "") //postgresql://username:password@host:port/dbname
	databaseUsername := flag.String("username", "postgres", "Input database username")
	databasePassword := flag.String("password", "123456", "Input database password")
	databaseHost := flag.String("db-host", "localhost", "Input database host")
	databasePort := flag.Int("db-port", 5432, "Input database port")
	databaseName := flag.String("db-name", "trello", "Input database name")
	flag.Parse()
	allowOrigins = append(allowOrigins, fmt.Sprint("http://89.208.199.114:", *clientPort))
	fmt.Println(allowOrigins)

	pool, err := database.GetPool(*databaseUsername, *databasePassword, *databaseName, *databaseHost, *databasePort)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	err = database.InitDatabase(pool)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		err = domain.GetEchoError(err)
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	Router(e)

	fileUtil := files.NewFileUtil()

	sessionRepo := postgreSessionsRepository.NewPostgreSessionsRepository(pool)
	sessionHandler := SessionsDelivery.NewSessionHandler(sessionRepo)

	userRepo := postgreUserRepository.NewPostgreUserRepository(pool, fileUtil)
	userUC := usecase.NewUserUsecase(userRepo)

	UserHttpDelivery.NewUserHandler(e, userUC, sessionHandler)
	FilesHttpDelivery.NewFilesHandler(e, userUC, fileUtil, sessionHandler)
	e.Logger.Fatal(e.Start(fmt.Sprint(":", *serverPort)))
}

func Router(e *echo.Echo) {
	e.File("/api/", "docs/index.html")
}
