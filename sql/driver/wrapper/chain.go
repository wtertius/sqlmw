package wrapper

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"fmt"
	"hash/crc32"
	"sync"

	"github.com/wtertius/sqlmw/sql/driver"
	"github.com/wtertius/sqlmw/sql/driver/dsn"
)

type Chain struct {
	origDriverName string
	driverName     string
	dsn            string
	host           string
	db             string
	user           string
	list           []driver.Wrapper
}

// NewChain creates new Chain wrapper by `driverName` and `dataSourceName` and `sql.Register` itself
// new driver with wrapper will have new unique name, you can get driver name by calling `Name()` method
// you can also use `NameAndDSN()` method, it can be useful for `sql.Open`, fox. ex:
// `sql.Open(wrapper.NewChain(driverName, dsn).Add(myWrapper1).Add(myWrapper2).NameAndDSN())`
func NewChain(driverName string, dataSourceName string) *Chain {
	host, db, user := dsn.Extract(dataSourceName)
	wrapper := &Chain{origDriverName: driverName, host: host, db: db, user: user, dsn: dataSourceName}

	newDriverName := fmt.Sprintf("%s_%d", driverName, crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s", wrapper))))
	register(driverName, newDriverName, wrapper)
	wrapper.driverName = newDriverName

	return wrapper
}

// DSN returns data source name
func (ch *Chain) DSN() string {
	return ch.dsn
}

// Name returns driver name of registered driver with this compound wrapper
func (ch *Chain) Name() string {
	return ch.driverName
}

// OrigName returns driver name which was wrapped by this wrapper
func (ch *Chain) OrigName() string {
	return ch.origDriverName
}

// NameAndDSN returns driver name and DSN, output can be used in `sql.Open` function, for ex.:
// `sql.Open(wrapper.NewCompound(driverName, dsn).Add(myWrapper1).Add(myWrapper2).NameAndDSN())`
func (ch *Chain) NameAndDSN() (driverName string, dataSourceName string) {
	return ch.driverName, ch.dsn
}

// Host returns db host from DSN
func (ch *Chain) Host() string {
	return ch.host
}

// DBX returns db name from DSN
func (ch *Chain) DB() string {
	return ch.db
}

// User returns user name from DSN
func (ch *Chain) User() string {
	return ch.user
}

// Add adds new `driver.Wrapper` to the chain
func (ch *Chain) Add(wrapper driver.Wrapper) *Chain {
	ch.list = append(ch.list, wrapper)
	return ch
}

// Before calls Before for all wrappers
func (ch Chain) Before(ctx context.Context) context.Context {
	for _, w := range ch.list {
		ctx = w.Before(ctx)
	}
	return ctx
}

// After calls After for all wrappers
func (ch Chain) After(ctx context.Context, err error) {
	for _, w := range ch.list {
		w.After(ctx, err)
	}
}

// ExecerContext calls ExecerContext for all wrappers
func (ch Chain) ExecerContext(execerContext sqldriver.ExecerContext) sqldriver.ExecerContext {
	for _, w := range ch.list {
		execerContext = w.ExecerContext(execerContext)
	}
	return execerContext
}

// QueryerContext calls QueryerContext for all wrappers
func (ch Chain) QueryerContext(queryerContext sqldriver.QueryerContext) sqldriver.QueryerContext {
	for _, w := range ch.list {
		queryerContext = w.QueryerContext(queryerContext)
	}
	return queryerContext
}

// ConnPrepareContext calls ConnPrepareContext for all wrappers
func (ch Chain) ConnPrepareContext(connPrepareContext sqldriver.ConnPrepareContext) sqldriver.ConnPrepareContext {
	for _, w := range ch.list {
		connPrepareContext = w.ConnPrepareContext(connPrepareContext)
	}
	return connPrepareContext
}

// ConnPrepare calls ConnPrepare for all wrappers
func (ch Chain) ConnPrepare(connPrepare driver.ConnPrepare) driver.ConnPrepare {
	for _, w := range ch.list {
		connPrepare = w.ConnPrepare(connPrepare)
	}
	return connPrepare
}

// ConnBeginTx calls ConnBeginTx for all wrappers
func (ch Chain) ConnBeginTx(connBeginTx sqldriver.ConnBeginTx) sqldriver.ConnBeginTx {
	for _, w := range ch.list {
		connBeginTx = w.ConnBeginTx(connBeginTx)
	}
	return connBeginTx
}

var mu sync.Mutex

func register(driverName, newDriverName string, w driver.Wrapper) {
	mu.Lock()
	defer mu.Unlock()

	for _, d := range sql.Drivers() {
		if d == newDriverName {
			return
		}
	}

	sql.Register(newDriverName, driver.WrapByName(driverName, w))
}
