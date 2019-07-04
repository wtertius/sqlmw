package driver

import (
	"context"
	"database/sql/driver"

	isqldata "github.com/wtertius/sqlmw/sql/driver/internal/sqldata"
	"github.com/wtertius/sqlmw/sql/driver/sqldata"
)

type Stmt struct {
	driver.Stmt
	Wrapper
	query string
}

func (stmt *Stmt) execContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if s, ok := stmt.Stmt.(driver.StmtExecContext); ok {
		return s.ExecContext(ctx, args)
	}

	values := make([]driver.Value, len(args))
	for _, arg := range args {
		values[arg.Ordinal-1] = arg.Value
	}

	return stmt.Exec(values)
}

func (stmt *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	var err error

	list := namedToInterface(args)

	ctx = sqldata.NewContext(ctx, isqldata.Stmt(stmt.query, list...), isqldata.Action(sqldata.ActionExec))

	ctx = stmt.Wrapper.Before(ctx)

	results, err := stmt.execContext(ctx, args)

	stmt.Wrapper.After(ctx, err)
	return results, err
}

func (stmt *Stmt) queryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if s, ok := stmt.Stmt.(driver.StmtQueryContext); ok {
		return s.QueryContext(ctx, args)
	}

	values := make([]driver.Value, len(args))
	for _, arg := range args {
		values[arg.Ordinal-1] = arg.Value
	}
	return stmt.Query(values)
}

func (stmt *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	var err error

	list := namedToInterface(args)

	ctx = sqldata.NewContext(ctx, isqldata.Stmt(stmt.query, list...), isqldata.Action(sqldata.ActionQuery))

	ctx = stmt.Wrapper.Before(ctx)

	rows, err := stmt.queryContext(ctx, args)

	stmt.Wrapper.After(ctx, err)
	return rows, err
}

func (stmt *Stmt) Close() error                                    { return stmt.Stmt.Close() }
func (stmt *Stmt) NumInput() int                                   { return stmt.Stmt.NumInput() }
func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) { return stmt.Stmt.Exec(args) }
func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error)  { return stmt.Stmt.Query(args) }
