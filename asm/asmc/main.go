package main

/*
#include <stdio.h>

int add(int a, int b) {
	return a + b;
}
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/labs/asm/pkg"
)

func main() {
	ret := pkg.CallCFunc(uintptr(unsafe.Pointer(C.add)), 123, 543)
	fmt.Println(ret)
}
