package market

import (
	"sync"
	"sync/atomic"
	"time"
)

type Item struct {
	ORDER_NO	uint64			// Order Number
	Code		int32			// Product Code
	Price		float64			// Price
	Count		uint32			// Count 
	Seller		int32			// Seller id
	Date		time.Time		// time on shelf
}

var (
	_items map[uint64]*Item				// map order->item
	_codes map[int32]map[uint64]*Item	// map to map 

	_lock sync.RWMutex
	_next_order_no uint64
)

func init() {
	_items = make(map[uint64]*Item)
	_codes = make(map[int32]map[uint64]*Item)
}

//--------------------------------------------------------- New Selling Order 
func NewSell(seller int32, code int32, price float64, count uint32) uint64 {
	nr := atomic.AddUint64(&_next_order_no, 1)

	_lock.Lock()
	defer _lock.Unlock()

	item := Item{ORDER_NO:nr, Code:code, Price:price, Count:count, Seller:seller, Date:time.Now()}
	_items[nr] = &item

	if _codes[code] == nil {
		_codes[code] = make(map[uint64]*Item)
	}
	_codes[code][nr] = &item

	return nr
}

//--------------------------------------------------------- Delete a Order
func DeleteOrder(order_no uint64) bool {
	_lock.Lock()
	defer _lock.Unlock()

	if item := _items[order_no]; item != nil {
		delete(_codes[item.Code], order_no)
		delete(_items, order_no)
		return true
	}

	return false
}

//--------------------------------------------------------- Get Product List
func List(start, count int, code int32) (ret []Item) {
	_lock.RLock()
	defer _lock.RUnlock()

	var idx int

	for _,v := range _codes[code] {
		if idx>=start {
			ret = append(ret, *v)
		}

		idx++

		if idx >= start+count {
			break
		}
	}

	return
}

//--------------------------------------------------------- Get Product Count
func Count(code int32) int {
	_lock.RLock()
	defer _lock.RUnlock()

	return len(_codes[code])
}
