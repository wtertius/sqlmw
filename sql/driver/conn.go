package driver

import (
	"context"
	"database/sql/driver"

	isqldata "github.com/wtertius/sqlmw/sql/driver/internal/sqldata"
	"github.com/wtertius/sqlmw/sql/driver/sqldata"
)

type conn interface {
	driver.Conn
	driver.ConnBeginTx
}

type conn_wr struct {
	conn
	Wrapper
}

func (c *conn_wr) cn() conn {
	return c.conn
}

func (c *conn_wr) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	var (
		stmt driver.Stmt
		err  error
	)

	if cp, ok := c.cn().(driver.ConnPrepareContext); ok {
		stmt, err = c.Wrapper.ConnPrepareContext(cp).PrepareContext(ctx, query)
	} else {
		stmt, err = c.Prepare(query)
	}

	if err != nil {
		return nil, err
	}

	return &Stmt{stmt, c.Wrapper, query}, nil
}

func (c *conn_wr) Prepare(query string) (driver.Stmt, error) {
	return c.Wrapper.ConnPrepare(c.cn()).Prepare(query)
}
func (c *conn_wr) Close() error              { return c.cn().Close() }
func (c *conn_wr) Begin() (driver.Tx, error) { return c.cn().Begin() }

func (c *conn_wr) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	txCtx := sqldata.NewContext(ctx, isqldata.Action(sqldata.ActionTx))
	txCtx = c.Wrapper.Before(txCtx)

	beginCtx := sqldata.NewContext(ctx, isqldata.Action(sqldata.ActionBegin))
	beginCtx = c.Wrapper.Before(beginCtx)

	tx, err := c.Wrapper.ConnBeginTx(c.cn()).BeginTx(beginCtx, opts)

	c.Wrapper.After(beginCtx, err)
	return tx, err

	return &Tx{Tx: tx, Wrapper: c.Wrapper, ctx: ctx, txCtx: txCtx}, err
}
