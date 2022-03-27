package store

import (
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/sqlparser"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/blocks"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/store/table"
	"github.com/zhukovaskychina/xmysql-server/util"
)

type TableTupleMeta struct {
	TableTuple
	DatabaseName   string
	TableName      string
	ColumnsMap     map[string]*FormColumnsWrapper //列
	IndexesMap     map[string]*IndexMeta
	PrimaryKeyMeta *IndexMeta
	Columns        []*FormColumnsWrapper
	Cfg            *conf.Cfg
	blockFile      *blocks.BlockFile
}

type IndexInfo struct {
	Type    string
	Name    string
	Primary bool
	Spatial bool
	Unique  bool
}

func NewTupleMeta(DatabaseName string, tableName string, Cfg *conf.Cfg) TableTuple {
	var indexMap = make(map[string]*IndexMeta)
	var columnsMap = make(map[string]*FormColumnsWrapper)
	var Columns = make([]*FormColumnsWrapper, 0)
	filePath := Cfg.DataDir + "/" + DatabaseName + "/"
	fmt.Println(filePath)
	var blockFile = blocks.NewBlockFileWithoutFileSize(filePath, tableName+".frm")

	return &TableTupleMeta{
		DatabaseName: DatabaseName,
		TableName:    tableName,
		ColumnsMap:   columnsMap,
		IndexesMap:   indexMap,
		Columns:      Columns,
		Cfg:          Cfg,
		blockFile:    blockFile,
	}
}

func (m *TableTupleMeta) GetTableName() string {
	return m.TableName
}

func (m *TableTupleMeta) GetDatabaseName() string {
	return m.DatabaseName
}

func (m *TableTupleMeta) WriteTupleColumn(columnDefs *sqlparser.ColumnDefinition) {
	frmColWrapper := NewFormColumnWrapper()
	frmColWrapper.InitializeFormWrapper(
		false,
		bool(columnDefs.Type.Autoincrement),
		bool(columnDefs.Type.NotNull),
		columnDefs.Type.Type,
		columnDefs.Name.String(),
		columnDefs.Type.Default,
		columnDefs.Type.Comment,
		int16(util.ReadUB2Byte2Int(columnDefs.Type.Length.Val)),
	)
	m.ColumnsMap[columnDefs.Name.String()] = frmColWrapper
	m.Columns = append(m.Columns, frmColWrapper)
}

func (m *TableTupleMeta) WriteTupleColumns(columnDefs []*sqlparser.ColumnDefinition) {
	for _, v := range columnDefs {
		m.WriteTupleColumn(v)
	}
}

func (m *TableTupleMeta) WriteIndexDefinitions(definitions []*sqlparser.IndexDefinition) {
	for _, v := range definitions {
		info := IndexInfo{
			Type:    v.Info.Type,
			Name:    v.Info.Name.String(),
			Primary: v.Info.Primary,
			Spatial: v.Info.Spatial,
			Unique:  v.Info.Unique,
		}
		var columns = make([]string, 0)
		for _, it := range v.Columns {
			columns = append(columns, it.Column.String())
		}

		var indexMeta = &IndexMeta{
			IndexInfo: info,
			Columns:   columns,
		}
		m.IndexesMap[v.Info.Name.String()] = indexMeta
	}
}

func (m *TableTupleMeta) ReadFrmBytes(form table.Form) {
	//读取form
	for _, v := range form.FieldBytes {
		currentFormCols := NewFormColumnWrapper()
		currentFormCols.ParseContent(v.FieldColumnsContent)
		m.Columns = append(m.Columns, currentFormCols)
		m.ColumnsMap[currentFormCols.FieldName] = currentFormCols
	}

	m.PrimaryKeyMeta = NewIndexMeta(form.ClusterIndex)

	for _, v := range form.SecondaryIndexes {
		currentIndex := NewIndexMeta(v.SecondaryIndexes)
		m.IndexesMap[currentIndex.Name] = currentIndex
	}
}

func (m *TableTupleMeta) FlushToDisk() {
	form := table.NewForm(m.DatabaseName, m.TableName)
	form.ColumnsLength = util.ConvertUInt4Bytes(uint32(len(m.Columns)))
	for _, v := range m.Columns {
		currentFields := table.FieldBytes{
			FieldColumnsOffset:  util.ConvertUInt4Bytes(uint32(len(v.ToBytes()))),
			FieldColumnsContent: v.ToBytes(),
		}
		form.FieldBytes = append(form.FieldBytes, currentFields)
	}
	form.SecondaryIndexesCount = byte(len(m.IndexesMap))
	for _, v := range m.IndexesMap {
		currentSecondIndex := table.SecondaryIndexes{
			SecondaryIndexesOffset: util.ConvertUInt4Bytes(uint32(len(v.ToBytes()))),
			SecondaryIndexes:       v.ToBytes(),
		}
		form.SecondaryIndexes = append(form.SecondaryIndexes, currentSecondIndex)
	}
	if m.PrimaryKeyMeta != nil {
		form.ClusterIndex = m.PrimaryKeyMeta.ToBytes()
		form.ClusterIndexOffSet = util.ConvertUInt4Bytes(uint32(len(m.PrimaryKeyMeta.ToBytes())))
	}

	bytes := form.ToBytes()
	m.blockFile.WriteContentToBlockFile(bytes)
}

//定义表元祖信息

type IndexMeta struct {
	IndexInfo
	Columns []string
}

func NewIndexMeta(content []byte) *IndexMeta {
	var cursor = 0
	cursor, types := util.ReadStringWithNull(content, cursor)
	cursor, indexName := util.ReadStringWithNull(content, cursor)
	cursor, primaryValue := util.ReadByte(content, cursor)
	cursor, Spatial := util.ReadByte(content, cursor)
	cursor, unique := util.ReadByte(content, cursor)
	indexInfo := IndexInfo{
		Type:    types,
		Name:    indexName,
		Primary: convertByteToBool(primaryValue),
		Spatial: convertByteToBool(Spatial),
		Unique:  convertByteToBool(unique),
	}
	var colNameArray = make([]string, 0)
	for {
		if cursor == len(content) {
			break
		}
		cursors, colName := util.ReadStringWithNull(content, cursor)
		cursor = cursors
		colNameArray = append(colNameArray, colName)
	}

	return &IndexMeta{
		IndexInfo: indexInfo,
		Columns:   colNameArray,
	}

}

//Type    string
//Name    string
//Primary bool
//Spatial bool
//Unique  bool
func (i *IndexMeta) ToBytes() []byte {
	var buff = make([]byte, 0)
	buff = append(buff, []byte(i.Type)...)
	buff = append(buff, []byte(i.Name)...)
	buff = append(buff, convertBoolToByte(i.Primary))
	buff = append(buff, convertBoolToByte(i.Spatial))
	buff = append(buff, convertBoolToByte(i.Unique))
	for _, v := range i.Columns {
		buff = append(buff, []byte(v)...)
		buff = append(buff, '0')
	}
	return buff
}
