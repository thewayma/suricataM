package check

import (
	"container/list"
	"github.com/thewayma/suricataM/comm/st"
	"sync"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func (this *SafeLinkedList) ToSlice() []*st.CheckerItem {
	this.RLock()
	defer this.RUnlock()
	sz := this.L.Len()
	if sz == 0 {
		return []*st.CheckerItem{}
	}

	ret := make([]*st.CheckerItem, 0, sz)
	for e := this.L.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*st.CheckerItem))
	}
	return ret
}

// @param limit 至多返回这些，如果不够，有多少返回多少
// @return bool isEnough
func (this *SafeLinkedList) HistoryData(limit int) ([]*st.HistoryData, bool) {
	if limit < 1 {
		// 其实limit不合法，此处也返回false吧，上层代码要注意
		// 因为false通常使上层代码进入异常分支，这样就统一了
		return []*st.HistoryData{}, false
	}

	size := this.Len()
	if size == 0 {
		return []*st.HistoryData{}, false
	}

	firstElement := this.Front()
	firstItem := firstElement.Value.(*st.CheckerItem)

	var vs []*st.HistoryData
	isEnough := true

    checkerType := firstItem.Type[0]
	if checkerType == 'G' || checkerType == 'g' {
		if size < limit {
			// 有多少获取多少
			limit = size
			isEnough = false
		}
		vs = make([]*st.HistoryData, limit)
		vs[0] = &st.HistoryData{Timestamp: firstItem.Timestamp, Value: firstItem.Value}
		i := 1
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			vs[i] = &st.HistoryData{
				Timestamp: nextElement.Value.(*st.CheckerItem).Timestamp,
				Value:     nextElement.Value.(*st.CheckerItem).Value,
			}
			i++
			currentElement = nextElement
		}
	} else {
		if size < limit+1 {
			isEnough = false
			limit = size - 1
		}

		vs = make([]*st.HistoryData, limit)

		i := 0
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			diffVal := currentElement.Value.(*st.CheckerItem).Value - nextElement.Value.(*st.CheckerItem).Value
			diffTs := currentElement.Value.(*st.CheckerItem).Timestamp - nextElement.Value.(*st.CheckerItem).Timestamp
			vs[i] = &st.HistoryData{
				Timestamp: currentElement.Value.(*st.CheckerItem).Timestamp,
				Value:     diffVal / float64(diffTs),
			}
			i++
			currentElement = nextElement
		}
	}

	return vs, isEnough
}

func (this *SafeLinkedList) PushFront(v interface{}) *list.Element {
	this.Lock()
	defer this.Unlock()
	return this.L.PushFront(v)
}

func (this *SafeLinkedList) PushFrontAndMaintain(v *st.CheckerItem, maxCount int) bool {
	this.Lock()
	defer this.Unlock()

	sz := this.L.Len()
	if sz > 0 {
		// 新push上来的数据有可能重复了，或者timestamp不对，这种数据要丢掉
		if v.Timestamp <= this.L.Front().Value.(*st.CheckerItem).Timestamp || v.Timestamp <= 0 {
			return false
		}
	}

	this.L.PushFront(v)

	sz++
	if sz <= maxCount {
		return true
	}

	del := sz - maxCount
	for i := 0; i < del; i++ {
		this.L.Remove(this.L.Back())
	}

	return true
}

func (this *SafeLinkedList) Front() *list.Element {
	this.RLock()
	defer this.RUnlock()
	return this.L.Front()
}

func (this *SafeLinkedList) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.L.Len()
}
