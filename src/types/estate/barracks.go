package estate

import (
	"types/solider"
)

const (
	BARRACKS_STATUS_NORMAL     = 0
	BARRACKS_STATUS_UPGRADING  = 1
	BARRACKS_STATUS_RECRUITING = 2
)

type BarracksProperty struct {
	Type         string
	CurrentLevel uint8
	Status       uint8
	Recruting    []solider.SoliderCD // barracks only manages solider creating
}
