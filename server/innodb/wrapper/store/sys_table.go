package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
)

//INNODB_SYS_TABLES
type MemoryInnodbSysTable struct {
	//系统表元祖信息
	TableTuple TableTuple
	//数据字典，参数传入
	DictionarySys    *DictionarySys
	sysTableIterator Iterator
}

func NewMemoryInnodbSysTable(sys *DictionarySys) schemas.Table {
	var memoryInnodbSysTable = new(MemoryInnodbSysTable)
	memoryInnodbSysTable.DictionarySys = sys
	memoryInnodbSysTable.TableTuple = NewSysTableTuple()
	sysTableIterator, _ := sys.SysTable.BTree.Iterate()
	memoryInnodbSysTable.sysTableIterator = sysTableIterator
	return memoryInnodbSysTable
}

func (m MemoryInnodbSysTable) TableName() string {
	return common.INNODB_SYS_TABLES
}

func (m MemoryInnodbSysTable) TableId() uint64 {
	return 1
}

func (m MemoryInnodbSysTable) SpaceId() uint32 {
	return 0
}

func (m MemoryInnodbSysTable) ColNums() int {
	return m.TableTuple.GetColumnLength()
}

func (m MemoryInnodbSysTable) RowIter() (innodb.RowIterator, error) {
	return NewMemoryIterator(m.sysTableIterator), nil
}

func (m MemoryInnodbSysTable) InsertWithoutKey(row innodb.Row) error {
	panic("implement me")
}

func (m MemoryInnodbSysTable) InsertWithKey(row innodb.Row, key innodb.Value) error {
	panic("implement me")
}

func (m MemoryInnodbSysTable) InsertReturnKey(row innodb.Row) (key innodb.Value, err error) {
	panic("implement me")
}

/****###############

#####################
*/

type MemoryInnodbSysColumns struct {
	//系统表元祖信息
	TableTuple TableTuple
	//数据字典，参数传入
	DictionarySys    *DictionarySys
	sysTableIterator Iterator
}

func NewMemoryInnodbSysColumns(sys *DictionarySys) schemas.Table {
	var memoryInnodbSysTable = new(MemoryInnodbSysColumns)
	memoryInnodbSysTable.DictionarySys = sys
	memoryInnodbSysTable.TableTuple = NewSysTableTuple()
	sysTableIterator, _ := sys.SysColumns.BTree.Iterate()
	memoryInnodbSysTable.sysTableIterator = sysTableIterator
	return memoryInnodbSysTable
}

func (m MemoryInnodbSysColumns) TableName() string {
	return common.INNODB_SYS_COLUMNS
}

func (m MemoryInnodbSysColumns) TableId() uint64 {
	return 1
}

func (m MemoryInnodbSysColumns) SpaceId() uint32 {
	return 0
}

func (m MemoryInnodbSysColumns) ColNums() int {
	return m.TableTuple.GetColumnLength()
}

func (m MemoryInnodbSysColumns) RowIter() (innodb.RowIterator, error) {
	return NewMemoryIterator(m.sysTableIterator), nil
}

func (m MemoryInnodbSysColumns) InsertWithoutKey(row innodb.Row) error {
	panic("implement me")
}

func (m MemoryInnodbSysColumns) InsertWithKey(row innodb.Row, key innodb.Value) error {
	panic("implement me")
}

func (m MemoryInnodbSysColumns) InsertReturnKey(row innodb.Row) (key innodb.Value, err error) {
	panic("implement me")
}

type MemoryInnodbTableSpaces struct {
	TableTuple TableTuple
}

func (m MemoryInnodbTableSpaces) TableName() string {
	panic("implement me")
}

func (m MemoryInnodbTableSpaces) TableId() uint64 {
	panic("implement me")
}

func (m MemoryInnodbTableSpaces) SpaceId() uint32 {
	panic("implement me")
}

func (m MemoryInnodbTableSpaces) ColNums() int {
	panic("implement me")
}

func (m MemoryInnodbTableSpaces) RowIter() (innodb.RowIterator, error) {
	panic("implement me")
}

func (m MemoryInnodbTableSpaces) InsertWithoutKey(row innodb.Row) error {
	panic("implement me")
}

func (m MemoryInnodbTableSpaces) InsertWithKey(row innodb.Row, key innodb.Value) error {
	panic("implement me")
}

func (m MemoryInnodbTableSpaces) InsertReturnKey(row innodb.Row) (key innodb.Value, err error) {
	panic("implement me")
}

func NewMemoryInnodbTableSpaces() schemas.Table {
	var memoryInnodbTableSpaces = new(MemoryInnodbTableSpaces)
	memoryInnodbTableSpaces.TableTuple = NewSysSpacesTuple()
	return memoryInnodbTableSpaces
}

type MemoryInnodbDataFiles struct {
	TableTuple TableTuple
}

func (m MemoryInnodbDataFiles) TableName() string {
	return common.INNODB_SYS_DATAFILES
}

func (m MemoryInnodbDataFiles) TableId() uint64 {
	return 14
}

func (m MemoryInnodbDataFiles) SpaceId() uint32 {
	return 0
}

func (m MemoryInnodbDataFiles) ColNums() int {
	return m.TableTuple.GetColumnLength()
}

func (m MemoryInnodbDataFiles) RowIter() (innodb.RowIterator, error) {
	panic("implement me")
}

func (m MemoryInnodbDataFiles) InsertWithoutKey(row innodb.Row) error {
	panic("implement me")
}

func (m MemoryInnodbDataFiles) InsertWithKey(row innodb.Row, key innodb.Value) error {
	panic("implement me")
}

func (m MemoryInnodbDataFiles) InsertReturnKey(row innodb.Row) (key innodb.Value, err error) {
	panic("implement me")
}
