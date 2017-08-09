package funcs

import (
	"fmt"
	"github.com/thewayma/suricataM/comm/st"
	"github.com/toolkits/nux"
	"sync"
)

const (
	historyCount int = 2
)

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

func UpdateCpuStat() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		return err
	}

	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		procStatHistory[i] = procStatHistory[i-1]
	}

	procStatHistory[0] = ps
	return nil
}

func deltaTotal() uint64 {
	if procStatHistory[1] == nil {
		return 0
	}
	return procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total
}

func CpuIdle() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Idle-procStatHistory[1].Cpu.Idle) * invQuotient
}

func CpuUser() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.User-procStatHistory[1].Cpu.User) * invQuotient
}

func CpuNice() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Nice-procStatHistory[1].Cpu.Nice) * invQuotient
}

func CpuSystem() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.System-procStatHistory[1].Cpu.System) * invQuotient
}

func CpuIowait() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Iowait-procStatHistory[1].Cpu.Iowait) * invQuotient
}

func CpuIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Irq-procStatHistory[1].Cpu.Irq) * invQuotient
}

func CpuSoftIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.SoftIrq-procStatHistory[1].Cpu.SoftIrq) * invQuotient
}

func CpuSteal() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Steal-procStatHistory[1].Cpu.Steal) * invQuotient
}

func CpuGuest() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Guest-procStatHistory[1].Cpu.Guest) * invQuotient
}

func CurrentCpuSwitches() uint64 {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[0].Ctxt
}

func CpuPrepared() bool {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[1] != nil
}

func CpuMetrics() []*st.MetricData {
	if !CpuPrepared() {
		return []*st.MetricData{}
	}

	cpuIdleVal := CpuIdle()
	idle := st.GaugeValue("cpu.idle", cpuIdleVal)
	busy := st.GaugeValue("cpu.busy", 100.0-cpuIdleVal)
	user := st.GaugeValue("cpu.user", CpuUser())
	nice := st.GaugeValue("cpu.nice", CpuNice())
	system := st.GaugeValue("cpu.system", CpuSystem())
	iowait := st.GaugeValue("cpu.iowait", CpuIowait())
	irq := st.GaugeValue("cpu.irq", CpuIrq())
	softirq := st.GaugeValue("cpu.softirq", CpuSoftIrq())
	steal := st.GaugeValue("cpu.steal", CpuSteal())
	guest := st.GaugeValue("cpu.guest", CpuGuest())
	switches := st.CounterValue("cpu.switches", CurrentCpuSwitches())
	return []*st.MetricData{idle, busy, user, nice, system, iowait, irq, softirq, steal, guest, switches}
}

func LoadAvgMetrics() []*st.MetricData {
	load, err := nux.LoadAvg()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return []*st.MetricData{
		st.GaugeValue("load.1min", load.Avg1min),
		st.GaugeValue("load.5min", load.Avg5min),
		st.GaugeValue("load.15min", load.Avg15min),
	}
}
