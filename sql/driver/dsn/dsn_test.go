package dsn

import "testing"

func Test_Extract(t *testing.T) {
	tests := []struct {
		name     string
		dsn      string
		wantHost string
		wantDb   string
		wantUser string
	}{
		{"Connection URI", "postgresql://", "", "", ""},
		{"Connection URI", "postgresql://localhost", "localhost", "", ""},
		{"Connection URI", "postgresql://localhost:5433", "localhost:5433", "", ""},
		{"Connection URI", "postgresql://localhost/mydb", "localhost", "mydb", ""},
		{"Connection URI", "postgresql://user@localhost", "localhost", "", "user"},
		{"Connection URI", "postgresql://user:secret@localhost", "localhost", "", "user"},
		{"Connection URI", "postgresql://other@localhost/otherdb?connect_timeout=10&application_name=myapp", "localhost", "otherdb", "other"},
		{"Connection URI", "postgresql://user@host:123/somedb?target_session_attrs=any&application_name=myapp", "host:123", "somedb", "user"},
		{"Connection URI", "sqlserver://user:password@host.o3.ru:1433?database=database_name", "host.o3.ru:1433", "database_name", "user"},
		{"SQL Server connection string", "server=host;user id=user;password=password;database=database;app name=My-API", "host", "database", "user"},
		{"SQL Server connection string", "server=host;database=database;app name=My-API", "host", "database", ""},
		{"SQL Server connection string", "server=host;app name=My-API", "host", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotDb, gotUser := Extract(tt.dsn)
			if gotHost != tt.wantHost {
				t.Errorf("Extract() got host = `%v`, want `%v`", gotHost, tt.wantHost)
			}
			if gotDb != tt.wantDb {
				t.Errorf("Extract() got db = `%v`, want `%v`", gotDb, tt.wantDb)
			}
			if gotUser != tt.wantUser {
				t.Errorf("Extract() got user = `%v`, want `%v`", gotUser, tt.wantUser)
			}
		})
	}
}
