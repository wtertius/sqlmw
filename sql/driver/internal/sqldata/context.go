package sqldata

type ActionEnum string

type Key int

var DataKey Key

type Option func(*Data)

type Data struct {
	// This group of props is for external using
	Operation string
	Handler   string

	// This group of props is for internal using
	Action ActionEnum
	Stmt   string
	Args   []interface{}

	// Changer adds callback to modify Data
	Changer func(*Data)
}

func Action(action ActionEnum) Option {
	return func(d *Data) {
		d.Action = action
	}
}

func Stmt(stmt string, args ...interface{}) Option {
	return func(d *Data) {
		d.Stmt = stmt
		d.Args = args
	}
}
