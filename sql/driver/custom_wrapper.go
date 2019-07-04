package driver

import (
	"context"
	"database/sql/driver"
)

type CustomWrapper struct {
	BeforeFunc             func(ctx context.Context) context.Context
	AfterFunc              func(ctx context.Context, err error)
	ExecerContextFunc      func(execerContext driver.ExecerContext) driver.ExecerContext
	QueryerContextFunc     func(queryerContext driver.QueryerContext) driver.QueryerContext
	ConnPrepareContextFunc func(connPrepareContext driver.ConnPrepareContext) driver.ConnPrepareContext
	ConnPrepareFunc        func(conn ConnPrepare) ConnPrepare
	ConnBeginTxFunc        func(connBeginTx driver.ConnBeginTx) driver.ConnBeginTx
}

func (w *CustomWrapper) Before(ctx context.Context) context.Context {
	if w == nil || w.BeforeFunc == nil {
		return ctx
	}

	return w.BeforeFunc(ctx)
}

func (w *CustomWrapper) After(ctx context.Context, err error) {
	if w == nil || w.AfterFunc == nil {
		return
	}

	w.AfterFunc(ctx, err)
}

func (w *CustomWrapper) ExecerContext(execerContext driver.ExecerContext) driver.ExecerContext {
	if w == nil || w.ExecerContextFunc == nil {
		return execerContext
	}

	return w.ExecerContextFunc(execerContext)
}

func (w *CustomWrapper) QueryerContext(queryerContext driver.QueryerContext) driver.QueryerContext {
	if w == nil || w.QueryerContextFunc == nil {
		return queryerContext
	}

	return w.QueryerContextFunc(queryerContext)
}

func (w *CustomWrapper) ConnPrepareContext(connPrepareContext driver.ConnPrepareContext) driver.ConnPrepareContext {
	if w == nil || w.ConnPrepareContextFunc == nil {
		return connPrepareContext
	}

	return w.ConnPrepareContextFunc(connPrepareContext)
}

func (w *CustomWrapper) ConnPrepare(conn ConnPrepare) ConnPrepare {
	if w == nil || w.ConnPrepareFunc == nil {
		return conn
	}

	return w.ConnPrepareFunc(conn)
}

func (w *CustomWrapper) ConnBeginTx(connBeginTx driver.ConnBeginTx) driver.ConnBeginTx {
	if w == nil || w.ConnBeginTxFunc == nil {
		return connBeginTx
	}

	return w.ConnBeginTxFunc(connBeginTx)
}

type execerContext struct {
	execContext func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error)
}

func (execer *execerContext) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return execer.execContext(ctx, query, args)
}

func ExecerContextFunc(execContext func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error)) driver.ExecerContext {
	return &execerContext{execContext: execContext}
}

type execer struct {
	exec func(query string, args []driver.Value) (driver.Result, error)
}

func (execer *execer) Exec(query string, args []driver.Value) (driver.Result, error) {
	return execer.exec(query, args)
}

func ExecerFunc(exec func(query string, args []driver.Value) (driver.Result, error)) driver.Execer {
	return &execer{exec: exec}
}

type queryerContext struct {
	queryContext func(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error)
}

func (queryer *queryerContext) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return queryer.queryContext(ctx, query, args)
}

func QueryerContextFunc(queryContext func(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error)) driver.QueryerContext {
	return &queryerContext{queryContext: queryContext}
}

type queryer struct {
	query func(query string, args []driver.Value) (driver.Rows, error)
}

func (queryer *queryer) Query(query string, args []driver.Value) (driver.Rows, error) {
	return queryer.query(query, args)
}

func QueryerFunc(query func(query string, args []driver.Value) (driver.Rows, error)) driver.Queryer {
	return &queryer{query: query}
}

type connPrepareContext struct {
	prepareContext func(ctx context.Context, query string) (driver.Stmt, error)
}

func (connPrepare *connPrepareContext) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return connPrepare.prepareContext(ctx, query)
}

func ConnPrepareContextFunc(prepareContext func(ctx context.Context, query string) (driver.Stmt, error)) driver.ConnPrepareContext {
	return &connPrepareContext{prepareContext: prepareContext}
}

type connPrepare struct {
	prepare func(query string) (driver.Stmt, error)
}

func (c *connPrepare) Prepare(query string) (driver.Stmt, error) {
	return c.prepare(query)
}

func ConnPrepareFunc(prepare func(query string) (driver.Stmt, error)) ConnPrepare {
	return &connPrepare{prepare: prepare}
}

type connBeginTx struct {
	beginTx func(ctx context.Context, opts driver.TxOptions) (driver.Tx, error)
}

func (connBeginTx *connBeginTx) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return connBeginTx.beginTx(ctx, opts)
}

func ConnBeginTxFunc(beginTx func(ctx context.Context, opts driver.TxOptions) (driver.Tx, error)) driver.ConnBeginTx {
	return &connBeginTx{beginTx: beginTx}
}
