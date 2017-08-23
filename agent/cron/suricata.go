package cron

import (
	"bytes"
	"github.com/thewayma/suricataM/agent/funcs"
	"github.com/thewayma/suricataM/agent/g"
	. "github.com/thewayma/suricataM/comm/log"
	"github.com/thewayma/suricataM/comm/st"
	"os/exec"
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

func scStartString() string {
	var comm bytes.Buffer

	comm.WriteString("nohup ")
	comm.WriteString(g.Config().Suricata.Bin)
	comm.WriteString(" -c ")
	comm.WriteString(g.Config().Suricata.Conf)

	for _, v := range g.Config().Suricata.Ifaces {
		comm.WriteString(" -i ")
		comm.WriteString(v)
	}
	comm.WriteString(" &> /dev/null &")

	str := comm.String()
	Log.Info("Suricata Start Command String: %s", str)

	return str
}

func syncControlCommand() {
	Log.Trace("Agent Starts syncControlCommand, Interval=%d", g.Config().Suricata.Interval)

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
			Log.Error("Agent <= Heartbeat, Pull Policy %s", err)
			continue
		}

		switch resp.Command {
		case "engine-start":
			Log.Trace("Agent <= Heartbeat, Start Suricata")
			res, err := exec.Command("/bin/sh", "-c", scStartString()).Output()
			if err != nil {
				Log.Error("Agent <= Heartbeat, Start Suricata Failure: %s\n", err)
			} else {
				Log.Trace("Agent <= Heartbeat, Start Suricata Success: %s\n", res)
			}

		case "engine-shutdown":
			Log.Trace("Agent <= Heartbeat, ShutDown Suricata: %s", funcs.ShutDown())

		case "engine-restart":
			Log.Trace("Agent <= Heartbeat, Restart Suricata")
			Log.Trace("     Restart_1. ShutDown: %s", funcs.ShutDown())
			res, err := exec.Command("/bin/sh", "-c", scStartString()).Output()
			if err != nil {
				Log.Error("     ReStart_2. Start Suricata Failure: %s\n", err)
			} else {
				Log.Trace("     ReStart_2. Start Suricata Success: %s\n", res)
			}

		case "reload-rules":
			Log.Trace("Agent <= Heartbeat, Reload Suricata Rules: %s", funcs.ReloadRules())

		case "rules-update":
			Log.Trace("Agent <= Heartbeat, Update Suricata Rules")

		default:
			Log.Trace("Agent <= Heartbeat, nothing to do right now!!!")
		}
	}
}
