package pages

import (
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/util"
)

//go:generate mockgen -source=inode_page.go -destination ./inode_page_mock.go -package pages

///对于一个新的segment，总是优先填满32个frag page数组，之后才会为其分配完整的Extent，可以利用碎片页，并避免小表占用太多空间。
//尽量获得hint page;
//如果segment上未使用的page太多，则尽量利用segment上的page
import (
	"fmt"
)

//实现INODE
type INodePage struct {
	AbstractPage
	INodePageList DESListNode   //12 byte 存储上一个和下一个INode的页面指针 38-50
	INodeEntries  []*INodeEntry //16320 byte 用于存储具体的段信息，每个INode 192 byte，一共85 个
	EmptySpace    []byte        //6 byte

}

type FragmentArrayEntry struct {
	PageNo []byte //4个字节,用于记载离散页面号码
}

func (fragmentArray FragmentArrayEntry) Byte() []byte {
	return fragmentArray.PageNo
}

//192个字节
type INodeEntry struct {
	SegmentId           []byte               //8个字节，该结构体对应的段的编号（ID） 若值为0，则表示该SLot未被泗洪
	NotFullNUsed        []byte               //4个字节，在Notfull链表中已经使用了多少个页面
	FreeListBaseNode    []byte               //16个字节，Free链表
	NotFullListBaseNode []byte               //16个字节，NotFull链表
	FullListBaseNode    []byte               //16个字节，Full链表
	MagicNumber         []byte               //4个字节 0x5D669D2
	FragmentArrayEntry  []FragmentArrayEntry //一共32个array，每个ArrayEntry为零散的页面号
}

//每当创建一个新的索引，构建一个新的Btree，先为非叶子节点的额segment段分配一个inodeentry，再创建一个rootpage，
//并将该色门头的位置记录到rootpage中，然后再分配leafsegment的inode entry，并记录到rootpage中
func NewINodeEntry(SegmentId uint64) *INodeEntry {
	framentArray := make([]FragmentArrayEntry, 32)
	for i := 0; i < 32; i++ {
		framentArray[i] = FragmentArrayEntry{PageNo: util.AppendByte(4)}
	}

	return &INodeEntry{
		SegmentId:           util.ConvertULong8Bytes(SegmentId),
		MagicNumber:         util.ConvertUInt4Bytes(0x5D669D2),
		NotFullNUsed:        util.ConvertUInt4Bytes(0),
		FreeListBaseNode:    util.AppendByte(16),
		NotFullListBaseNode: util.AppendByte(16),
		FullListBaseNode:    util.AppendByte(16),
		FragmentArrayEntry:  framentArray,
	}
}

func (ientry *INodeEntry) GetCloseZeroFrag() int {

	var result = -1
	for i := 0; i < 32; i++ {
		if (util.ReadUB4Byte2UInt32(ientry.FragmentArrayEntry[i].PageNo)) == 0 {
			result = i
			break
		}
	}
	return result
}

//构造INode
func NewINodePage(pageNumber uint32) INodePage {
	var iPage = new(INodePage)
	var fileHeader = new(FileHeader)

	fileHeader.FilePageOffset = util.ConvertUInt4Bytes(pageNumber) //第一个byte
	//写入FSP文件头
	fileHeader.FilePageType = util.ConvertInt2Bytes(common.FILE_PAGE_INODE)
	fileHeader.FilePagePrev = util.ConvertInt4Bytes(0)
	fileHeader.FilePageOffset = util.ConvertUInt4Bytes(uint32(0))
	fileHeader.WritePageOffset(0)
	fileHeader.WritePagePrev(0)
	fileHeader.WritePageFileType(common.FILE_PAGE_INODE)
	fileHeader.WritePageNext(3)
	fileHeader.WritePageLSN(0)
	fileHeader.WritePageFileFlushLSN(0)
	fileHeader.WritePageArch(0)
	fileHeader.WritePageSpaceCheckSum(nil)
	iPage.FileHeader = *fileHeader
	iPage.INodePageList = DESListNode{
		PreNodePageNumber:  util.AppendByte(4),
		PreNodeOffset:      util.AppendByte(2),
		NextNodePageNumber: util.AppendByte(4),
		NextNodeOffSet:     util.AppendByte(2),
	}
	iPage.INodeEntries = make([]*INodeEntry, 85)
	for k, v := range iPage.INodeEntries {
		v = NewINodeEntry(0)
		iPage.INodeEntries[k] = v
	}
	iPage.FileTrailer = NewFileTrailer()
	iPage.EmptySpace = make([]byte, 6)
	return *iPage
}

func (ibuf *INodePage) SerializeBytes() []byte {
	var buff = make([]byte, 0)
	buff = append(buff, ibuf.FileHeader.GetSerialBytes()...)
	buff = append(buff, util.AppendByte(12)...)
	//SegmentId           []byte               //8个字节，该结构体对应的段的编号（ID） 若值为0，则表示该SLot未被泗洪
	//NotFullNUsed        []byte               //4个字节，在Notfull链表中已经使用了多少个页面
	//FreeListBaseNode    []byte               //16个字节，Free链表
	//NotFullListBaseNode []byte               //16个字节，NotFull链表
	//FullListBaseNode    []byte               //16个字节，Full链表
	//MagicNumber         []byte               //4个字节 0x5D669D2
	//FragmentArrayEntry  []FragmentArrayEntry //一共32个array，每个ArrayEntry为零散的页面号

	for _, v := range ibuf.INodeEntries {
		buff = append(buff, v.SegmentId...)
		buff = append(buff, v.NotFullNUsed...)
		buff = append(buff, v.FreeListBaseNode...)
		buff = append(buff, v.NotFullListBaseNode...)
		buff = append(buff, v.FullListBaseNode...)
		buff = append(buff, v.MagicNumber...)
		buff = append(buff, util.AppendByte(32*4)...)
	}
	buff = append(buff, ibuf.EmptySpace...)
	buff = append(buff, ibuf.FileTrailer.FileTrailer...)
	fmt.Println(len(buff))
	return buff
}

func (ibuf *INodePage) GetSerializeBytes() []byte {
	return ibuf.SerializeBytes()
}
