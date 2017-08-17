package main

import (
	"flag"
	"github.com/thewayma/suricataM/agent/cron"
	"github.com/thewayma/suricataM/agent/funcs"
	"github.com/thewayma/suricataM/agent/g"
	"github.com/thewayma/suricataM/agent/http"
	_ "github.com/thewayma/suricataM/comm/log"
    "bytes"
    "fmt"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)
	g.InitRpcClients()

	funcs.GenerateCollectorFuncs()

	cron.PreCollect()
	cron.Collect()
	cron.ReportAgentStatus()
	cron.SyncSuricata()

	go http.Start()



    var comm bytes.Buffer
    comm.WriteString(g.Config().Suricata.Bin)
    comm.WriteString(" -c ")
    comm.WriteString(g.Config().Suricata.Conf)

    for _, v := range g.Config().Suricata.Ifaces {
        comm.WriteString(" -i ")
        comm.WriteString(v)
    }
	fmt.Println(comm.String())


	select {}
}
