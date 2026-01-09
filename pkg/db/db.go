package db

import (
	"context"
	"database/sql"
)

// DB defines an sql database interface for the service.
type DB interface {
	SQLDB() *sql.DB
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	PingContext(context.Context) error
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	markUnretriable(error) error
}
