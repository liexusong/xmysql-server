package innodb

//定义迭代
type RowIterator interface {
	Open() error

	Next() (Row, error)

	Close() error
}
