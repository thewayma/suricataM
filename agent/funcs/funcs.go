package funcs

import (
	"github.com/thewayma/suricataM/agent/g"
	. "github.com/thewayma/suricataM/comm/st"
)

type FuncsAndInterval struct {
	Fs       []func() []*MetricData
	Interval int
}

var CollectorFuncs []FuncsAndInterval

func GenerateCollectorFuncs() {
	interval := g.Config().Transfer.Interval
	CollectorFuncs = []FuncsAndInterval{
		{
			Fs: []func() []*MetricData{
				//GetUptime,
				CpuMetrics,
				LoadAvgMetrics,
				MemMetrics,
				DiskIOMetrics,
				IOStatsMetrics,
			},
			Interval: interval,
		},
	}
}
