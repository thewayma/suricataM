package tx

import (
	"github.com/thewayma/suricataM/transporter/g"
	nlist "github.com/toolkits/container/list"
)

//!< 半异步队列, 由半异步发送任务消费
func initSendQueues() {
	cfg := g.Config()
	for node := range cfg.Judge.Cluster {
		Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		JudgeQueues[node] = Q
	}
    /*
	for node, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
			GraphQueues[node+addr] = Q
		}
	}

	if cfg.Tsdb.Enabled {
		TsdbQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}
    */
}
