#!/bin/sh
$GOPATH/bin/hub &
$GOPATH/bin/event &
$GOPATH/bin/stats &
sleep 3
$GOPATH/bin/agent
