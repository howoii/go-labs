// show using cgo by c source code
package main

/*
#include "lib.h"
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// SayHello use 'export' do define C function in GO
//export SayHello
func SayHello(s string) {
	fmt.Println(s)
}

func CString(s string) *C.char {
	sp := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return (*C.char)(unsafe.Pointer(sp.Data))
}

type File struct {
	fd int
}

func OpenFile(filename string) (*File, error) {
	// errno returned as the second return value
	fd, err := C.Open(C.CString(filename))
	if err != nil {
		return nil, err
	}
	return &File{fd: int(fd)}, nil
}

func (f *File) Write(msg string) (int, error) {
	n, err := C.Write(C.int(f.fd), CString(msg), C.int(len(msg)))
	return int(n), err
}

func (f *File) Close() error {
	_, err := C.Close(C.int(f.fd))
	return err
}

func main() {
	C.SayHello("Hello, World!")
	f, err := OpenFile("cgo.log")
	if err != nil {
		fmt.Printf("open file failed: %v\n", err)
	}
	_, err = f.Write("hello log\n")
	if err != nil {
		fmt.Printf("write failed: %v\n", err)
	}
	f.Close()
}
