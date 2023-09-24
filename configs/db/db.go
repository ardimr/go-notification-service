package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type DBInterface interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type DB struct {
	*sql.DB
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func NewDBConnection(host, port, user, password, dbName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", "host="+host+" port="+port+" user="+user+" password="+password+" dbname="+dbName+" sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewDB(host, port, user, password, dbName string) (*DB, error) {
	db, err := NewDBConnection(host, port, user, password, dbName)
	return &DB{
		DB:       db,
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbName:   dbName,
	}, err
}

func (db *DB) reconnect() {
	db.DB, _ = NewDBConnection(db.host, db.port, db.user, db.password, db.dbName)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db.DB == nil {
		db.reconnect()
	}

	rows, err := db.DB.QueryContext(ctx, query, args...)
	return rows, err
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if db.DB == nil {
		db.reconnect()
	}

	return db.DB.QueryRowContext(ctx, query, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db.DB == nil {
		db.reconnect()
	}

	return db.DB.ExecContext(ctx, query, args...)
}

func (db *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if db.DB == nil {
		db.reconnect()
	}

	return db.DB.PrepareContext(ctx, query)
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if db.DB == nil {
		db.reconnect()
	}

	return db.DB.BeginTx(ctx, opts)
}
