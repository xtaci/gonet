package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

import (
	"cfg"
)

//----------------------------------------------- handle unix signals
func SignalProc() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP)

	for {
		msg := <-ch
		log.Println("Recevied signal:", msg)
		cfg.Reload()
	}
}
