package driver

import (
	"database/sql/driver"
)

type conn_qe_nvc interface {
	conn_qe
	driver.NamedValueChecker
}

type conn_qe_nvc_wr struct {
	*conn_qe_wr
}

func (c *conn_qe_nvc_wr) cn() conn_qe_nvc {
	return c.conn.(conn_qe_nvc)
}

func (c *conn_qe_nvc_wr) CheckNamedValue(v *driver.NamedValue) error {
	return c.cn().CheckNamedValue(v)
}
