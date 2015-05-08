package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

import (
	"cfg"
	"gamedata"
	. "helper"
)

//----------------------------------------------- handle unix signals
func SignalProc() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP: // reload config
			log.Println("\033[043;1m[SIGHUP]\033[0m")
			cfg.Reload()
			gamedata.Reload()
		case syscall.SIGTERM: // server close
			close(die)
			log.Println("\033[043;1m[SIGTERM]\033[0m")
			INFO("waiting for agents close, please wait...")
			wg.Wait()
			INFO("all work done. bye bye!!!")
			os.Exit(0)
		}
	}
}
