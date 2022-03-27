package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/util"
	"strings"
)

/***
########################################################################################################################
**/

//用于描述非叶子节点的记录

type ClusterInternalRowHeader struct {
	innodb.FieldDataHeader
	deleteFlag bool
	minRecFlag bool   //B+树非叶子项都会增加该标记
	nOwned     uint16 //槽位中最大的值有值，该槽位其他的行皆为0
	heapNo     uint16 //表示当前记录在页面中的相对位置
	recordType uint8  //0 表示普通记录，1表示B+树非叶子节点的目录项记录，2表示Infimum，3表示Supremum
	nextRecord uint16 //表示下一条记录相对位置
	Content    []byte //长度5个字节+长度列表，都是bit
}

func NewClusterUnLeafRowHeader() innodb.FieldDataHeader {
	var clr = new(ClusterInternalRowHeader)
	clr.deleteFlag = false
	clr.minRecFlag = false
	clr.nOwned = 1
	clr.heapNo = 0
	clr.nextRecord = 0
	clr.Content = []byte{util.ConvertBits2Byte("00000000")}
	clr.Content = append(clr.Content, util.ConvertBits2Bytes("0000000000000011")...)
	clr.Content = append(clr.Content, util.ConvertBits2Bytes("0000000000000000")...)
	return clr
}

func (cldr *ClusterInternalRowHeader) SetDeleteFlag(delete bool) {
	if delete {
		cldr.Content[0] = util.ConvertValueOfBitsInBytes(cldr.Content[0], common.DELETE_OFFSET, common.COMMON_TRUE)

	} else {
		cldr.Content[0] = util.ConvertValueOfBitsInBytes(cldr.Content[0], common.DELETE_OFFSET, common.COMMON_FALSE)
	}
	cldr.deleteFlag = delete
}

func (cldr *ClusterInternalRowHeader) GetDeleteFlag() bool {
	value := util.ReadBytesByIndexBit(cldr.Content[0], common.DELETE_OFFSET)
	if value == "1" {
		return true
	} else {
		return false
	}
}
func (cldr *ClusterInternalRowHeader) GetRecMinFlag() bool {
	value := util.ReadBytesByIndexBit(cldr.Content[0], common.DELETE_OFFSET)
	if value == "1" {
		return true
	} else {
		return false
	}
}
func (cldr *ClusterInternalRowHeader) SetRecMinFlag(flag bool) {
	if flag {
		cldr.Content[0] = util.ConvertValueOfBitsInBytes(cldr.Content[0], common.MIN_REC_OFFSET, common.COMMON_TRUE)

	} else {
		cldr.Content[0] = util.ConvertValueOfBitsInBytes(cldr.Content[0], common.MIN_REC_OFFSET, common.COMMON_FALSE)
	}
	cldr.minRecFlag = flag
}
func (cldr *ClusterInternalRowHeader) SetNOwned(size byte) {
	cldr.Content[0] = util.ConvertBits2Byte(util.WriteBitsByStart(cldr.Content[0], util.TrimLeftPaddleBitString(size, 4), 4, 8))
	cldr.nOwned = uint16(size)
}
func (cldr *ClusterInternalRowHeader) GetNOwned() byte {
	return util.LeftPaddleBitString(util.ReadBytesByIndexBitByStart(cldr.Content[0], 4, 8), 4)
}
func (cldr *ClusterInternalRowHeader) GetHeapNo() uint16 {
	var heapNo = make([]string, 0)
	heapNo = append(heapNo, "0")
	heapNo = append(heapNo, "0")
	heapNo = append(heapNo, "0")
	heapNo = append(heapNo, util.ConvertByte2BitsString(cldr.Content[1])...)
	heapNo = append(heapNo, util.ConvertByte2BitsString(cldr.Content[2])[0:5]...)
	return util.ReadUB2Byte2Int(util.ConvertBits2Bytes(strings.Join(heapNo, "")))
}
func (cldr *ClusterInternalRowHeader) SetHeapNo(heapNo uint16) {
	var result = util.ConvertUInt2Bytes(heapNo)
	resultArray := util.ConvertBytes2BitStrings(result)
	//取值
	cldr.Content[1] = util.ConvertString2Byte(strings.Join(resultArray[3:11], ""))
	cldr.Content[2] = util.ConvertString2Byte(util.WriteBitsByStart(cldr.Content[2], resultArray[11:16], 0, 5))
	cldr.nOwned = uint16(heapNo)
}
func (cldr *ClusterInternalRowHeader) GetRecordType() uint8 {
	var recordType = make([]string, 0)
	recordType = append(recordType, "0")
	recordType = append(recordType, "0")
	recordType = append(recordType, "0")
	recordType = append(recordType, "0")
	recordType = append(recordType, "0")
	recordType = append(recordType, util.ConvertByte2BitsString(cldr.Content[2])[5:8]...)
	return uint8(util.ReadUB2Byte2Int(util.ConvertBits2Bytes(strings.Join(recordType, ""))))
}
func (cldr *ClusterInternalRowHeader) SetRecordType(recordType uint8) {
	resultArray := util.ConvertByte2BitsString(recordType)
	cldr.Content[2] = util.ConvertString2Byte(util.WriteBitsByStart(cldr.Content[2], resultArray[5:8], 5, 8))
	cldr.recordType = recordType
}
func (cldr *ClusterInternalRowHeader) GetNextRecord() uint16 {
	return util.ReadUB2Byte2Int(cldr.Content[3:5])
}
func (cldr *ClusterInternalRowHeader) SetNextRecord(nextRecord uint16) {
	cldr.Content[3] = util.ConvertUInt2Bytes(nextRecord)[0]
	cldr.Content[4] = util.ConvertUInt2Bytes(nextRecord)[1]
}

