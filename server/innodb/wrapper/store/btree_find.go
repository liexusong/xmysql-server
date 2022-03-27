package store

import (
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
)

func (self *BTree) Has(key innodb.Value) (has bool, err error) {
	a, i, err := self.getStart(key)
	if err != nil {
		return false, err
	}
	empty, err := self.empty(a)
	if err != nil {
		return false, err
	}
	if empty {
		return false, nil
	}
	err = self.doKey(a, i, func(akey innodb.Value) error {

		hasValue, errEqual := key.Equal(akey)
		if errEqual != nil {
			return errEqual
		}
		has = hasValue.Raw().(bool)
		return nil
	})
	if err != nil {
		return false, err
	}
	return has, nil
}

func (self *BTree) empty(a uint32) (empty bool, err error) {
	err = self.do(
		a,
		func(n *Index) error {
			empty = n.GetRecordSize() == 0
			return nil
		},
		func(n *Index) error {
			empty = n.GetRecordSize() == 0
			return nil
		},
	)
	if err != nil {
		return false, err
	}
	return empty, nil
}

func (self *BTree) Iterate() (kvi Iterator, err error) {
	return self.Range(nil, nil)
}

func (self *BTree) Values() (it RowItemsIterator, err error) {
	kvi, err := self.Iterate()
	if err != nil {
		return nil, err
	}
	it = func() (value innodb.Row, err error, _it RowItemsIterator) {
		_, value, err, kvi = kvi()
		if err != nil {
			return nil, err, nil
		}
		if kvi == nil {
			return nil, nil, nil
		}
		return value, nil, it
	}
	return it, nil
}

func doIter(run func() (Iterator, error), do func(key innodb.Value, value innodb.Row) error) error {
	kvi, err := run()
	if err != nil {
		return err
	}
	var key innodb.Value
	var value innodb.Row
	for key, value, err, kvi = kvi(); kvi != nil; key, value, err, kvi = kvi() {
		e := do(key, value)
		if e != nil {
			return e
		}
	}
	return err
}

func (self *BTree) Count(key innodb.Value) (count int, err error) {
	kvi, err := self.UnsafeRange(key, key)
	if err != nil {
		return 0, err
	}
	count = 0
	for _, _, err, kvi = kvi(); kvi != nil; _, _, err, kvi = kvi() {
		count++
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (self *BTree) UnsafeRange(from, to innodb.Value) (kvi Iterator, err error) {
	bi, err := self.rangeIterator(from, to)
	if err != nil {
		return nil, err
	}
	return self._rangeUnsafe(bi)
}

// Iterate over all of the key/values pairs between [from, to]
// inclusive. See DoIterate() for usage details.
func (self *BTree) DoRange(from, to innodb.Value, do func(key innodb.Value, value innodb.Row) error) error {
	return doIter(
		func() (Iterator, error) { return self.Range(from, to) },
		do,
	)
}

func (self *BTree) Range(from, to innodb.Value) (kvi Iterator, err error) {
	bi, err := self.rangeIterator(from, to)
	if err != nil {
		return nil, err
	}
	return self._range(bi)
}

func (self *BTree) Find(key innodb.Value) (kvi Iterator, err error) {
	return self.Range(key, key)
}
func (self *BTree) DoFind(key innodb.Value, do func(key innodb.Value, value innodb.Row) error) error {
	return doIter(
		func() (Iterator, error) { return self.Find(key) },
		do,
	)
}
