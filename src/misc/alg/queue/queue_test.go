package queue

import "testing"
import "fmt"

func TestQueue(t *testing.T) {
	q := New(10)

	for i := 0; i < 10; i++ {
		q.Enqueue(i + 10)
	}

	q.Dequeue()
	q.Enqueue(999)

	fmt.Println("testing All")
	s := q.All()

	for k := range s {
		fmt.Println(s[k])
	}

	fmt.Println("testing dequeue")
	for {
		if v, ok := q.Dequeue(); ok {
			fmt.Println(v)
		} else {
			break
		}
	}
}

func BenchmarkQueue(b *testing.B) {
	q := New(b.N)

	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}

	q.All()
	for {
		if _, ok := q.Dequeue(); ok {
		} else {
			break
		}
	}
}
