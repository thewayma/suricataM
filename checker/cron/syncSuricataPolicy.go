package cron

import (
	"github.com/thewayma/suricataM/checker/g"
	. "github.com/thewayma/suricataM/comm/log"
	"github.com/thewayma/suricataM/comm/st"
	"time"
)

func SyncSuricataPolicy() {
	duration := time.Duration(g.Config().Hbs.PolicyPollInterval) * time.Second
	for {
		syncStrategies()
		time.Sleep(duration)
	}
}

func syncStrategies() {
	req := st.StrategiesRequest{
		IP:       g.IP(),
		Hostname: g.Hostname(),
	}
	var resp st.StrategiesResponse
	err := g.HbsClient.Call("Agent.FetechPolicyConfig", req, &resp)
	if err != nil {
		Log.Error("Checker <= Hearbeat, RPC Agent.FetechPolicyConfig Err: %s", err)
		return
	}

	//rebuildStrategyMap(&strategiesResponse)
}
