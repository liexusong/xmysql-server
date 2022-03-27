package store

import (
	"fmt"
	"github.com/smartystreets/assertions"
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/blocks"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
	"github.com/zhukovaskychina/xmysql-server/util"
	"path"
)

/*************

系统表空间和独立空间的前3个页面一致，
但是3-7的页面是系统表空间独有的。

//////////////////////////
// type   //  描述
/////////////////////////
// SYS	  //  Insert Buffer Header  存储ChangeBuffer的头部信息
/////////////////////////
// INDEX  //  Insert Buffer Root   存储ChangeBuffer的根页面
///////////////////////////
// TRX_SYS//  Transaction System 事物系统的相关信息
////////////////////////////
// SYS    //  First Rollback Segment  第一个回滚段的信息
////////////////////////////
// SYS    //  Data Dictionary Header 数据字典头部信息
////////////////////////////
对于一个新的segment，总是优先填满32个frag page数组，之后才会为其分配完整的Extent，可以利用碎片页，并避免小表占用太多空间。
尽量获得hint page;
如果segment上未使用的page太多，则尽量利用segment上的page。

***********/
type SysTableSpace struct {
	conf      *conf.Cfg
	blockFile *blocks.BlockFile

	Fsp        *Fsp   //0 号页面
	IBuf       *IBuf  //1 号页面
	FirstInode *INode //2 号页面

	DataDict *DataDictWrapper //7号页面

	SysTables *Index //8号页面

	SysTablesIds *Index //9号页面

	SysColumns *Index //10号页面

	SysIndexes *Index //11号页面

	SysFields *Index //12号页面

	SysForeign *Index //13号页面

	SysForeignCols *Index //14号页面

	SysTableSpaces *Index //15号页面

	SysDataFiles *Index //16号页面

	SysVirtual *Index //17号页面

	InodeMap map[uint32]*INode

	FirstIBDExtent *FirstIBDExtent //第一个区

	OtherExtents []Extent //其他区

	BtreeMaps map[string]*BTree //每个Btree中有两个段，索引段之叶子段和非叶子段

	DictionarySys *DictionarySys
}

//初始化数据库
func NewSysTableSpace(cfg *conf.Cfg) TableSpace {
	tableSpace := new(SysTableSpace)
	tableSpace.conf = cfg
	filePath := path.Join(cfg.BaseDir, "/", "ibdata1")
	isFlag, _ := util.PathExists(filePath)
	blockfile := blocks.NewBlockFile(cfg.BaseDir, "ibdata1", 256*64*16384)
	tableSpace.blockFile = blockfile
	if !isFlag {
		tableSpace.blockFile.CreateFile()
		tableSpace.initHeadPage()
		tableSpace.initFreeLimit()
		tableSpace.initSysTableDataDict()
		tableSpace.initDatabaseDictionary()
		tableSpace.flushToDisk()
		tableSpace.initAllSysTables()
		//	tableSpace.flushToDisk()
	} else {
		tableSpace.initAfterRebootMySQL()
		tableSpace.initDatabaseDictionary()
	}
	return tableSpace
}

/******************

加载页面，第二次启动加载已知的数据字典表，复盘段信息
*****/

func (sysTable *SysTableSpace) initAfterRebootMySQL() {
	//加载fsp
	fspBytes, _ := sysTable.LoadPageByPageNumber(0)

	sysTable.Fsp = NewFspByLoadBytes(fspBytes)

	inodeBytes, _ := sysTable.LoadPageByPageNumber(2)

	sysTable.FirstInode = NewINodeByByte(inodeBytes)

	dataDictBytes, _ := sysTable.LoadPageByPageNumber(7)

	sysTable.DataDict = NewDataDictWrapperByBytes(dataDictBytes)

	sysTableCluster, _ := sysTable.LoadPageByPageNumber(8)
	tuple := NewSysTableTuple()
	sysTable.SysTables = NewPageIndexByLoadBytesWithTuple(sysTableCluster, tuple)

	assertions.ShouldEqual(len(sysTable.SysTables.SlotRowData.FullRowList()), 3)
	//sysTableIds, _ := sysTable.LoadPageByPageNumber(9)
	//sysTable.SysTablesIds = NewPageIndexByLoadBytesWithTuple(sysTableIds)
	//
	columnTuple := NewSysColumnsTuple()
	sysColumns, _ := sysTable.LoadPageByPageNumber(10)

	sysTable.SysColumns = NewPageIndexByLoadBytesWithTuple(sysColumns, columnTuple)
	fmt.Println("")
	//
	//SysIndexes, _ := sysTable.LoadPageByPageNumber(11)
	//
	//sysTable.SysIndexes = NewPageIndexByLoadBytes(SysIndexes)
	//
	//sysFields, _ := sysTable.LoadPageByPageNumber(12)
	//
	//sysTable.SysFields = NewPageIndexByLoadBytes(sysFields)
	//
	//sysForeign, _ := sysTable.LoadPageByPageNumber(13)
	//sysTable.SysForeign = NewPageIndexByLoadBytes(sysForeign)
	//
	//sysForeignCols, _ := sysTable.LoadPageByPageNumber(14)
	//sysTable.SysForeignCols = NewPageIndexByLoadBytes(sysForeignCols)
	//
	//sysTableSpaces, _ := sysTable.LoadPageByPageNumber(15)
	//sysTable.SysTableSpaces = NewPageIndexByLoadBytes(sysTableSpaces)

}

