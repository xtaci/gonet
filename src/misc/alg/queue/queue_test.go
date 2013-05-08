package queue

import "testing"
import "fmt"

func TestQueue(t *testing.T) {
	q := New(10)

	for i:=0;i<10;i++ {
		q.Enqueue(i+10)
	}

	v := q.Dequeue()
	fmt.Println(v)

	for v = q.Dequeue() ; v!= nil; v = q.Dequeue() {
		fmt.Println(v)
	}
}
