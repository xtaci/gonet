package stats_tbl

import (
	"fmt"
	"testing"
)

func TestSTRSet(t *testing.T) {
	/*
		for i := int32(1); i <= 1; i++ {
			StrValue := "3"
			Key := fmt.Sprintf("%v#level", 5)
			SetUpdate(Key, StrValue)
		}
	*/
	//UserLogin(33323, "111")
}

func TestIntSet(t *testing.T) {
	for i := int32(1); i <= 30000; i++ {
		IntValue := int32(100)
		Key := fmt.Sprintf("%v#gems_buy_10%v_food", i, "%")
		SetAdds(Key, IntValue, "cn")
	}
	for i := int32(1); i <= 30000; i++ {
		IntValue := int32(100)
		Key := fmt.Sprintf("%v#gems_buy_10%v_gold", i, "%")
		SetAdds(Key, IntValue, "en")
	}
	for i := int32(1); i <= 30000; i++ {
		IntValue := int32(100)
		Key := fmt.Sprintf("%v#shop_gems_buy_10%v_popu", i, "%")
		SetAdds(Key, IntValue, "en")
	}
}
