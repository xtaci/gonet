package protos

import (
	"sync"
)

var (
	// 各个服务器的Forward消息队列
	_servers    map[int32]chan []byte
	_serverlock sync.RWMutex
)

func AddServer(hostid int32, forward chan []byte) {
	_serverlock.Lock()
	defer _serverlock.Unlock()
	_servers[hostid] = forward
}

func RemoveServer(hostid int32) {
	_serverlock.Lock()
	defer _serverlock.Unlock()
	delete(_servers, hostid)
}

func ForwardChan(hostid int32) chan []byte {
	_serverlock.RLock()
	defer _serverlock.RUnlock()
	return _servers[hostid]
}

func init() {
	_servers = make(map[int32]chan []byte)
}
