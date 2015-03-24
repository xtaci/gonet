package consistent_hash

import (
	"sync"
)

type ConsistentHashing struct {
	split_points []uint32
	keys         map[uint32]string
	sync.Mutex
}

func (ch *ConsistentHashing) Init() {
	ch.keys = make(map[uint32]string)
}

//----------------------------------------------- add a key with hashcode to the circle
func (ch *ConsistentHashing) AddNode(key string, hashcode uint32) bool {
	ch.Lock()
	defer ch.Unlock()

	// test hashcode existence
	if _, ok := ch.keys[hashcode]; ok {
		return false
	}
	ch.keys[hashcode] = key

	// hashcode in the middle of the circle
	for i := 0; i < len(ch.split_points); i++ {
		if hashcode < ch.split_points[i] {
			ch.split_points = append(ch.split_points[:i], append([]uint32{hashcode}, ch.split_points[i:]...)...)
			return true
		}
	}

	// largest hashcode or empty circle
	ch.split_points = append(ch.split_points, hashcode)
	return true
}

//----------------------------------------------- remove a node from the circle
func (ch *ConsistentHashing) RemoveNode(hashcode uint32) bool {
	ch.Lock()
	defer ch.Unlock()

	if _, ok := ch.keys[hashcode]; ok {
		delete(ch.keys, hashcode)
		for i := 0; i < len(ch.split_points); i++ {
			if ch.split_points[i] == hashcode { // node found!
				ch.split_points = append(ch.split_points[:i], ch.split_points[i+1:]...)
				return true
			}
		}
	}
	return false
}

//----------------------------------------------- get the node by a given hashcode
func (ch *ConsistentHashing) GetNode(hashcode uint32) (key string, ok bool) {
	ch.Lock()
	defer ch.Unlock()

	// if empty circle
	if len(ch.split_points) == 0 {
		return "", false
	}

	// find nearest node
	for i := range ch.split_points {
		if ch.split_points[i] >= hashcode {
			return ch.keys[ch.split_points[i]], true
		}
	}

	// hashcode is larger than the largest node, return to the first node
	return ch.keys[ch.split_points[0]], true
}
