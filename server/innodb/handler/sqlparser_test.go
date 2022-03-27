package handler

import (
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
	"testing"
)

func TestProcessMySQLPacketFromClient(t *testing.T) {
	//sql := "SELECT * FROM table WHERE a = 'abc'"
	//sql := "show variables like '%innodb_old_blocks'"
	//sql:="show schemas"
	sql := "DESCRIBE student"
	stmt, err := sqlparser.Parse(sql)

	if err != nil {
		// Do something with the err
		fmt.Println(err)
	}

	// Otherwise do something with stmt
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		_ = stmt
	case *sqlparser.Insert:
	case *sqlparser.OtherRead:
		{
			fmt.Println("other read")
		}
	case *sqlparser.Show:
		{
			fmt.Println(stmt.Type)

			fmt.Println(stmt)

		}

	}
}

func TestCreateTableWithIndex(t *testing.T) {
	//sql:="CREATE TABLE index1(id INT,name VARCHAR(20),sex int, INDEX(id))"
	sql := "CREATE TABLE index5(id INT PRIMARY KEY,name VARCHAR(20), sex CHAR(4), INDEX index5_ns(name,sex),UNIQUE INDEX uqx(name))"
	//sql:="create index index_test using hash on test1(id)"
	stmt, _ := sqlparser.Parse(sql)

	switch stmt := stmt.(type) {
	case *sqlparser.DDL:
		{
			action := stmt.Action
			switch action {
			case sqlparser.CreateStr:
				fmt.Println(stmt)
				//sqlparser.IndexDefinition{}
				fmt.Println(stmt)
			case sqlparser.DropStr:

			case sqlparser.RenameStr:

			case sqlparser.AlterStr:
				fmt.Println(stmt)
			case sqlparser.CreateVindexStr:
				fmt.Println(stmt)
			case sqlparser.AddColVindexStr:
				fmt.Println(stmt)
			case sqlparser.DropColVindexStr:

			default:
			}
		}
	}
}
