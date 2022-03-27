package store

import (
	"bytes"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
)
import "github.com/zhukovaskychina/xmysql-server/util"

type Fsp struct {
	fspHrdBinaryPage *pages.FspHrdBinaryPage
	xdesEntryMap     map[uint64]XDESEntryWrapper //用于存储和描述XDES的相关信息，方便操作

}

func NewFsp(fspHrdBinaryPage pages.FspHrdBinaryPage) *Fsp {
	return &Fsp{fspHrdBinaryPage: &fspHrdBinaryPage}
}

func NewFspByLoadBytes(content []byte) *Fsp {

	var fspBinary = new(pages.FspHrdBinaryPage)
	fspBinary.FileHeader = pages.NewFileHeader()
	fspBinary.FileTrailer = pages.NewFileTrailer()

	fspBinary.LoadFileHeader(content[0:38])
	fspBinary.LoadFileTrailer(content[16384-8 : 16384])
	fspBinary.EmptySpace = content[16384-8-5986 : 16384-8]
	//初始化
	fspBinary.FileSpaceHeader = pages.FileSpaceHeader{
		SpaceId:                 content[38:42],
		NotUsed:                 content[42:46],
		Size:                    content[46:50],
		FreeLimit:               content[50:54],
		SpaceFlags:              content[54:58],
		FragNUsed:               content[58:62],
		BaseNodeForFreeList:     content[62:78],
		BaseNodeForFragFreeList: content[78:94],
		BaseNodeForFullFragList: content[94:110],
		NextUnusedSegmentId:     content[110:118],
		SegFullINodesList:       content[118:134],
		SegFreeINodesList:       content[134:150],
	}
	fspBinary.XDESEntrys = make([]pages.XDESEntry, 256)
	//复盘XDESEntry，一共255个
	for k, _ := range fspBinary.XDESEntrys {
		fspBinary.XDESEntrys[k] = pages.XDESEntry{
			XDesId:       content[150+k*20 : 158+k*20],
			XDesFlstNode: content[158+k*20 : 170+k*20],
			XDesState:    content[170+k*20 : 174+k*20],
			XDesBitMap:   content[174+k*20 : 190+k*20],
		}
	}

	return NewFsp(*fspBinary)
}

//页面最小
func (fsp *Fsp) SetFreeLimit(freePageNo uint32) {
	fsp.fspHrdBinaryPage.FileSpaceHeader.FreeLimit = util.ConvertUInt4Bytes(freePageNo)
}

func (fsp *Fsp) SetXDesEntryInfo() {

}

func (fsp *Fsp) GetNextSegmentId() []byte {
	var buff = fsp.fspHrdBinaryPage.FileSpaceHeader.NextUnusedSegmentId
	fsp.ChangeNextSegmentId()
	return buff
}

func (fsp *Fsp) ChangeNextSegmentId() {
	//fsp.fspHrdBinaryPage.FileSpaceHeader.NextUnusedSegmentId
	segmentId := util.ReadUB8Byte2Long(fsp.fspHrdBinaryPage.FileSpaceHeader.NextUnusedSegmentId)
	segmentId = segmentId + 1
	fsp.fspHrdBinaryPage.FileSpaceHeader.NextUnusedSegmentId = util.ConvertULong8Bytes(segmentId)
}

//计算表空间，以Page页面数量计算
func (fsp *Fsp) GetFspSize() uint32 {
	return util.ReadUB4Byte2UInt32(fsp.fspHrdBinaryPage.FileSpaceHeader.Size)
}

//在空闲的Extent上最小的尚未被初始化的Page的PageNumber
func (fsp *Fsp) GetFspFreeLimit() uint32 {
	return util.ReadUB4Byte2UInt32(fsp.fspHrdBinaryPage.FileSpaceHeader.FreeLimit)
}

//获取FreeFrag链表中已经使用的数量
func (fsp *Fsp) GetFragNUsed() uint32 {

	return 0
}

//当一个Extent中所有page都未被使用时，放到该链表上，可以用于随后的分配
func (fsp *Fsp) GetFspFreeExtentListInfo() CommonNodeInfo {
	segFullInodeList := fsp.fspHrdBinaryPage.FileSpaceHeader.BaseNodeForFreeList
	return CommonNodeInfo{
		NodeInfoLength:     util.ReadUB4Byte2UInt32(segFullInodeList[0:4]),
		PreNodePageNumber:  util.ReadUB4Byte2UInt32(segFullInodeList[4:8]),
		PreNodeOffset:      util.ReadUB2Byte2Int(segFullInodeList[8:10]),
		NextNodePageNumber: util.ReadUB4Byte2UInt32(segFullInodeList[10:14]),
		NextNodeOffset:     util.ReadUB2Byte2Int(segFullInodeList[14:16]),
	}
}

//给fsp设置free链表信息
func (fsp *Fsp) SetFspFreeExtentListInfo(info CommonNodeInfo) {
	fsp.fspHrdBinaryPage.FileSpaceHeader.BaseNodeForFreeList = info.ToBytes()
}

//Extent中所有的page都被使用掉时，会放到该链表上，当有Page从该Extent释放时，则移回FREE_FRAG链表
func (fsp *Fsp) GetFspFullFragListInfo() CommonNodeInfo {
	segFullInodeList := fsp.fspHrdBinaryPage.FileSpaceHeader.BaseNodeForFullFragList
	return CommonNodeInfo{
		NodeInfoLength:     util.ReadUB4Byte2UInt32(segFullInodeList[0:4]),
		PreNodePageNumber:  util.ReadUB4Byte2UInt32(segFullInodeList[4:8]),
		PreNodeOffset:      util.ReadUB2Byte2Int(segFullInodeList[8:10]),
		NextNodePageNumber: util.ReadUB4Byte2UInt32(segFullInodeList[10:14]),
		NextNodeOffset:     util.ReadUB2Byte2Int(segFullInodeList[14:16]),
	}
}

