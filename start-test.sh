#!/bin/sh 
killall -9 agent
killall -9 stats
killall -9 hub

killall -0 agent
while [ $? -ne 1 ]; do
	sleep 1
	killall -0 agent
done

killall -0 stats
while [ $? -ne 1 ]; do
	sleep 1
	killall -0 stats
done

killall -0 hub
while [ $? -ne 1 ]; do
	sleep 1
	killall -0 hub
done

$GOPATH/bin/hub &
$GOPATH/bin/stats &
sleep 1
$GOPATH/bin/agent &
