package db

import (
	"github.com/jackc/pgx"
)

type Database interface {

	Execute(handler func (tx Transaction))

}

type DatabaseConfig struct {
	Host              string // host (e.g. localhost) or path to unix domain socket directory (e.g. /private/tmp)
	Port              uint16 // default: 5432
	Database          string
	User              string // default: OS user name
	Password          string
							 // Run-time parameters to set on connection as session default values
							 // (e.g. search_path or application_name)
	RuntimeParams     map[string]string
}

type db struct {

	poolConfig *pgx.ConnPoolConfig
	connectionPool *pgx.ConnPool

}

func (self *db) connection() *pgx.ConnPool {
	if self.connectionPool == nil {
		var pool, err = pgx.NewConnPool(*self.poolConfig)
		if err != nil {
			panic(err)
		}
		self.connectionPool = pool
	}
	return self.connectionPool
}

func (self *db) Execute(handler func (tx Transaction)) {
	var tx txWrapper
	defer func() {
		err := recover()
		tx.Rollback()
		if err != nil {
			panic(err)
		}
	}()
	tx.conn = self.connection()
	handler(&tx)
	tx.Commit()
}

func NewDatabase(config *DatabaseConfig) Database {
	cfg := pgx.ConnConfig{
		Host : config.Host,
		Port : config.Port,
		Database : config.Database,
		User : config.User,
		Password : config.Password,
	}
	return &db{poolConfig: &pgx.ConnPoolConfig{ConnConfig:cfg}}
}