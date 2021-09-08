package main

import (
	"context"
	"fmt"
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

func TestPanic(t *testing.T) {
	PanicFunc()
	PanicGoRoutine()
}
