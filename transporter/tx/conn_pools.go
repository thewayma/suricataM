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

	CheckerConnPools = CreateSafeRpcConnPools(cfg.Checker.MaxConcurrentConns, cfg.Checker.MaxIdle,
		cfg.Checker.ConnTimeout, cfg.Checker.CallTimeout, checkerInstances.ToSlice())
}

func DestroyConnPools() {
	CheckerConnPools.Destroy()
}
