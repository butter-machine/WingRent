package postgres

import (
	"database/sql"
	"fmt"
	"wingrent/database"
)

type Postgres struct {
	database.Database
	database.IDatabase
}


func Initialize(username, password, database string, host string, port int) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	p := Postgres{}
	p.Conn = conn
	err = p.Conn.Ping()
	if err != nil {
		return nil, err
	}

	return &p, err
}

func (p *Postgres) Close()  {
	p.Conn.Close()
}

func (p *Postgres) GetConnection() *sql.DB {
	return p.Conn
}
