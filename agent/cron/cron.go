package cron

import (
	"github.com/thewayma/suricataM/agent/funcs"
	"github.com/thewayma/suricataM/agent/g"
	. "github.com/thewayma/suricataM/comm/st"
	"time"
)

func PreCollect() {
	go func() {
		for {
			funcs.UpdateCpuStat()
			funcs.UpdateDiskStats()
			time.Sleep(time.Second) //!< TODO: hardcode
		}
	}()
}

func Collect() {
	if !g.Config().Transporter.Enabled {
		return
	}

	if len(g.Config().Transporter.Addrs) == 0 {
		return
	}

	for _, v := range funcs.CollectorFuncs {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*MetricData) {
	t := time.NewTicker(time.Second * time.Duration(sec)).C
	for {
		<-t
		ip := g.IP()

		metrics := []*MetricData{}

		for _, fn := range fns {
			items := fn()
			if items == nil {
				continue
			}

			if len(items) == 0 {
				continue
			}

			for _, mv := range items {
				metrics = append(metrics, mv)
			}
		}

		dt := g.Config().DefaultTags
		now := time.Now().Unix()
		for j := 0; j < len(metrics); j++ { //!< Metric, Endpoint等在GaugeValue构造填充
			metrics[j].Step = sec
			//metrics[j].Endpoint = fmt.Sprintf("%s_%s", hostname, ip)
			metrics[j].Endpoint = ip
			metrics[j].Timestamp = now

			if len(dt) > 0 { //!< Attach DefaultTags
				for k, v := range dt {
					metrics[j].Tags[k] = v
				}
			}
		}

		g.SendToTransporter(metrics)
	}
}
