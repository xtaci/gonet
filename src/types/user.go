package types

type User struct {
	MQ                 chan string
	Id                 int
	Name               string
	OwnerId            int
	IsCapital          int
	Achievement        int
	X                  int
	Y                  int
	Gold               int
	Wood               int
	Food               int
	Iron               int
	Stone              int
	Scout              int
	Swordsman          int
	CrossbowArcher     int
	Squire             int
	Templar            int
	Paladin            int
	ArcherCavalry      int
	RoyalKnight        int
	ActionEventsCount  int
	RecruitEventsCount int
	DealsCount         int
	LockVersion        int
	LastMoveTime       int
	Durability         int
	ArcaneMage         int
	BattleMage         int
	HolyMage           int
	IsAutoFix          int
	ReviveTime         int
	ItemWarehouseLv    int
	ItemTransportLv    int
	Skeleton           int
	GhostRider         int
	Ram                int
	Zeppelin           int
	SteelGolem         int
	Cruiser            int
}
