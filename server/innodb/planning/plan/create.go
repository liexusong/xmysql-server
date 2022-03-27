package plan

import (
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
)

type DatabaseCreator interface {
	CreateDatabase(databaseName string) (schemas.Database, error)
}

type CreateTable struct {
	database schemas.Database
	conf     *conf.Cfg
	stmt     *sqlparser.DDL
}

func NewCreateTable(creator schemas.Database, conf *conf.Cfg, stmt *sqlparser.DDL) *CreateTable {
	return &CreateTable{
		conf:     conf,
		database: creator,
		stmt:     stmt,
	}
}

func (t *CreateTable) Columns() []string {
	return nil
}

func (t *CreateTable) RowIter() (innodb.RowIterator, error) {

	if _, err := t.database.CreateTable(t.conf, t.stmt); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return nil, nil
}
