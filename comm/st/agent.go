package st

import (
	"sync"
)

type AgentControlCommandRequest struct {
	IP string
}

type AgentControlCommandResponse struct {
	Command string
}

type AgentMetricCommandRequest struct {
	IP string
}

type cpuMetric struct {
	Idle     string `json:"idle"`
	Busy     string `json:"busy"`
	User     string `json:"user"`
	Nice     string `json:"nice"`
	System   string `json:"system"`
	Iowait   string `json:"iowait"`
	Irq      string `json:"irq"`
	Softirq  string `json:"softirq"`
	Steal    string `json:"steal"`
	Guest    string `json:"guest"`
	Switches string `json:"switches"`
}

type memMetric struct {
	Memtotal         string `json:"memtotal"`
	Memused          string `json:"memused"`
	Memfree          string `json:"memfree"`
	Swaptotal        string `json:"swaptotal"`
	Swapused         string `json:"swapused"`
	Swapfree         string `json:"swapfree"`
	Memfree_percent  string `json:"memfree_percent"`
	Memused_percent  string `json:"memused_percent"`
	Swapfree_percent string `json:"swapfree_percent"`
	Swapused_percent string `json:"swapused_percent"`
}

type diskMetric struct {
	Read_requests       string `json:"read_requests"`
	Read_merged         string `josn:"read_merged"`
	Read_sectors        string `json:"read_sectors"`
	Msec_read           string `json:"msec_read"`
	Write_requests      string `json:"write_requests"`
	Write_merged        string `json:"write_merged"`
	Write_sectors       string `json:"write_sectors"`
	Msec_write          string `json:"msec_write"`
	Ios_in_progress     string `json:"ios_in_progress"`
	Msec_total          string `json:"msec_total"`
	Msec_weighted_total string `json:"msec_weighted_total"`
	Read_bytes          string `json:"read_bytes"`
	Write_bytes         string `json:"write_bytes"`
	Avgrq_sz            string `json:"avgrq_sz"`
	Avgqu_sz            string `json:"avgqu_sz"`
	Await               string `json:"await"`
	Svctm               string `json:"svctm"`
	Util                string `json:"util"`
}

type loadMetric struct {
	Min1   string `json:"min1"`
	Mins5  string `json:"mins5"`
	Mins15 string `json:"mins15"`
}

type AgentMetricCommandResponse struct {
	Cpu  cpuMetric
	Mem  memMetric
	Disk diskMetric
	Load loadMetric
}

type ignoreMetric struct {
	sync.RWMutex
	Item map[string]bool
}

var (
	IgnoreMetric = &ignoreMetric{
		Item: make(map[string]bool),
	}

	MetricMapper = map[string]string{
		/*cpu*/
		"Idle":     "cpu.idle",
		"Busy":     "cpu.busy",
		"User":     "cpu.user",
		"Nice":     "cpu.nice",
		"System":   "cpu.system",
		"Iowait":   "cpu.iowait",
		"Irq":      "cpu.irq",
		"Softirq":  "cpu.softirq",
		"Steal":    "cpu.steal",
		"Guest":    "cpu.guest",
		"Switches": "cpu.switches",
		/*mem*/
		"Memtotal":         "mem.memtotal",
		"Memused":          "mem.memused",
		"Memfree":          "mem.memfree",
		"Swaptotal":        "mem.swaptotal",
		"Swapused":         "mem.swapused",
		"Swapfree":         "mem.swapfree",
		"Memfree_percent":  "mem.memfree.percent",
		"Memused_percent":  "mem.memused.percent",
		"Swapfree_percent": "mem.swapfree.percent",
		"Swapused_percent": "mem.swapused.percent",
		/*disk, io*/
		"Read_requests":       "disk.io.read_requests",
		"Read_merged":         "disk.io.read_merged",
		"Read_sectors":        "disk.io.read_sectors",
		"Msec_read":           "disk.io.msec_read",
		"Write_requests":      "disk.io.write_requests",
		"Write_merged":        "disk.io.write_merged",
		"Write_sectors":       "disk.io.write_sectors",
		"Msec_write":          "disk.io.msec_write",
		"Ios_in_progress":     "disk.io.ios_in_progress",
		"Msec_total":          "disk.io.msec_total",
		"Msec_weighted_total": "disk.io.msec_weighted_total",
		"Read_bytes":          "disk.io.read_bytes",
		"Write_bytes":         "disk.io.write_bytes",
		"Avgrq_sz":            "disk.io.avgrq_sz",
		"Avgqu_sz":            "disk.io.avgqu_sz",
		"Await":               "disk.io.await",
		"Svctm":               "disk.io.svctm",
		"Util":                "disk.io.util",
		/*load*/
		"Min1":   "load.1min",
		"Mins5":  "load.5min",
		"Mins15": "load.15min",
	}
)

func init() {
	IgnoreMetric.Lock()
	defer IgnoreMetric.Unlock()

	for _, k := range MetricMapper {
		IgnoreMetric.Item[k] = true
	}
}
