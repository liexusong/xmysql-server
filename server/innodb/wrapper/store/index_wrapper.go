package store

import (
	"errors"
	"fmt"
	"github.com/smartystreets/assertions"
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"math"
	"sort"
)
import "github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"
import "github.com/zhukovaskychina/xmysql-server/util"

//数据页面的包装
type Index struct {
	PageWrapper
	IndexPage   *pages.IndexPage
	SlotRowData *SlotRows
	Tuple       TableTuple //表元祖信息，用于构建row
}

//定义槽位行
type SlotRow struct {
	IsSupremum bool //是否最大
	IsInfimum  bool //是否最小
	MaxRow     innodb.Row
	Rows       innodb.Rows
	RowOffSet  uint16 //以最大值为例子，其他值之和
}

/**
如果是最小槽位，则返回1
如果是最大槽位，则返回值不能超过7
如果是中间槽位，则返回不能超过8
*/
func (sr *SlotRow) CalculateSlotSize() int {

	if sr.IsInfimum {
		return 1
	}
	if sr.IsSupremum {
		return len(sr.Rows)
	}
	return len(sr.Rows) + 1
}

type SlotRows []*SlotRow

func NewSlotRows() SlotRows {
	var slotRows = make([]*SlotRow, 0)
	ir := NewInfimumRow()
	isr := SlotRow{
		IsSupremum: false,
		IsInfimum:  true,
		MaxRow:     ir,
		Rows:       make([]innodb.Row, 0),
		RowOffSet:  0,
	}
	slotRows = append(slotRows, &isr)

	sr := NewSupremumRow()
	ssr := SlotRow{
		IsSupremum: true,
		IsInfimum:  false,
		MaxRow:     sr,
		Rows:       make([]innodb.Row, 0),
		RowOffSet:  0,
	}
	slotRows = append(slotRows, &ssr)
	return slotRows
}

func NewSlotRowsWithContent(infimum_supremum_content []byte) SlotRows {
	var slotRows = make([]*SlotRow, 0)
	ir := NewInfimumRowByContent(infimum_supremum_content[0:13])
	isr := SlotRow{
		IsSupremum: false,
		IsInfimum:  true,
		MaxRow:     ir,
		Rows:       make([]innodb.Row, 0),
		RowOffSet:  0,
	}
	slotRows = append(slotRows, &isr)

	sr := NewSupremumRowByContent(infimum_supremum_content[13:26])
	ssr := SlotRow{
		IsSupremum: true,
		IsInfimum:  false,
		MaxRow:     sr,
		Rows:       make([]innodb.Row, 0),
		RowOffSet:  0,
	}
	slotRows = append(slotRows, &ssr)
	return slotRows
}

func AddRow(sr **SlotRows, row innodb.Row) {
	//v:=len(*sr)
	psr := *sr
	index := psr.search(row)

	indexPointer := (*psr)[index]

	if indexPointer.IsSupremum && indexPointer.IsInfimum {
		if indexPointer.CalculateSlotSize() < 7 {
			//此处需要排序
			indexPointer.Rows = append(indexPointer.Rows, row)

			indexPointer.Rows = append(indexPointer.Rows, indexPointer.MaxRow)
			sort.Sort(indexPointer.Rows)
			indexPointer.MaxRow = indexPointer.Rows[len(indexPointer.Rows)-1]
			indexPointer.RowOffSet = getRowsAllength(indexPointer.Rows) + (indexPointer.MaxRow).GetRowLength()
			indexPointer.MaxRow.SetNOwned(byte(len(indexPointer.Rows)))

		} else {
			//重排
			//此处需要排序
			indexPointer.Rows = append(indexPointer.Rows, row)
			sort.Sort(indexPointer.Rows)
			ReAssignSlotRows(&psr)

		}
	}

	if indexPointer.IsSupremum {
		if indexPointer.CalculateSlotSize() < 7 {
			//此处需要排序
			indexPointer.Rows = append(indexPointer.Rows, row)
			sort.Sort(indexPointer.Rows)
			indexPointer.RowOffSet = getRowsAllength(indexPointer.Rows) + (indexPointer.MaxRow).GetRowLength()
			il := len(indexPointer.Rows)
			indexPointer.MaxRow.SetNOwned(byte(il))
		} else {
			//重排
			//此处需要排序
			indexPointer.Rows = append(indexPointer.Rows, row)
			sort.Sort(indexPointer.Rows)
			ReAssignSlotRows(&psr)

		}
	}
	*sr = psr

	//计算heapNo和
	var heapNo uint16 = 38 + 56 + 26
	for i := 0; i < (*sr).GetNDirs(); i++ {
		slotRow := (*sr).GetDirRows(uint16(i))

		for _, v := range slotRow.Rows {
			v.SetHeapNo(heapNo)
			v.SetNextRowOffset(v.GetRowLength())
			heapNo = v.GetRowLength() + heapNo
		}
		slotRow.MaxRow.SetHeapNo(heapNo)
		slotRow.RowOffSet = slotRow.MaxRow.GetHeapNo()
		heapNo = slotRow.MaxRow.GetRowLength() + heapNo
		if i == (*sr).GetNDirs()-1 {
			slotRow.MaxRow.SetNextRowOffset(0)
			break
		} else {
			slotRow.MaxRow.SetNextRowOffset(slotRow.MaxRow.GetRowLength())
		}

	}

}
func (sr *SlotRows) getSupremumMaxRow() innodb.Row {
	return (*sr)[len(*sr)-1].MaxRow
}
func (sr *SlotRows) getInfimumMaxRow() innodb.Row {
	return (*sr)[0].MaxRow
}