//FREE_FRAG链表的Base Node，通常这样的Extent中的Page可能归属于不同的segment，用于segment frag array page的分配（见下文）
func (fsp *Fsp) GetFspFreeFrageExtentListInfo() CommonNodeInfo {
	segFullInodeList := fsp.fspHrdBinaryPage.FileSpaceHeader.BaseNodeForFragFreeList
	return CommonNodeInfo{
		NodeInfoLength:     util.ReadUB4Byte2UInt32(segFullInodeList[0:4]),
		PreNodePageNumber:  util.ReadUB4Byte2UInt32(segFullInodeList[4:8]),
		PreNodeOffset:      util.ReadUB2Byte2Int(segFullInodeList[8:10]),
		NextNodePageNumber: util.ReadUB4Byte2UInt32(segFullInodeList[10:14]),
		NextNodeOffset:     util.ReadUB2Byte2Int(segFullInodeList[14:16]),
	}
}

//segInodeFull链表的基节点
//已被完全用满的Inode Page链表
func (fsp *Fsp) GetFullINodeBaseInfo() CommonNodeInfo {
	segFullInodeList := fsp.fspHrdBinaryPage.FileSpaceHeader.SegFullINodesList
	return CommonNodeInfo{
		NodeInfoLength:     util.ReadUB4Byte2UInt32(segFullInodeList[0:4]),
		PreNodePageNumber:  util.ReadUB4Byte2UInt32(segFullInodeList[4:8]),
		PreNodeOffset:      util.ReadUB2Byte2Int(segFullInodeList[8:10]),
		NextNodePageNumber: util.ReadUB4Byte2UInt32(segFullInodeList[10:14]),
		NextNodeOffset:     util.ReadUB2Byte2Int(segFullInodeList[14:16]),
	}
}

//至少存在一个空闲Inode Entry的Inode Page被放到该链表上
func (fsp *Fsp) GetFreeSegINodeBaseInfo() CommonNodeInfo {
	segFullInodeList := fsp.fspHrdBinaryPage.FileSpaceHeader.SegFullINodesList
	return CommonNodeInfo{
		NodeInfoLength:     util.ReadUB4Byte2UInt32(segFullInodeList[0:4]),
		PreNodePageNumber:  util.ReadUB4Byte2UInt32(segFullInodeList[4:8]),
		PreNodeOffset:      util.ReadUB2Byte2Int(segFullInodeList[8:10]),
		NextNodePageNumber: util.ReadUB4Byte2UInt32(segFullInodeList[10:14]),
		NextNodeOffset:     util.ReadUB2Byte2Int(segFullInodeList[14:16]),
	}
}

func (fsp *Fsp) GetSerializeBytes() []byte {
	return fsp.fspHrdBinaryPage.GetSerializeBytes()
}

//type DESListNode struct {
//	PreNodePageNumber  []byte //4个字节	表示指向前一个INode页面号
//	PreNodeOffset      []byte //2个字节 65536-1
//	NextNodePageNumber []byte //4个字节  表示指向后一个INode页面号
//	NextNodeOffSet     []byte //2个字节	65536-1
//}
//
////XDES entry,每个Entry 占用40个字节
////一个XDES-ENtry 对应一个extent
////
//type XDESEntry struct {
//	XDesId       []byte //8 个 byte 每个段都有唯一的编号，分配段的号码
//	XDesFlstNode []byte //12 个长度 XDesEntry链表
//	XDesState    []byte //4个字节长度，根据该Extent状态信息，包括：XDES_FREE,FREE_FRAG,FULL_FRAG,FSEG
//	XDesBitMap   []byte //16个字节，一共128个bit，用两个bit表示Extent中的一个page，一个bit表示该page是否空闲的（XDES_FREE_BIT）,另一个保留位
//}

type XDESEntryWrapper struct {
	XDesId             uint64         //段ID
	PreNodePageNumber  uint32         //前一个Extent链表
	PreNodeOffset      uint16         //偏移量
	NextNodePageNumber uint32         //下一个Extent
	NextNodeoffset     uint16         //偏移量
	XdesDescPageMap    map[uint8]bool //2个bit 表示一个Page，2个表示一个page，1个表示是否空闲，1个空
}

func (xdes *XDESEntryWrapper) ToBytes() []byte {

	var buff = make([]byte, 0)
	buff = append(buff, util.ConvertULong8Bytes(xdes.XDesId)...)
	buff = append(buff, util.ConvertUInt4Bytes(xdes.PreNodePageNumber)...)
	buff = append(buff, util.ConvertUInt2Bytes(xdes.PreNodeOffset)...)
	var buffer bytes.Buffer
	for k, v := range xdes.XdesDescPageMap {
		bitKey := util.ConvertByte2Bits(k)
		buffer.WriteString(util.Substr(bitKey, 4, 8))
		if v {
			buffer.WriteString("1")
		} else {
			buffer.WriteString("0")
		}
		buffer.WriteString("0")
	}
	buff = append(buff, util.ConvertBits2Bytes(buffer.String())...)

	return buff
}
