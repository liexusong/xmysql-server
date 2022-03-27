package rows

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/wrapper/store"
)

type SysTableClusterRow struct {
	innodb.Row
	header  innodb.FieldDataHeader
	value   innodb.FieldDataValue
	FrmMeta *store.SysTableTuple
}

//系统表
func NewSysTableClusterRow(FrmMeta *store.SysTableTuple) innodb.Row {
	var row = new(SysTableClusterRow)
	row.FrmMeta = FrmMeta
	row.header = store.NewClusterLeafRowHeader(FrmMeta)
	row.value = store.NewClusterSysIndexLeafRowData(FrmMeta)
	return row
}

func (s *SysTableClusterRow) Less(than innodb.Row) bool {
	thanPk := than.GetPrimaryKey()
	thisPk := s.GetPrimaryKey()
	lessValue, _ := thanPk.LessThan(thisPk)

	return lessValue.Raw().(bool)
}

func (s *SysTableClusterRow) ToByte() []byte {
	var buff = make([]byte, 0)
	buff = append(buff, s.header.ToByte()...)
	buff = append(buff, s.value.ToByte()...)
	return buff
}

func (s *SysTableClusterRow) IsInfimumRow() bool {
	return false
}

func (s *SysTableClusterRow) IsSupremumRow() bool {
	return false
}

func (s *SysTableClusterRow) GetPageNumber() uint32 {
	panic("implement me")
}

func (s *SysTableClusterRow) WriteWithNull(content []byte) {
	s.value.WriteBytesWithNull(content)
}

func (s *SysTableClusterRow) GetRowLength() uint16 {
	panic("implement me")
}

func (s *SysTableClusterRow) GetPrimaryKey() innodb.Value {
	return NewDBRowIdWithBytes(s.value.ReadBytesWithNullWithPosition(0))
}

func (s *SysTableClusterRow) GetFieldLength() int {
	panic("implement me")
}

type SysTableInternalRow struct {
	innodb.Row
	header  innodb.FieldDataHeader
	value   innodb.FieldDataValue
	FrmMeta *store.TableTupleMeta
}

func (s SysTableInternalRow) Less(than innodb.Row) bool {
	panic("implement me")
}

func (s SysTableInternalRow) ToByte() []byte {
	panic("implement me")
}

func (s SysTableInternalRow) IsInfimumRow() bool {
	panic("implement me")
}

func (s SysTableInternalRow) IsSupremumRow() bool {
	panic("implement me")
}

func (s SysTableInternalRow) GetPageNumber() uint32 {
	panic("implement me")
}

func (s SysTableInternalRow) WriteWithNull(content []byte) {
	panic("implement me")
}

func (s SysTableInternalRow) GetRowLength() uint16 {
	panic("implement me")
}

func (s SysTableInternalRow) GetPrimaryKey() innodb.Value {
	panic("implement me")
}

func (s SysTableInternalRow) GetFieldLength() int {
	panic("implement me")
}
