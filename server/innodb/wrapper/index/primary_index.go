package index

import (
	"github.com/google/btree"
	"github.com/zhukovaskychina/xmysql-server/server"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
)

type PrimaryIndex struct {
	//索引名称
	IndexName string

	//Btree
	Btree *btree.BTree
}

func (p *PrimaryIndex) Find(session server.MySQLServerSession, searchValueStart innodb.Value, searchValueEnd innodb.Value) innodb.Cursor {
	panic("implement me")
}

func (p PrimaryIndex) GetRow(session server.MySQLServerSession, primaryKey innodb.Value) innodb.Row {
	panic("implement me")
}
