package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
)
import "github.com/zhukovaskychina/xmysql-server/util"

//用于管理数据文件中的segment,用于存储各种INodeEntry
//第三个page的类型FIL_PAGE_INODE
//每个Inode页面可以存储85个记录，
//PreNodePageNumber  []byte //4个字节	表示指向前一个INode页面号
//PreNodeOffset      []byte //2个字节 65536-1
//NextNodePageNumber []byte //4个字节  表示指向后一个INode页面号
//NextNodeOffSet     []byte //2个字节	65536-1
type INode struct {
	INodePage          pages.INodePage
	SegMap             map[uint64]*INodeEntryWrapper
	ts                 TableSpace
	PreNodePageNumber  uint32
	PreNodeOffset      uint16
	NextNodePageNumber uint32
	NextNodeOffset     uint16

	SegFreeExtentMap map[uint64]*ExtentList
	SegFullExtentMap map[uint64]*ExtentList
	SegNotExtentMap  map[uint64]*ExtentList
}

//NotFullNUsed        []byte               //4个字节，在Notfull链表中已经使用了多少个页面
//FreeListBaseNode    []byte               //16个字节，Free链表 segment上所有page均空闲的extent链表
//NotFullListBaseNode []byte               //16个字节，NotFull链表 至少有一个page分配给当前Segment的Extent链表，全部用完时，转移到FSEG_FULL上，全部释放时，则归还给当前表空间FSP_FREE链表
//FullListBaseNode    []byte               //16个字节，Full链表 segment上page被完全使用的extent链表
//MagicNumber         []byte               //4个字节 0x5D669D2
//FragmentArrayEntry  []FragmentArrayEntry //一共32个array，每个ArrayEntry为零散的页面号
//
type INodeEntryWrapper struct {
	SegmentId          uint64
	NotFullInUsedPages uint32

	FreeListLength          uint32
	FreeFirstNodePageNumber uint32
	FreeFirstNodeOffset     uint16
	FreeLastNodePageNumber  uint32
	FreeLastNodeOffset      uint16

	NotFullListLength          uint32
	NotFullFirstNodePageNumber uint32
	NotFullFirstNodeOffset     uint16
	NotFullLastNodePageNumber  uint32
	NotFullLastNodeOffset      uint16

	FullListLength          uint32
	FullFirstNodePageNumber uint32
	FullFirstNodeOffset     uint16
	FullLastNodePageNumber  uint32
	FullLastNodeOffset      uint16

	fragmentArray []uint32
}

func NewINode(pageNo uint32) *INode {
	ipage := pages.NewINodePage(pageNo)
	return &INode{
		INodePage: ipage,
	}
}
func NewINodeByByte(content []byte) *INode {
	var inodePage = new(pages.INodePage)
	inodePage.FileHeader = pages.NewFileHeader()
	inodePage.FileTrailer = pages.NewFileTrailer()

	inodePage.LoadFileHeader(content[0:38])
	inodePage.LoadFileTrailer(content[16384-8 : 16384])

	//PreNodePageNumber  []byte //4个字节	表示指向前一个INode页面号
	//PreNodeOffset      []byte //2个字节 65536-1
	//NextNodePageNumber []byte //4个字节  表示指向后一个INode页面号
	//NextNodeOffSet     []byte //2个字节	65536-1
	inodePage.INodePageList.PreNodePageNumber = content[38:42]
	inodePage.INodePageList.PreNodeOffset = content[42:44]
	inodePage.INodePageList.NextNodePageNumber = content[44:48]
	inodePage.INodePageList.NextNodeOffSet = content[48:50]

	inodePage.EmptySpace = content[16384-8-6 : 16384-8]

	inodePage.INodeEntries = make([]*pages.INodeEntry, 85)

	//NotFullNUsed        []byte               //4个字节，在Notfull链表中已经使用了多少个页面
	//FreeListBaseNode    []byte               //16个字节，Free链表 segment上所有page均空闲的extent链表
	//NotFullListBaseNode []byte               //16个字节，NotFull链表 至少有一个page分配给当前Segment的Extent链表，全部用完时，转移到FSEG_FULL上，全部释放时，则归还给当前表空间FSP_FREE链表
	//FullListBaseNode    []byte               //16个字节，Full链表 segment上page被完全使用的extent链表
	//MagicNumber         []byte               //4个字节 0x5D669D2
	//FragmentArrayEntry  []FragmentArrayEntry //一共32个array，每个ArrayEntry为零散的页面号
	//for k, v := range inodePage.INodeEntries {
	//	v.FragmentArrayEntry = make([]pages.FragmentArrayEntry, 32)
	//	inodePage.INodeEntries[k] = &pages.INodeEntry{
	//		SegmentId:           content[50+192*k : 58+192*k],
	//		NotFullNUsed:        content[58+192*k : 62+192*k],
	//		FreeListBaseNode:    content[62+192*k : 78+192*k],
	//		NotFullListBaseNode: content[78+192*k : 94+192*k],
	//		FullListBaseNode:    content[94+192*k : 110+192*k],
	//		MagicNumber:         content[110+192*k : 114+192*k],
	//		FragmentArrayEntry:  parseFragmentArray(content[114+192*k : 242+192*k]),
	//	}
	//}
	for k := 0; k < 85; k++ {
		//	FragmentArrayEntry := make([]pages.FragmentArrayEntry, 32)
		inodePage.INodeEntries[k] = &pages.INodeEntry{
			SegmentId:           content[50+192*k : 58+192*k],
			NotFullNUsed:        content[58+192*k : 62+192*k],
			FreeListBaseNode:    content[62+192*k : 78+192*k],
			NotFullListBaseNode: content[78+192*k : 94+192*k],
			FullListBaseNode:    content[94+192*k : 110+192*k],
			MagicNumber:         content[110+192*k : 114+192*k],
			FragmentArrayEntry:  parseFragmentArray(content[114+192*k : 242+192*k]),
		}
	}
	return &INode{
		INodePage: *inodePage,
	}
}

