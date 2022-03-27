package store

import (
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/segs"
	"github.com/zhukovaskychina/xmysql-server/util"
)

type DataDictWrapper struct {
	PageWrapper
	DataHrdPage   *pages.DataDictionaryHeaderSysPage
	SegmentHeader *segs.SegmentHeader
	MaxRowId      uint64 //没有主键的，也不运训null的unique主键的，则分配一个rowid

	MaxTableId uint64 //tableId

	MaxIndexId uint64 //索引ID

	MaxSpaceId uint32 //空间ID

	SysTableClusterRoot   uint32 //SYS_TABLES_CLUSTER 根页面
	SysTableIdsIndexRoot  uint32 //SYS_TABLES_IDS 二级索引的根页面号
	SysColumnsIndexRoot   uint32 //SYS_COLUMNS 根页面号
	SysIndexesClusterRoot uint32 //SYS_INDEXS
	SysFieldsClusterRoot  uint32 //SYS_FIELDS
}

func NewDataDictWrapper() *DataDictWrapper {
	var fspBinary = pages.NewDataDictHeaderPage()

	return &DataDictWrapper{
		DataHrdPage: fspBinary,
	}
}

func (d *DataDictWrapper) GetMaxTableId() uint64 {
	d.MaxTableId++
	return d.MaxTableId
}
func (d *DataDictWrapper) GetMaxIndexId() uint64 {
	d.MaxIndexId++
	return d.MaxIndexId
}
func (d *DataDictWrapper) GetMaxSpaceId() uint32 {
	d.MaxSpaceId++
	return d.MaxSpaceId
}

func (d *DataDictWrapper) SetDataDictSegments() {
	d.SegmentHeader = segs.NewSegmentHeader(0, 2, 50)
	d.DataHrdPage.SegmentHeader = d.SegmentHeader.GetBytes()
}

//
func NewDataDictWrapperByBytes(content []byte) *DataDictWrapper {
	return &DataDictWrapper{
		DataHrdPage: pages.ParseDataDictHrdPage(content),
	}
}

func (d *DataDictWrapper) processDataDicts() {
	d.MaxRowId = util.ReadUB8Byte2Long(d.DataHrdPage.DataDictHeader.MaxRowId)
	d.MaxTableId = util.ReadUB8Byte2Long(d.DataHrdPage.DataDictHeader.MaxTableId)
	d.MaxIndexId = util.ReadUB8Byte2Long(d.DataHrdPage.DataDictHeader.MaxIndexId)
	d.MaxSpaceId = util.ReadUB4Byte2UInt32(d.DataHrdPage.DataDictHeader.MaxSpaceId)

	d.SysTableClusterRoot = util.ReadUB4Byte2UInt32(d.DataHrdPage.DataDictHeader.SysTableRootPage)
	d.SysTableIdsIndexRoot = util.ReadUB4Byte2UInt32(d.DataHrdPage.DataDictHeader.SysTablesIDSRootPage)
	d.SysColumnsIndexRoot = util.ReadUB4Byte2UInt32(d.DataHrdPage.DataDictHeader.SysColumnsRootPage)
	d.SysIndexesClusterRoot = util.ReadUB4Byte2UInt32(d.DataHrdPage.DataDictHeader.SysIndexesRootPage)
	d.SysFieldsClusterRoot = util.ReadUB4Byte2UInt32(d.DataHrdPage.DataDictHeader.SysFieldsRootPage)
}

func (d *DataDictWrapper) GetSerializeBytes() []byte {

	return d.DataHrdPage.GetSerializeBytes()
}

/**
定义数据字典结构
**/

type DictionarySys struct {
	currentRowId uint64

	currentTableId uint64 //tableId

	currentIndexId uint64 //索引ID

	currentSpaceId uint32 //空间ID

	SysTable *DictTable

	SysColumns *DictTable

	SysIndex *DictTable

	SysFields *DictTable

	DataDict *DataDictWrapper //7号页面

	sysTableTuple   TableTuple
	sysColumnsTuple TableTuple
	sysIndexTuple   TableTuple
	sysFieldsTuple  TableTuple
}

