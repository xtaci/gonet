package timer

import "testing"
import "time"
import "fmt"

func TestTimer(t *testing.T) {
	ch := make(chan int32, 100)

	now := time.Now().Unix()

	for i := 0; i < 100; i++ {
		Add(1, now+int64(i), ch)
	}

	count := 0
	for {
		fmt.Println("timer event :", <-ch)
		count++
		if count == 100 {
			break
		}
	}
}
