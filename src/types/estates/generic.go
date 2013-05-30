package estates

const (
	COLLECTION    = "ESTATES"
	STATUS_NORMAL = byte(0)
	STATUS_CD     = 1
)

//----------------------------------------------- Generic Move event records
type Move struct {
	OID uint32
	X   uint16
	Y   uint16
}

//----------------------------------------------- Generic Cooldown event records
type CD struct {
	OID     uint32
	Timeout int64
}

type Manager struct {
	UserId    int32
	Version   uint32
	Basic     *BasManager
	Offensive *OffManager
	Defensive *DefManager
	Res       *ResManager
}

func (m *Manager) CheckAllCD() int {
	ops := 0
	ops += m.Basic.CheckCD()
	ops += m.Offensive.CheckCD()
	ops += m.Defensive.CheckCD()
	ops += m.Res.CheckCD()

	return ops
}
