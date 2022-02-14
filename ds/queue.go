// Package ds Queue implementation of queue in net/http/transport
// it's not concurrency safe
package ds

type Queue struct {
	headPos int
	head    []int
	tail    []int
}

func (q *Queue) Len() int {
	return len(q.head) - q.headPos + len(q.tail)
}

func (q *Queue) PushBack(v int) {
	q.tail = append(q.tail, v)
}

// PopFront return element, isValid
func (q *Queue) PopFront() (int, bool) {
	if q.headPos >= len(q.head) {
		// the head part is empty now
		if len(q.tail) == 0 {
			return 0, false
		}
		q.head, q.headPos, q.tail = q.tail, 0, q.head[:0]
	}

	front := q.head[q.headPos]
	q.headPos++
	return front, true
}

func (q *Queue) Peek() (int, bool) {
	if q.headPos < len(q.head) {
		return q.head[q.headPos], true
	}

	if len(q.tail) > 0 {
		return q.tail[0], true
	}

	return 0, false
}
