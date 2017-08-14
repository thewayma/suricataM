package check

import (
	"container/list"
	"github.com/thewayma/suricataM/comm/st"
	"sync"
)

type CheckerItemMap struct {
	sync.RWMutex
	M map[string]*SafeLinkedList //!< 因为policychecker需要检测过去N个监控项, 故这个状态维护在此SafeLinkList上, 以PK()作为map的key索引
}

func NewCheckerItemMap() *CheckerItemMap {
	return &CheckerItemMap{M: make(map[string]*SafeLinkedList)}
}

func (this *CheckerItemMap) Get(key string) (*SafeLinkedList, bool) {
	this.RLock()
	defer this.RUnlock()
	val, ok := this.M[key]
	return val, ok
}

func (this *CheckerItemMap) Set(key string, val *SafeLinkedList) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = val
}

func (this *CheckerItemMap) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.M)
}

func (this *CheckerItemMap) Delete(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, key)
}

func (this *CheckerItemMap) BatchDelete(keys []string) {
	count := len(keys)
	if count == 0 {
		return
	}

	this.Lock()
	defer this.Unlock()
	for i := 0; i < count; i++ {
		delete(this.M, keys[i])
	}
}

func (this *CheckerItemMap) CleanStale(before int64) {
	keys := []string{}

	this.RLock()
	for key, L := range this.M {
		front := L.Front()
		if front == nil {
			continue
		}

		if front.Value.(*st.CheckerItem).Timestamp < before {
			keys = append(keys, key)
		}
	}
	this.RUnlock()

	this.BatchDelete(keys)
}

func (this *CheckerItemMap) PushFrontAndMaintain(key string, val *st.CheckerItem, maxCount int, now int64) {
	if linkedList, exists := this.Get(key); exists {
		needCheck := linkedList.PushFrontAndMaintain(val, maxCount)
		if needCheck {
			Checker(linkedList, val, now)
		}
	} else {
		NL := list.New()
		NL.PushFront(val)
		safeList := &SafeLinkedList{L: NL}
		this.Set(key, safeList)
		Checker(safeList, val, now)
	}
}

var HistoryBigMap = make(map[string]*CheckerItemMap)

func InitHistoryBigMap() {
	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			HistoryBigMap[arr[i]+arr[j]] = NewCheckerItemMap()
		}
	}
}

func init() {
	InitHistoryBigMap()
}
