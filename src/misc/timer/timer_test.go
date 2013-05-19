package timer

import "testing"
import "time"
import "fmt"

func TestTimer(t *testing.T) {
	ch := make(chan uint32, 100)

	now := time.Now().Unix()
	Add(now+1, ch)
	Add(now+2, ch)
	Add(now+3, ch)
	Add(now+4, ch)
	Add(now+5, ch)

	count := 0
	for {
		fmt.Println("timer event :", <-ch)
		count++
		if count == 5 {
			break
		}
	}
}
