package cron

import (
	"github.com/thewayma/suricataM/checker/g"
	. "github.com/thewayma/suricataM/comm/log"
	"github.com/thewayma/suricataM/comm/st"
	"time"
)

var lastStrategiesVersion = "0.0.0"

func SyncSuricataPolicy() {
	Log.Trace("Checker <= Hearbeat, Start SyncSuricataPolicy")
	duration := time.Duration(g.Config().Hbs.PolicyPollInterval) * time.Second

	for {
		syncStrategies()
		time.Sleep(duration)
	}
}

func syncStrategies() {
	req := st.StrategiesRequest{
		Ip:       g.IP(),
		Hostname: g.Hostname(),
	}
	var resp st.StrategiesResponse
	err := g.HbsClient.Call("Agent.FetechPolicyConfig", req, &resp)
	if err != nil {
		Log.Error("Checker <= Hearbeat, RPC Agent.FetechPolicyConfig Err: %s", err)
		return
	}

	if lastStrategiesVersion != resp.Version {
		Log.Trace("Checker <= Hearbeat, Strategies need to be updated, lastVersion=%s, newVersion=%s", lastStrategiesVersion, resp.Version)
		lastStrategiesVersion = resp.Version
		rebuildStrategyMap(&strategiesResponse)
	} else {
		Log.Trace("Checker <= Hearbeat, Strategies dnt need to be updated, lastVersion=newVersion=%s", resp.Version)
	}
}

func rebuildStrategyMap(resp *st.StrategiesResponse) {
	key := "allStrategy"
	m := make(map[string][]st.Strategy)

	for _, v := range resp.Strategies {
		if _, exist := m[key]; exist {
			m[key] = append(m[key], v)
		} else {
			m[key] = []st.Strategy{v}
		}
	}

	g.StrategyMap.ReInit(m)
}