func NewDictionarySys() *DictionarySys {
	var dictSys = new(DictionarySys)
	dictSys.currentRowId = 0
	dictSys.currentTableId = 0
	dictSys.currentIndexId = 0
	dictSys.currentSpaceId = 0

	dictSys.sysTableTuple = NewSysTableTuple()
	dictSys.sysColumnsTuple = NewSysColumnsTuple()
	dictSys.sysFieldsTuple = NewSysFieldsTuple()
	dictSys.sysIndexTuple = NewSysIndexTuple()
	return dictSys
}

func NewDictionarySysByWrapper(dt *DataDictWrapper) *DictionarySys {
	var dictSys = new(DictionarySys)
	dictSys.currentRowId = dt.MaxRowId
	dictSys.currentTableId = dt.MaxTableId
	dictSys.currentIndexId = dt.MaxIndexId
	dictSys.currentSpaceId = dt.MaxSpaceId

	dictSys.sysTableTuple = NewSysTableTuple()
	dictSys.sysColumnsTuple = NewSysColumnsTuple()
	dictSys.sysFieldsTuple = NewSysFieldsTuple()
	dictSys.sysIndexTuple = NewSysIndexTuple()
	return dictSys
}

func (dictSys *DictionarySys) initDictionary(space *SysTableSpace) {

	dictSys.SysTable = NewDictTableWithRootIndex("INNODB_SYS_TABLE", "SYS_TABLE", 8, space.SysTables, space, dictSys.sysTableTuple)
	dictSys.SysColumns = NewDictTableWithRootIndex("INNODB_SYS_COLUMNS", "SYS_COLUMNS", 10, space.SysColumns, space, dictSys.sysColumnsTuple)
	//dictSys.SysIndex = NewDictTableWithRootIndex("INNODB_SYS_INDEX", "SYS_INDEX", 10, space.SysIndexes, space, nil)
	//dictSys.SysFields = NewDictTableWithRootIndex("INNODB_SYS_FIELDS", "SYS_FIELDS", 11, space.SysFields, space, nil)
	dictSys.currentSpaceId = space.DataDict.MaxSpaceId
	dictSys.currentTableId = space.DataDict.MaxTableId
	dictSys.currentRowId = space.DataDict.MaxRowId
	dictSys.currentIndexId = space.DataDict.MaxIndexId
}

func (dictSys *DictionarySys) CreateTable(databaseName string, tuple *TableTupleMeta) (err error) {
	//插入到SYS_TABLE中

	currentSysTableRow := NewClusterSysIndexLeafRow(dictSys.sysTableTuple, false)
	currentSysColumnRow := NewClusterSysIndexLeafRow(dictSys.sysColumnsTuple, false)
	currentSysIndexRow := NewClusterSysIndexLeafRow(dictSys.sysIndexTuple, false)
	currentSysFieldsRow := NewClusterSysIndexLeafRow(dictSys.sysFieldsTuple, false)

	dictSys.initSysTableRow(databaseName, tuple, currentSysTableRow)
	dictSys.initSysColumns(databaseName, tuple, currentSysColumnRow)
	dictSys.initSysIndex(databaseName, tuple, currentSysIndexRow)
	dictSys.initSysFields(databaseName, tuple, currentSysFieldsRow)

	err = dictSys.SysTable.AddDictRow(currentSysTableRow)
	if err != nil {
		return err
	}
	err = dictSys.SysTable.AddDictRow(currentSysColumnRow)
	if err != nil {
		return err
	}
	err = dictSys.SysTable.AddDictRow(currentSysIndexRow)
	if err != nil {
		return err
	}
	err = dictSys.SysTable.AddDictRow(currentSysFieldsRow)
	if err != nil {
		return err
	}
	return nil
}

//创建系统文件表
func (dictSys *DictionarySys) createSysDataFilesTable(databaseName string, tuple TableTuple) (err error) {
	//currentSysTableRow := NewClusterSysIndexLeafRow(dictSys.sysTableTuple)
	return dictSys.createSystemTable(databaseName, tuple)
}

//创建空间表
func (dictSys *DictionarySys) createSysTableSpacesTable(databaseName string, tuple TableTuple) (err error) {
	//currentSysTableRow := NewClusterSysIndexLeafRow(dictSys.sysTableTuple)

	return dictSys.createSystemTable(databaseName, tuple)

}

