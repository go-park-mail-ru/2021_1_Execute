package main

import (
	"2021_1_Execute/database"
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/files"
	FilesHttpDelivery "2021_1_Execute/internal/files/delivery/http"
	SessionsDelivery "2021_1_Execute/internal/session/delivery"
	"2021_1_Execute/internal/session/repository/postgre_sessions_repository"
	"2021_1_Execute/internal/tasks/repository/postgre_task_repository"
	"log"

	BoardsHttpDelivery "2021_1_Execute/internal/boards_and_rows/delivery"
	"2021_1_Execute/internal/boards_and_rows/repository/postgre_board_repository"
	UserHttpDelivery "2021_1_Execute/internal/users/delivery/http"

	BoardUC "2021_1_Execute/internal/boards_and_rows/usecase"
	TasksHttpDelivery "2021_1_Execute/internal/tasks/delivery"
	TaskUC "2021_1_Execute/internal/tasks/usecase"
	"2021_1_Execute/internal/users/repository/postgre_user_repository"
	UserUC "2021_1_Execute/internal/users/usecase"
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
		echoErr := domain.GetEchoError(err)
		e.DefaultHTTPErrorHandler(echoErr, c)
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} err=${error}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	Router(e)

	fileUtil := files.NewFileUtil()

	sessionRepo := postgre_sessions_repository.NewPostgreSessionsRepository(pool)
	sessionHandler := SessionsDelivery.NewSessionHandler(sessionRepo)

	userRepo := postgre_user_repository.NewPostgreUserRepository(pool, fileUtil)
	userUC := UserUC.NewUserUsecase(userRepo)

	tasksRepo := postgre_task_repository.NewPostgreTaskRepository(pool)
	boardRepo := postgre_board_repository.NewPostgreBoardRepository(pool)

	taskUC := TaskUC.NewTasksUsecase(tasksRepo, boardRepo)
	boardUC := BoardUC.NewBoardsUsecase(boardRepo, userUC, taskUC)

	UserHttpDelivery.NewUserHandler(e, userUC, sessionHandler)
	FilesHttpDelivery.NewFilesHandler(e, userUC, fileUtil, sessionHandler)
	BoardsHttpDelivery.NewBoardsHandler(e, boardUC, sessionHandler, taskUC)
	TasksHttpDelivery.NewTasksHandler(e, sessionHandler, taskUC)

	e.Logger.Fatal(e.Start(fmt.Sprint(":", *serverPort)))
}

func Router(e *echo.Echo) {
	e.File("/api/", "docs/index.html")
}
