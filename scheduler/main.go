package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Printf("num of CPU: %d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		// watcher goroutine
		time.Sleep(1 * time.Second)
		runtime.Gosched()
		fmt.Println("watcher goroutine got scheduled")
		wg.Done()
	}()

	go func() {
		// loop without function call
		for {
		}
	}()
	wg.Wait()
}
