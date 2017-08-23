package g

import (
	. "github.com/thewayma/suricataM/comm/log"
	. "github.com/thewayma/suricataM/comm/st"
	"math/rand"
	"sync"
	"time"
)

var (
	ClientsLock *sync.RWMutex = new(sync.RWMutex)
	Clients                   = make(map[string]*RpcClient)
)

func initTransporterClient(addr string) *RpcClient {
	var c *RpcClient = &RpcClient{
		Peer:      "Agent => Transporter",
		RpcServer: addr,
		Timeout:   time.Duration(Config().Transporter.Timeout) * time.Millisecond,
	}
	ClientsLock.Lock()
	defer ClientsLock.Unlock()
	Clients[addr] = c

	return c
}

func getTransporterClient(addr string) *RpcClient {
	ClientsLock.RLock()
	defer ClientsLock.RUnlock()

	if c, ok := Clients[addr]; ok {
		return c
	}
	return nil
}

func updateMetrics(c *RpcClient, metrics []*MetricData, resp *TransporterResponse) bool {
	err := c.Call("Transporter.Update", metrics, resp)
	if err != nil {
		Log.Error("Agent => Transporter Transporter.Update RPC fail, Rpc Client:%v, Error Code:%s", c, err)
		return false
	}
	return true
}

func SendMetrics(metrics []*MetricData, resp *TransporterResponse) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(Config().Transporter.Addrs)) {
		addr := Config().Transporter.Addrs[i]

		c := getTransporterClient(addr)
		if c == nil {
			c = initTransporterClient(addr)
		}

		if updateMetrics(c, metrics, resp) {
			break
		}
	}
}

func SendToTransporter(m []*MetricData) {
	if len(m) == 0 {
		return
	}

	Log.Trace("Agent => Transporter, Total=%d, MetricData[0]=%v\n", len(m), m[0])

	var resp TransporterResponse
	SendMetrics(m, &resp)

	Log.Trace("Agent <= Transporter, TransporterResponse=%v", resp)
}