func (cldr *ClusterInternalRowHeader) GetRowHeaderLength() uint16 {
	return uint16(len(cldr.Content))
}

func (cldr *ClusterInternalRowHeader) ToByte() []byte {
	return cldr.Content
}

type ClusterUnLeafRowData struct {
	innodb.FieldDataValue
	PrimaryKeyMeta *TableTupleMeta
	Content        []byte
}

func NewClusterUnLeafRowData(PrimaryKeyMeta *TableTupleMeta) innodb.FieldDataValue {
	var clusterLeafRowData = new(ClusterUnLeafRowData)
	clusterLeafRowData.Content = make([]byte, 0)
	clusterLeafRowData.PrimaryKeyMeta = PrimaryKeyMeta
	return clusterLeafRowData
}
func (cld *ClusterUnLeafRowData) WriteBytesWithNull(content []byte) {
	cld.Content = util.WriteWithNull(cld.Content, content)
}

func (cld *ClusterUnLeafRowData) GetPrimaryKey() innodb.Value {
	//return cld.ReadBytesWithNullWithPosition(0)
	return nil
}

func (cld *ClusterUnLeafRowData) GetRowDataLength() uint16 {
	return uint16(len(cld.Content))
}

func (cld *ClusterUnLeafRowData) ToByte() []byte {
	return cld.Content
}

func (cld *ClusterUnLeafRowData) ReadBytesWithNullWithPosition(index int) []byte {
	//var cursor = 0
	//var i = 0

	//for {
	//	cursor, result = util.ReadBytesWithNull(cld.Content, cursor)
	//	util.ReadUB4Byte2UInt32(cld.Content)
	//	if i == index {
	//		break
	//	}
	//}
	//_,result:=util.ReadBytesWithNull(cld.Content[0+5*index:5*index+5],0)
	return cld.Content[0+5*index : 5*index+5][0:4]
}

//大致为  页面号/主键
type ClusterInternalRow struct {
	innodb.Row
	header    innodb.FieldDataHeader
	value     innodb.FieldDataValue
	TupleMeta TableTuple
}

func NewClusterUnLeafRow(meta *TableTupleMeta) innodb.Row {
	return &ClusterInternalRow{
		Row:       nil,
		header:    NewClusterUnLeafRowHeader(),
		value:     NewClusterUnLeafRowData(meta),
		TupleMeta: meta,
	}
}

func (row *ClusterInternalRow) Less(than innodb.Row) bool {

	if than.IsSupremumRow() {
		return true
	}

	if than.IsInfimumRow() {
		return false
	}

	//thanPrimaryKey := row.GetPrimaryKey()
	//thisPrimaryKey := row.GetPrimaryKey()
	//
	//switch row.PrimaryKeyMeta.PrimaryKeyType {
	//case common.COLUMN_TYPE_TINY:
	//	{
	//
	//	}
	//case common.COLUMN_TYPE_STRING:
	//	{
	//		fmt.Println(string(thanPrimaryKey))
	//		fmt.Println(string(thisPrimaryKey))
	//
	//	}
	//case common.COLUMN_TYPE_VARCHAR:
	//	{
	//
	//	}
	//case common.COLUMN_TYPE_LONG:
	//	{
	//
	//	}
	//case common.COLUMN_TYPE_INT24:
	//	{
	//		var that = util.ReadUB4Byte2UInt32(thanPrimaryKey)
	//		var this = util.ReadUB4Byte2UInt32(thisPrimaryKey)
	//		if that > this {
	//			return true
	//		} else {
	//			return false
	//		}
	//	}
	//}
	return false

}

func (row *ClusterInternalRow) ToByte() []byte {
	return append(row.header.ToByte(), row.value.ToByte()...)
}

func (row *ClusterInternalRow) WriteWithNull(content []byte) {

	row.value.WriteBytesWithNull(content)
}

func (row *ClusterInternalRow) GetRowLength() uint16 {

	return row.header.GetRowHeaderLength() + row.value.GetRowDataLength()
}

func (row *ClusterInternalRow) GetPrimaryKey() innodb.Value {
	return nil
}
func (row *ClusterInternalRow) GetPageNumber() uint32 {
	return 0
}

func (row *ClusterInternalRow) IsSupremumRow() bool {
	return false
}

func (row *ClusterInternalRow) IsInfimumRow() bool {
	return false
}
