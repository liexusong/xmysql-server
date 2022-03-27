package engine

import (
	"context"
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server"
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/planning/planner"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
	"time"
)

//定义执行器

type XMySQLExecutor struct {
	infosSchemaManager schemas.InfoSchemas
	conf               *conf.Cfg
}

func NewXMySQLExecutor(infosSchemaManager schemas.InfoSchemas, conf *conf.Cfg) *XMySQLExecutor {
	var xMySQLExecutor = new(XMySQLExecutor)
	xMySQLExecutor.infosSchemaManager = infosSchemaManager
	xMySQLExecutor.conf = conf
	return xMySQLExecutor
}
func (e *XMySQLExecutor) ExecuteWithQuery(mysqlSession server.MySQLServerSession, query string, databaseName string) <-chan *Result {
	results := make(chan *Result)
	ctx := &ExecutionContext{
		Context:     context.Background(),
		statementId: 0,
		QueryId:     0,
		Results:     results,
		Cfg:         nil,
	}
	go e.executeQuery(ctx, mysqlSession, query, databaseName, results)
	return results
}

func (e *XMySQLExecutor) executeQuery(ctx *ExecutionContext, mysqlSession server.MySQLServerSession, query string, databaseName string, results chan *Result) {
	stmt, _ := sqlparser.Parse(query)
	defer close(results)
	e.recover(query, results)

	switch stmt := stmt.(type) {

	case sqlparser.SelectStatement:
		{

			e.executeSelectStatement(ctx, (stmt).(*sqlparser.Select), databaseName)
		}

	case *sqlparser.DDL:
		{
			action := stmt.Action
			fmt.Println(action)
			switch action {
			//CREATE, ALTER, DROP, RENAME or TRUNCATE statement
			case "create":
				{

					e.executeCreateTableStatement(ctx, databaseName, stmt)
				}

			}
		}
	case *sqlparser.DBDDL:
		{
			action := stmt.Action
			fmt.Println(action)
			switch action {
			case "create":
				{
					e.executeCreateDatabaseStatement(ctx, stmt)
				}
			case "drop":
				{

				}
			}
		}
	case *sqlparser.Set:
		{
			for _, v := range stmt.Exprs {
				fmt.Println(v.Name.String())
				//mysqlSession.SetParamByName(v.Name.String(), v.Expr)

				fmt.Println("-------------------0000000000000000----------------------------9991", time.Now())
				time.Sleep(time.Duration(10) * time.Second)
				fmt.Println("-------------------0000000000000000----------------------------999000", time.Now())

				ctx.Send(&Result{
					StatementID: 0,
					Rows:        nil,
					Err:         nil,
					ResultType:  common.RESULT_TYPE_SET,
				})

			}

		}
	}

}

func (e *XMySQLExecutor) executeSelectStatement(ctx *ExecutionContext, stmt *sqlparser.Select, databaseName string) error {

	currentPlanner := planner.NewPlanner(e.conf, e.infosSchemaManager)

	selectNode, errorSelect := currentPlanner.PlanSelect(databaseName, stmt)

	if errorSelect != nil {
		return errorSelect
	}

	currentRowIter, errs := selectNode.RowIter()

	if errs != nil {
		return errs
	}
	var dataRows []innodb.Row
	for {
		cr, err := currentRowIter.Next()
		if err != nil {
			break
		}
		dataRows = append(dataRows, cr)
	}
	ctxErr := ctx.Send(&Result{
		StatementID: 0,
		Rows:        dataRows,
		Err:         nil,
		ResultType:  "",
	})

	return ctxErr
}

func (e *XMySQLExecutor) buildWhereConditions(where *sqlparser.Where) {

}

func (e *XMySQLExecutor) executeInsertStatement(ctx *ExecutionContext, stmt *sqlparser.SelectStatement) {

}

func (e *XMySQLExecutor) executeSetStatement(ctx *ExecutionContext, stmt *sqlparser.Set) {
	//currentPlanner := planner.NewPlanner(e.infosSchemaManager)

}

func (e *XMySQLExecutor) executeCreateTableStatement(ctx *ExecutionContext, databaseName string, stmt *sqlparser.DDL) {
	currentPlanner := planner.NewPlanner(e.conf, e.infosSchemaManager)

	createTableIterator, createrror := currentPlanner.PlanCreateTable(databaseName, stmt)
	if createrror != nil {

	} else {
		createTableIterator.RowIter()
	}

	ctx.Send(&Result{
		StatementID: 0,
		Rows:        nil,
		Err:         createrror,
		ResultType:  common.RESULT_TYPE_DDL,
	})

}

func (e *XMySQLExecutor) executeCreateDatabaseStatement(ctx *ExecutionContext, stmt *sqlparser.DBDDL) {
	currentPlanner := planner.NewPlanner(e.conf, e.infosSchemaManager)

	currentPlanner.PlanCreateDatabase(stmt)

	ctx.Send(&Result{
		StatementID: 0,
		Rows:        nil,
		Err:         nil,
		ResultType:  common.RESULT_TYPE_DDL,
	})

}

//本处代码参考influxdb
//用于获取查询中的异常
func (e *XMySQLExecutor) recover(query string, results chan *Result) {
	if err := recover(); err != nil {
		//atomic.AddInt64(&e.stats.RecoveredPanics, 1) // Capture the panic in _internal stats.
		//e.Logger.Error(fmt.Sprintf("%s [panic:%s] %s", query.String(), err, debug.Stack()))
		//
		results <- &Result{
			StatementID: -1,
			Err:         fmt.Errorf("%s [panic:%s]", query, err),
		}
		//
		//if willCrash {
		//	e.Logger.Error(fmt.Sprintf("\n\n=====\nAll goroutines now follow:"))
		//	buf := debug.Stack()
		//	e.Logger.Error(fmt.Sprintf("%s", buf))
		//	os.Exit(1)
		//}
	}
}
