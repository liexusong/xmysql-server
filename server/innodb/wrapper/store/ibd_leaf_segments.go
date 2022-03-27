package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/segs"
	"github.com/zhukovaskychina/xmysql-server/util"
)

//必须先从segment内分配extent和page，创建segment核心是从inode page 中分配空闲的inode
// u
// 段空间，叶子，非叶子,rollback,undo
// 每个索引有两个segment，一个leaf，一个non-leaf
// 每个表的段数，就是索引的*2，段空间
// 得注意系统表和非系统表之间的差别
// 也就是每一个IBD对象会有多个Segments
type DataSegment struct {
	SegmentHeader     *segs.SegmentHeader
	iNodePageNo       uint32
	segmentId         uint64
	SegmentType       bool
	IndexName         string
	spaceId           uint32
	extents           []Extent //区
	fsp               *Fsp
	inode             *INode
	index             *uint32
	currentTableSpace TableSpace
	inodeEntry        *pages.INodeEntry
}

func NewLeafSegment(spaceId uint32, pageNumber uint32, offset uint16, indexName string, space TableSpace) Segment {
	var segment = new(DataSegment)
	segment.SegmentHeader = segs.NewSegmentHeader(spaceId, pageNumber, offset)
	currentINodeBytes, _ := space.LoadPageByPageNumber(pageNumber)
	segment.IndexName = indexName
	currentInode := NewINodeByByte(currentINodeBytes)
	segment.inode = currentInode
	return segment
}

func (d *DataSegment) NewSegmentByBytes(bytes []byte, spaceId uint32) Segment {
	panic("implement me")
}

func (d *DataSegment) AllocatePage() *Index {

	panic("implement me")
}

func (d *DataSegment) AllocateLeafPage() *Index {
	*d.index = (*d.index) + 1
	return NewPageIndex(*d.index)
}

func (d *DataSegment) AllocateInternalPage() *Index {
	panic("implement me")
}

func (d *DataSegment) GetSegmentHeader() *segs.SegmentHeader {
	return d.SegmentHeader
}

//叶子段，非不定长
func NewDataSegmentWithTableSpaceAtInit(spaceId uint32, pageNumber uint32, indexName string, offset uint16, space TableSpace) Segment {
	var segment = new(DataSegment)
	segment.SegmentHeader = segs.NewSegmentHeader(spaceId, pageNumber, offset)
	segment.IndexName = indexName
	segment.spaceId = spaceId
	segment.currentTableSpace = space
	segIdBytes := space.GetFirstFsp().GetNextSegmentId()
	segment.inode = space.GetFirstINode()
	inodeEntry := segment.inode.AllocateINodeEntry(util.ReadUB8Byte2Long(segIdBytes))
	segment.inodeEntry = inodeEntry
	return segment
}

func NewDataSegmentWithTableSpace(spaceId uint32, pageNumber uint32, indexName string, offset uint16, space TableSpace) Segment {
	var segment = new(DataSegment)
	segment.SegmentHeader = segs.NewSegmentHeader(spaceId, pageNumber, offset)
	segment.IndexName = indexName
	segment.spaceId = spaceId
	segment.currentTableSpace = space
	segIdBytes := space.GetFirstFsp().GetNextSegmentId()
	inodeByte, _ := space.LoadPageByPageNumber(pageNumber)
	segment.inode = NewINodeByByte(inodeByte)
	inodeEntry := segment.inode.AllocateINodeEntry(util.ReadUB8Byte2Long(segIdBytes))
	segment.inodeEntry = inodeEntry
	return segment
}

func (d *DataSegment) AllocateNewExtent() Extent {
	currentExtent := d.currentTableSpace.GetFspFreeExtentList().DequeFirstElement()
	d.inode.GetFreeExtentList()
	d.inode.SegFreeExtentMap[d.segmentId].AddExtent(currentExtent)
	return currentExtent
}

func (d *DataSegment) GetNotFullNUsedSize() uint32 {
	panic("implement me")
}

func (d *DataSegment) AllocateDiscretePage(pageNumber uint32) {
	index := d.inodeEntry.GetCloseZeroFrag()
	if index == -1 {
		return
	}
	d.inodeEntry.FragmentArrayEntry = append(d.inodeEntry.FragmentArrayEntry[:index], d.inodeEntry.FragmentArrayEntry[index+1:]...)
	d.inodeEntry.FragmentArrayEntry = append(d.inodeEntry.FragmentArrayEntry, pages.FragmentArrayEntry{PageNo: util.ConvertUInt4Bytes(pageNumber)})

}

