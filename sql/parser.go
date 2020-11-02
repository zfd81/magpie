package sql

import (
	"fmt"

	"vitess.io/vitess/go/vt/sqlparser"
)

func Parse(sql string) (*MQS, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	return convert(stmt)
}

func convert(stmt sqlparser.Statement) (*MQS, error) {
	mqs := &MQS{}
	switch n := stmt.(type) {
	case *sqlparser.Select:
		return mqs, convertSelect(mqs, n)
	case *sqlparser.Insert:
		return mqs, nil
	case *sqlparser.Delete:
		return mqs, nil
	case *sqlparser.Update:
		return mqs, nil
	default:
		return mqs, fmt.Errorf("unsupported syntax: %#v", n)
	}
}

func convertSelect(mqs *MQS, s *sqlparser.Select) error {
	tables, err := tableExprsToTables(s.From)
	if err != nil {
		return err
	}
	mqs.From = tables
	if s.Where != nil {
		conditions, err := whereToFilter(s.Where)
		if err != nil {
			return err
		}
		mqs.Where = conditions
	}
	fields, err := selectExprsToFields(s.SelectExprs)
	if err != nil {
		return err
	}
	mqs.Select = fields
	return nil
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
	var fields []*Field
	for _, e := range se {
		field, err := selectExprToField(e)
		if err != nil {
			return fields, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func selectExprToField(se sqlparser.SelectExpr) (*Field, error) {
	switch e := se.(type) {
	case *sqlparser.AliasedExpr:
		field := &Field{}
		expr := e.Expr
		et1, ok := expr.(*sqlparser.ColName)
		if ok {
			field.Name = et1.Name.String()
		}
		et2, ok := expr.(*sqlparser.BinaryExpr)
		if ok {
			buf := sqlparser.NewTrackedBuffer(nil)
			et2.Format(buf)
			field.Expr = buf.String()
			field.As = e.As.String()
		}

		et3, ok := expr.(*sqlparser.ComparisonExpr)
		if ok {
			buf := sqlparser.NewTrackedBuffer(nil)
			et3.Format(buf)
			field.Expr = buf.String()
			field.As = e.As.String()
		}

		return field, nil
	default:
		return nil, fmt.Errorf("unsupported syntax: %#v", e)
	}
}
