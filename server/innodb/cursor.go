package innodb

type Cursor interface {
	GetCurrentRow() Row

	Next() Row
}
