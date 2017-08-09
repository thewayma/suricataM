package tx

import (
	."github.com/thewayma/suricataM/comm/pool"
	"github.com/thewayma/suricataM/transporter/g"
	"github.com/toolkits/container/set"
)

func initConnPools() {
	cfg := g.Config()

	// judge
	judgeInstances := set.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}
	JudgeConnPools = CreateSafeRpcConnPools(cfg.Judge.MaxConns, cfg.Judge.MaxIdle,
		cfg.Judge.ConnTimeout, cfg.Judge.CallTimeout, judgeInstances.ToSlice())

        /*
	// tsdb
	if cfg.Tsdb.Enabled {
		TsdbConnPoolHelper = NewTsdbConnPoolHelper(cfg.Tsdb.Address, cfg.Tsdb.MaxConns, cfg.Tsdb.MaxIdle, cfg.Tsdb.ConnTimeout, cfg.Tsdb.CallTimeout)
	}

	// graph
	graphInstances := set.NewSafeSet()
	for _, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			graphInstances.Add(addr)
		}
	}
	GraphConnPools = CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
		cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphInstances.ToSlice())
        */

}

func DestroyConnPools() {
	JudgeConnPools.Destroy()
	//GraphConnPools.Destroy()
	//TsdbConnPoolHelper.Destroy()
}
