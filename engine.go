package qb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/serenize/snaker"
)

// NewEngine generates a new engine and returns it as an engine pointer
func NewEngine(driver string, dsn string) (*Engine, error) {
	conn, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// set name mapper function
	conn.MapperFunc(func(name string) string {
		return snaker.CamelToSnake(name)
	})

	return &Engine{
		driver: driver,
		dsn:    dsn,
		db:     conn,
	}, err
}

// Engine is the generic struct for handling db connections
type Engine struct {
	driver string
	dsn    string
	db     *sqlx.DB
}

// Exec executes insert & update type queries and returns sql.Result and error
func (e *Engine) Exec(query *QueryElem) (sql.Result, error) {
	stmt, err := e.db.Prepare(query.SQL())
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(query.Bindings()...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// QueryRow wraps *sql.DB.QueryRow()
func (e *Engine) QueryRow(query *QueryElem) *sql.Row {
	return e.db.QueryRow(query.SQL(), query.Bindings()...)
}

// Query wraps *sql.DB.Query()
func (e *Engine) Query(query *QueryElem) (*sql.Rows, error) {
	return e.db.Query(query.SQL(), query.Bindings()...)
}

// Get maps the single row to a model
func (e *Engine) Get(query *QueryElem, model interface{}) error {
	return e.db.Get(model, query.SQL(), query.Bindings()...)
}

// Select maps multiple rows to a model array
func (e *Engine) Select(query *QueryElem, model interface{}) error {
	return e.db.Select(model, query.SQL(), query.Bindings()...)
}

// DB returns sql.DB of wrapped engine connection
func (e *Engine) DB() *sqlx.DB {
	return e.db
}

// Ping pings the db using connection and returns error if connectivity is not present
func (e *Engine) Ping() error {
	return e.db.Ping()
}

// Driver returns the driver as string
func (e *Engine) Driver() string {
	return e.driver
}

// Dsn returns the connection dsn
func (e *Engine) Dsn() string {
	return e.dsn
}
