package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/labs/rudp/rudp"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	var buf [1420]byte

	serverMsg := "Server Message [%d]"
	loop := 0
	for {
		loop++
		msg := fmt.Sprintf(serverMsg, loop)
		_, err := c.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[%s] Send: %s\n", c.RemoteAddr(), msg)
		n, err := c.Read(buf[:])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[%s] Recv: %s\n", c.RemoteAddr(), string(buf[:n]))
	}
}
