package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"testing"
)

func PanicFunc() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p)
		}
	}()

	panic("oops!")
}

func PanicGoRoutine() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println(p)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("waiting to panic")
		defer cancel()
		panic("oops! goroutine!")
	}()

	<-ctx.Done()
	fmt.Println("main goroutine finished without panic")
}

func doSomethingWillPanic() {
	var m map[int]int
	m[1] = 1
}

func PanicStack() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				err := fmt.Errorf("%v", p)
				fmt.Println(err, string(debug.Stack()))
			}
			wg.Done()
		}()
		doSomethingWillPanic()
	}()
	wg.Wait()
}

func TestPanic(t *testing.T) {
	PanicStack()
}
