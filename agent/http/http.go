package http

import (
	"github.com/thewayma/suricataM/agent/g"
	. "github.com/thewayma/suricataM/comm/log"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	configEngine()
}

func Start() {
	if !g.Config().Http.Enabled {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	Log.Info("Agent HTTP listening: %v", addr)
	log.Fatalln(s.ListenAndServe())
}