func parseFragmentArray(content []byte) []pages.FragmentArrayEntry {
	var buff = make([]pages.FragmentArrayEntry, 0)
	for i := 0; i < 32; i++ {
		buff = append(buff, pages.FragmentArrayEntry{PageNo: content[i*4 : i*4+4]})
	}
	return buff
}

func (iNode *INode) AllocateINodeEntry(segmentId uint64) *pages.INodeEntry {
	nodeEntry := pages.NewINodeEntry(segmentId)
	var index = iNode.getCloseZeroSeg()
	INodeEntries := append(iNode.INodePage.INodeEntries[:index], nodeEntry)
	INodeEntries = append(INodeEntries, iNode.INodePage.INodeEntries[index+1:]...)
	iNode.INodePage.INodeEntries = INodeEntries
	return nodeEntry
}

func (iNode *INode) getCloseZeroSeg() int {

	var result = 0
	for i := 0; i < 85; i++ {
		if (util.ReadUB8Byte2Long(iNode.INodePage.INodeEntries[i].SegmentId)) == 0 {
			result = i
			break
		}
	}
	return result
}

//根据
func (iNode *INode) GetInodeEntryBySegmentId(segmentId uint64) (*pages.INodeEntry, bool) {

	for _, v := range iNode.INodePage.INodeEntries {
		flags := util.ReadUB8Byte2Long(v.SegmentId) == segmentId
		if flags {
			return v, true
		}
	}
	return nil, false
}

func (iNode *INode) GetINodeRootPageBySegId(segmentId uint64) (uint32, bool) {
	for _, v := range iNode.INodePage.INodeEntries {
		flags := util.ReadUB8Byte2Long(v.SegmentId) == segmentId
		if flags {
			return util.ReadUB4Byte2UInt32(v.FragmentArrayEntry[0].Byte()), true
		}
	}

	return 0, false
}

//获取所有FreeExtent链表
func (iNode *INode) GetFreeExtentList() {
	var SegFreeExtentMap = make(map[uint64]*ExtentList)
	for k, v := range iNode.SegMap {
		var extentList = NewExtentList("FREE_EXTENT")
		//获取当前free的区链表
		nextNodePgNo := v.FreeLastNodePageNumber
		//下一个offset,最大值256
		nextOffset := v.FreeLastNodeOffset

		//只处理当前256MB的数据,后面拓展
		for {
			//递归完成
			nextXDesEntryPage, _ := iNode.ts.LoadPageByPageNumber(nextNodePgNo)
			nextXDESEntry := nextXDesEntryPage[nextOffset : nextOffset+40]
			//如果是free，则将该区加入到链表中
			if util.ReadUB4Byte2UInt32(nextXDESEntry[20:24]) == uint32(common.XDES_FREE) {
				currentPageNodeOffsetNo := (nextOffset - 150) / 40
				u := uint32(currentPageNodeOffsetNo) + nextNodePgNo
				extentList.AddExtent(NewOtherIBDExtent(u))
				nextNodePgNo = util.ReadUB4Byte2UInt32(nextXDESEntry[14:18])
				nextOffset = util.ReadUB2Byte2Int(nextXDESEntry[18:20])
				if nextOffset == 0 {
					break
				}
			}
		}
		SegFreeExtentMap[k] = extentList
	}
	iNode.SegFreeExtentMap = SegFreeExtentMap
}