func (sr *SlotRows) setSupremumMaxRow(supremumRow innodb.Row) {
	(*sr)[1].MaxRow = supremumRow
}
func (sr *SlotRows) setInfimumMaxRow(infimumRow innodb.Row) {
	(*sr)[0].MaxRow = infimumRow
}

/***
当前页面，获取所有记录链表
的二进制码，
槽位的具体内容
所有的记录二进制大小


***/
func (sr *SlotRows) GetRowDataAndSlotBytes() (rowdata []byte, slotdata []byte, recordSize uint16) {
	//获取所有记录的字节数
	var buff = make([]byte, 0)
	var rows = sr.FullRowList()
	for _, v := range rows {
		buff = append(buff, (v).ToByte()...)
	}
	//计算所有的槽位
	var slotBuff = make([]byte, 0)
	var setBegin uint16 = 0
	for i := 0; i < len(*sr); i++ {
		set := sr.GetDirRows(uint16(i)).RowOffSet
		set = set + setBegin
		slotBuff = append(slotBuff, util.ConvertUInt2Bytes(set)...)
	}
	return buff, slotBuff, uint16(len(rows) - 2)
}

func (sr *SlotRows) GetNDirs() int {
	return len(*sr)
}

//根据槽位获取槽记录值
func (sr *SlotRows) GetDirRows(index uint16) *SlotRow {

	return (*sr)[index]

}

func getRowsAllength(rows []innodb.Row) uint16 {
	var offset uint16 = 0
	for _, v := range rows {
		offset = offset + (v).GetRowLength()
	}
	return offset
}

func ReAssignSlotRowsByRows(psr **SlotRows, rowList []innodb.Row) *SlotRows {
	sr := *psr

	infimumSlotRows := SlotRow{
		IsSupremum: false,
		IsInfimum:  true,
		MaxRow:     rowList[0],
		Rows:       make([]innodb.Row, 0),
	}

	infimumSlotRows.MaxRow.SetNOwned(1)
	supremumSlotRows := SlotRow{
		IsSupremum: true,
		IsInfimum:  false,
		MaxRow:     rowList[len(rowList)-1],
		Rows:       make([]innodb.Row, 0),
	}

	rowList = rowList[1 : len(rowList)-1]

	resultSR := make([]*SlotRow, 0)
	resultSR = append(resultSR, &infimumSlotRows)

	restSize := math.Mod(float64(len(rowList)), 8)

	supremumSlotRows.Rows = append(supremumSlotRows.Rows, rowList[len(rowList)-int(restSize):]...)
	supremumSlotRows.MaxRow.SetNOwned(byte(len(supremumSlotRows.Rows)))
	size := (float64(len(rowList)) - restSize) / 8

	slotRowBuffer := make([]*SlotRow, 0)
	for i := 0; i < int(size); i++ {
		sriter := &SlotRow{
			IsInfimum:  false,
			IsSupremum: false,
			MaxRow:     rowList[7+i*8],
			Rows:       rowList[0+i*8 : 7+i*8],
			RowOffSet:  getRowsAllength(rowList[0:7+i*8]) + (rowList[7+i*8]).GetRowLength(),
		}
		sriter.MaxRow.SetNOwned(byte(len(sriter.Rows)))
		slotRowBuffer = append(slotRowBuffer, sriter)
	}
	resultSR = append(resultSR, slotRowBuffer...)
	resultSR = append(resultSR, &supremumSlotRows)

	sr = (*SlotRows)(&resultSR)
	*psr = sr
	return sr
}

