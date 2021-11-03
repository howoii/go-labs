package lru

import (
	"testing"
)

const count = 100

var testCases = []struct {
	keySet interface{}
	keyGet interface{}
	result bool
}{
	{"key", "key", true},
	{"keyA", "keY", false},
	{"keyB", "key", false},
	{"keyA", "keyB", true},
}

func TestCache_Add(t *testing.T) {
	c := New(2, nil)
	for _, v := range testCases {
		c.Add(v.keySet, 111)
		val, ok := c.Get(v.keyGet)
		if ok != v.result {
			t.Errorf("%v, cache hit = %v; want %v", v.keyGet, ok, v.result)
		} else if ok && val != 111 {
			t.Errorf("%v cache value = %v; want 111", v.keyGet, val)
		}
	}
}

func BenchmarkCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := New(10, nil)
		for k := 0; k < count; k++ {
			c.Add(k, k)
		}
	}
}

func BenchmarkCachePool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := New(10, &Option{UsePool: true})
		for k := 0; k < count; k++ {
			c.Add(k, k)
		}
	}
}
