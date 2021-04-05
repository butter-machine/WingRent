package database

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

type Database struct {
	Conn *sql.DB
}

type IDatabase interface {
	Close()
	GetOrCreate(m *interface{}, fieldNames []string, tableName string) (*interface{}, bool, error)
	GetConnection() *sql.DB
}

var DBSingleton IDatabase
