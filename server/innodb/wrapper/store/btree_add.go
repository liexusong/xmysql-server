package store

import (
	"errors"
	"fmt"
	"github.com/zhukovaskychina/xmysql-server/server/common"
	"github.com/zhukovaskychina/xmysql-server/server/innodb"
	"github.com/zhukovaskychina/xmysql-server/util"
	"log"
)

func (self *BTree) Add(key innodb.Value, value innodb.Row) (err error) {

	if key == nil {
		return errors.New("key 为null")
	}

	_, _, err = self.add(self.rootPageNo, key, value)
	if err != nil {
		return err
	}
	//self.meta.itemCount += cntDelta
	//self.meta.root = root
	return nil
}

func (self *BTree) add(root uint32, key innodb.Value, value innodb.Row) (cntDelta, newRoot uint32, err error) {
	a, b, err := self.insert(root, key, value)
	if err != nil {
		return 0, 0, err
	} else if b == 0 {
		return 1, a, nil
	}

	//newRoot, err = self.newInternal()
	//if err != nil {
	//	return 0, 0, err
	//}
	//err = self.doInternal(newRoot, func(nPage *Index) error {
	//	err := self.firstKey(a, func(akey []byte) error {
	//		return nPage.putKP(self.varchar, akey, a)
	//	})
	//	if err != nil {
	//		return err
	//	}
	//	return self.firstKey(b, func(bkey []byte) error {
	//		return n.putKP(self.varchar, bkey, b)
	//	})
	//})
	//if err != nil {
	//	return 0, 0, err
	//}
	return 1, newRoot, nil
}

/* right is only set on split left is always set.
 * - When split is false left is the pointer to block
 * - When split is true left is the pointer to the new left block
 */
func (self *BTree) insert(n uint32, key innodb.Value, value innodb.Row) (a, b uint32, err error) {
	var leafOrInternal string
	err = self.blockFile.Do(n, func(content []byte) error {
		//此处需要判断是否index
		if util.ReadUB2Byte2Int(content[24:26]) == common.FILE_PAGE_INDEX {
			//判断是叶子还是非叶子
			//currentIndex := NewPageIndexByLoadBytes(content)
			currentIndex := NewPageIndexByLoadBytesWithTuple(content, self.Tuple)
			if currentIndex.PageLeafOrInternal() == common.PAGE_LEAF {
				leafOrInternal = common.PAGE_LEAF
			} else {
				leafOrInternal = common.PAGE_INTERNAL
			}
			return nil
		} else {
			return errors.New("非Index页面")
		}
	})
	if err != nil {
		return 0, 0, err
	}
	if leafOrInternal == common.PAGE_INTERNAL {
		return self.internalInsert(n, key, value)
	} else {
		return self.leafInsert(n, key, value)
	}
}

/* - first find the child to insert into
 * - do the child insert
 * - if there was a split:
 *    - if the block is full, split this block
 *    - else insert the new key/pointer into this block
 */
func (self *BTree) internalInsert(n uint32, key innodb.Value, value innodb.Row) (a, b uint32, err error) {

	var ptr uint32
	//尽可能的查找子节点
	err = self.doInternal(n, func(nIndex *Index) (err error) {

		currentRow, found := nIndex.FindByKey(key)

		if found {

		}

		ptr = currentRow.GetPageNumber()

		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	//下一层节点可能分裂后，产生的页面号
	//两种情况
	//一，没有分裂过的页面
	//二，分裂过的页面
	p, q, err := self.insert(ptr, key, value)
	if err != nil {
		return 0, 0, err
	}
	var must_split bool = false

	//加载非叶子节点
	//当前非叶子节点，需要将生成的子页面号+key值组合，落到
	err = self.doInternal(n, func(m *Index) error {
		//*m.ptr(i) = p
		//err := self.firstKey(p, func(key innodb.Value) error {
		//	return m.updateK(self.varchar, i, key)
		//})
		//if err != nil {
		//	return err
		//}

		//获取p页面的第一个记录的值
		//需要更新
		self.firstKey(p, func(key innodb.Value) error {

			return nil
		})
		//如果分裂了子页面，子页面可能是叶子也可能是非叶子节点
		if q != 0 {
			return self.firstKey(q, func(key innodb.Value) error {
				if m.IsFull(nil) {
					must_split = true
				}
				//n 页面需要将q页面的页面号和key值加入到n的记录

				return nil
			})
		}
		return nil
	})
	if err != nil {
		self.doInternal(n, func(n *Index) (err error) {
			//	log.Println(n.Debug(self.varchar))
			return nil
		})
		log.Printf("n: %v, p: %v, q: %v", n, p, q)
		log.Println(err)
		return 0, 0, err
	}
	if must_split {
		a, b, err = self.internalSplit(n)
		if err != nil {
			return 0, 0, err
		}
	} else {
		a = n
		b = 0
	}
	return a, b, nil
}

func (self *BTree) leafInsert(n uint32, key innodb.Value, value innodb.Row) (a, b uint32, err error) {

	var mustSplit bool = false

	err = self.doLeaf(n, func(nIndex *Index) error {

		if nIndex.GetRecordSize() <= 0 {
			nIndex.AddRow(value)
			self.blockFile.WriteFileBySeekStart(uint64(n*common.PAGE_SIZE), nIndex.IndexPage.GetSerializeBytes())
			return nil
		}
		//去重判断
		_, found := nIndex.FindByKey(key)
		if found {
			return errors.New("主键重复")
		}

		if nIndex.IsFull(value) {
			mustSplit = true
		} else {
			nIndex.AddRow(value)
		}
		fmt.Println(len(nIndex.IndexPage.SerializeBytes()))
		self.blockFile.WriteFileBySeekStart(uint64(n*common.PAGE_SIZE), nIndex.IndexPage.GetSerializeBytes())
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	if mustSplit {
		return self.leafSplit(n, key, value)
	}
	return n, 0, nil
}

/***
   先分裂当前叶子节点，
	然后待插入值，判断左右，如果小于b的最小值，那么则插入a节点
***/
func (self *BTree) leafSplit(n uint32, key innodb.Value, value innodb.Row) (a, b uint32, err error) {

	a = n

	//申请叶子页面
	bIndex := self.dataSegment.AllocateLeafPage()

	b = bIndex.GetPageNumber()

	err = self.doLeaf(a, func(aIndex *Index) (err error) {
		aIndex.SetNextPageNo(b)
		aIndex.BalancePage(bIndex)
		if value.Less(bIndex.GetMiniumRecord()) {
			aIndex.AddRow(value)
		} else {
			bIndex.AddRow(value)
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}

/* 非叶子节点分裂
 * - first assert that the key to be inserted is not already in the block.
 * - Make a new block
 * - balance the two blocks.
 * - insert the new key/pointer combo into the correct block
 *
 * Note. in the varchar case, the key is not the key but a pointer to a
 * key. This complicates the bytes.Compare line significantly.
 */
func (self *BTree) internalSplit(n uint32) (a, b uint32, err error) {
	// log.Println("internalSplit", n, key)
	a = n
	//重新申请页面
	//
	bIndex := self.indexSegment.AllocateInternalPage()
	b = bIndex.GetPageNumber()
	err = self.doInternal(a, func(n *Index) error {
		return self.doInternal(b, func(m *Index) (err error) {
			n.BalancePage(m)
			return nil
		})
	})
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}
