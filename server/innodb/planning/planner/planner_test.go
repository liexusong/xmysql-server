package planner

import (
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/wrapper/store"
	"testing"
)

func TestPlannerSelectInnodbSysTable(t *testing.T) {
	t.Parallel()

	t.Run("select * from information_schemas.innodb_sys_tables", func(t *testing.T) {
		t.Parallel()
		conf := conf.NewCfg()
		conf.DataDir = "/Users/zhukovasky/xmysql/data"
		conf.BaseDir = "/Users/zhukovasky/xmysql"
		schemaManager := store.NewInfoSchemaManager(conf)
		planner := NewPlanner(conf, schemaManager)
		sql := "select * from innodb_sys_tables sysTables"
		stmt, _ := sqlparser.Parse(sql)

		switch stmt := stmt.(type) {

		case sqlparser.SelectStatement:
			{
				selectNode, errorSelect := planner.PlanSelect(common.INFORMATION_SCHEMAS, stmt.(*sqlparser.Select))

				if errorSelect != nil {
					t.FailNow()
				}

				currentRowIter, errs := selectNode.RowIter()
				currentRowIter.Open()
				if errs != nil {
					t.FailNow()
				}
				var dataRows []innodb.Row
				for {
					cr, err := currentRowIter.Next()
					if err != nil {
						break
					}
					tableName := cr.ReadValueByIndex(4).Raw().(string)
					fmt.Println(tableName)
					dataRows = append(dataRows, cr)
				}
			}
		}

	})

	t.Run("select * from information_schemas.innodb_sys_columns", func(t *testing.T) {
		t.Parallel()
		conf := conf.NewCfg()
		conf.DataDir = "/Users/zhukovasky/xmysql/data"
		conf.BaseDir = "/Users/zhukovasky/xmysql"
		schemaManager := store.NewInfoSchemaManager(conf)
		planner := NewPlanner(conf, schemaManager)
		sql := "select * from innodb_sys_columns sysTables"
		stmt, _ := sqlparser.Parse(sql)

		switch stmt := stmt.(type) {

		case sqlparser.SelectStatement:
			{
				selectNode, errorSelect := planner.PlanSelect(common.INFORMATION_SCHEMAS, stmt.(*sqlparser.Select))

				if errorSelect != nil {
					t.FailNow()
				}

				currentRowIter, errs := selectNode.RowIter()
				currentRowIter.Open()
				if errs != nil {
					t.FailNow()
				}
				var dataRows []innodb.Row
				for {
					cr, err := currentRowIter.Next()
					if err != nil {
						break
					}
					tableName := cr.ReadValueByIndex(4).Raw().(string)
					i := cr.ReadValueByIndex(0).Raw().(int64)
					fmt.Println(tableName)
					fmt.Println(i)
					dataRows = append(dataRows, cr)
				}
				currentRowIter.Close()
			}
		}

	})
}
