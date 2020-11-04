package sql

type Row []interface{}

func (r *Row) Append(data interface{}) {
	*r = append(*r, data)
}

type Expression interface {
	Expr() string
	// Eval evaluates the given row and returns a result.
	//Eval(*Context, Row) (interface{}, error)
}

type Statement interface {
	Type() string
}

func (*InsertStatement) Type() string {
	return "Insert"
}

func (*DeleteStatement) Type() string {
	return "Delete"
}

func (*UpdateStatement) Type() string {
	return "Update"
}

func (*SelectStatement) Type() string {
	return "Select"
}

type InsertStatement struct {
	Table   string
	Columns []string
	Rows    []*Row
}

type DeleteStatement struct {
	Table string
	Where []*Condition
}

type UpdateStatement struct {
	Table  string
	Fields []*Field
	Where  []*Condition
}

//magpie query statement
type SelectStatement struct {
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
