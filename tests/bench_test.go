package main

import "testing"

func BenchmarkAppend(b *testing.B) {
	src := make([]int, 50)
	for i := 0; i < b.N; i++ {
		dst := make([]int, 1)
		dst = append(dst, src...)
	}
}

func BenchmarkCopy(b *testing.B) {
	src := make([]int, 50)
	for i := 0; i < b.N; i++ {
		dst := make([]int, 1)
		tmp := make([]int, len(src)+len(dst))
		copy(tmp, dst)
		copy(tmp[len(dst):], src)
		dst = tmp
	}
}

func BenchmarkBad(b *testing.B) {
	const size = 50
	for i := 0; i < b.N; i++ {
		data := make([]int, 0)
		for k := 0; k < size; k++ {
			data = append(data, k)
		}
	}
}

func BenchmarkGood(b *testing.B) {
	const size = 50
	for i := 0; i < b.N; i++ {
		data := make([]int, 0, size)
		for k := 0; k < size; k++ {
			data = append(data, k)
		}
	}
}
