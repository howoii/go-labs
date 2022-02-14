package ds

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := Queue{}
	for i := 0; i < 10; i++ {
		q.PushBack(i)
	}
	t.Logf("len: %d", q.Len())

	for q.Len() > 0 {
		t.Log(q.PopFront())
	}
}
