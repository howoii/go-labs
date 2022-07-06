package main

import (
	"fmt"
	"log"
	"net"

	"github.com/labs/rudp/rudp"
)

func main() {
	l, err := rudp.Listen(":8848")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("server listen on: %s\n", l.Addr().String())
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serveConnection(c)
	}
}

func serveConnection(c net.Conn) {
	for {
		_, err := c.Write([]byte("nice to meet you"))
		if err != nil {
			log.Fatal(err)
		}
		var buf [1420]byte
		_, err = c.Read(buf[:])
		if err != nil {
			log.Fatal(err)
		}
	}
}
