package driver

import (
	"context"
	"database/sql/driver"

	isqldata "github.com/wtertius/sqlmw/sql/driver/internal/sqldata"
	"github.com/wtertius/sqlmw/sql/driver/sqldata"
)

type conn_qe interface {
	conn
	driver.QueryerContext
	driver.ExecerContext
}

type conn_qe_wr struct {
	*conn_wr
}

func (c *conn_qe_wr) cn() conn_qe {
	return c.conn.(conn_qe)
}

func (c *conn_qe_wr) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	var err error

	list := namedToInterface(args)

	ctx = sqldata.NewContext(ctx, isqldata.Stmt(query, list...), isqldata.Action(sqldata.ActionQuery))

	ctx = c.Wrapper.Before(ctx)

	results, err := c.Wrapper.QueryerContext(c.cn()).QueryContext(ctx, query, args)

	c.Wrapper.After(ctx, err)
	return results, err
}

func (c *conn_qe_wr) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	var err error

	list := namedToInterface(args)

	ctx = sqldata.NewContext(ctx, isqldata.Stmt(query, list...), isqldata.Action(sqldata.ActionExec))

	ctx = c.Wrapper.Before(ctx)

	results, err := c.Wrapper.ExecerContext(c.cn()).ExecContext(ctx, query, args)

	c.Wrapper.After(ctx, err)
	return results, err
}
