package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"
	"time"

	"github.com/labs/rudp/rudp"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8848")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	hello := &rudp.HelloMessage{
		SeqID: 0,
	}
	_, err = conn.Write(hello.Marshal())
	if err != nil {
		fmt.Printf("send hello failed: %v\n", err)
	}

	time.Sleep(1)

	ack := &rudp.AckMessage{
		AckID: 2,
	}
	_, err = conn.Write(ack.Marshal())
	if err != nil {
		fmt.Printf("send ack failed: %v\n", err)
	}

	go readLoop(conn)
	go writeLoop(conn)

	select {}
}

func readLoop(conn *net.UDPConn) {
	buf := make([]byte, 1420)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Recv: %s(%n)\n", string(buf[:n]), n)
	}
}

func writeLoop(conn *net.UDPConn) {
	var buf [1024]byte
	seqID := uint32(1)
	for {
		n, err := os.Stdin.Read(buf[:])
		if n > 0 {
			msg := rudp.DataMessage{
				SeqID: seqID,
				Data:  buf[:n],
			}
			_, err := conn.Write(msg.Marshal())
			if err != nil {
				log.Fatal(err, string(debug.Stack()))
			}
		}
		if err != nil {
			log.Fatal(err, string(debug.Stack()))
		}
	}
}
