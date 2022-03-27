package plan

import (
	"errors"
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/planning"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
	"io"
)

type Filter struct {
	cond      sqlparser.Expr
	tableName string

	child planning.PlanNode
}

func NewFilter(cond sqlparser.Expr, tableName string, child planning.PlanNode) *Filter {
	return &Filter{
		cond:      cond,
		tableName: tableName,
		child:     child,
	}
}

func (f Filter) Columns() []string {
	panic("implement me")
}

func (f Filter) RowIter() (innodb.RowIterator, error) {
	iter, err := f.child.RowIter()
	if err != nil {
		return nil, fmt.Errorf("failed to get row iter: %w", err)
	}

	iter = &filterIter{
		cond: f.cond,
		iter: iter,
	}

	return iter, nil
}

type filterIter struct {
	cond sqlparser.Expr
	iter innodb.RowIterator
}

func NewFilterIter(cond sqlparser.Expr, iter innodb.RowIterator) innodb.RowIterator {
	var filterIterOne = new(filterIter)
	filterIterOne.iter = iter
	filterIterOne.cond = cond
	return filterIterOne
}

func (f *filterIter) Open() error {
	return f.iter.Open()
}

func (f *filterIter) Next() (innodb.Row, error) {
	for {
		row, err := f.iter.Next()
		switch {
		case errors.Is(err, io.EOF):
			return nil, err
		case err != nil:
			return nil, fmt.Errorf("failed to get next row: %w", err)
		}

		value, err := f.cond.Eval()
		if err != nil {
			return nil, err
		}

		isTrue, ok := value.Raw().(bool)
		if !ok {
			return nil, fmt.Errorf("argument must be type boolean, not type %T", isTrue)
		}

		if isTrue {
			return row, nil
		}
		return row, nil
	}
}

func (f *filterIter) Close() error {
	return f.iter.Close()
}

type TableFilter struct {
	TableName  string
	Cost       float64
	LeftValue  *ColNameVariable
	RightValue innodb.Value
}
type TableName struct {
	Name, Qualifier string
}
type ColNameVariable struct {
	Name      string
	Qualifier TableName
}
