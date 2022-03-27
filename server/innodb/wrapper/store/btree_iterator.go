package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
)

type Iterator func() (innodb.Value, innodb.Row, error, Iterator)

type RowItemsIterator func() (innodb.Row, error, RowItemsIterator)

type bpt_iterator func() (pageNo uint32, idxRecord int, err error, bi bpt_iterator)
