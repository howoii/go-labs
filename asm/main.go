package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

var name = "hello, world"

// sayHello this function is call by Print in assembly
func sayHello(msg string) {
	println(msg)
}

func getGID() int64 {
	buf := make([]byte, 0, 64)
	n := runtime.Stack(buf, false)
	stk := strings.TrimPrefix(string(buf[:n]), "goroutine ")

	id, err := strconv.Atoi(strings.Fields(stk)[0])
	if err != nil {
		panic(fmt.Errorf("can not gt goroutine id: %w", err))
	}
	return int64(id)
}

func GetGoroutineID() int64 {
	g := getg()
	return reflect.ValueOf(g).FieldByName("goid").Int()
}

// Print defined in assembly
func Print()
func SyscallDarwin(fd int, msg string) int
func getg() interface{}

func main() {
	Print()
	SyscallDarwin(syscall.Stdout, "hello, world\n")

	println(getGID())
	println(GetGoroutineID())
}
