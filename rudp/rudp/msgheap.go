package rudp

import (
	"container/heap"
)

type IntHeap []uint32

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(uint32))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type MHeap struct {
	h    IntHeap
	dmap map[uint32]*DataMessage
}

func newMHeap() *MHeap {
	m := &MHeap{
		dmap: make(map[uint32]*DataMessage),
	}
	heap.Init(&m.h)
	return m
}

func (m *MHeap) Push(msg *DataMessage) {
	heap.Push(&m.h, msg.SeqID)
	m.dmap[msg.SeqID] = msg
}

func (m *MHeap) Pop() *DataMessage {
	i := heap.Pop(&m.h)
	if i != nil {
		key := i.(uint32)
		msg := m.dmap[key]
		delete(m.dmap, key)
		return msg
	}
	return nil
}

func (m *MHeap) Top() *DataMessage {
	top := m.h[0]
	return m.dmap[top]
}

func (m *MHeap) Has(msg *DataMessage) bool {
	_, ok := m.dmap[msg.SeqID]
	return ok
}

func (m *MHeap) Len() int {
	return len(m.dmap)
}