//获取所有FreeExtent链表
func (d *DataSegment) GetFreeExtentList() *ExtentList {

	var extentList = NewExtentList("FREE_EXTENT")

	//加载inode所在的页面
	currentINodeBytes, _ := d.currentTableSpace.LoadPageByPageNumber(util.ReadUB4Byte2UInt32(d.SegmentHeader.PageNumberINodeEntry))

	currentINodeEntryOffset := util.ReadUB2Byte2Int(d.SegmentHeader.ByteOffsetINodeEntry)

	currentINodePageEntry := currentINodeBytes[currentINodeEntryOffset : currentINodeEntryOffset+192]

	freeBaseNodeInfo := currentINodePageEntry[12:28]

	//获取当前free的区链表
	nextNodePgNo := util.ReadUB4Byte2UInt32(freeBaseNodeInfo[10:14])
	//下一个offset,最大值256
	nextOffset := util.ReadUB2Byte2Int(freeBaseNodeInfo[14:16])

	//只处理当前256MB的数据,后面拓展
	for {
		//递归完成
		nextXDesEntryPage, _ := d.currentTableSpace.LoadPageByPageNumber(nextNodePgNo)
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
	return extentList
}

//获取所有FULLExtent链表
func (d *DataSegment) GetFullExtentList() *ExtentList {

	var extentList = NewExtentList("FULL_EXTENT")

	//加载inode所在的页面
	currentINodeBytes, _ := d.currentTableSpace.LoadPageByPageNumber(util.ReadUB4Byte2UInt32(d.SegmentHeader.PageNumberINodeEntry))

	currentINodeEntryOffset := util.ReadUB2Byte2Int(d.SegmentHeader.ByteOffsetINodeEntry)

	currentINodePageEntry := currentINodeBytes[currentINodeEntryOffset : currentINodeEntryOffset+192]

	freeBaseNodeInfo := currentINodePageEntry[44:60]

	//获取当前free的区链表
	nextNodePgNo := util.ReadUB4Byte2UInt32(freeBaseNodeInfo[10:14])
	//下一个offset,最大值256
	nextOffset := util.ReadUB2Byte2Int(freeBaseNodeInfo[14:16])

	//只处理当前256MB的数据,后面拓展
	for {
		//递归完成
		nextXDesEntryPage, _ := d.currentTableSpace.LoadPageByPageNumber(nextNodePgNo)
		nextXDESEntry := nextXDesEntryPage[nextOffset : nextOffset+40]
		//如果是free，则将该区加入到链表中
		if util.ReadUB4Byte2UInt32(nextXDESEntry[20:24]) == uint32(common.XDES_FULL_FRAG) {
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
	return extentList
}
func (d *DataSegment) GetNotFullExtentList() *ExtentList {

	var extentList = NewExtentList("NOT_FULL_EXTENT")

	//加载inode所在的页面
	currentINodeBytes, _ := d.currentTableSpace.LoadPageByPageNumber(util.ReadUB4Byte2UInt32(d.SegmentHeader.PageNumberINodeEntry))

	currentINodeEntryOffset := util.ReadUB2Byte2Int(d.SegmentHeader.ByteOffsetINodeEntry)

	currentINodePageEntry := currentINodeBytes[currentINodeEntryOffset : currentINodeEntryOffset+192]

	freeBaseNodeInfo := currentINodePageEntry[28:44]

	//获取当前free的区链表
	nextNodePgNo := util.ReadUB4Byte2UInt32(freeBaseNodeInfo[10:14])
	//下一个offset,最大值256
	nextOffset := util.ReadUB2Byte2Int(freeBaseNodeInfo[14:16])

	//只处理当前256MB的数据,后面拓展
	for {
		//递归完成
		nextXDesEntryPage, _ := d.currentTableSpace.LoadPageByPageNumber(nextNodePgNo)
		nextXDESEntry := nextXDesEntryPage[nextOffset : nextOffset+40]
		//如果是free，则将该区加入到链表中
		if util.ReadUB4Byte2UInt32(nextXDESEntry[20:24]) == uint32(common.XDES_FREE_FRAG) {
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
	return extentList
}
