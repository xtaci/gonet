package stats_client

import (
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

import (
	"cfg"
	. "helper"
	"misc/packet"
	"misc/timer"
)

const (
	STATS_QUEUE_SIZE     = 100000
	STATS_COLLECT_PERIOD = 30 //secs
)

var _conn net.Conn
var _seq_id uint64
var _wait_ack map[uint64]chan []byte
var _wait_ack_lock sync.Mutex

var (
	Accum_queue  chan SET_ADDS_REQ
	Statsd_queue chan SET_ADDS_REQ
	Update_queue chan SET_UPDATE_REQ
)

func init() {
	_wait_ack = make(map[uint64]chan []byte)
	Accum_queue = make(chan SET_ADDS_REQ, STATS_QUEUE_SIZE)
	Statsd_queue = make(chan SET_ADDS_REQ, STATS_QUEUE_SIZE)
	Update_queue = make(chan SET_UPDATE_REQ, STATS_QUEUE_SIZE)

	go stats_sender()
}

//----------------------------------------------- connect to Stats server
func DialStats() {
	INFO("Connecting to Stats server")
	config := cfg.Get()

	addr, err := net.ResolveUDPAddr("udp", config["stats_service"])
	if err != nil {
		ERR(err)
		os.Exit(-1)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		ERR(err)
		os.Exit(-1)
		return
	}

	_conn = conn
	INFO("Stats Service Connected")
}

//------------------------------------------------ send msg to stats server
func Send(data []byte) bool {
	// send the packet
	_, err := _conn.Write(data)
	if err != nil {
		ERR("Error send packet to Stats Server:", err)
		return false
	}
	return true
}

func stats_sender() {
	_accum_buffer := make(map[string]map[string]int32)
	_update_buffer := make(map[string]map[string]string)
	_statsd_buffer := make(map[string]int32)

	stats_timer := make(chan int32, 100)
	stats_timer <- 1
	for {
		select {
		case req := <-Accum_queue:
			if _, ok := _accum_buffer[req.F_lang]; !ok {
				val := make(map[string]int32)
				val[req.F_key] = 0
				_accum_buffer[req.F_lang] = val
			}
			val := _accum_buffer[req.F_lang]
			val[req.F_key] += req.F_value
			_accum_buffer[req.F_lang] = val
		case req := <-Update_queue:
			if _, ok := _update_buffer[req.F_lang]; !ok {
				val := make(map[string]string)
				val[req.F_key] = ""
				_update_buffer[req.F_lang] = val
			}
			val := _update_buffer[req.F_lang]
			val[req.F_key] = req.F_value
			_update_buffer[req.F_lang] = val
		case req := <-Statsd_queue:
			_statsd_buffer[req.F_key] += req.F_value
		case <-stats_timer:
			INFO("Stats Buffer:", len(_accum_buffer), len(_update_buffer), len(_statsd_buffer))
			// 累计
			accum := SET_ADDS_REQ{}
			for accum.F_lang, _ = range _accum_buffer {
				for accum.F_key, accum.F_value = range _accum_buffer[accum.F_lang] {
					Send(packet.Pack(Code["set_adds_req"], &accum, nil))
				}
			}
			_accum_buffer = make(map[string]map[string]int32)

			// 更新
			update := SET_UPDATE_REQ{}
			for update.F_lang, _ = range _update_buffer {
				for update.F_key, update.F_value = range _update_buffer[update.F_lang] {
					Send(packet.Pack(Code["set_update_req"], &update, nil))
				}
			}
			_update_buffer = make(map[string]map[string]string)

			// FINI
			config := cfg.Get()
			period := STATS_COLLECT_PERIOD
			if config["stats_collect_period"] != "" {
				period, _ = strconv.Atoi(config["stats_collect_period"])
			}

			timer.Add(0, time.Now().Unix()+int64(period), stats_timer)
			runtime.GC()
		}
	}
}

func checkErr(err error) {
	if err != nil {
		ERR(err)
		panic("error occured in protocol module")
	}
}
