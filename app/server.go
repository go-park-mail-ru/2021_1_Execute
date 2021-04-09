package main

import (
	"2021_1_Execute/internal/files"
	FilesHttpDelivery "2021_1_Execute/internal/files/delivery/http"
	SessionsDelivery "2021_1_Execute/internal/session/delivery"
	"2021_1_Execute/internal/session/repository/postgreSessionsRepository"
	"io/ioutil"
	"log"

	UserHttpDelivery "2021_1_Execute/internal/users/delivery/http"
	"2021_1_Execute/internal/users/repository/postgreUserRepository"
	"2021_1_Execute/internal/users/usecase"
	"context"
	"flag"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
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

	DBUri := "postgresql://postgres:123456@localhost:5432/trello"
	config, err := pgxpool.ParseConfig(DBUri)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	initFile, err := ioutil.ReadFile("database/trello.sql")
	if err != nil {
		log.Fatal(err)
	}
	initComands := string(initFile)

	_, err = pool.Exec(ctx, initComands)

	e := echo.New()

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
