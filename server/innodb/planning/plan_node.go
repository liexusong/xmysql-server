package planning

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
)

type PlanNode interface {
	Columns() []string

	RowIter() (innodb.RowIterator, error)
}
