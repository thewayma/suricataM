package tx

import (
	"github.com/thewayma/suricataM/transporter/g"
	nlist "github.com/toolkits/container/list"
)

//!< 半异步队列, 由半异步发送任务消费
func initSendQueues() {
	cfg := g.Config()
	for node := range cfg.Checker.Cluster {
		Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		CheckerQueues[node] = Q
	}

	if cfg.InfluxDB.Enabled {
		InfluxDBQueue = nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
	}
}
