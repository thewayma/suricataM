package cron

import (
	"bytes"
	"github.com/thewayma/suricataM/agent/funcs"
	"github.com/thewayma/suricataM/agent/g"
	. "github.com/thewayma/suricataM/comm/log"
	"github.com/thewayma/suricataM/comm/st"
	"io"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"time"
)

func SyncSuricata() {
	if !g.Config().Heartbeat.Enabled {
		return
	}

	if g.Config().Heartbeat.HttpAddr == "" {
		return
	}

	go syncControlCommand()
	go syncMonitorMetric()
}

func ParseStructByReflect(u interface{}) {
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			Log.Trace("field=%s, type=%s, value=%v, tag=%s\n", t.Field(i).Name, v.Field(i).Type().Kind(), v.Field(i).Interface(), t.Field(i).Tag)

			if v.Field(i).Type().Kind() == reflect.Struct {
				sft := v.Field(i).Type()
				sfv := v.Field(i)
				for j := 0; j < sfv.NumField(); j++ {
					Log.Trace("    structfield=%s, type=%s, value=%v, tag=%s\n",
						sft.Field(j).Name, sfv.Field(j).Type().Kind(), sfv.Field(j).Interface(), sft.Field(j).Tag)

					key, ok := st.MetricMapper[sft.Field(j).Name]
					if !ok {
						continue
					}
					val := false
					if sfv.Field(j).Interface() == "off" {
						val = true
					}

					st.IgnoreMetric.Lock()
					st.IgnoreMetric.Item[key] = val
					st.IgnoreMetric.Unlock()
				}
			}
		}
	}
}

func syncMonitorMetric() {
	interval := g.Config().Suricata.MonitorMetricPollInterval
	Log.Trace("Agent Starts syncMonitorMetric, MonitorMetricPollInterval=%d", interval)
	duration := time.Duration(interval) * time.Second

	for {
		time.Sleep(duration)

		req := st.AgentMetricCommandRequest{
			IP:       g.IP(),
			Hostname: g.Hostname(),
		}

		var resp st.AgentMetricCommandResponse
		err := g.HbsClient.Call("Agent.FetchMonitorConfig", req, &resp)
		if err != nil {
			Log.Error("Agent <= Heartbeat, Pull Policy %s", err.Error())
			continue
		}

		ParseStructByReflect(resp)
	}

}

func scStartString() string {
	var comm bytes.Buffer

	comm.WriteString("nohup ")
	comm.WriteString(g.Config().Suricata.Bin)
	comm.WriteString(" --dpdkintel ")
	comm.WriteString(" -c ")
	comm.WriteString(g.Config().Suricata.Conf)
	/*
	   for _, v := range g.Config().Suricata.Ifaces {
	       comm.WriteString(" -i ")
	       comm.WriteString(v)
	   }
	*/
	comm.WriteString(" &> /dev/null &")

	str := comm.String()
	Log.Info("Suricata Start Command String: %s", str)

	return str
}

func replaceRuleFile() error {
	rulefile_bak := g.Config().Suricata.RulesDir + "/jcloud_suricata.rule.bak"
	rulefile_org := g.Config().Suricata.RulesDir + "/jcloud_suricata.rule"
	str := "mv " + rulefile_bak + " " + rulefile_org

	_, err := exec.Command("/bin/sh", "-c", str).Output()
	if err != nil {
		Log.Error("Agent <= Heartbeat, mv %s %s, err=%s", rulefile_bak, rulefile_org, err.Error())
		return err
	}

	Log.Trace("Agent <= Heartbeat, mv %s %s, Success", rulefile_bak, rulefile_org)
	return nil
}

func downloadSCRuleFile() error {
	rulefile := g.Config().Suricata.RulesDir + "/jcloud_suricata.rule.bak"
	f, err := os.OpenFile(rulefile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		Log.Error("Agent <= Heartbeat, Fail to New File %s,  err=%s", rulefile, err.Error())
		return err
	}

	url := "http://" + g.Config().Heartbeat.HttpAddr + "/manage/opt/download_rules_file"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("Agent <= Heartbeat, Heartbeat Wrong URL: %s", url)
		return err
	}
	req.Header.Set("JCloudSec_NIDS_Agent", g.IP())
	//!< TODO: 鉴权措施

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Log.Error("Agent <= Heartbeat, Http Get Request Err: %s", err.Error())
		return err
	}

	if resp.StatusCode == 200 {
		written, err := io.Copy(f, resp.Body)
		if err != nil {
			Log.Error("Agent <= Heartbeat, Copy File Err: %s", err.Error())
			return err
		}

		Log.Trace("Agent <= Heartbeat, Fetch Suricata Rulefile From %s, and Copy to %s, FileLen=%d", url, rulefile, written)
	} else if resp.StatusCode == 304 {
		Log.Trace("Agent <= Heartbeat, Resp.StatusCode=304, No Need to Update Suricata Rulefile")
	} else {
		Log.Trace("Agent <= Heartbeat, Resp.StatusCode=%s, Dont know how to deal with that!", resp.StatusCode)
	}

	return nil
}

func syncControlCommand() {
	interval := g.Config().Suricata.ControlCommandPollInterval
	Log.Trace("Agent Starts syncControlCommand, ControlCommandPollInterval=%d", interval)
	duration := time.Duration(interval) * time.Second

	for {
		time.Sleep(duration)

		req := st.AgentControlCommandRequest{
			IP:       g.IP(),
			Hostname: g.Hostname(),
		}

		var resp st.AgentControlCommandResponse
		err := g.HbsClient.Call("Agent.FetchOpt", req, &resp)
		if err != nil {
			Log.Error("Agent <= Heartbeat, Pull Policy %s", err.Error())
			continue
		}

		switch resp.Command {
		case "engine-start":
			Log.Trace("Agent <= Heartbeat, Start Suricata")
			res, err := exec.Command("/bin/sh", "-c", scStartString()).Output()
			if err != nil {
				Log.Error("Agent <= Heartbeat, Start Suricata Failure: %s\n", err.Error())
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
				Log.Error("     ReStart_2. Start Suricata Failure: %s\n", err.Error())
			} else {
				Log.Trace("     ReStart_2. Start Suricata Success: %s\n", res)
			}

		case "reload-rules":
			Log.Trace("Agent <= Heartbeat, Reload Suricata Rules: %s", funcs.ReloadRules())

		case "rules-update":
			Log.Trace("Agent <= Heartbeat, Update Suricata Rules")
			downloadSCRuleFile()
			replaceRuleFile()
			Log.Trace("Agent <= Heartbeat, Reload Suricata Rules: %s", funcs.ReloadRules())

		default:
			Log.Trace("Agent <= Heartbeat, nothing to do right now!!!")
		}
	}
}
