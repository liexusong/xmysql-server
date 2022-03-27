package store

import "github.com/zhukovaskychina/xmysql-server/server/common"

type SysTableTuple struct {
	TableTuple
	TableName string
	Columns   []*FormColumnsWrapper
}

func (s *SysTableTuple) GetTableName() string {
	return s.TableName
}

func (s *SysTableTuple) GetDatabaseName() string {
	return "information_schema"
}

func (s *SysTableTuple) GetColumnLength() int {
	return len(s.Columns)
}

func (s *SysTableTuple) GetUnHiddenColumnsLength() int {
	var result = 0

	for _, column := range s.Columns {
		if !column.IsHidden {
			result = result + 1
		}
	}

	return result
}

func (s *SysTableTuple) GetColumnInfos(index byte) *FormColumnsWrapper {
	return s.Columns[index]
}

/**
获取可变长度变量列表
**/
func (s *SysTableTuple) GetVarColumns() []*FormColumnsWrapper {
	var formColumnsWrapperCols = make([]*FormColumnsWrapper, 0)
	for i := 0; i < len(s.Columns); i++ {
		if s.Columns[i].FieldType == "VARCHAR" {
			formColumnsWrapperCols = append(formColumnsWrapperCols, s.Columns[i])
		}
	}
	return formColumnsWrapperCols
}

func NewSysTableTuple() TableTuple {
	var sysTableTuple = new(SysTableTuple)
	sysTableTuple.TableName = common.INNODB_SYS_TABLES
	sysTableTuple.Columns = make([]*FormColumnsWrapper, 0)

	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TRX_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_POINTER",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TABLE_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "NAME",
		FieldLength:       655,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "FLAG",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "N_COLS",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "SPACE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "FILE_FORMAT",
		FieldLength:       10,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "ROW_FORMAT",
		FieldLength:       12,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ZIP_PAGE_SIZE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "SPACE_TYPE",
		FieldLength:       10,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})

	return sysTableTuple
}

func NewSysColumnsTuple() TableTuple {
	var sysTableTuple = new(SysTableTuple)
	sysTableTuple.TableName = common.INNODB_SYS_COLUMNS
	sysTableTuple.Columns = make([]*FormColumnsWrapper, 0)

	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TRX_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_POINTER",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TABLE_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "NAME",
		FieldLength:       655,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "POS",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "MTYPE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "PRTYPE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "LEN",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})

	return sysTableTuple
}

func NewSysIndexTuple() TableTuple {
	var sysTableTuple = new(SysTableTuple)
	sysTableTuple.TableName = common.INNODB_SYS_INDEXES
	sysTableTuple.Columns = make([]*FormColumnsWrapper, 0)

	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TRX_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_POINTER",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TABLE_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "NAME",
		FieldLength:       655,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "FLAG",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "N_COLS",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "SPACE",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "FILE_FORMAT",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "ROW_FORMAT",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ZIP_PAGE_SIZE",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "SPACE_TYPE",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})

	return sysTableTuple
}

func NewSysFieldsTuple() TableTuple {
	var sysTableTuple = new(SysTableTuple)
	sysTableTuple.TableName = common.INNODB_SYS_INDEXES
	sysTableTuple.Columns = make([]*FormColumnsWrapper, 0)

	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TRX_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_POINTER",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TABLE_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "NAME",
		FieldLength:       655,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "FLAG",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "N_COLS",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "SPACE",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "FILE_FORMAT",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "ROW_FORMAT",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ZIP_PAGE_SIZE",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "SPACE_TYPE",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})

	return sysTableTuple
}

func NewSysSpacesTuple() TableTuple {
	var sysTableTuple = new(SysTableTuple)
	sysTableTuple.TableName = common.INNODB_SYS_TABLESPACES
	sysTableTuple.Columns = make([]*FormColumnsWrapper, 0)

	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TRX_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_POINTER",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "SPACE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "NAME",
		FieldLength:       655,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "FLAG",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "FILE_FORMAT",
		FieldLength:       10,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "ROW_FORMAT",
		FieldLength:       22,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "PAGE_SIZE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ZIP_PAE_SIZE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "SPACE_TYPE",
		FieldLength:       10,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})

	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "FS_BLOCK_SIZE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "FILE_SIZE",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "ALLOCATED_SIZ",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	return sysTableTuple

}

//datafile文件元祖
func NewSysDataFilesTuple() TableTuple {
	var sysTableTuple = new(SysTableTuple)
	sysTableTuple.TableName = common.INNODB_SYS_DATAFILES
	sysTableTuple.Columns = make([]*FormColumnsWrapper, 0)

	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "TRX_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_ID",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          true,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "BIGINT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "ROW_POINTER",
		FieldLength:       21,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "INT",
		FieldTypeIntValue: common.COLUMN_TYPE_INT24,
		FieldName:         "SPACE",
		FieldLength:       11,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})
	sysTableTuple.Columns = append(sysTableTuple.Columns, &FormColumnsWrapper{
		IsHidden:          false,
		AutoIncrement:     false,
		NotNull:           true,
		ZeroFill:          false,
		FieldType:         "VARCHAR",
		FieldTypeIntValue: common.COLUMN_TYPE_VARCHAR,
		FieldName:         "PATH",
		FieldLength:       4000,
		FieldCommentValue: "",
		FieldDefaultValue: nil,
	})

	return sysTableTuple

}

func (s *SysTableTuple) GetPrimaryKeyColumn() *FormColumnsWrapper {
	return s.Columns[0]
}
