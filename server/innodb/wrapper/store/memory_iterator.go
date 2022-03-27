package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"io"
)

type MemoryIterator struct {
	innodb.RowIterator
	index        int
	currentRow   innodb.Row
	iterator     Iterator
	tmpIteraotor Iterator
}

func NewMemoryIterator(iterator Iterator) innodb.RowIterator {
	var memoryIterator = new(MemoryIterator)
	memoryIterator.iterator = iterator
	memoryIterator.index = 0
	return memoryIterator
}

func (m *MemoryIterator) Open() error {

	_, currentRow, err, tbIterator := m.iterator()
	if err != nil {
		return err
	}
	m.tmpIteraotor = tbIterator
	m.currentRow = currentRow
	return nil
}

func (m *MemoryIterator) Next() (innodb.Row, error) {
	if m.tmpIteraotor == nil {
		return nil, io.EOF
	}
	var resultError error
	_, m.currentRow, resultError, m.tmpIteraotor = m.tmpIteraotor()

	return m.currentRow, resultError
}

func (m *MemoryIterator) Close() error {
	m.tmpIteraotor = nil
	m.index = 0
	return nil
}
