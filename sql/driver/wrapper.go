package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type Wrapper interface {
	// Before calls before sql query
	Before(ctx context.Context) context.Context
	// After calls after sql query
	After(ctx context.Context, err error)

	ExecerContext(driver.ExecerContext) driver.ExecerContext
	QueryerContext(driver.QueryerContext) driver.QueryerContext
	ConnPrepareContext(driver.ConnPrepareContext) driver.ConnPrepareContext
	ConnPrepare(ConnPrepare) ConnPrepare
	ConnBeginTx(driver.ConnBeginTx) driver.ConnBeginTx
}

type ConnPrepare interface {
	// Prepare returns a prepared statement, bound to this connection.
	Prepare(query string) (driver.Stmt, error)
}

func Wrap(driver driver.Driver, wrapper Wrapper) driver.Driver {
	return &Driver{driver, wrapper}
}

func WrapByName(driverName string, wrapper Wrapper) driver.Driver {
	db, err := sql.Open(driverName, "")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}()
	return Wrap(db.Driver(), wrapper)
}
