package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/segs"
	"github.com/zhukovaskychina/xmysql-server/util"
)

type InternalSegment struct {
	Segment
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

func NewInternalSegment(spaceId uint32, pageNumber uint32, offset uint16, indexName string, space TableSpace) Segment {
	var segment = new(InternalSegment)
	segment.SegmentHeader = segs.NewSegmentHeader(spaceId, pageNumber, offset)
	currentINodeBytes, _ := space.LoadPageByPageNumber(pageNumber)
	segment.IndexName = indexName
	currentInode := NewINodeByByte(currentINodeBytes)
	segment.inode = currentInode
	return segment
}

func NewInternalSegmentWithTableSpaceAtInit(spaceId uint32, pageNumber uint32, indexName string, offset uint16, space TableSpace) Segment {
	var segment = new(InternalSegment)
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

func (i *InternalSegment) AllocatePage() *Index {
	panic("implement me")
}

func (i *InternalSegment) AllocateLeafPage() *Index {
	panic("implement me")
}

func (i *InternalSegment) AllocateInternalPage() *Index {
	panic("implement me")
}

func (i *InternalSegment) AllocateNewExtent() Extent {
	panic("implement me")
}

func (i *InternalSegment) GetNotFullNUsedSize() uint32 {
	panic("implement me")
}

func (i *InternalSegment) GetFreeExtentList() *ExtentList {
	panic("implement me")
}

func (i *InternalSegment) GetFullExtentList() *ExtentList {
	panic("implement me")
}

func (i *InternalSegment) GetNotFullExtentList() *ExtentList {
	panic("implement me")
}

func (i *InternalSegment) AllocateDiscretePage(pageNumber uint32) {
	index := i.inodeEntry.GetCloseZeroFrag()
	if index == -1 {
		return
	}
	i.inodeEntry.FragmentArrayEntry = append(i.inodeEntry.FragmentArrayEntry[:index], i.inodeEntry.FragmentArrayEntry[index+1:]...)
	i.inodeEntry.FragmentArrayEntry = append(i.inodeEntry.FragmentArrayEntry, pages.FragmentArrayEntry{PageNo: util.ConvertUInt4Bytes(pageNumber)})
}

func (i *InternalSegment) GetSegmentHeader() *segs.SegmentHeader {
	return i.SegmentHeader
}

//type INodeEntry struct {
//	SegmentId           []byte               //8个字节，该结构体对应的段的编号（ID） 若值为0，则表示该SLot未被泗洪
//	NotFullNUsed        []byte               //4个字节，在Notfull链表中已经使用了多少个页面
//	FreeListBaseNode    []byte               //16个字节，Free链表
//	NotFullListBaseNode []byte               //16个字节，NotFull链表
//	FullListBaseNode    []byte               //16个字节，Full链表
//	MagicNumber         []byte               //4个字节 0x5D669D2
//	FragmentArrayEntry  []FragmentArrayEntry //一共32个array，每个ArrayEntry为零散的页面号
//}

func (i *InternalSegment) NewSegmentByBytes(bytes []byte, spaceId uint32) Segment {
	var segment = new(InternalSegment)

	segment.SegmentHeader = segs.NewSegmentHeader(spaceId, 0, 0)

	return segment
}
