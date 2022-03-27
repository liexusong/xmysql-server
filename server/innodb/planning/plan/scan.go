package plan

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
)

type Scan struct {
	table schemas.Table
}

func NewScan(table schemas.Table) *Scan {
	return &Scan{
		table: table,
	}
}

func (s *Scan) Columns() []string {
	panic("implement me")
}

func (s *Scan) RowIter() (innodb.RowIterator, error) {
	return s.table.RowIter()
}
