package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
)

type TableSpace interface {
	LoadPageByPageNumber(pageNo uint32) ([]byte, error)

	GetSegINodeFullList() *INodeList

	GetSegINodeFreeList() *INodeList

	GetFspFreeExtentList() *ExtentList

	GetFspFreeFragExtentList() *ExtentList

	GetFspFullFragExtentList() *ExtentList

	GetFirstFsp() *Fsp

	GetFirstINode() *INode

	GetDictTable() *DictionarySys

	GetSpaceId() uint32
}

/**
*
*
******/

type RowIter struct {
	index int
	rows  []innodb.Row
}

func (r RowIter) Next() (innodb.Row, error) {
	panic("implement me")
}

func (r RowIter) Close() error {
	panic("implement me")
}

type TableWrapper struct {
	Tuple *TableTupleMeta
	Index *IndexMeta
	Name  string
}

func (t *TableWrapper) InsertWithoutKey(row innodb.Row) error {
	panic("implement me")
}

func (t *TableWrapper) InsertWithKey(row innodb.Row, key innodb.Value) error {
	panic("implement me")
}

func (t *TableWrapper) InsertReturnKey(row innodb.Row) (key innodb.Value, err error) {
	panic("implement me")
}

func (t *TableWrapper) TableName() string {
	return t.Name
}

func (t TableWrapper) TableId() uint64 {
	panic("implement me")
}

func (t TableWrapper) SpaceId() uint32 {
	panic("implement me")
}

func (t TableWrapper) ColNums() int {
	panic("implement me")
}

func (t TableWrapper) RowIter() (innodb.RowIterator, error) {
	panic("implement me")
}

func NewTableWrapper(conf *conf.Cfg, databaseName string, tableName string, colDefs []*sqlparser.ColumnDefinition, indexDefs []*sqlparser.IndexDefinition) {
	tableTuple := NewTupleMeta(databaseName, tableName, conf).(*TableTupleMeta)
	tableTuple.WriteIndexDefinitions(indexDefs)
	tableTuple.WriteTupleColumns(colDefs)
	tableTuple.FlushToDisk()
	NewTableSpaceFile(conf, databaseName, tableName, 0, false)
}

func CreateTable(conf *conf.Cfg, databaseName string, tableName string, colDefs []*sqlparser.ColumnDefinition, indexDefs []*sqlparser.IndexDefinition) {
	tableTuple := NewTupleMeta(databaseName, tableName, conf).(*TableTupleMeta)
	tableTuple.WriteIndexDefinitions(indexDefs)
	tableTuple.WriteTupleColumns(colDefs)
	tableTuple.FlushToDisk()
	NewTableSpaceFile(conf, databaseName, tableName, 0, false)
}

//系统表wrapper

//系统表空间Space
type SysTableSpaceWrapper struct {
	SysTableSpace *SysTableSpace
	tableName     string
}

func (s SysTableSpaceWrapper) InsertWithoutKey(row innodb.Row) error {
	//s.SysTableSpace.BtreeMaps[common.INNODB_SYS_TABLESPACES].Insert(row)
	return nil
}

func (s SysTableSpaceWrapper) InsertWithKey(row innodb.Row, key innodb.Value) error {
	panic("implement me")
}

func (s SysTableSpaceWrapper) InsertReturnKey(row innodb.Row) (key innodb.Value, err error) {
	//	s.SysTableSpace.BtreeMaps[""].Insert(row)
	return row.GetPrimaryKey(), err
}

//
func NewSysTableSpaceWrapper(SysTableSpace *SysTableSpace) schemas.Table {
	var sysTableSpaceWrapper = new(SysTableSpaceWrapper)
	sysTableSpaceWrapper.SysTableSpace = SysTableSpace
	sysTableSpaceWrapper.tableName = common.INNODB_SYS_TABLESPACES
	return sysTableSpaceWrapper
}

func (s *SysTableSpaceWrapper) TableName() string {
	return s.tableName
}

func (s *SysTableSpaceWrapper) TableId() uint64 {
	return 13
}

func (s *SysTableSpaceWrapper) SpaceId() uint32 {
	return 0
}

func (s *SysTableSpaceWrapper) ColNums() int {
	panic("implement me")
}

func (s *SysTableSpaceWrapper) RowIter() (innodb.RowIterator, error) {
	//var sqlRows = make([]innodb.Row, 0)
	////s.SysTableSpace.BtreeMaps[common.INNODB_SYS_TABLESPACES].Ascend(func(i innodb.Row) bool {
	////	sqlRows = append(sqlRows, i)
	////	return true
	////})
	//m := &MemoryIterator{
	//	Rows: sqlRows,
	//}
	//return m, nil
	return nil, nil
}
