package g

import (
	. "github.com/thewayma/suricataM/comm/st"
	"sync"
	"time"
)

type SafeStrategyMap struct {
	sync.RWMutex
	M map[string][]Strategy
}

func (this *SafeStrategyMap) ReInit(m map[string][]Strategy) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeStrategyMap) Get() map[string][]Strategy {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

type SafeEventMap struct {
	sync.RWMutex
	M map[string]*Event
}

func (this *SafeEventMap) Get(key string) (*Event, bool) {
	this.RLock()
	defer this.RUnlock()
	event, exists := this.M[key]
	return event, exists
}

func (this *SafeEventMap) Set(key string, event *Event) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = event
}

var (
	HbsClient   *RpcClient
	StrategyMap = &SafeStrategyMap{M: make(map[string][]Strategy)}
	LastEvents  = &SafeEventMap{M: make(map[string]*Event)}
)

func InitHbsClient() {
	HbsClient = &RpcClient{
		Peer:      "Checker => HeartBeat",
		RpcServer: Config().Hbs.Server,
		Timeout:   time.Duration(Config().Hbs.Timeout) * time.Millisecond,
	}
}
