package engine

import (
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	planner3 "github.com/zhukovaskychina/xmysql-server/server/innodb/planning/planner"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
)

type StatementExecutor struct {
	DatabaseName string
	planner      *planner3.Planner
}

func (s *StatementExecutor) ExecuteStatement(ctx *ExecutionContext, stmt *sqlparser.Statement) error {
	panic("implement me")
	//此处需要改写

	switch stmt := (*stmt).(type) {
	case sqlparser.SelectStatement:
		{
			fmt.Println(stmt)
		}
	}
	return nil
}

func (s *StatementExecutor) executeSelectStatement(ctx *ExecutionContext, stmt *sqlparser.SelectStatement) error {
	var err error
	var Result *Result

	rowIter, err := s.planner.PlanSelect(s.DatabaseName, (*stmt).(*sqlparser.Select))

	if err != nil {
		Result.Err = err
		ctx.Send(Result)
		return err
	}

	rowiterator, err := rowIter.RowIter()
	if err != nil {
		Result.Err = err
		ctx.Send(Result)
		return err
	}

	var Rows []innodb.Row
	for {
		var row innodb.Row
		row, err = rowiterator.Next()
		if err != nil {
			break
		}
		Rows = append(Rows, row)
	}
	Result.Err = err
	Result.Rows = Rows
	ctx.Send(Result)
	return nil
}
