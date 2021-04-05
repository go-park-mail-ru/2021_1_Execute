package postgreRepo

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

type PostgreRepo struct {
	Host           string
	Port           uint16
	DatabaseName   string
	User           string
	MaxConnections int
	pool           *pgx.ConnPool
}

func (repo *PostgreRepo) ConfigConnection(host, dbName, user, passw string, port uint16, maxConnections int) error {
	repo.Host = host
	repo.DatabaseName = dbName
	repo.User = user
	repo.Port = port
	repo.MaxConnections = maxConnections
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     repo.Host,
			Port:     repo.Port,
			Database: repo.DatabaseName,
			User:     repo.User,
			Password: passw,
		},
		MaxConnections: repo.MaxConnections,
	}
	pool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return errors.Wrap(err, "Error while creating ConnPool")
	}
	repo.pool = pool
	return nil
}

func (repo *PostgreRepo) CloseConnection() {
	repo.pool.Close()
}
