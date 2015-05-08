package protos

import (
	"sync"
)

import (
	. "types"
)

var (
	// 各个服务器的Forward消息队列
	_servers    map[int32]chan IPCObject
	_serverlock sync.Mutex
)

func AddServer(hostid int32, forward chan IPCObject) {
	_serverlock.Lock()
	defer _serverlock.Unlock()
	_servers[hostid] = forward
}

func RemoveServer(hostid int32) {
	_serverlock.Lock()
	defer _serverlock.Unlock()
	delete(_servers, hostid)
}

func ForwardChan(hostid int32) chan IPCObject {
	_serverlock.Lock()
	defer _serverlock.Unlock()
	return _servers[hostid]
}

func AllServers() []int32 {
	_serverlock.Lock()
	defer _serverlock.Unlock()

	_all := make([]int32, len(_servers))
	idx := 0
	for k := range _servers {
		_all[idx] = k
		idx++
	}

	return _all
}

func init() {
	_servers = make(map[int32]chan IPCObject)
}
