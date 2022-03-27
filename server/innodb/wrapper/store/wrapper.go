package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/util"
)

type PageWrapper interface {
}

type DatabaseWrapper struct {
	DatabaseName string
	Cache        map[string]*TableWrapper
}

func (databaseWrapper *DatabaseWrapper) Name() string {
	return databaseWrapper.DatabaseName
}

func (databaseWrapper *DatabaseWrapper) GetTable(name string) (schemas.Table, error) {
	return databaseWrapper.Cache[name], nil
}

func (databaseWrapper *DatabaseWrapper) ListTables() []schemas.Table {
	var list = make([]schemas.Table, 0)

	for _, v := range databaseWrapper.Cache {
		list = append(list, v)
	}

	return list
}

func (databaseWrapper *DatabaseWrapper) DropTable(name string) error {
	panic("implement me")
}

type CommonNodeInfo struct {
	NodeInfoLength     uint32 //节点数量
	PreNodePageNumber  uint32 //上一个节点的页面号
	PreNodeOffset      uint16 //上一个节点的偏移量
	NextNodePageNumber uint32
	NextNodeOffset     uint16
}

func (ci CommonNodeInfo) ToBytes() []byte {
	var buff = make([]byte, 0)
	buff = append(buff, util.ConvertUInt4Bytes(ci.NodeInfoLength)...)
	buff = append(buff, util.ConvertUInt4Bytes(ci.PreNodePageNumber)...)
	buff = append(buff, util.ConvertUInt2Bytes(ci.PreNodeOffset)...)
	buff = append(buff, util.ConvertUInt4Bytes(ci.NextNodePageNumber)...)
	buff = append(buff, util.ConvertUInt2Bytes(ci.NextNodeOffset)...)
	return nil
}
