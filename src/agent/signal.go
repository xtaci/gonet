package main

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

import (
	"cfg"
	"gamedata"
	. "helper"
)

//----------------------------------------------- handle unix signals
func SignalProc() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP: // 重载入 config/gamedata等配置数据
			log.Println("\033[043;1m[SIGHUP]\033[0m")
			cfg.Reload()
			gamedata.Reload()
		case syscall.SIGTERM: // 关闭agent
			atomic.StoreInt32(&SIGTERM, 1)
			log.Println("\033[043;1m[SIGTERM]\033[0m")
			INFO("waiting for agents close, please wait...")
			wg.Wait()
			INFO("all work done. bye bye!!!")
			os.Exit(0)
		}
	}
}