//重新计算槽位对应的记录
//第一个和最后一个是最大和最小

func ReAssignSlotRows(psr **SlotRows) *SlotRows {
	sr := *psr
	rowList := sr.FullRowList()
	infimumSlotRows := SlotRow{
		IsSupremum: false,
		IsInfimum:  true,
		MaxRow:     rowList[0],
		Rows:       make([]innodb.Row, 0),
	}

	supremumSlotRows := SlotRow{
		IsSupremum: true,
		IsInfimum:  false,
		MaxRow:     rowList[len(rowList)-1],
		Rows:       make([]innodb.Row, 0),
	}

	rowList = rowList[1 : len(rowList)-1]

	resultSR := make([]*SlotRow, 0)
	resultSR = append(resultSR, &infimumSlotRows)

	restSize := math.Mod(float64(len(rowList)), 8)

	supremumSlotRows.Rows = append(supremumSlotRows.Rows, rowList[len(rowList)-int(restSize):]...)

	size := (float64(len(rowList)) - restSize) / 8

	slotRowBuffer := make([]*SlotRow, 0)
	for i := 0; i < int(size); i++ {
		sriter := &SlotRow{
			IsInfimum:  false,
			IsSupremum: false,
			MaxRow:     rowList[7+i*8],
			Rows:       rowList[0+i*8 : 7+i*8],
			RowOffSet:  getRowsAllength(rowList[0:7+i*8]) + (rowList[7+i*8]).GetRowLength(),
		}
		slotRowBuffer = append(slotRowBuffer, sriter)
	}
	resultSR = append(resultSR, slotRowBuffer...)
	resultSR = append(resultSR, &supremumSlotRows)

	sr = (*SlotRows)(&resultSR)

	sr.FullRowList()

	*psr = sr
	return sr
}

func (sr *SlotRows) FullRowList() []innodb.Row {
	var rows = make([]innodb.Row, 0)

	for _, v := range *sr {
		rows = append(rows, v.Rows...)
		rows = append(rows, v.MaxRow)
	}
	return rows
}

func (sr *SlotRows) GetMaxRows() innodb.Row {
	size := len(sr.FullRowList())
	return sr.FullRowList()[size-2]
}

func (sr *SlotRows) GetRowListWithoutInfiuAndSupremum() []innodb.Row {
	var rows = make([]innodb.Row, 0)

	for _, v := range *sr {
		rows = append(rows, v.Rows...)
		rows = append(rows, v.MaxRow)
	}

	return rows[1 : len(rows)-1]
}

//根据槽位返回数组下标
func (sr *SlotRows) search(row innodb.Row) int {

	lowIndex := 0
	highIndex := len(*sr) - 1
	result := 0
	for lowIndex < highIndex {

		mid := (lowIndex + highIndex) >> 1

		maxRowpointer := (*sr)[mid].MaxRow

		if lowIndex == highIndex-1 {
			result = highIndex
			break
		}

		if (maxRowpointer).Less(row) {
			lowIndex = mid
		} else if row.Less(maxRowpointer) {
			highIndex = mid
		} else {
			highIndex = mid
		}
	}
	return result
}

func NewPageIndex(pageNumber uint32) *Index {

	var slotRowData = NewSlotRows()
	return &Index{
		IndexPage:   pages.NewIndexPage(pageNumber, 0),
		SlotRowData: &slotRowData,
	}
}

