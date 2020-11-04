package sql

import (
	"fmt"

	"vitess.io/vitess/go/vt/sqlparser"
)

func Parse(sql string) (Statement, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	return convert(stmt)
}

func convert(stmt sqlparser.Statement) (Statement, error) {
	switch n := stmt.(type) {
	case *sqlparser.Select:
		return convertSelect(n)
	case *sqlparser.Insert:
		return convertInsert(n)
	case *sqlparser.Delete:
		return convertDelete(n)
	case *sqlparser.Update:
		return convertUpdate(n)
	default:
		return nil, fmt.Errorf("unsupported syntax: %#v", n)
	}
}

func convertSelect(s *sqlparser.Select) (Statement, error) {
	tables, err := tableExprsToTables(s.From)
	if err != nil {
		return nil, err
	}
	stmt := &SelectStatement{}
	stmt.From = tables
	if s.Where != nil {
		conditions, err := whereToFilter(s.Where)
		if err != nil {
			return nil, err
		}
		stmt.Where = conditions
	}
	fields, err := selectExprsToFields(s.SelectExprs)
	if err != nil {
		return nil, err
	}
	stmt.Select = fields
	return stmt, nil
}

func tableExprsToTables(te sqlparser.TableExprs) ([]*TableItem, error) {
	if len(te) == 0 {
		return nil, fmt.Errorf("zero tables in FROM")
	}
	var tables []*TableItem
	for _, t := range te {
		table, err := tableExprToTable(t)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func tableExprToTable(te sqlparser.TableExpr) (*TableItem, error) {
	tableItem := &TableItem{}
	switch t := (te).(type) {
	case *sqlparser.AliasedTableExpr:
		switch e := t.Expr.(type) {
		case sqlparser.TableName:
			tableItem.Name = e.Name.String()
			return tableItem, nil
		case *sqlparser.Subquery:
			return nil, nil
		default:
			return nil, fmt.Errorf("unsupported syntax: %#v", t)
		}
	default:
		return nil, fmt.Errorf("unsupported syntax: %#v", t)
	}
}

func whereToFilter(w *sqlparser.Where) ([]*Condition, error) {
	conditions := []*Condition{}
	err := exprToExpression(&conditions, w.Expr)
	if err != nil {
		return nil, err
	}
	return conditions, nil
}

func exprToExpression(conditions *[]*Condition, e sqlparser.Expr) error {
	switch v := e.(type) {
	default:
		return fmt.Errorf("unsupported syntax: %#v", e)
	case *sqlparser.ComparisonExpr:
		condition := &Condition{}
		left := v.Left.(*sqlparser.ColName)
		condition.Name = left.Name.String()
		condition.Operator = v.Operator
		right := v.Right.(*sqlparser.SQLVal)
		condition.Value = right.Val
		*conditions = append(*conditions, condition)
		return nil
	case *sqlparser.AndExpr:
		err := exprToExpression(conditions, v.Left)
		if err != nil {
			return err
		}
		err = exprToExpression(conditions, v.Right)
		if err != nil {
			return err
		}
		return nil
	}
}

func selectExprsToFields(se sqlparser.SelectExprs) ([]*Field, error) {
	fields := make([]*Field, len(se))
	for i, e := range se {
		field, err := selectExprToField(e)
		if err != nil {
			return fields, err
		}
		fields[i] = field
	}
	return fields, nil
}

func selectExprToField(se sqlparser.SelectExpr) (*Field, error) {
	switch e := se.(type) {
	case *sqlparser.AliasedExpr:
		field := &Field{}
		expr := e.Expr
		switch f := expr.(type) {
		case *sqlparser.ColName:
			field.Name = f.Name.String()
		default:
			buf := sqlparser.NewTrackedBuffer(nil)
			f.Format(buf)
			field.Expr = buf.String()
			field.As = e.As.String()
		}
		return field, nil
	default:
		return nil, fmt.Errorf("unsupported syntax: %#v", e)
	}
}

func convertInsert(i *sqlparser.Insert) (Statement, error) {
	stmt := &InsertStatement{}
	stmt.Table = i.Table.Name.String()
	cols := i.Columns
	if cols != nil {
		for _, c := range cols {
			stmt.Columns = append(stmt.Columns, c.String())
		}
	}
	rows, ok := i.Rows.(sqlparser.Values)
	if ok {
		for _, r := range rows {
			row := &Row{}
			for _, f := range r {
				switch t := f.(type) {
				case *sqlparser.SQLVal:
					row.Append(string(t.Val))
				default:
					buf := sqlparser.NewTrackedBuffer(nil)
					t.Format(buf)
					row.Append(buf.String())
				}
			}
			stmt.Rows = append(stmt.Rows, row)
		}
	}
	return stmt, nil
}

func convertDelete(d *sqlparser.Delete) (Statement, error) {
	tables, err := tableExprsToTables(d.TableExprs)
	if err != nil {
		return nil, err
	}
	stmt := &DeleteStatement{}
	stmt.Table = tables[0].Name

	if d.Where != nil {
		conditions, err := whereToFilter(d.Where)
		if err != nil {
			return nil, err
		}
		stmt.Where = conditions
	}

	return stmt, nil
}

func convertUpdate(u *sqlparser.Update) (Statement, error) {
	tables, err := tableExprsToTables(u.TableExprs)
	if err != nil {
		return nil, err
	}
	stmt := &UpdateStatement{}
	stmt.Table = tables[0].Name
	fields, err := updateExprsToFields(u.Exprs)
	if err != nil {
		return nil, err
	}
	stmt.Fields = fields
	if u.Where != nil {
		conditions, err := whereToFilter(u.Where)
		if err != nil {
			return nil, err
		}
		stmt.Where = conditions
	}
	return stmt, nil
}

func updateExprsToFields(e sqlparser.UpdateExprs) ([]*Field, error) {
	fields := make([]*Field, len(e))
	for i, updateExpr := range e {
		field := &Field{}
		field.Name = updateExpr.Name.Name.String()

		buf := sqlparser.NewTrackedBuffer(nil)
		updateExpr.Expr.Format(buf)
		field.Expr = buf.String()
		fields[i] = field
	}
	return fields, nil
}
