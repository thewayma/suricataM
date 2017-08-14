package rx

import (
	"github.com/thewayma/suricataM/checker/check"
	"github.com/thewayma/suricataM/checker/g"
	. "github.com/thewayma/suricataM/comm/log"
	"github.com/thewayma/suricataM/comm/st"
	"log"
	"net"
	"net/rpc"
	"time"
)

type Checker struct{}

func (this *Checker) Send(items []*st.CheckerItem, resp *st.SimpleRpcResponse) error {
	remain := g.Config().Remain
	now := time.Now().Unix()

	Log.Trace("Checker <= Transporter, Len=%d, CheckerItem[0]=%v", len(items), items[0])

	for _, item := range items {
		pk := item.PrimaryKey()
		check.HistoryBigMap[pk[0:2]].PushFrontAndMaintain(pk, item, remain, now)
	}
	return nil
}

func RpcStart() {
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
		Log.Trace("Checker <= Transporter, rpc listening=%v", addr)
	}

	rpc.Register(new(Checker))

	for {
		conn, err := listener.Accept()
		if err != nil {
			Log.Error("checker <= transporter, listener.Accept occur error: %s", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
