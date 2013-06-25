package main

import (
	"agent/gsdb"
	"agent/ipc"
	"math/rand"
	"sync"
	"testing"
	. "types"
)

func _simple_receiver(sess *Session, wg *sync.WaitGroup) {
	for {
		obj := <-sess.MQ
		IPCRequestProxy(sess, &obj)
		wg.Done()
	}
}

func TestIPC(t *testing.T) {
	// fake 2 user
	var sess1 Session
	sess1.User = &User{Id: 1}
	sess1.MQ = make(chan IPCObject, 10)
	var sess2 Session
	sess2.User = &User{Id: 2}
	sess2.MQ = make(chan IPCObject, 10)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go _simple_receiver(&sess1, wg)
	go _simple_receiver(&sess2, wg)

	gsdb.RegisterOnline(&sess1, sess1.User.Id)
	gsdb.RegisterOnline(&sess2, sess2.User.Id)

	if !ipc.Send(1, 2, ipc.SERVICE_PING, false, "ABC") {
		t.Fatal("ipc.Send failed")
	}
	wg.Wait()
}

func BenchmarkIPC(b *testing.B) {
	wg := &sync.WaitGroup{}

	for i := 1; i <= b.N; i++ {
		var sess Session
		sess.User = &User{Id: int32(i)}
		sess.MQ = make(chan IPCObject, 10)
		gsdb.RegisterOnline(&sess, sess.User.Id)
		go _simple_receiver(&sess, wg)
	}

	for i := 1; i <= b.N; i++ {
		src := rand.Int31n(int32(b.N)) + 1
		dest := rand.Int31n(int32(b.N)) + 1
		if src != dest {
			wg.Add(2)
			ipc.Send(src, dest, ipc.SERVICE_PING, false, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		}
	}
	wg.Wait()
}
