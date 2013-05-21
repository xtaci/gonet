package estate

const (
	DEFENSE_STATUS_NORMAL    = 0
	DEFENSE_STATUS_UPGRADING = 1
)

type DefenseProperty struct {
	Type         string
	CurrentLevel uint8
	Status       uint8
}
