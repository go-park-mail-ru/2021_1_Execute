package database

import (
	"context"
	"io/ioutil"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

const PathToInitDBFile = "database/trello.sql"

func GetPool(username, password, dbname, host string, port int) (*pgxpool.Pool, error) {
	DBUri := "postgresql://" + username + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/" + dbname
	config, err := pgxpool.ParseConfig(DBUri)
	if err != nil {
		return nil, err
	}

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
