package g

import (
	. "github.com/thewayma/suricataM/comm/st"
	"time"
)

var HbsClient *RpcClient

func InitRpcClients() {
	if Config().Heartbeat.Enabled {
		HbsClient = &RpcClient{
			Peer:      "Agent => Heartbeat",
			RpcServer: Config().Heartbeat.Addr,
			Timeout:   time.Duration(Config().Heartbeat.Timeout) * time.Millisecond,
		}
	}
}
