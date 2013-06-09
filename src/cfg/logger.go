package cfg

import (
	"log"
	"os"
)

func StartLogger(logfile string) {
	bl := []byte(logfile)

	var err error
	var f *os.File

	if bl[0] == '/' { // start with slash, just open
		f, err = os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	} else {
		path := os.Getenv("GOPATH") + "/" + logfile
		f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	}

	if err != nil {
		log.Println("cannot open logfile %v\n", err)
		os.Exit(-1)
	}

	var r Repeater

	config := Get()
	switch config["log_output"] {
	case "both":
		r.out1 = os.Stdout
		r.out2 = f
	case "file":
		r.out2 = f
	}
	log.SetOutput(&r)
}
