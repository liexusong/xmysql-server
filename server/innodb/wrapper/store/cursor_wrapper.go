package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/server/innodb/schemas"
)

type PrimaryKeyIndexCursor struct {
	btree *BTree

	currentRow innodb.Row

	index int
}

func NewPrimaryKeyIndexCursor(table schemas.Table) innodb.Cursor {
	//此处加载btree

	return &PrimaryKeyIndexCursor{
		btree:      nil,
		currentRow: nil,
		index:      0,
	}
}

func (p *PrimaryKeyIndexCursor) GetCurrentRow() innodb.Row {
	return p.currentRow
}

func (p *PrimaryKeyIndexCursor) Next() innodb.Row {

	//var first, _ = p.btree.GetRootNode().Scan(  Ascend, p.currentRow, nil, false, false, func(i innodb.Row) bool {
	//	if i == nil {
	//		return false
	//	}
	//	if p.index == 1 {
	//		p.index++
	//		p.currentRow = i
	//		return false
	//	}
	//	return false
	//})
	//if first {
	//	return p.currentRow
	//}

	return nil
}

func (p *PrimaryKeyIndexCursor) GetIndexPages() int {
	//return p.btree.Len()
	return 0
}

type RangeCursor struct {
	Start      innodb.Value
	End        innodb.Value
	btree      *BTree
	currentRow innodb.Row

	startRow innodb.Row
	endRow   innodb.Row
}

func (r *RangeCursor) GetCurrentRow() innodb.Row {
	panic("implement me")
}

func (r *RangeCursor) Next() innodb.Row {

	//r.btree.AscendRange(r.startRow, r.endRow, func(i innodb.Row) bool {
	//	return false
	//})
	return nil
}
