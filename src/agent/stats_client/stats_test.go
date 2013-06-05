package stats_client

import (
	"fmt"
	"testing"
)

func TestStatsFunc(t *testing.T) {
	DialStats()
	fmt.Println(Ping())
}
