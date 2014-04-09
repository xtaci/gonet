package main

import (
	"agent/hub_client"
	"agent/stats_client"
	"cfg"
	. "helper"
)

//---------------------------------------------------------- 服务器启动流程
func startup() {
	INFO("Starting GS.")
	// start logger
	config := cfg.Get()
	if config["gs_log"] != "" {
		cfg.StartLogger(config["gs_log"])
	}

	// dial HUB
	hub_client.DialHub()

	// signal
	go SignalProc()

	// sys routine
	go SysRoutine()

	// stats
	go stats_client.DialStats()
}
