package timer

import "testing"
import "time"
import "fmt"

func TestTimer(t *testing.T) {
	ch := make(chan int32, 10)

	now := time.Now().Unix()

	for i := 0; i < 10; i++ {
		Add(int32(i), now+int64(i), ch)
	}

	count := 0
	for {
		fmt.Println("timer event :", <-ch)
		count++
		if count == 10 {
			break
		}
	}
}

func BenchmarkTimer(b *testing.B) {
	ch := make(chan int32, 100)
	now := time.Now().Unix()
	for i := 0; i < b.N; i++ {
		Add(1, now+100000, ch)
	}
}
