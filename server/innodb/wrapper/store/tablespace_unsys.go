package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/blocks"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
	"github.com/zhukovaskychina/xmysql-server/util"
	"path"
)

//分区表实际是由多个Tablespace组成的，每个Tablespace有独立的”.ibd”文件和Space_id，
//其中”.ibd”文件的名字会以分区名加以区分，但给用户返回的是一个统一的逻辑表。
//初始化表空间的rootPage

//FSP_SIZE：表空间大小，以Page数量计算
//FSP_FREE_LIMIT：目前在空闲的Extent上最小的尚未被初始化的Page的`Page Number
//FSP_FREE：空闲extent链表，链表中的每一项为代表extent的xdes，所谓空闲extent是指该extent内所有page均未被使用
//FSP_FREE_FRAG：free frag extent链表，链表中的每一项为代表extent的xdes，所谓free frag extent是指该extent内有部分page未被使用
//FSP_FULL_FRAG：full frag extent链表，链表中的每一项为代表extent的xdes，所谓full frag extent是指该extent内所有Page均已被使用
//FSP_SEG_ID：下次待分配的segment id，每次分配新segment时均会使用该字段作为segment id，并将该字段值+1写回
//FSP_SEG_INODES_FULL：full inode page链表，链表中的每一项为inode page，该链表中的每个inode page内的inode entry都已经被使用
//FSP_SEG_INODES_FREE：free inode page链表，链表中的每一项为inode page，该链表中的每个inode page内上有空闲inode entry可分配

//Innodb的逻辑存储形式，表空间
//表空间下面分为段，区，page
//每个索引两个段，叶子段和非叶子段
//回滚段
//每个表文件都对应一个表空间
//非系统表文件加载段的时候，需要从SYS_INDEXES 中查找并且加载出来
type UnSysTableSpace struct {
	TableSpace
	conf         *conf.Cfg
	tableName    string
	spaceId      uint32
	dataBaseName string
	isSys        bool
	filePath     string
	blockFile    *blocks.BlockFile

	Fsp        *Fsp   //0 号页面
	IBuf       *IBuf  //1 号页面
	FirstInode *INode //2 号页面

	InodeMap map[uint32]*INode

	FirstIBDExtent *FirstIBDExtent //第一个区

	BtreeMaps map[string]*BTree //每个Btree中有两个段，索引段之叶子段和非叶子段

	tableMeta *TableTupleMeta //表元祖信息

}

/***
FSP HEADER PAGE

FSP header page是表空间的root page，存储表空间关键元数据信息。由page file header、fsp header、xdes entries三大部分构成

**/
func NewTableSpaceFile(cfg *conf.Cfg, databaseName string, tableName string, spaceId uint32, isSys bool) TableSpace {
	tableSpace := new(UnSysTableSpace)
	filePath := path.Join(cfg.DataDir, "/", databaseName)
	isFlag, _ := util.PathExists(filePath)
	if !isFlag {
		util.CreateDataBaseDir(cfg.DataDir, databaseName)
	}
	tableName = tableName + ".ibd"
	blockfile := blocks.NewBlockFile(filePath, tableName, 256*64*16384)
	tableSpace.blockFile = blockfile
	tableSpace.spaceId = spaceId
	tableSpace.conf = cfg
	tableSpace.isSys = isSys
	tableSpace.tableName = tableName
	tableSpace.initHeadPage()
	tableSpace.loadInodePage()
	return tableSpace
}

/**
初始化创建表文件的时候，首先会创建fsp,ibuf,inode
当重启加载的时候，可以
*/
func (tableSpace *UnSysTableSpace) initHeadPage() {

	//如果存在该文件，则是加载前3个页面
	if tableSpace.blockFile.Exists() {
		tableSpace.blockFile.OpenFile()
		fsp, err := tableSpace.blockFile.ReadPageByNumber(0)
		ibuf, err := tableSpace.blockFile.ReadPageByNumber(1)
		inode, err := tableSpace.blockFile.ReadPageByNumber(2)
		if err != nil {
			panic("err")
		}
		tableSpace.Fsp = NewFspByLoadBytes(fsp)
		tableSpace.IBuf = NewIBufByLoadBytes(ibuf)
		tableSpace.FirstInode = NewINodeByByte(inode)

		//此处需要初始化段，段目前暂时做索引段之叶子段，索引段和非叶子段

		//读取当前inode里面的段

	} else {
		tableSpace.blockFile.CreateFile()
		//初始化FspHrdPage
		fspHrdPage := pages.NewFspHrdPage(tableSpace.spaceId)
		//初始化BitMapPage
		ibuBufMapPage := pages.NewIBufBitMapPage(tableSpace.spaceId)
		//初始化INodePage
		iNodePage := pages.NewINodePage(tableSpace.spaceId)
		tableSpace.Fsp = NewFsp(fspHrdPage)
		tableSpace.IBuf = NewIBuf(tableSpace.spaceId)
		tableSpace.FirstInode = NewINode(2)
		tableSpace.blockFile.WriteContentByPage(0, fspHrdPage.GetSerializeBytes())
		tableSpace.blockFile.WriteContentByPage(1, ibuBufMapPage.GetSerializeBytes())
		tableSpace.blockFile.WriteContentByPage(2, iNodePage.GetSerializeBytes())
	}
}

func (tableSpace *UnSysTableSpace) loadInodePage() {
	tableSpace.InodeMap = make(map[uint32]*INode)
	nextPageNum := tableSpace.FirstInode.NextNodePageNumber
	tableSpace.InodeMap[2] = tableSpace.FirstInode
	for true {
		if nextPageNum == 0 {
			break
		}
		inode, _ := tableSpace.LoadPageByPageNumber(nextPageNum)
		tableSpace.InodeMap[nextPageNum] = NewINodeByByte(inode)
		nextPageNum = tableSpace.InodeMap[nextPageNum].NextNodePageNumber
	}
}

func (tableSpace *UnSysTableSpace) loadAllBtrees() {
	//
	tableSpace.BtreeMaps = make(map[string]*BTree)
	//加载clusterIndex btree

	//加载主键的段

	//tableSpace.BtreeMaps["primary"] = btree.NewTree(4, nil, nil)
	//加载secondaryIndex btree
	//for k,v:=range tableSpace.tableMeta.SecondaryIndices{
	//	btree.NewTree(4,nil,nil)
	//}

}

func (tableSpace *UnSysTableSpace) LoadPageByPageNumber(pageNumber uint32) ([]byte, error) {
	return tableSpace.blockFile.ReadPageByNumber(pageNumber)
}
