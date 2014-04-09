package main

const (
	TCP_TIMEOUT           = 60 // TCP读超时(sec)
	MAX_DELAY_IN          = 60 // 数据包等待进入队列最长时间(sec)
	DEFAULT_INQUEUE_SIZE  = 64 // 默认输入数据包队列
	DEFAULT_OUTQUEUE_SIZE = 15 // 默认输出数据包队列
)

const (
	DEFAULT_MQ_SIZE   = 32  // 默认玩家IPC消息队列大小
	DEFAULT_FLUSH_OPS = 128 // 默认刷入操作数
	CUSTOM_TIMER      = 60  // 玩家定时器间隔
)

const (
	PACKET_EXPIRE = 120 // 包过期时间
	PACKET_ERROR  = 5   // 时间误差
)

const (
	SYS_MQ_SIZE    = 65535 // 系统进程的消息队列大小
	CACHE_CLEAN    = 3600  // 玩家信息CACHE过期间隔
	ALLIANCE_CACHE = 3600  // 联盟CACHE过期间隔
	GC_INTERVAL    = 300   // 主动垃圾回收间隔
)
