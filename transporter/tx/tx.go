package tx

import (
	. "github.com/thewayma/suricataM/comm/log"
	. "github.com/thewayma/suricataM/comm/pool"
	. "github.com/thewayma/suricataM/comm/st"
	"github.com/thewayma/suricataM/transporter/g"
	rings "github.com/toolkits/consistent/rings"
	nlist "github.com/toolkits/container/list"
)

const (
	DefaultSendQueueMaxSize = 102400 //10.24w
)

var (
	MinStep int //最小上报周期,单位sec
)

// 服务节点的一致性哈希环
// pk -> node
var (
	CheckerNodeRing *rings.ConsistentHashNodeRing
)

// 发送缓存队列
// node -> queue_of_data
var (
	CheckerQueues = make(map[string]*nlist.SafeListLimited)
	//TsdbQueue   *nlist.SafeListLimited
)

// 连接池
// node_address -> connection_pool
var (
	CheckerConnPools *SafeRpcConnPools
	//TsdbConnPools      *TsdbConnPoolHelper
)

func Start() {
	MinStep = g.Config().MinStep
	if MinStep < 1 {
		MinStep = 30
	}

	initConnPools()
	initSendQueues()
	initNodeRings()
	startSendTasks()
	Log.Trace("Tx.RPC Init Finished")
}

// 将数据 打入 某个Checker的发送缓存队列
func Push2CheckerSendQueue(items []*MetricData) {
	for _, item := range items {
		pk := item.PK()
		node, err := CheckerNodeRing.GetNode(pk)
		if err != nil {
			Log.Error("CheckerNodeRing.GetNode, pk=%s, err=%s", err)
			continue
		}

		// align ts
		step := int(item.Step)
		if step < MinStep {
			step = MinStep
		}
		ts := alignTs(item.Timestamp, int64(step))

		checkerItem := &CheckerItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: ts,
			Type:      item.Type,
			Tags:      item.Tags,
		}
		Q := CheckerQueues[node]
		isSuccess := Q.PushFront(checkerItem)

		// statistics
		if !isSuccess {
			//proc.SendToCheckerDropCnt.Incr()
		}
	}
}

/*
// 将原始数据入到tsdb发送缓存队列
func Push2TsdbSendQueue(items []*MetricData) {
	for _, item := range items {
		tsdbItem := convert2TsdbItem(item)
		isSuccess := TsdbQueue.PushFront(tsdbItem)

		if !isSuccess {
			//proc.SendToTsdbDropCnt.Incr()
		}
	}
}

// 转化为tsdb格式
func convert2TsdbItem(d *MetricData) *TsdbItem {
	t := TsdbItem{Tags: make(map[string]string)}
	for k, v := range d.Tags {
		t.Tags[k] = v
	}
	t.Tags["endpoint"] = d.Endpoint
	t.Metric = d.Metric
	t.Timestamp = d.Timestamp
	t.Value = d.Value
	return &t
}
*/

func alignTs(ts int64, period int64) int64 {
	return ts - ts%period
}