func NewPageIndexWithTuple(pageNumber uint32, tuple TableTuple) *Index {

	var slotRowData = NewSlotRows()
	return &Index{
		IndexPage:   pages.NewIndexPage(pageNumber, 0),
		SlotRowData: &slotRowData,
		Tuple:       tuple,
	}
}

func NewPageIndexByLoadBytes(content []byte) *Index {

	var indexPage = new(pages.IndexPage)
	indexPage.FileHeader = pages.NewFileHeader()
	indexPage.FileTrailer = pages.NewFileTrailer()

	indexPage.LoadFileHeader(content[0:38])
	indexPage.LoadFileTrailer(content[16384-8 : 16384])

	indexPage.ParsePageHeader(content[38 : 38+56])
	indexPage.ParseInfimumSupermum(content[38+56 : 38+56+26])
	indexPage.ParsePageSlots(content)
	indexPage.ParseUserRecordsAndFreeSpace(content)

	var index = new(Index)
	index.IndexPage = indexPage

	indexPageType := index.GetIndexPageType()
	slotRows := NewSlotRowsWithContent(indexPage.InfimumSupermum)
	index.SlotRowData = &slotRows
	index.parseIndexBytes2SlotRows(content, indexPageType)

	return &Index{
		IndexPage:   indexPage,
		SlotRowData: &slotRows,
	}
}

func NewPageIndexByLoadBytesWithTuple(content []byte, tuple TableTuple) *Index {

	var indexPage = new(pages.IndexPage)
	indexPage.FileHeader = pages.NewFileHeader()
	indexPage.FileTrailer = pages.NewFileTrailer()

	indexPage.LoadFileHeader(content[0:38])
	indexPage.LoadFileTrailer(content[16384-8 : 16384])

	indexPage.ParsePageHeader(content[38 : 38+56])
	indexPage.ParseInfimumSupermum(content[38+56 : 38+56+26])
	indexPage.ParsePageSlots(content)
	indexPage.ParseUserRecordsAndFreeSpace(content)

	var index = new(Index)
	index.IndexPage = indexPage
	index.Tuple = tuple
	indexPageType := index.GetIndexPageType()
	slotRows := NewSlotRowsWithContent(indexPage.InfimumSupermum)
	index.SlotRowData = &slotRows
	index.parseIndexBytes2SlotRows(content, indexPageType)

	return index
}

func (i *Index) GetPageNo() []byte {
	return i.IndexPage.FileHeader.FilePageOffset
}

func (i *Index) GetPageNumber() uint32 {
	return util.ReadUB4Byte2UInt32(i.IndexPage.FileHeader.FilePageOffset)
}

//叶子段
func (i *Index) SetPageBtrTop(btrtop []byte) {
	i.IndexPage.PageHeader.PageBtrSegLeaf = btrtop
}

//PageBtrSegTop
//非叶子段
func (i *Index) SetPageBtrSegs(btrsegs []byte) {
	i.IndexPage.PageHeader.PageBtrSegTop = btrsegs
}

//索引页的跟页面记录两个结构
//获取Inode的表空间号，INode的页面号，以及Inode的entry offset
//返回spaceId，inode page id,inode entry offset
func (i *Index) GetCurrentINodePageNumber() (uint32, uint32, uint16) {

	var buff = i.IndexPage.PageHeader.PageBtrSegTop

	return util.ReadUB4Byte2UInt32(buff[0:4]), util.ReadUB4Byte2UInt32(buff[4:8]), util.ReadUB2Byte2Int(buff[8:10])
}

//获取页面类型
func (i *Index) GetFilePageType() uint16 {
	return util.ReadUB2Byte2Int(i.IndexPage.FileHeader.FilePageType)
}

//根据26个字段判断当前页面属于叶子还是非叶子节点
func (i *Index) GetIndexPageType() string {

	firstByte := i.IndexPage.InfimumSupermum[0]
	return util.ConvertByte2BitsString(firstByte)[3]
}

