package driver

import (
	"database/sql/driver"
)

type conn_nvc interface {
	conn
	driver.NamedValueChecker
}

type conn_nvc_wr struct {
	*conn_wr
}

func (c *conn_nvc_wr) cn() conn_nvc {
	return c.conn.(conn_nvc)
}

func (c *conn_nvc_wr) CheckNamedValue(v *driver.NamedValue) error {
	return c.cn().CheckNamedValue(v)
}
