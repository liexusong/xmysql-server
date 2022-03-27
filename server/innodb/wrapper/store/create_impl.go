package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
)

type CreateTableImpl struct {
	//plan.TableCreator
}

func (c CreateTableImpl) CreateTableWithDatabase(database schemas.Database, stmt *sqlparser.DDL) (schemas.Table, error) {
	panic("implement me")
}
