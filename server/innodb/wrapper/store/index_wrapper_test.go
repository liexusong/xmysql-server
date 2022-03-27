package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/util"
	"testing"
)

func TestIndex_AddRows(t *testing.T) {

	tuple := NewSysTableTuple()
	index := NewPageIndexWithTuple(10, tuple)
	currentSysTableRow := NewClusterSysIndexLeafRow(tuple, false)
	initSysTableRow("test", tuple, currentSysTableRow)
	index.AddRow(currentSysTableRow)
	assert.Equal(t, index.GetRecordSize(), 3)

	content := index.ToBytes()

	assert.Equal(t, len(content), 16384)
	//currentContent := NewPageIndexByLoadBytes(content)
	currentContent := NewPageIndexByLoadBytesWithTuple(content, tuple)
	assert.Equal(t, currentContent.GetPageNumber(), index.GetPageNumber())

	assert.Equal(t, currentContent.IndexPage.PageDirectory, index.IndexPage.PageDirectory)
}

func TestIndex_AddRow(t *testing.T) {
	index := NewPageIndex(10)
	tuple := NewSysTableTuple()
	currentSysTableRow := NewClusterSysIndexLeafRow(tuple, false)
	initSysTableRow("test", tuple, currentSysTableRow)
	index.AddRow(currentSysTableRow)
}

func Test_AddRow(t *testing.T) {
	currentIndex := NewPageIndex(10)

	slotRows := NewSlotRowsWithContent(currentIndex.IndexPage.InfimumSupermum)

	fmt.Println(slotRows.FullRowList())
}

func TestIndex_AddRowsWithRow(t *testing.T) {

	tuple := NewSysTableTuple()
	index := NewPageIndexWithTuple(10, tuple)
	currentSysTableRow := NewClusterSysIndexLeafRow(tuple, false)
	initSysTableRow("test", tuple, currentSysTableRow)
	index.AddRow(currentSysTableRow)
	assert.Equal(t, index.GetRecordSize(), 1)

	content := index.ToBytes()

	assert.Equal(t, len(content), 16384)
	//currentContent := NewPageIndexByLoadBytes(content)
	currentContent := NewPageIndexByLoadBytesWithTuple(content, tuple)
	assert.Equal(t, currentContent.GetPageNumber(), index.GetPageNumber())

	assert.Equal(t, currentContent.IndexPage.PageDirectory, index.IndexPage.PageDirectory)
}

func initSysTableRow(databaseName string, tuple TableTuple, currentSysTableRow innodb.Row) {
	//rowId
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(1), 0)
	//transaction_id
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(1), 1)
	//rowpointer
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(1), 2)
	//tableId
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(1), 3)
	//tableName
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte(databaseName+"/"+tuple.GetTableName()), 4)
	//flag
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte{0, 0, 0, 0, 0, 0, 0, 0}, 5)
	//N_COLS
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(uint64(uint32(tuple.GetColumnLength()))), 6)

	//space_id
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(1), 7)

	//FileFormat
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte("Antelope"), 8)
	//RowFormat
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte("Redundant"), 9)
	//ZipPageSize
	currentSysTableRow.WriteBytesWithNullWithsPos(util.ConvertULong8Bytes(0), 10)
	//SpaceType
	currentSysTableRow.WriteBytesWithNullWithsPos([]byte("space"), 11)

	fmt.Println(currentSysTableRow.ToByte())

}

func TestdoKeyAt(t *testing.T) {

}
