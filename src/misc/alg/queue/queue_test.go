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
	v := q.Dequeue()
	fmt.Println(v)

	for v = q.Dequeue(); v != nil; v = q.Dequeue() {
		fmt.Println(v)
	}
}
