package driver

import (
	"context"
	"database/sql/driver"

	isqldata "github.com/wtertius/sqlmw/sql/driver/internal/sqldata"
	"github.com/wtertius/sqlmw/sql/driver/sqldata"
)

type Tx struct {
	driver.Tx
	Wrapper
	ctx   context.Context
	txCtx context.Context
}

func (tx *Tx) Commit() error {
	var err error

	tx.ctx = sqldata.NewContext(tx.ctx, isqldata.Action(sqldata.ActionCommit))

	tx.ctx = tx.Wrapper.Before(tx.ctx)

	err = tx.Tx.Commit()

	// finish commit query
	tx.Wrapper.After(tx.ctx, err)

	// finish TX
	tx.Wrapper.After(tx.txCtx, err)

	return err
}

func (tx *Tx) Rollback() error {
	var err error

	tx.ctx = sqldata.NewContext(tx.ctx, isqldata.Action(sqldata.ActionRollback))

	tx.ctx = tx.Wrapper.Before(tx.ctx)

	err = tx.Tx.Rollback()

	// finish rollback query
	tx.Wrapper.After(tx.ctx, err)

	// finish Tx
	tx.Wrapper.After(tx.txCtx, err)

	return err
}
