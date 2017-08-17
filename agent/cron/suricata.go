package cron

import (
	"github.com/thewayma/suricataM/agent/funcs"
	"github.com/thewayma/suricataM/agent/g"
	. "github.com/thewayma/suricataM/comm/log"
	"github.com/thewayma/suricataM/comm/st"
	"time"
)

func SyncSuricata() {
	if !g.Config().Heartbeat.Enabled {
		return
	}

	if g.Config().Heartbeat.Addr == "" {
		return
	}

	go syncControlCommand()
}

func syncControlCommand() {
	// TODO: 现用g.Config().Suricata.Interval统一命令控制, 策略规则的更新间隔, 将来可以在配置文件中分开处理
	duration := time.Duration(g.Config().Suricata.Interval) * time.Second

	for {
		time.Sleep(duration)

		req := st.AgentControlCommandRequest{
			IP: g.IP(),
		}

		var resp st.AgentControlCommandResponse
		err := g.HbsClient.Call("Agent.FetchOpt", req, &resp)
		if err != nil {
			Log.Error("Agent <= Heartbeat Pull Policy %s", err)
			continue
		}

		switch resp.Command {
		case "engine-start":
			Log.Trace("Agent <= Heartbeat, Start Suricata")

		case "engine-shutdown":
			Log.Trace("Agent <= Heartbeat, Shutdown Suricata")
			funcs.ShutDown()

		case "engine-restart":
			Log.Trace("Agent <= Heartbeat, Restart Suricata")
			funcs.ShutDown()

		case "reload-rules":
			Log.Trace("Agent <= Heartbeat, Reload Suricata Rules")
			funcs.ReloadRules()

		case "rules-update":
			Log.Trace("Agent <= Heartbeat, Update Suricata Rules")

		default:
			Log.Trace("Agent <= Heartbeat, Nothing to do right now")
		}
	}
}
