package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Conn *sql.DB
}

func NewDB(username, password, host, port, dbname string) (*DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	db := &DB{
		Conn: conn,
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
