package queue

type Queue struct {
	capacity int
	size     int
	front    int
	rear     int
	elements []interface{}
}

func New(max int) *Queue {
	queue := &Queue{capacity: max, size: 0, front: 0, rear: -1}
	queue.elements = make([]interface{}, max)
	return queue
}

//----------------------------------------------- Enqueue
func (q *Queue) Enqueue(elem interface{}) bool {
	if q.size < q.capacity {
		q.size++
		q.rear++

		if q.rear == q.capacity {
			q.rear = 0
		}

		q.elements[q.rear] = elem

		return true
	}

	return false
}

//----------------------------------------------- Dequeue
func (q *Queue) Dequeue() (interface{}, bool) {
	if q.size > 0 {
		ret := q.elements[q.front]

		q.size--
		q.front++

		if q.front == q.capacity {
			q.front = 0
		}
		return ret, true
	}

	return nil, false
}

//----------------------------------------------- return queue
func (q *Queue) All() (all []interface{}) {
	all = make([]interface{}, q.size)

	count := q.size
	idx := q.front

	for count > 0 {
		all[q.size-count] = q.elements[idx]

		idx++
		if idx >= q.capacity {
			idx = 0
		}

		count--
	}

	return
}
