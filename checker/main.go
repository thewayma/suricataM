package main

import (
	"flag"
	"github.com/thewayma/suricataM/checker/cron"
	"github.com/thewayma/suricataM/checker/g"
	"github.com/thewayma/suricataM/checker/rx"
	_ "github.com/thewayma/suricataM/comm/log"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)

	g.InitRedisConnPool()
	g.InitHbsClient()

	go rx.RpcStart()
	go cron.SyncSuricataPolicy()
	//go cron.CleanStale()

	select {}
}
