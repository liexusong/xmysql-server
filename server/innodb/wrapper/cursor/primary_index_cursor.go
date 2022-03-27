package cursor

import (
	"github.com/zhukovaskychina/xmysql-server/server"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
)

type PrimaryIndexCursor struct {
	Cursor
	serverSession server.MySQLServerSession

	table schemas.Table
}

func NewPrimaryIndexCursor(serverSession server.MySQLServerSession, table schemas.Table) Cursor {
	var primaryIndexCursor = new(PrimaryIndexCursor)
	primaryIndexCursor.serverSession = serverSession
	primaryIndexCursor.table = table
	return primaryIndexCursor
}

func (p *PrimaryIndexCursor) Open() error {
	panic("implement me")
}

func (p *PrimaryIndexCursor) Close() error {
	panic("implement me")
}

func (p *PrimaryIndexCursor) GetRow() innodb.Row {
	panic("implement me")
}

func (p *PrimaryIndexCursor) Next() bool {
	panic("implement me")
}
