package sqldata

import (
	"context"

	"github.com/wtertius/sqlmw/sql/driver/internal/sqldata"
)

type Option = sqldata.Option

type Data = sqldata.Data

const (
	ActionQuery    sqldata.ActionEnum = "query"
	ActionExec     sqldata.ActionEnum = "exec"
	ActionTx       sqldata.ActionEnum = "tx"
	ActionBegin    sqldata.ActionEnum = "begin"
	ActionCommit   sqldata.ActionEnum = "commit"
	ActionRollback sqldata.ActionEnum = "rollback"
)

// NewContext gets previous sqldata.Data from context, sets new opts for it and stores it in context
func NewContext(ctx context.Context, opts ...Option) context.Context {
	data := FromContext(ctx)
	for _, o := range opts {
		o(&data)
	}
	if data.Changer != nil {
		data.Changer(&data)
	}
	return context.WithValue(ctx, sqldata.DataKey, data)
}

func FromContext(ctx context.Context) Data {
	d, _ := ctx.Value(sqldata.DataKey).(Data)
	return d
}

// Operation gives short description of the operation which is executed
func Operation(o string) Option {
	return func(d *Data) {
		d.Operation = o
	}
}

// Handler sets the name of rpc handler in which the operation is executed
func Handler(h string) Option {
	return func(d *Data) {
		d.Handler = h
	}
}

// Changer adds callback to modify sqldata.Data
func Changer(fn func(*Data)) Option {
	return func(d *Data) {
		d.Changer = fn
	}
}

// WithoutArgs cleans Args from sqldata.Data, useful if you do not want to store arguments
// to restore writing args, just set <nil> instead of Changer:
// sqldata.NewContext(ctx, sqldata.Changer(nil))
func WithoutArgs() Option {
	return Changer(func(data *Data) {
		data.Args = nil
	})
}
