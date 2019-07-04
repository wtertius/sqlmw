package dsn

import (
	"net/url"
	"strings"
)

func Extract(dsn string) (host string, db string, user string) {
	// for pgx dsn string
	if len(dsn) >= 9 && dsn[0] == 0 {
		dsn = dsn[9:]
	}
	if u, err := url.Parse(dsn); err == nil {
		host = u.Host
		db = strings.Trim(u.Path, "/")
		if u.User != nil {
			user = u.User.Username()
		}
		if db == "" && len(u.Query()["database"]) > 0 {
			db = u.Query()["database"][0]
		}
		if host != "" {
			return
		}
		db = ""
		user = ""
	}
	if u, err := url.ParseQuery(dsn); err == nil {
		if len(u["server"]) > 0 {
			host = u["server"][0]
		}
		if len(u["user id"]) > 0 {
			user = u["user id"][0]
		}
		if len(u["database"]) > 0 {
			db = u["database"][0]
		}
		if host != "" {
			return
		}
		db = ""
		user = ""
	}
	return
}
