package schemas

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
)

type Table interface {
	TableName() string

	TableId() uint64

	SpaceId() uint32

	ColNums() int

	RowIter() (innodb.RowIterator, error)

	InsertWithoutKey(row innodb.Row) error

	InsertWithKey(row innodb.Row, key innodb.Value) error

	InsertReturnKey(row innodb.Row) (key innodb.Value, err error)
}