//增加记录
func (i *Index) AddRow(row innodb.Row) {
	//初始化
	if i.GetRecordSize() == 0 {

		//解析infimumRow
		infimumRow := NewInfimumRowByContent(i.IndexPage.InfimumSupermum[0:13])
		//解析supreRow
		supremumRow := NewSupremumRowByContent(i.IndexPage.InfimumSupermum[13:26])

		infimumRow.SetHeapNo(38 + 56 + 26)
		infimumRow.SetNextRowOffset(13)
		supremumRow.SetHeapNo(38 + 56 + 26 + 13)

		i.SlotRowData.setInfimumMaxRow(infimumRow)

		i.SlotRowData.setSupremumMaxRow(supremumRow)

	}

	AddRow(&i.SlotRowData, row)
	infimumSupremumBytes := make([]byte, 0)
	infimumSupremumBytes = append(infimumSupremumBytes, i.SlotRowData.getInfimumMaxRow().ToByte()...)
	infimumSupremumBytes = append(infimumSupremumBytes, i.SlotRowData.getSupremumMaxRow().ToByte()...)
	i.IndexPage.InfimumSupermum = infimumSupremumBytes
	rowData, slotData, recordSize := i.SlotRowData.GetRowDataAndSlotBytes()

	i.IndexPage.PageHeader.PageNRecs = util.ConvertUInt2Bytes(recordSize)
	i.IndexPage.PageHeader.PageNDirSlots = util.ConvertUInt2Bytes(uint16(i.SlotRowData.GetNDirs()))
	i.IndexPage.UserRecords = rowData
	i.IndexPage.PageDirectory = slotData

	i.IndexPage.FreeSpace = util.AppendByte(common.PAGE_SIZE -
		common.PAGE_FILE_HEADER_SIZE -
		common.PAGE_PAGE_HEADER_SIZE -
		common.PAGE_INFIMUMSUPERUM_SIZE -
		common.PAGE_FILE_TRAILER_SIZE -
		len(rowData) -
		len(slotData))
}

func (i *Index) AddRows(rows []innodb.Row) {
	for _, v := range rows {
		i.AddRow(v)
	}
}

//判断是否还有足够空间插入
func (i *Index) IsFullParams(row innodb.Row) bool {
	size := len(i.IndexPage.FreeSpace)
	return len(row.ToByte()) < size
}

func (i *Index) IsFull(row innodb.Row) bool {

	var rest = common.PAGE_SIZE - common.PAGE_FILE_HEADER_SIZE -
		common.PAGE_FILE_TRAILER_SIZE - common.PAGE_PAGE_HEADER_SIZE - common.PAGE_INFIMUMSUPERUM_SIZE
	//len(i.IndexPage.UserRecords)+len(i.IndexPage.PageDirectory)

	return rest-len(i.IndexPage.UserRecords)-len(i.IndexPage.PageDirectory) < int(row.GetRowLength())
}

func (i *Index) Delete(key innodb.Row, row innodb.Row) {

}

//根据字节码转换行
//需要根据页面类型，判断是否是leaf还是internal
//
func (i *Index) parseIndexBytes2SlotRows(content []byte, pageType string) {
	//获取当前页面的所有记录
	recordSize := i.GetRecordSize()

	//槽位数量
	//slotCnt := i.GetSlotNDirs()
	//

	//userRecords := i.IndexPage.UserRecords

	//解析infimumRow
	infimumRow := NewInfimumRowByContent(i.IndexPage.InfimumSupermum[0:13])
	//解析supreRow
	supremumRow := NewSupremumRowByContent(i.IndexPage.InfimumSupermum[13:26])
	//
	//
	//
	//AddRow(&i.SlotRowData, infimumRow)
	//AddRow(&i.SlotRowData, supremumRow)
	//根据infimumRow可以获取到下一条记录的位置

	//获取

	if recordSize == 0 {
		return
	}
	nextOffset := infimumRow.GetNextRowOffset()

	var startOffset uint16 = infimumRow.GetHeapNo()

	var endOffset uint16 = supremumRow.GetHeapNo()
	for nextOffset != 0 {

		startOffset = startOffset + nextOffset

		if startOffset == endOffset {
			break
		}
		if startOffset == uint16(16384) {
			break
		}

		//获得了当前记录的开始字节段
		prepareContent := content[startOffset:]
		//需要解析出当前的不定长字段的大小，以及

		var currentRow innodb.Row
		if pageType == common.PAGE_INTERNAL {
			fmt.Println("=================================" + pageType)
		} else {
			switch currentTuple := i.Tuple.(type) {

			case *SysTableTuple:
				{
					currentRow = NewClusterSysIndexLeafRowWithContent(prepareContent, currentTuple)
				}
			case *TableTupleMeta:
				{
					currentRow = NewClusterLeafRow(prepareContent, currentTuple)
				}
			default:
				{
					fmt.Println("default")
				}
			}

		}
		if currentRow == nil {
			break
		}
		nextOffset = currentRow.GetNextRowOffset()
		AddRow(&i.SlotRowData, currentRow)
	}

	//slotDirSizeArrays := i.getSlotDirs()
	//
	//
	//
	//for i := 0; i < len(slotDirSizeArrays); i++ {
	//	//特殊值处理
	//	if i == 0 {
	//		slotDirSizeArrays[0:13]
	//		NewInfimumRow()
	//	}
	//
	//}
}

