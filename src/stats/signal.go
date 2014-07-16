package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	//    "sync/atomic"
)

import (
	"cfg"
)

//----------------------------------------------- handle unix signals
func SignalProc() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			log.Println("Recevied signal:", msg)
			cfg.Reload()
			/*
			   case syscall.SIGTREM:
			       atomic.StoreInt32(&SIGTREM,1)
			       wg.Wait()
			       os.Exit(-1)
			*/
		}
	}
}