func (sysTable *SysTableSpace) initHeadPage() {
	//初始化FspHrdPage
	fspHrdPage := pages.NewFspHrdPage(0)
	sysTable.Fsp = NewFsp(fspHrdPage)
	sysTable.IBuf = NewIBuf(0)
	sysTable.FirstInode = NewINode(2)
	sysTable.DataDict = NewDataDictWrapper()

	sysTable.SysTables = NewPageIndex(8)
	sysTable.SysTablesIds = NewPageIndex(9)
	sysTable.SysColumns = NewPageIndex(10)
	sysTable.SysIndexes = NewPageIndex(11)
	sysTable.SysFields = NewPageIndex(12)
	sysTable.SysTableSpaces = NewPageIndex(13)
	sysTable.SysDataFiles = NewPageIndex(14)
	sysTable.SysVirtual = NewPageIndex(15)

}

//初始化数据字典
func (sysTable *SysTableSpace) initDatabaseDictionary() {
	sysTable.DictionarySys = NewDictionarySysByWrapper(sysTable.DataDict)
	sysTable.DictionarySys.initDictionary(sysTable)
}

//初始化各种系统表

//SysTableSpaces
//SysDataFiles
func (sysTable *SysTableSpace) initAllSysTables() {

	sysTable.DictionarySys.createSysTableSpacesTable(common.INFORMATION_SCHEMAS, NewSysSpacesTuple())

	sysTable.DictionarySys.createSysDataFilesTable(common.INFORMATION_SCHEMAS, NewSysDataFilesTuple())

}

func (sysTable *SysTableSpace) initFreeLimit() {
	sysTable.Fsp.SetFreeLimit(64)
}

//初始化FspExtent信息
func (sysTable *SysTableSpace) initFspExtents() {
	sysTable.Fsp.SetFspFreeExtentListInfo(CommonNodeInfo{NodeInfoLength: 2, PreNodePageNumber: 0, PreNodeOffset: 0, NextNodePageNumber: 0, NextNodeOffset: 0})
	sysTable.Fsp.SetFreeLimit(128)
}

///

func (sysTable SysTableSpace) initSysDict() {

}

//初始化字典段
//
func (sysTable *SysTableSpace) initSysTableClusterSegments() {
	dataDictSegs := NewDataSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLES_CLUSTER", 0, sysTable)
	internalSegs := NewInternalSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLES_CLUSTER", 1, sysTable)

	sysTable.SysTables.SetPageBtrTop(dataDictSegs.GetSegmentHeader().GetBytes())
	sysTable.SysTables.SetPageBtrSegs(internalSegs.GetSegmentHeader().GetBytes())
	internalSegs.AllocateDiscretePage(8)
}

//初始化字典段
//
func (sysTable *SysTableSpace) initSysTableIdsSegments() {
	dataDictSegs := NewDataSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLE_IDS", 2, sysTable)
	internalSegs := NewInternalSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLE_IDS", 3, sysTable)
	sysTable.SysTablesIds.SetPageBtrTop(dataDictSegs.GetSegmentHeader().GetBytes())
	sysTable.SysTablesIds.SetPageBtrSegs(internalSegs.GetSegmentHeader().GetBytes())
	internalSegs.AllocateDiscretePage(9)
}

func (sysTable *SysTableSpace) initSysTableColumnsSegments() {
	dataDictSegs := NewDataSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLE_COLUMNS", 4, sysTable)
	internalSegs := NewInternalSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLE_COLUMNS", 5, sysTable)
	sysTable.SysColumns.SetPageBtrTop(dataDictSegs.GetSegmentHeader().GetBytes())
	sysTable.SysColumns.SetPageBtrSegs(internalSegs.GetSegmentHeader().GetBytes())
	internalSegs.AllocateDiscretePage(10)
}

