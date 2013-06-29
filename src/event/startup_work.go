package main

import (
	"db/event_tbl"
	"event/core"
)

func startup_work() {
	events := event_tbl.GetAll()

	for _, v := range events {
		core.Load(&v)
	}
}
