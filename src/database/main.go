package postgresRepo

import (
	"github.com/jackc/pgx"
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
		return err
	}
	repo.pool = pool
	return nil
}

func (repo *PostgresRepo) GetConnection() (*pgx.Conn, error) {
	conn, err := repo.pool.Acquire()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (repo *PostgresRepo) CloseConnection() {
	repo.pool.Close()
}

func (repo *PostgresRepo) RunQuery(sql string) (*pgx.Rows, error) {
	conn, err := repo.GetConnection()
	if err != nil {
		return nil, err
	}
	//defer conn.Close()
	rows, err := conn.Query(sql)
	if err != nil {
		return nil, err
	}
	//fmt.Println(rows)
	return rows, nil
}

// func main() {
// 	var repo PostgresRepo
// 	err := repo.ConfigConnection("localhost", "testdb", "test", "test", 5432, 10)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer repo.CloseConnection()
// 	// poolConfig := pgx.ConnPoolConfig{
// 	// 	ConnConfig: pgx.ConnConfig{
// 	// 		Host:     "localhost",
// 	// 		Port:     5432,
// 	// 		Database: "testdb",
// 	// 		User:     "test",
// 	// 		Password: "test",
// 	// 	},
// 	// 	MaxConnections: 10,
// 	// }

// 	// pool, err := pgx.NewConnPool(poolConfig)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// }
// 	// defer pool.Close()

// 	// conn, err := pool.Acquire()
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// }
// 	// defer conn.Close()
// 	for i := 0; i < 11; i++ {
// 		fmt.Println("Hi")
// 		rows, err := repo.RunQuery("select * from boards")
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		defer rows.Close()
// 		fmt.Println(rows)
// 		for rows.Next() {
// 			answ, err := rows.Values()
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			fmt.Println(answ)
// 		}
// 	}
// }
