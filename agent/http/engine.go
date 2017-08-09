package http

import (
	"github.com/thewayma/suricataM/agent/funcs"
	. "github.com/thewayma/suricataM/comm/st"
	"net/http"
)

/* Agent 开放接口
1. Get操作
    /engine/version,        GetVersion
    /engine/runningmode,    GetRunningMode
    /engine/capturemode,    GetCaptureMode
    /engine/uptime,         GetUptime
    /engine/allportstats,   GetAllPortStats

2. Set操作
    /engine/shutdown,       ShutDown
    /engine/reloadrules,    ReloadRules
*/
func configEngine() {
	http.HandleFunc("/engine/version", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.GetVersion())
	})

	http.HandleFunc("/engine/runningmode", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.GetRunningMode())
	})

	http.HandleFunc("/engine/capturemode", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.GetCaptureMode())
	})

	http.HandleFunc("/engine/uptime", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.GetUptime())
	})

	http.HandleFunc("/engine/allportstat", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.GetAllPortStats())
	})

	http.HandleFunc("/engine/shutdown", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.ShutDown())
	})

	http.HandleFunc("/engine/reloadrule", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.ReloadRules())
	})
}
