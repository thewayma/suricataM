package cron

import (
	"fmt"
	"time"

	"github.com/thewayma/suricataM/agent/funcs"
	"github.com/thewayma/suricataM/agent/g"
	. "github.com/thewayma/suricataM/comm/log"
	. "github.com/thewayma/suricataM/comm/st"
)

type AgentReportRequest struct {
	Hostname        string
	IP              string
	AgentVersion    string
	SuricataVersion string
	Uptime          int64
}

func (this *AgentReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname=%s, IP=%s, agentVersion=%s, engineVersion=%s, engineUptime=%d>",
		this.Hostname,
		this.IP,
		this.AgentVersion,
		this.SuricataVersion,
		this.Uptime,
	)
}

func ReportAgentStatus() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.RpcAddr != "" {
		go reportAgentStatus(time.Duration(g.Config().Heartbeat.Interval) * time.Second)
	}
}

func reportAgentStatus(interval time.Duration) {
	for {
		req := AgentReportRequest{
			Hostname:        g.Hostname(),
			IP:              g.IP(),
			AgentVersion:    g.VERSION,
			Uptime:          funcs.GetUptime(),
			SuricataVersion: funcs.GetVersion(),
		}

		var resp SimpleRpcResponse
		err := g.HbsClient.Call("Agent.ReportStatus", req, &resp)
		if err != nil || resp.Code != 0 {
			Log.Error("Agent <= Heartbeat, Agent.ReportStatus fail:%s, Request=%v, Response=%v", err, req, resp)
		}

		time.Sleep(interval)
	}
}
