package types

//------------------------------------------------ 状态机定义
const (
	UNKNOWN = byte(iota)
	OFF_FREE
	OFF_RAID
	OFF_PROT
	ON_FREE
	ON_PROT
)