func (dictSys *DictionarySys) createSysTableWithData(databaseName string, tuple TableTuple, row innodb.Row) (err error) {
	//插入到SYS_TABLE中

	err = dictSys.SysTable.AddDictRow(row)

	if err != nil {
		return err
	}

	//err = dictSys.SysTable.AddDictRow(currentSysColumnRow)
	//if err != nil {
	//	return err
	//}
	//err = dictSys.SysTable.AddDictRow(currentSysIndexRow)
	//if err != nil {
	//	return err
	//}
	//err = dictSys.SysTable.AddDictRow(currentSysFieldsRow)
	//if err != nil {
	//	return err
	//}
	return nil
}

/**
创建系统表
@param	 databaseName 数据库名称
@param   tuple 元祖信息

每创建一个表，需要记录Columns
***/
func (dictSys *DictionarySys) createSystemTable(databaseName string, tuple TableTuple) (err error) {
	//插入到SYS_TABLE中

	currentSysTableRow := NewClusterSysIndexLeafRow(dictSys.sysTableTuple, false)
	dictSys.initSysTableRow(databaseName, tuple, currentSysTableRow)
	dictSys.currentRowId++
	err = dictSys.SysTable.AddDictRow(currentSysTableRow)
	if err != nil {
		return err
	}
	//currentSysColumnRow := NewClusterSysIndexLeafRow(dictSys.sysColumnsTuple)
	//currentSysIndexRow := NewClusterSysIndexLeafRow(dictSys.sysIndexTuple)
	//currentSysFieldsRow := NewClusterSysIndexLeafRow(dictSys.sysFieldsTuple)

	//dictSys.initSysColumns(databaseName, tuple, currentSysColumnRow)
	//dictSys.initSysIndex(databaseName, tuple, currentSysIndexRow)
	//dictSys.initSysFields(databaseName, tuple, currentSysFieldsRow)

	///插入SysColumns
	columnlength := tuple.GetColumnLength()

	for i := 0; i < columnlength; i++ {

		currentColumn := tuple.GetColumnInfos(byte(i))

		isHidden := currentColumn.IsHidden
		if isHidden {
			continue
		}

		currentColumnTableRow := NewClusterSysIndexLeafRow(dictSys.sysColumnsTuple, false)

		//
		//rowId
		currentColumnTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 0)
		//transaction_id
		currentColumnTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 1)
		//rowpointer
		currentColumnTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 2)
		//tableId
		currentColumnTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentTableId), 3)

		//name
		currentColumnTableRow.WriteBytesWithNullWithsPos([]byte(currentColumn.FieldName), 4)

		//pos
		currentColumnTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(uint64(uint32(i))), 5)

		//mtype
		currentColumnTableRow.WriteBytesWithNullWithsPos(util.ConvertUInt4Bytes(uint32(uint64(uint32(i)))), 6)

		//prtype
		currentColumnTableRow.WriteBytesWithNullWithsPos(util.ConvertUInt4Bytes(uint32(uint64(uint32(i)))), 7)

		//len
		currentColumnTableRow.WriteBytesWithNullWithsPos(util.ConvertUInt4Bytes(uint32(uint64(uint32(i)))), 8)
		//插入columns表
		err = dictSys.SysColumns.AddDictRow(currentColumnTableRow)
		dictSys.currentRowId++
		fmt.Println("currentRowId", dictSys.currentRowId)
	}
	dictSys.currentTableId++
	//err = dictSys.SysTable.AddDictRow(currentSysColumnRow)
	//if err != nil {
	//	return err
	//}
	//err = dictSys.SysTable.AddDictRow(currentSysIndexRow)
	//if err != nil {
	//	return err
	//}
	//err = dictSys.SysTable.AddDictRow(currentSysFieldsRow)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (dictSys *DictionarySys) initSysTableRow(databaseName string, tuple TableTuple, currentSysTableRow innodb.Row) {
	//rowId
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 0)
	//transaction_id
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 1)
	//rowpointer
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 2)
	//tableId
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentTableId), 3)
	//tableName
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte(databaseName+"/"+tuple.GetTableName()), 4)
	//flag
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte{0, 0, 0, 0}, 5)
	//N_COLS
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertUInt4Bytes(uint32(tuple.GetColumnLength())), 6)

	//space_id
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertUInt4Bytes(dictSys.currentSpaceId), 7)

	//FileFormat
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte("Antelope"), 8)
	//RowFormat
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte("Redundant"), 9)
	//ZipPageSize
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertUInt4Bytes(0), 10)
	//SpaceType
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte("space"), 11)
}

