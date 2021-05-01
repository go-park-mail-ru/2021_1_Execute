package database

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

const PathToInitDBFile = "database/trello.sql"
const PathToDropDBFile = "database/drop.sql"

func GetPool(username, password, dbname, host string, port int) (*pgxpool.Pool, error) {
	DBUri := "postgresql://" + username + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname
	config, err := pgxpool.ParseConfig(DBUri)
	if err != nil {
		return nil, err
	}

	logger := &logrus.Logger{
		Out:          os.Stderr,
		Hooks:        make(logrus.LevelHooks),
		Formatter:    new(logrus.TextFormatter),
		Level:        logrus.DebugLevel,
		ReportCaller: false,
	}
	config.ConnConfig.Logger = logrusadapter.NewLogger(logger)
	config.ConnConfig.LogLevel = pgx.LogLevelDebug

	ctx := context.Background()

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func InitDatabase(pool *pgxpool.Pool) error {
	initFile, err := ioutil.ReadFile(PathToInitDBFile)
	if err != nil {
		return err
	}
	initComands := string(initFile)

	ctx := context.Background()
	_, err = pool.Exec(ctx, initComands)
	return err
}

func DropDatabase(pool *pgxpool.Pool) error {
	dropFile, err := ioutil.ReadFile(PathToDropDBFile)
	if err != nil {
		return err
	}
	initComands := string(dropFile)

	ctx := context.Background()
	_, err = pool.Exec(ctx, initComands)
	return err
}
