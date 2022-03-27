package optimizer

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
)

type BtreeIndexCursor struct {
}

func (b BtreeIndexCursor) GetCurrentRow() innodb.Row {
	panic("implement me")
}

func (b BtreeIndexCursor) Next() innodb.Row {
	panic("implement me")
}
