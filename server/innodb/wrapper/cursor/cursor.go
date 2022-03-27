package cursor

import "github.com/zhukovaskychina/xmysql-server/server/innodb"

type Cursor interface {

	//打开游标
	Open() error

	//获取当前行
	GetRow() innodb.Row

	//获取下一个
	Next() bool

	//关闭游标
	Close() error
}
