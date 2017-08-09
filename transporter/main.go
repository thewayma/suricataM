package main

import (
	"flag"
	_"github.com/thewayma/suricataM/comm/log"
	"github.com/thewayma/suricataM/transporter/g"
	"github.com/thewayma/suricataM/transporter/rx"
	"github.com/thewayma/suricataM/transporter/tx"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)

	tx.Start()
	go rx.RpcServer()

	select {}
}
