package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbInstance *DB
	once       sync.Once
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

func DBInstance(dbName string) *DB {
	once.Do(func() {
		db, err := NewDB("root", "", "localhost", "3306", dbName)
		if err != nil {
			panic(err)
		}
		dbInstance = db
	})
	return dbInstance
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
