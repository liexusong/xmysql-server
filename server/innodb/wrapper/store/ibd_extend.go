package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
)

/////////////////////////////////////////////////////////////
///
/// 每个Extent的大小均为1MB 256*16384 Byte
///
//////////////////////////////////////////////////////////////

type Extent interface {

	//获取区的类型，返回first和other
	ExtentType() string

	//获取
	//	GetSegmentById(segmentId uint64) *segment.Segment
	//释放页面
	FreePage(pageNumber uint32)

	AllocateNewPage(pageType int)

	GetExtentId() uint32
}

//每个区有64个page
//每次生成一个区就有64个页面
//第一个区有Fsp,BitMap,Inode，管理包括他自身的25
//Inode会有多个
type FirstIBDExtent struct {
	Extent
	fspHrdBinaryPage *Fsp   //FileSpace Header,用于存储表空间的元数据信息，
	IBufBitMapPage   *IBuf  //Insert Buffer Bookkeeping
	INodePage        *INode //Index Node Information

	//剩下的Page
	Pages []*PageWrapper

	ExtentId int64 //	extentId
}

//其他区XdesPage,BufBitMap
type OtherFirstIBDExtent struct {
	Extent
	XDesPage  *XDesPageWrapper
	IBufPages *IBuf
	Pages     []*PageWrapper

	ExtentId int64
}

//其他区
type OtherIBDExtent struct {
	Extent
	extentNumber uint32 //区的首个页面号码，会管理后面64个页面
	Pages        []*PageWrapper
}

//64 Page
func NewFirstIBDExtent(spaceId uint32, initPageNumber uint32) Extent {
	fspHrdPage := pages.NewFspHrdPage(spaceId)
	firstIBDExtent := new(FirstIBDExtent)
	firstIBDExtent.fspHrdBinaryPage = NewFsp(fspHrdPage)
	firstIBDExtent.INodePage = NewINode(initPageNumber + 2)
	firstIBDExtent.IBufBitMapPage = NewIBuf(initPageNumber + 1)
	firstIBDExtent.Pages = make([]*PageWrapper, 61)
	firstIBDExtent.ExtentId = 0
	return firstIBDExtent
}

//判断当前区是否满了，
func (firstIBDExtent *FirstIBDExtent) IsFull() bool {
	return false
}

func (firstIBDExtent *FirstIBDExtent) GetExtentId() uint32 {
	return 0
}

func (firstIBDExtent *FirstIBDExtent) FreePage(pageNumber uint32) {

}

func NewOtherFirstIBDExtent(initPageNumber uint32) Extent {
	otherIBDExtent := new(OtherFirstIBDExtent)
	otherIBDExtent.XDesPage = NewXDesWrapper(initPageNumber)
	otherIBDExtent.IBufPages = NewIBuf(initPageNumber + 1)
	otherIBDExtent.ExtentId = int64(initPageNumber >> 6)
	return otherIBDExtent
}

func NewOtherIBDExtent(extentPageNumber uint32) Extent {
	var otherIBDExtent = new(OtherIBDExtent)
	otherIBDExtent.extentNumber = extentPageNumber
	return otherIBDExtent
}
