package timer

import "testing"
import "time"
import "fmt"

func TestTimer(t *testing.T) {
	ch := make(chan int32, 100)

	now := time.Now().Unix()
	Add(1, now+1, ch)
	Add(2, now+2, ch)
	Add(3, now+3, ch)
	Add(4, now+4, ch)
	Add(5, now+5, ch)
	// Add(10, now+10, ch)
	//	Add(now+60, ch)

	count := 0
	for {
		fmt.Println("timer event :", <-ch)
		count++
		if count == 5 {
			break
		}
	}
}