func (sysTable *SysTableSpace) initSysTableIndexesSegments() {
	dataDictSegs := NewDataSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLE_INDEXES", 6, sysTable)
	internalSegs := NewInternalSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLE_INDEXES", 7, sysTable)
	sysTable.SysIndexes.SetPageBtrTop(dataDictSegs.GetSegmentHeader().GetBytes())
	sysTable.SysIndexes.SetPageBtrSegs(internalSegs.GetSegmentHeader().GetBytes())
	internalSegs.AllocateDiscretePage(11)
}

func (sysTable *SysTableSpace) initSysTableFieldsSegments() {
	dataDictSegs := NewDataSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLE_FIELDS", 8, sysTable)
	internalSegs := NewInternalSegmentWithTableSpaceAtInit(0, 2, "SYS_TABLE_FIELDS", 9, sysTable)
	sysTable.SysFields.SetPageBtrTop(dataDictSegs.GetSegmentHeader().GetBytes())
	sysTable.SysFields.SetPageBtrSegs(internalSegs.GetSegmentHeader().GetBytes())
	internalSegs.AllocateDiscretePage(12)
}

//初始化数据字典表
func (sysTable *SysTableSpace) initSysTableDataDict() {
	sysTable.initSysTableClusterSegments()
	sysTable.initSysTableIdsSegments()
	sysTable.initSysTableColumnsSegments()
	sysTable.initSysTableIndexesSegments()
	sysTable.initSysTableFieldsSegments()

}