func (i *Index) PageLeafOrInternal() string {

	flags := util.ConvertByte2BitsString(i.IndexPage.InfimumSupermum[0])[3]
	return flags
}

func (i *Index) GetRowByIndex(index int) (row innodb.Row, found bool) {
	var resultRow innodb.Row
	found = true
	if index <= 0 {
		return nil, false
	}
	if index <= len(i.SlotRowData.GetRowListWithoutInfiuAndSupremum()) {
		resultRow = i.SlotRowData.GetRowListWithoutInfiuAndSupremum()[index-1]
		found = true
	} else {
		found = false
	}

	return resultRow, found

}

//平衡
func (i *Index) BalancePage(index *Index) error {
	recordSize := i.GetRecordSize()
	rows, _ := i.GetRowsByIndex(recordSize + 1)

	index.AddRows(rows)
	index.AddRow(NewSupremumRow())
	i.TruncateByIndex(recordSize)

	return nil
}

func (i *Index) GetRowsByIndex(index int) (rows []innodb.Row, found bool) {

	for k, v := range i.SlotRowData.FullRowList() {
		if k == index {
			break
		}
		rows = append(rows, v)
	}
	return rows, false
}

func (i *Index) Truncate() {
	i.IndexPage.UserRecords = make([]byte, 0)
	i.IndexPage.PageHeader.PageNRecs = util.ConvertUInt2Bytes(0)
	i.IndexPage.PageDirectory = make([]byte, 0)
	i.IndexPage.FreeSpace = util.AppendByte(common.PAGE_SIZE - common.PAGE_FILE_HEADER_SIZE - common.PAGE_INFIMUMSUPERUM_SIZE - common.PAGE_FILE_TRAILER_SIZE)
}

func (i *Index) TruncateByIndex(index int) {
	rows := i.SlotRowData.FullRowList()
	rows = rows[0:index]
	ReAssignSlotRowsByRows(&i.SlotRowData, rows)
	rowData, slotData, recordSize := i.SlotRowData.GetRowDataAndSlotBytes()
	i.IndexPage.PageHeader.PageNRecs = util.ConvertUInt2Bytes(recordSize)
	i.IndexPage.UserRecords = rowData
	i.IndexPage.PageDirectory = slotData
	i.IndexPage.FreeSpace = util.AppendByte(common.PAGE_SIZE - common.PAGE_FILE_HEADER_SIZE - common.PAGE_INFIMUMSUPERUM_SIZE - common.PAGE_FILE_TRAILER_SIZE - len(rowData) - len(slotData))

}

//根据Key值查找
//如果没有则返回false，同时返回该非叶子记录的行，里面包括了，子页面的页面号
func (i *Index) Find(rows innodb.Row) (row innodb.Row, found bool) {
	fullList := i.SlotRowData.GetRowListWithoutInfiuAndSupremum()
	if rows == nil {
		return fullList[0], false
	}

	index := sort.Search(len(fullList), func(i int) bool {
		return rows.Less(fullList[i])
	})

	if index > 0 && !(fullList[index-1]).Less(rows) {
		return fullList[index-1], true
	}

	return fullList[index-1], false
}

