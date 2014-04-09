package main

const (
	TCP_TIMEOUT           = 60 // tcp read timeout
	MAX_DELAY_IN          = 60 // packet process wait
	DEFAULT_INQUEUE_SIZE  = 64 // default input buffer size
	DEFAULT_OUTQUEUE_SIZE = 15 // default output buffer size
)

const (
	DEFAULT_MQ_SIZE   = 32  // size of user's message queue
	DEFAULT_FLUSH_OPS = 128 // max ops before flush
	CUSTOM_TIMER      = 60  // a timer for each user
)

const (
	PACKET_EXPIRE = 120 // packet before this duration is illegal
	PACKET_ERROR  = 5   // due to time error of client and server
)

const (
	SYS_MQ_SIZE = 65535 // size of sys routine's message queue
	GC_INTERVAL = 300   // voluntary GC interval
)