func (dictSys *DictionarySys) initSysColumns(databaseName string, tuple TableTuple, currentSysColumnRow innodb.Row) {

	//rowId
	currentSysColumnRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 0)
	//transaction_id
	currentSysColumnRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 1)
	//rowpointer
	currentSysColumnRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentRowId), 2)
	//tableId
	currentSysColumnRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(dictSys.currentTableId), 3)

	//tableName
	currentSysColumnRow.WriteWithNull([]byte(databaseName + "/" + tuple.GetTableName()))
	//flag
	currentSysColumnRow.WriteWithNull([]byte{0, 0, 0, 0})
	//N_COLS
	currentSysColumnRow.WriteWithNull(util.ConvertUInt4Bytes(uint32(tuple.GetColumnLength())))

	//space_id
	currentSysColumnRow.WriteWithNull(util.ConvertUInt4Bytes(dictSys.currentSpaceId))

	//FileFormat
	currentSysColumnRow.WriteWithNull([]byte("Antelope"))
	//RowFormat
	currentSysColumnRow.WriteWithNull([]byte("Redundant"))
	//ZipPageSize
	currentSysColumnRow.WriteWithNull(util.ConvertUInt4Bytes(0))
	//SpaceType
	currentSysColumnRow.WriteWithNull([]byte("space"))
}

func (dictSys *DictionarySys) initSysIndex(databaseName string, tuple TableTuple, currentSysIndexRow innodb.Row) {
	//rowId
	currentSysIndexRow.WriteWithNull(util.ConvertULong8Bytes(dictSys.currentRowId))
	//transaction_id
	currentSysIndexRow.WriteWithNull(util.ConvertULong8Bytes(dictSys.currentRowId))
	//rowpointer
	currentSysIndexRow.WriteWithNull(util.ConvertULong8Bytes(dictSys.currentRowId))
	//tableId
	currentSysIndexRow.WriteWithNull(util.ConvertULong8Bytes(dictSys.currentTableId))
	//tableName
	currentSysIndexRow.WriteWithNull([]byte(databaseName + "/" + tuple.GetTableName()))
	//flag
	currentSysIndexRow.WriteWithNull([]byte{0, 0, 0, 0})
	//N_COLS
	currentSysIndexRow.WriteWithNull(util.ConvertUInt4Bytes(uint32(tuple.GetColumnLength())))

	//space_id
	currentSysIndexRow.WriteWithNull(util.ConvertUInt4Bytes(dictSys.currentSpaceId))

	//FileFormat
	currentSysIndexRow.WriteWithNull([]byte("Antelope"))
	//RowFormat
	currentSysIndexRow.WriteWithNull([]byte("Redundant"))
	//ZipPageSize
	currentSysIndexRow.WriteWithNull(util.ConvertUInt4Bytes(0))
	//SpaceType
	currentSysIndexRow.WriteWithNull([]byte("space"))
}

func (dictSys *DictionarySys) initSysFields(databaseName string, tuple TableTuple, currentSysFieldsRow innodb.Row) {
	//rowId
	currentSysFieldsRow.WriteWithNull(util.ConvertULong8Bytes(dictSys.currentRowId))
	//transaction_id
	currentSysFieldsRow.WriteWithNull(util.ConvertULong8Bytes(dictSys.currentRowId))
	//rowpointer
	currentSysFieldsRow.WriteWithNull(util.ConvertULong8Bytes(dictSys.currentRowId))
	//tableId
	currentSysFieldsRow.WriteWithNull(util.ConvertULong8Bytes(dictSys.currentTableId))
	//tableName
	currentSysFieldsRow.WriteWithNull([]byte(databaseName + "/" + tuple.GetDatabaseName()))
	//flag
	currentSysFieldsRow.WriteWithNull([]byte{0, 0, 0, 0})
	//N_COLS
	currentSysFieldsRow.WriteWithNull(util.ConvertUInt4Bytes(uint32(tuple.GetColumnLength())))

	//space_id
	currentSysFieldsRow.WriteWithNull(util.ConvertUInt4Bytes(dictSys.currentSpaceId))

	//FileFormat
	currentSysFieldsRow.WriteWithNull([]byte("Antelope"))
	//RowFormat
	currentSysFieldsRow.WriteWithNull([]byte("Redundant"))
	//ZipPageSize
	currentSysFieldsRow.WriteWithNull(util.ConvertUInt4Bytes(0))
	//SpaceType
	currentSysFieldsRow.WriteWithNull([]byte("space"))
}