//获取所有FULLExtent链表
func (iNode *INode) GetFullExtentList() {
	var SegFreeExtentMap = make(map[uint64]*ExtentList)
	for k, v := range iNode.SegMap {

		var extentList = NewExtentList("FULL_EXTENT")
		//获取当前free的区链表
		nextNodePgNo := v.FreeLastNodePageNumber
		//下一个offset,最大值256
		nextOffset := v.FreeLastNodeOffset

		//只处理当前256MB的数据,后面拓展
		for {
			//递归完成
			nextXDesEntryPage, _ := iNode.ts.LoadPageByPageNumber(nextNodePgNo)
			nextXDESEntry := nextXDesEntryPage[nextOffset : nextOffset+40]
			//如果是free，则将该区加入到链表中
			if util.ReadUB4Byte2UInt32(nextXDESEntry[20:24]) == uint32(common.XDES_FREE) {
				currentPageNodeOffsetNo := (nextOffset - 150) / 40
				u := uint32(currentPageNodeOffsetNo) + nextNodePgNo
				extentList.AddExtent(NewOtherIBDExtent(u))
				nextNodePgNo = util.ReadUB4Byte2UInt32(nextXDESEntry[14:18])
				nextOffset = util.ReadUB2Byte2Int(nextXDESEntry[18:20])
				if nextOffset == 0 {
					break
				}
			}
		}
		SegFreeExtentMap[k] = extentList
	}
	iNode.SegFullExtentMap = SegFreeExtentMap
}
func (iNode *INode) GetNotFullExtentList() {
	var SegFreeExtentMap = make(map[uint64]*ExtentList)
	for k, v := range iNode.SegMap {

		var extentList = NewExtentList("NOT_FULL_EXTENT")
		//获取当前free的区链表
		nextNodePgNo := v.FreeLastNodePageNumber
		//下一个offset,最大值256
		nextOffset := v.FreeLastNodeOffset

		//只处理当前256MB的数据,后面拓展
		for {
			//递归完成
			nextXDesEntryPage, _ := iNode.ts.LoadPageByPageNumber(nextNodePgNo)
			nextXDESEntry := nextXDesEntryPage[nextOffset : nextOffset+40]
			//如果是free，则将该区加入到链表中
			if util.ReadUB4Byte2UInt32(nextXDESEntry[20:24]) == uint32(common.XDES_FREE) {
				currentPageNodeOffsetNo := (nextOffset - 150) / 40
				u := uint32(currentPageNodeOffsetNo) + nextNodePgNo
				extentList.AddExtent(NewOtherIBDExtent(u))
				nextNodePgNo = util.ReadUB4Byte2UInt32(nextXDESEntry[14:18])
				nextOffset = util.ReadUB2Byte2Int(nextXDESEntry[18:20])
				if nextOffset == 0 {
					break
				}
			}
		}
		SegFreeExtentMap[k] = extentList
	}
	iNode.SegNotExtentMap = SegFreeExtentMap
}

func (iNode *INode) GetSerializeBytes() []byte {
	return iNode.INodePage.GetSerializeBytes()
}

//
//SegmentId           []byte               //8个字节，该结构体对应的段的编号（ID） 若值为0，则表示该SLot未被泗洪
//NotFullNUsed        []byte               //4个字节，在Notfull链表中已经使用了多少个页面
//FreeListBaseNode    []byte               //16个字节，Free链表
//NotFullListBaseNode []byte               //16个字节，NotFull链表
//FullListBaseNode    []byte               //16个字节，Full链表
//MagicNumber         []byte               //4个字节 0x5D669D2
//FragmentArrayEntry  []FragmentArrayEntry //一共32个array，每个ArrayEntry为零散的页面号

func (iNode *INode) GetSegmentByOffset(offset uint16, segmentType int) Segment {

	//entry := iNode.INodePage.INodeEntries[offset]
	//currentINodeEntry:=iNode.INodePage.INodeEntries[offset]

	return nil
}