//根据Key值查找
//如果没有则返回false，同时返非叶子记录的行，里面包括了，子页面的页面号
//暂时搁置二分查找逻辑
func (i *Index) FindByKey(targetKey innodb.Value) (row innodb.Row, found bool) {
	fullList := i.SlotRowData.GetRowListWithoutInfiuAndSupremum()
	if targetKey == nil {
		return fullList[0], false
	}
	if len(fullList) == 0 {
		return nil, false
	}

	//index := sort.Search(len(fullList), func(i int) bool {
	//	return rows.Less(*fullList[i])
	//})
	//
	//if index > 0 && !(*fullList[index-1]).Less(rows) {
	//	return *fullList[index-1], true
	//}

	return fullList[0], false
}

func (i *Index) FindReturnIndex(rows innodb.Row) (rowIndex int, found bool) {

	fullList := i.SlotRowData.GetRowListWithoutInfiuAndSupremum()

	index := sort.Search(len(fullList), func(i int) bool {
		return rows.Less(fullList[i])
	})

	if index > 0 && !(fullList[index-1]).Less(rows) {
		return index - 1, true
	}

	return index, false
}

func (i *Index) GetRecordByIndex(index int) innodb.Row {
	fullList := i.SlotRowData.GetRowListWithoutInfiuAndSupremum()
	return fullList[index-1]
}

///获取所有记录
func (i *Index) GetRecordSize() int {
	return int(util.ReadUB2Byte2Int(i.IndexPage.PageHeader.PageNRecs))
}

//获取槽位的大小
func (i *Index) GetSlotNDirs() int {
	return int(util.ReadUB2Byte2Int(i.IndexPage.PageHeader.PageNDirSlots))
}

//根据槽位的偏移量获取字节
func (i *Index) getSlotMaxRowSize(currentIdx int) uint16 {
	slotDirsCnt := i.GetSlotNDirs()
	assertions.ShouldBeLessThan(currentIdx, slotDirsCnt)
	directory := i.getSlotDirs()
	return directory[currentIdx]
}

//获取槽位下标
func (i *Index) getSlotDirs() []uint16 {
	pageDirectory := i.IndexPage.PageDirectory
	var buff = make([]uint16, 0)
	for i := 0; i < (len(pageDirectory) / 2); i = i + 2 {
		buff = append(buff, util.ReadUB2Byte2Int(pageDirectory[i:i+2]))
	}
	return buff
}

//获取第几个记录
func (n *Index) doValueAt(i int, do func(row innodb.Row) error) error {
	value := n.GetRecordByIndex(i)
	return do(value)
}

func (n *Index) doKeyAt(i int, do func(key innodb.Value) error) error {
	row, _ := n.GetRowByIndex(i)
	if row == nil {
		return errors.New("没有记录")
	}
	return do(row.GetPrimaryKey())
}

func (n *Index) GetPrePageNo() uint32 {
	return util.ReadUB4Byte2UInt32(n.IndexPage.FileHeader.FilePagePrev)
}

func (n *Index) GetNextPageNo() uint32 {
	return util.ReadUB4Byte2UInt32(n.IndexPage.FileHeader.FilePageNext)
}

func (n *Index) SetPrePageNo(prePageNo uint32) {
	n.IndexPage.FileHeader.FilePagePrev = util.ConvertUInt4Bytes(prePageNo)
}

func (n *Index) SetNextPageNo(nextPageNo uint32) {
	n.IndexPage.FileHeader.FilePageNext = util.ConvertUInt4Bytes(nextPageNo)
}

func (n *Index) GetMiniumRecord() innodb.Row {

	return n.GetRecordByIndex(0)
}

//ph.PageBtrSegLeaf = buff[36:46]  //B+树中叶子节点端的头部信息，尽在B+树中的跟页面中定义
//ph.PageBtrSegTop = buff[46:56]   //B+树中非叶子节点端的头部信息，尽在B+树中的跟页面中定义
func (i *Index) GetSegLeaf() []byte {
	return i.IndexPage.PageHeader.PageBtrSegLeaf
}

func (i *Index) SetSegLeaf(segLeaf []byte) {
	i.IndexPage.PageHeader.PageBtrSegLeaf = segLeaf
}
func (i *Index) SetSegTop(segTop []byte) {
	i.IndexPage.PageHeader.PageBtrSegTop = segTop
}
func (i *Index) GetSegTop() []byte {
	return i.IndexPage.PageHeader.PageBtrSegTop
}

func (i *Index) ToBytes() []byte {
	return i.IndexPage.GetSerializeBytes()
}
