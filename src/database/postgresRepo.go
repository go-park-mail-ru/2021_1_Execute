package postgresRepo

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

type PostgresRepo struct {
	Host           string
	Port           uint16
	DatabaseName   string
	User           string
	MaxConnections int
	pool           *pgx.ConnPool
}

func (repo *PostgresRepo) ConfigConnection(host, dbName, user, passw string, port uint16, maxConnections int) error {
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

func (repo *PostgresRepo) GetConnection() (*pgx.Conn, error) {
	conn, err := repo.pool.Acquire()
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating Conn")
	}
	return conn, nil
}

func (repo *PostgresRepo) CloseConnection() {
	repo.pool.Close()
}

// func main() {
// 	var repo PostgresRepo
// 	err := repo.ConfigConnection("localhost", "testdb", "test", "test", 5432, 10)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer repo.CloseConnection()
// 	user := &api.UserRegistrationRequest{
// 		Email:    "ggwp@",
// 		Username: "ggwpomg",
// 		Password: "123456",
// 	}

// 	usr, err := repo.CreateUser(user)

// 	fmt.Println("OUT:\n\t", usr, "\n\t", err)
// }
