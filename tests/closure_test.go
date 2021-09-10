package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestClosureVariable(t *testing.T) {
	v := 1
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Wait()
		fmt.Println(v)
	}()

	v = 2
	wg.Done()
}
