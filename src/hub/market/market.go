package market

import (
	"sync"
	"sync/atomic"
)

type Item struct {
	ORDER_NO	uint64
	Code		int32
	Price		float64
	Count		uint32
	Seller		int32
}

var _items map[uint64]*Item
var _lock sync.Mutex
var _next_order_no uint64

func init() {
	_items = make(map[uint64]*Item)
}

func Sell(seller int32, code int32, price float64, count uint32) uint64 {
	nr := atomic.AddUint64(&_next_order_no, 1)

	_lock.Lock()
	defer _lock.Unlock()

	_items[nr] = &Item{ORDER_NO:nr, Code:code, Price:price, Count:count, Seller:seller}

	return nr
}

func Buy(order_no uint64) bool {
	_lock.Lock()
	defer _lock.Unlock()

	if _items[order_no] != nil {
		delete(_items, order_no)
		return true
	}

	return false
}
