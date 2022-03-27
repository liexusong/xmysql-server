package plan

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
)

type InsertTableIterator struct {
	db   schemas.Database
	stmt *sqlparser.Insert
}

func NewInsertTableIterator(db schemas.Database, stmt *sqlparser.Insert) *InsertTableIterator {

	return &InsertTableIterator{
		db:   db,
		stmt: stmt,
	}
}

func (i InsertTableIterator) Columns() []string {
	panic("implement me")
}

func (i InsertTableIterator) RowIter() (innodb.RowIterator, error) {
	panic("implement me")
}
