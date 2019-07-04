//go:generate reform

package test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/wtertius/pp"
	sqlmwdriver "github.com/wtertius/sqlmw/sql/driver"
	"github.com/wtertius/sqlmw/sql/driver/wrapper"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pgsqlCreateEmptyUserTableQuery = `
		CREATE table IF NOT EXISTS users(id int, name text);
		TRUNCATE table users;
	`

func pgsqlCreateEmptyUserTable(pgDB *sql.DB) error {
	_, err := pgDB.ExecContext(context.Background(), pgsqlCreateEmptyUserTableQuery)

	return err
}

func createEmptyPeopleTable() error {
	_, err := pgDB.ExecContext(context.Background(), `
		CREATE table IF NOT EXISTS people(
		  id int not null,
		  name text not null,
		  email text,
		  created_at timestamp(9) with time zone not null,
		  updated_at timestamp(9) with time zone
		);
		TRUNCATE table people;
	`)

	return err
}

func newQueryWrap(queries *[]string) sqlmwdriver.Wrapper {
	return &sqlmwdriver.CustomWrapper{
		BeforeFunc: func(ctx context.Context) context.Context {
			_, file, line, ok := runtime.Caller(3)
			pp.Printf("line := %s:%s\t%t\n", file, line, ok)
			return ctx
		},
		ExecerContextFunc: func(execerContext driver.ExecerContext) driver.ExecerContext {
			return sqlmwdriver.ExecerContextFunc(func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
				*queries = append(*queries, query)
				return execerContext.ExecContext(ctx, query, args)
			})
		},
		QueryerContextFunc: func(queryerContext driver.QueryerContext) driver.QueryerContext {
			return sqlmwdriver.QueryerContextFunc(func(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
				*queries = append(*queries, query)
				return queryerContext.QueryContext(ctx, query, args)
			})
		},
		ConnPrepareContextFunc: func(connPrepareContext driver.ConnPrepareContext) driver.ConnPrepareContext {
			return sqlmwdriver.ConnPrepareContextFunc(func(ctx context.Context, query string) (driver.Stmt, error) {
				*queries = append(*queries, query)
				return connPrepareContext.PrepareContext(ctx, query)
			})
		},
		ConnPrepareFunc: func(connPrepare sqlmwdriver.ConnPrepare) sqlmwdriver.ConnPrepare {
			return sqlmwdriver.ConnPrepareFunc(func(query string) (driver.Stmt, error) {
				*queries = append(*queries, query)
				return connPrepare.Prepare(query)
			})
		},
	}
}

func Test_PGSQL_SelectInSliceOfInt(t *testing.T) {
	queries := make([]string, 0, 3)
	queriesExpected := make([]string, 0, 3)

	ctx := context.Background()

	chain := wrapper.NewChain("postgres", fmt.Sprintf("postgres://postgres:password@%s:5432?sslmode=disable", pgHOST))
	chain.Add(newQueryWrap(&queries))

	pgDB, err := sql.Open(chain.NameAndDSN())
	if err != nil {
		log.Fatal(err)
	}

	require.NoError(t, pgsqlCreateEmptyUserTable(pgDB))
	queriesExpected = append(queriesExpected, pgsqlCreateEmptyUserTableQuery)

	id1, id2 := 1, 2
	query := "INSERT INTO users (id) VALUES ($1), ($2)"
	_, err = pgDB.ExecContext(ctx, query, id1, id2)
	require.NoError(t, err)

	queriesExpected = append(queriesExpected, query)

	userIds := make([]int, 0, 2)
	query = "SELECT id FROM users"
	queriesExpected = append(queriesExpected, query)

	rows, err := pgDB.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)
		if err != nil {
			log.Fatal(err)
		}

		userIds = append(userIds, userId)
	}

	assert.NoError(t, err)
	require.Equal(t, 2, len(userIds), "expected count of results: 2, got: %d", len(userIds))
	assert.Equal(t, id1, userIds[0])
	assert.Equal(t, id2, userIds[1])

	assert.Equal(t, queriesExpected, queries)
}
