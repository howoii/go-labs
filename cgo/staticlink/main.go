// show using cgo by static linking
package main

import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

//#cgo CFLAGS: -I./file
//#cgo LDFLAGS: -L/Users/haiwei.zhang/Documents/go/labs/cgo/staticlink/file -lfile
//#include "file.h"
import "C"

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
