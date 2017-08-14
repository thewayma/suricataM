package rx

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	. "github.com/thewayma/suricataM/comm/log"
	. "github.com/thewayma/suricataM/comm/st"
	"github.com/thewayma/suricataM/transporter/g"
	"github.com/thewayma/suricataM/transporter/tx"
)

type Transfer struct{}

func RpcServer() {
	if !g.Config().Rpc.Enabled {
		return
	}

	addr := g.Config().Rpc.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		Log.Info("transporter <= agent, net.ListenTCP.Addr=%v", addr)
	}

	server := rpc.NewServer()
	server.Register(new(Transfer))

	for {
		conn, err := listener.Accept()
		if err != nil {
			Log.Error("transporter <= agent, listener.Accept occur error:", err)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (t *Transfer) Update(args []*MetricData, reply *TransporterResponse) error {
	return RecvMetric(args, reply, "rpc")
}

func RecvMetric(items []*MetricData, reply *TransporterResponse, from string) error {
	start := time.Now()
	reply.Invalid = 0

	Log.Trace("Transporter <= Agent, Total=%d, MetricData[0]=%v", len(items), items[0])

	//!< sanity check已前移至agent上
	cfg := g.Config()
	if cfg.Checker.Enabled {
		tx.Push2CheckerSendQueue(items)
	}

	if cfg.Tsdb.Enabled {
		tx.Push2TsdbSendQueue(items)
	}

	reply.Message = "ok"
	reply.Total = len(items)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	Log.Trace("Transporter => Agent, TransferResp=%v", reply)

	return nil
}
