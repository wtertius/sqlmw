package driver

import "database/sql/driver"

type Driver struct {
	driver.Driver
	Wrapper
}

func (drv *Driver) Open(name string) (driver.Conn, error) {
	c, err := drv.Driver.Open(name)
	if err != nil {
		return c, err
	}
	if cc, ok := c.(conn_qe_nvc); ok {
		return &conn_qe_nvc_wr{&conn_qe_wr{&conn_wr{conn: cc, Wrapper: drv.Wrapper}}}, nil
	}
	if cc, ok := c.(conn_qe); ok {
		return &conn_qe_wr{&conn_wr{conn: cc, Wrapper: drv.Wrapper}}, nil
	}
	if cc, ok := c.(conn_nvc); ok {
		return &conn_nvc_wr{&conn_wr{conn: cc, Wrapper: drv.Wrapper}}, nil
	}
	if cc, ok := c.(conn); ok {
		return &conn_wr{conn: cc, Wrapper: drv.Wrapper}, nil
	}
	return c, nil
}
