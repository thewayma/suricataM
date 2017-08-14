package tx

import (
	. "github.com/thewayma/suricataM/comm/pool"
	"github.com/thewayma/suricataM/transporter/g"
	"github.com/toolkits/container/set"
)

func initConnPools() {
	cfg := g.Config()

	checkerInstances := set.NewStringSet()
	for _, instance := range cfg.Checker.Cluster {
		checkerInstances.Add(instance)
	}
	CheckerConnPools = CreateSafeRpcConnPools(cfg.Checker.MaxConns, cfg.Checker.MaxIdle,
		cfg.Checker.ConnTimeout, cfg.Checker.CallTimeout, checkerInstances.ToSlice())

	/*
		// tsdb
		if cfg.Tsdb.Enabled {
			TsdbConnPoolHelper = NewTsdbConnPoolHelper(cfg.Tsdb.Address, cfg.Tsdb.MaxConns, cfg.Tsdb.MaxIdle, cfg.Tsdb.ConnTimeout, cfg.Tsdb.CallTimeout)
		}
	*/

}

func DestroyConnPools() {
	CheckerConnPools.Destroy()
	//TsdbConnPoolHelper.Destroy()
}
