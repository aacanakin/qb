package qb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/serenize/snaker"
	"log"
	"os"
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
		logger: DefaultLogger{LDefault, log.New(os.Stdout, "", -1)},
	}, err
}

// Engine is the generic struct for handling db connections
type Engine struct {
	driver  string
	dsn     string
	db      *sqlx.DB
	dialect Dialect
	logger  Logger
}

// SetDialects sets the dialect of engine lazily
func (e *Engine) SetDialect(dialect Dialect) {
	e.dialect = dialect
}

// Logger returns the active logger of engine
func (e *Engine) Logger() Logger {
	return e.logger
}

// SetLogger sets the logger of engine
func (e *Engine) SetLogger(logger Logger) {
	e.logger = logger
}

func (e *Engine) log(statement *Stmt) {
	logFlags := e.logger.LogFlags()
	if logFlags == LQuery || logFlags == (LQuery|LBindings) {
		e.logger.Println(statement.SQL())
	}
	if logFlags == LBindings || logFlags == (LQuery|LBindings) {
		e.logger.Println(statement.Bindings())
	}
	if logFlags != LDefault {
		e.logger.Println()
	}
}

// Exec executes insert & update type queries and returns sql.Result and error
func (e *Engine) Exec(builder Builder) (sql.Result, error) {
	statement := builder.Build(e.dialect)
	stmt, err := e.db.Prepare(statement.SQL())
	if err != nil {
		return nil, err
	}

	e.log(statement)

	res, err := stmt.Exec(statement.Bindings()...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// QueryRow wraps *sql.DB.QueryRow()
func (e *Engine) QueryRow(builder Builder) *sql.Row {
	statement := builder.Build(e.dialect)
	e.log(statement)
	return e.db.QueryRow(statement.SQL(), statement.Bindings()...)
}

// Query wraps *sql.DB.Query()
func (e *Engine) Query(builder Builder) (*sql.Rows, error) {
	statement := builder.Build(e.dialect)
	e.log(statement)
	return e.db.Query(statement.SQL(), statement.Bindings()...)
}

// Get maps the single row to a model
func (e *Engine) Get(builder Builder, model interface{}) error {
	statement := builder.Build(e.dialect)
	e.log(statement)
	return e.db.Get(model, statement.SQL(), statement.Bindings()...)
}

// Select maps multiple rows to a model array
func (e *Engine) Select(builder Builder, model interface{}) error {
	statement := builder.Build(e.dialect)
	e.log(statement)
	return e.db.Select(model, statement.SQL(), statement.Bindings()...)
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