func (sysTable *SysTableSpace) flushToDisk() {
	sysTable.blockFile.WriteContentByPage(0, sysTable.Fsp.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(1, sysTable.IBuf.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(2, sysTable.FirstInode.GetSerializeBytes())

	sysTable.blockFile.WriteContentByPage(7, sysTable.DataDict.DataHrdPage.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(8, sysTable.SysTables.IndexPage.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(9, sysTable.SysTablesIds.IndexPage.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(10, sysTable.SysColumns.IndexPage.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(11, sysTable.SysIndexes.IndexPage.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(12, sysTable.SysFields.IndexPage.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(13, sysTable.SysTableSpaces.IndexPage.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(14, sysTable.SysDataFiles.IndexPage.GetSerializeBytes())
	sysTable.blockFile.WriteContentByPage(15, sysTable.SysVirtual.IndexPage.GetSerializeBytes())
}

func (sysTable *SysTableSpace) initSysTableTable() {
}

func (sysTable *SysTableSpace) LoadPageByPageNumber(pageNo uint32) ([]byte, error) {
	return sysTable.blockFile.ReadPageByNumber(pageNo)
}

//获取所有的完全用满的InodePage链表
func (sysTable *SysTableSpace) GetSegINodeFullList() *INodeList {
	var inodeList = NewINodeList("FULL_LIST")

	fullNodeInfo := sysTable.Fsp.GetFullINodeBaseInfo()

	nextNodePageNo := fullNodeInfo.NextNodePageNumber

	//nextOffset := fullNodeInfo.NextNodeOffset

	//只处理当前256MB的数据,后面拓展
	for {
		//递归完成
		nextINodePage, _ := sysTable.LoadPageByPageNumber(nextNodePageNo)
		nextINode := NewINodeByByte(nextINodePage)
		inodeList.AddINode(nextINode)
		nextNodePageNo = nextINode.NextNodePageNumber
		if nextNodePageNo == 0 {
			break
		}
	}

	return inodeList
}

//至少存在一个空闲Inode Entry的Inode Page被放到该链表上
func (sysTable *SysTableSpace) GetSegINodeFreeList() *INodeList {
	var inodeList = NewINodeList("FREE_LIST")

	fullNodeInfo := sysTable.Fsp.GetFreeSegINodeBaseInfo()

	nextNodePageNo := fullNodeInfo.NextNodePageNumber

	//nextOffset := fullNodeInfo.NextNodeOffset

	//只处理当前256MB的数据,后面拓展
	for {
		//递归完成
		nextINodePage, _ := sysTable.LoadPageByPageNumber(nextNodePageNo)
		nextINode := NewINodeByByte(nextINodePage)
		inodeList.AddINode(nextINode)
		nextNodePageNo = nextINode.NextNodePageNumber
		if nextNodePageNo == 0 {
			break
		}
	}

	return inodeList
}

//获取所有FreeExtent链表
func (sysTable *SysTableSpace) GetFspFreeExtentList() *ExtentList {
	var extentList = NewExtentList("FREE_EXTENT")
	//获取当前free的区链表
	freeInitNode := sysTable.Fsp.GetFspFreeExtentListInfo()
	//下一个区的首页面，它的属性应该是xdes类型
	//在FSP页面里面，页面号应该是0，
	//
	nextNodePageNo := freeInitNode.NextNodePageNumber
	//下一个offset,最大值256
	nextOffset := freeInitNode.NextNodeOffset

	//只处理当前256MB的数据,后面拓展
	for {
		//递归完成
		nextXDesEntryPage, _ := sysTable.LoadPageByPageNumber(nextNodePageNo)
		nextXDESEntry := nextXDesEntryPage[nextOffset : nextOffset+40]
		//如果是free，则将该区加入到链表中
		if util.ReadUB4Byte2UInt32(nextXDESEntry[20:24]) == uint32(common.XDES_FREE) {
			currentPageNodeOffsetNo := (nextOffset - 150) / 40
			u := uint32(currentPageNodeOffsetNo) + nextNodePageNo
			extentList.AddExtent(NewOtherIBDExtent(u))
			nextNodePageNo = util.ReadUB4Byte2UInt32(nextXDESEntry[14:18])
			nextOffset = util.ReadUB2Byte2Int(nextXDESEntry[18:20])
			if nextOffset == 0 {
				break
			}
		}

	}

	if extentList.IsEmpty() {

	}

	return extentList
}

//获取所有的FreeFragExtent
func (sysTable *SysTableSpace) GetFspFreeFragExtentList() *ExtentList {
	var extentList = NewExtentList("FREE_FRAG")
	//获取当前free的区链表
	freeInitNode := sysTable.Fsp.GetFspFreeFrageExtentListInfo()
	//下一个区的首页面，它的属性应该是xdes类型
	//在FSP页面里面，页面号应该是0，
	//
	nextNodePageNo := freeInitNode.NextNodePageNumber
	//下一个offset,最大值256
	nextOffset := freeInitNode.NextNodeOffset

	//只处理当前256MB的数据,后面拓展
	for {
		//递归完成
		nextXDesEntryPage, _ := sysTable.LoadPageByPageNumber(nextNodePageNo)
		nextXDESEntry := nextXDesEntryPage[nextOffset : nextOffset+40]
		//如果是free，则将该区加入到链表中
		if util.ReadUB4Byte2UInt32(nextXDESEntry[20:24]) == uint32(common.XDES_FREE) {
			currentPageNodeOffsetNo := (nextOffset - 150) / 40
			u := uint32(currentPageNodeOffsetNo) + nextNodePageNo
			extentList.AddExtent(NewOtherIBDExtent(u))
			nextNodePageNo = util.ReadUB4Byte2UInt32(nextXDESEntry[14:18])
			nextOffset = util.ReadUB2Byte2Int(nextXDESEntry[18:20])
			if nextOffset == 0 {
				break
			}
		}

	}

	return extentList
}

//获取所有的FullFragExtent
//
func (sysTable *SysTableSpace) GetFspFullFragExtentList() *ExtentList {
	var extentList = NewExtentList("FULL_FRAG")
	//获取当前free的区链表
	freeInitNode := sysTable.Fsp.GetFspFullFragListInfo()
	//下一个区的首页面，它的属性应该是xdes类型
	//在FSP页面里面，页面号应该是0，
	//
	nextNodePageNo := freeInitNode.NextNodePageNumber
	//下一个offset,最大值256
	nextOffset := freeInitNode.NextNodeOffset

	//只处理当前256MB的数据,后面拓展
	for {
		//递归完成
		nextXDesEntryPage, _ := sysTable.LoadPageByPageNumber(nextNodePageNo)
		nextXDESEntry := nextXDesEntryPage[nextOffset : nextOffset+40]
		//如果是free，则将该区加入到链表中
		if util.ReadUB4Byte2UInt32(nextXDESEntry[20:24]) == uint32(common.XDES_FREE) {
			currentPageNodeOffsetNo := (nextOffset - 150) / 40
			u := uint32(currentPageNodeOffsetNo) + nextNodePageNo
			extentList.AddExtent(NewOtherIBDExtent(u))
			nextNodePageNo = util.ReadUB4Byte2UInt32(nextXDESEntry[14:18])
			nextOffset = util.ReadUB2Byte2Int(nextXDESEntry[18:20])
			if nextOffset == 0 {
				break
			}
		}

	}

	return extentList
}

func (sysTable *SysTableSpace) GetFirstFsp() *Fsp {

	return sysTable.Fsp
}

func (sysTable *SysTableSpace) GetFirstINode() *INode {
	return sysTable.FirstInode
}

//获取数据字典
func (sysTable *SysTableSpace) GetDictTable() *DictionarySys {
	dictionarySys := NewDictionarySys()
	return dictionarySys
}

//获取数据字典叶子段
func (sysTable *SysTableSpace) GetSysDictDataTableIndexSegments() Segment {
	sysTable.GetFirstFsp()
	return nil
}

//获取数据字典数据段
func (sysTable *SysTableSpace) GetSysDictDataTableDataSegments() Segment {

	return nil
}

func (sysTable *SysTableSpace) GetSpaceId() uint32 {
	return 0
}