type DictTable struct {
	TableName  string
	RootPageNo uint32
	IndexName  string
	BTree      *BTree
	tuple      TableTuple
}

//创建
func NewDictTable(TableName string, IndexName string, rootPageNo uint32, space *SysTableSpace, tuple *SysTableTuple) *DictTable {

	var dictTable = new(DictTable)

	dictTable.TableName = TableName
	dictTable.RootPageNo = rootPageNo

	rootIndexPageBytes, err := space.LoadPageByPageNumber(rootPageNo)
	currentIndex := NewPageIndexByLoadBytes(rootIndexPageBytes)

	segLeafBytes := currentIndex.GetSegLeaf()

	segInternalBytes := currentIndex.GetSegTop()

	if err != nil {
		panic("加载异常")
	}
	segInternalSpaceId := util.ReadUB4Byte2UInt32(segInternalBytes[0:4])
	segInternalPageNumber := util.ReadUB4Byte2UInt32(segInternalBytes[4:8])
	segInternalOffset := util.ReadUB2Byte2Int(segInternalBytes[8:10])

	segLeafSpaceId := util.ReadUB4Byte2UInt32(segLeafBytes[0:4])
	segLeafPageNumber := util.ReadUB4Byte2UInt32(segLeafBytes[4:8])
	segLeafOffset := util.ReadUB2Byte2Int(segLeafBytes[8:10])

	internalSeg := NewInternalSegment(segInternalSpaceId, segInternalPageNumber, segInternalOffset, IndexName, space)

	datasegs := NewLeafSegment(segLeafSpaceId, segLeafPageNumber, segLeafOffset, IndexName, space)

	dictTable.BTree = NewBtree(rootPageNo, IndexName, internalSeg, datasegs, currentIndex, space.blockFile, true, true)

	dictTable.tuple = tuple
	return dictTable
}

func NewDictTableWithRootIndex(TableName string, IndexName string, rootPageNo uint32, rootIndex *Index, space *SysTableSpace, tuple TableTuple) *DictTable {

	var dictTable = new(DictTable)

	dictTable.TableName = TableName
	dictTable.RootPageNo = rootPageNo

	currentIndex := rootIndex

	segLeafBytes := currentIndex.GetSegLeaf()

	segInternalBytes := currentIndex.GetSegTop()

	segInternalSpaceId := util.ReadUB4Byte2UInt32(segInternalBytes[0:4])
	segInternalPageNumber := util.ReadUB4Byte2UInt32(segInternalBytes[4:8])
	segInternalOffset := util.ReadUB2Byte2Int(segInternalBytes[8:10])

	segLeafSpaceId := util.ReadUB4Byte2UInt32(segLeafBytes[0:4])
	segLeafPageNumber := util.ReadUB4Byte2UInt32(segLeafBytes[4:8])
	segLeafOffset := util.ReadUB2Byte2Int(segLeafBytes[8:10])

	internalSeg := NewInternalSegment(segInternalSpaceId, segInternalPageNumber, segInternalOffset, IndexName, space)

	datasegs := NewLeafSegment(segLeafSpaceId, segLeafPageNumber, segLeafOffset, IndexName, space)

	dictTable.BTree = NewBtree(rootPageNo, IndexName, internalSeg, datasegs, currentIndex, space.blockFile, true, true)
	dictTable.BTree.Tuple = tuple
	dictTable.tuple = tuple
	return dictTable
}

func (d *DictTable) AddDictRow(rows innodb.Row) error {

	err := d.BTree.Add(rows.GetPrimaryKey(), rows)

	return err
}

func (d *DictTable) RemoveRow(rows innodb.Row) error {

	//d.BTree._range(func() (a uint32, idx int, err error, bi bpt_iterator) {
	//		return
	//});

	return nil
}
