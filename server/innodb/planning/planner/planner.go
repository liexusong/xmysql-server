package planner

import (
	"errors"
	"fmt"
	_ "github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/planning"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/planning/plan"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
	"github.com/zhukovaskychina/xmysql-server/util"
	"path"
	"strings"
)

type Planner struct {
	conf    *conf.Cfg
	schemas schemas.InfoSchemas
}

func NewPlanner(conf *conf.Cfg, schemas schemas.InfoSchemas) *Planner {
	return &Planner{conf: conf, schemas: schemas}
}

func (p *Planner) Plan(database string, stmt *sqlparser.Statement) (planning.PlanNode, error) {

	switch stmts := (*stmt).(type) {
	case sqlparser.SelectStatement:
		{

			return p.PlanSelect(database, stmts.(*sqlparser.Select))
		}
	}

	return nil, nil
}

func (p *Planner) PlanSelect(database string, stmt *sqlparser.Select) (planning.PlanNode, error) {
	var table schemas.Table
	var node planning.PlanNode
	var err error
	if table, node, err = p.planScan(database, stmt.From); err != nil {
		return nil, fmt.Errorf("failed to plan scan: %w", err)
	}

	if node, err = p.planFilter(table, stmt.Where, node); err != nil {
		return nil, fmt.Errorf("failed to plan filter: %w", err)
	}

	return node, err
}

func (p *Planner) planScan(database string, stmt sqlparser.TableExprs) (schemas.Table, planning.PlanNode, error) {
	if stmt == nil {
		return nil, nil, nil
	}

	tb, _ := stmt.Eval()
	tableName := tb.Raw().(string)
	table, err := p.getTable(database, strings.ToUpper(tableName))
	if err != nil {
		return nil, nil, err
	}
	return table, plan.NewScan(table), nil
}

func (p *Planner) planFilter(table schemas.Table, stmt *sqlparser.Where, child planning.PlanNode) (planning.PlanNode, error) {
	if stmt == nil {
		return child, nil
	}

	if table == nil {
		return nil, fmt.Errorf("table not specified")
	}

	return plan.NewFilter(stmt.Expr, table.TableName(), child), nil
}

func (p *Planner) PlanCreateTable(databaseName string, stmt *sqlparser.DDL) (planning.PlanNode, error) {
	isDataBaseExist := p.schemas.GetSchemaExist(databaseName)
	if !isDataBaseExist {
		return nil, errors.New("数据库不存在")
	}

	db, err := p.schemas.GetSchemaByName(databaseName)

	if err != nil {
		return nil, errors.New("数据库不存在")
	}

	return plan.NewCreateTable(db, p.conf, stmt), nil
}

func (p *Planner) PlanCreateDatabase(stmt *sqlparser.DBDDL) (planning.PlanNode, error) {

	dbName := stmt.DBName

	util.CreateDataBaseDir(p.conf.DataDir, dbName)
	util.CreateFile(path.Join(p.conf.DataDir, dbName), "db.opt")
	charsetOption := "default-character-set=utf-8" + "\n"
	collation := "default-collation=utf8_general_ci" + "\n"
	util.WriteToFileByAppendBytes(path.Join(p.conf.DataDir, dbName, "/"), "db.opt", []byte(charsetOption))
	util.WriteToFileByAppendBytes(path.Join(p.conf.DataDir, dbName, "/"), "db.opt", []byte(collation))

	return nil, nil
}

func (p *Planner) PlanInsert(databaseName string, stmt *sqlparser.Insert) (planning.PlanNode, error) {
	isDataBaseExist := p.schemas.GetSchemaExist(databaseName)
	if !isDataBaseExist {
		return nil, errors.New("数据库不存在")
	}

	db, err := p.schemas.GetSchemaByName(databaseName)
	if err != nil {
		return nil, errors.New("数据库不存在")
	}
	return plan.NewInsertTableIterator(db, stmt), nil
}

func (p *Planner) getTable(databaseName, tableName string) (schemas.Table, error) {
	if databaseName == "" {
		return nil, fmt.Errorf("database not specified")
	}

	if tableName == "" {
		return nil, fmt.Errorf("table not specified")
	}
	currentTable, errors := p.schemas.GetTableByName(databaseName, tableName)
	if errors != nil {
		return nil, errors
	}

	return currentTable, nil
}
