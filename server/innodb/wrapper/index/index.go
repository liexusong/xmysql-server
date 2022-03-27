package index

import (
	"github.com/zhukovaskychina/xmysql-server/server"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
)

type Index interface {

	//Range 查询
	Find(session server.MySQLServerSession, searchValueStart innodb.Value, searchValueEnd innodb.Value) innodb.Cursor

	//根据主键获取行记录
	GetRow(session server.MySQLServerSession, primaryKey innodb.Value) innodb.Row
}
