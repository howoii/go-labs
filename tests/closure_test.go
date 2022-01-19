package main

import (
	"fmt"
	"sync"
	"testing"
)

// result:
// 1: 2
// 2: 1
func TestClosureVariable(t *testing.T) {
	v := 1
	begin := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-begin
		fmt.Printf("1: %d\n", v)
		wg.Done()
	}()

	wg.Add(1)
	go func(i int) {
		<-begin
		fmt.Printf("2: %d\n", i)
		wg.Done()
	}(v)

	v = 2
	close(begin)
	wg.Wait()
}
