package sql

type Item interface {
	Name() string
	Expr() string
	As() string
}

//magpie query statement
type MQS struct {
	Select []*Field
	From   []*TableItem
	Where  []*Condition
}

type Field struct {
	Name string
	Expr string
	As   string
}

type TableItem struct {
	Name string
	Expr string
	As   string
}

type Condition struct {
	Name     string
	Operator string
	Value    []byte
}
