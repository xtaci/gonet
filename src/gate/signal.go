package main

import "os"
import "os/signal"
import "syscall"
import "log"

func SignalProc() {

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT)

	for {
		msg := <-ch
		log.Println(msg)
	}
}
