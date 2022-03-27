package store

type TableTuple interface {
	GetTableName() string

	GetDatabaseName() string

	GetColumnLength() int

	//获取非隐藏列
	GetUnHiddenColumnsLength() int

	GetColumnInfos(index byte) *FormColumnsWrapper

	//获取可变列链表
	GetVarColumns() []*FormColumnsWrapper

	//根据列下标，计算出可变列表的下标
	//
	GetVarDescribeInfoIndex(index byte) byte
}
