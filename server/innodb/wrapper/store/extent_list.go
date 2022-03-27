package store

import (
	"container/list"
)

//用于构造ibd的extent链表，用于管理各种链表
type ExtentList struct {
	list           *list.List
	extentListType string
}

func NewExtentList(extentListType string) *ExtentList {
	var extentList = new(ExtentList)
	extentList.extentListType = extentListType
	extentList.list = list.New()
	return extentList
}

func (el *ExtentList) AddExtent(extent Extent) {
	el.list.PushBack(extent)
}

func (el *ExtentList) DequeFirstElement() Extent {
	element := el.list.Front()

	extent := element.Value
	el.list.Remove(element)
	return extent.(Extent)
}

func (el *ExtentList) IsEmpty() bool {
	return el.list.Len() == 0
}

type FreeExtentList struct {
	list *list.List
}

func NewFreeExtentList() *FreeExtentList {
	var extentList = new(FreeExtentList)
	extentList.list = list.New()
	return extentList
}

//func (el *FreeExtentList) AddExtent(freeExtentFirstPageNo uint32) {
//	el.list.PushBack(extent)
//}
//
//func (el *FreeExtentList) RemoveExtent(extent Extent) {
//
//}
