package engine

import (
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server"
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/wrapper/store"
)

//SQL执行引擎
//默认一个实例
type XMySQLEngine struct {
	conf *conf.Cfg
	//定义查询线程
	QueryExecutor *XMySQLExecutor
	//定义purge线程
	//定义SchemaManager
	infoSchemaManager schemas.InfoSchemas
}

func NewXMySQLEngine(conf *conf.Cfg) *XMySQLEngine {
	var mysqlEngine = new(XMySQLEngine)
	mysqlEngine.conf = conf
	mysqlEngine.infoSchemaManager = store.NewInfoSchemaManager(conf)
	mysqlEngine.QueryExecutor = NewXMySQLExecutor(mysqlEngine.infoSchemaManager, conf)
	return mysqlEngine
}

//ast->planner->store->result->net
func (e *XMySQLEngine) ExecuteQuery(session server.MySQLServerSession, query string, databaseName string) {
	fmt.Println(query)
	results := e.QueryExecutor.ExecuteWithQuery(session, query, databaseName)

	select {
	case rs, okv := <-results:
		{
			if !okv {
				panic("出错了")
			}
			switch rs.ResultType {
			case common.RESULT_TYPE_QUERY:
				{
					session.SendHandleOk()

				}
			case common.RESULT_TYPE_DDL:
				{
					session.SendOK()
				}
			case common.RESULT_TYPE_SET:
				{
					session.SendOK()
				}
			}

			fmt.Println("=======结束了======")
			break
		}
	}

}
