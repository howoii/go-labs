package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"
)

func main() {
	addr, err := net.ResolveUDPAddr("", ":8848")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err, string(debug.Stack()))
	}
	fmt.Printf("local address: %s\n", conn.LocalAddr().String())
	var buf [1024]byte
	for {
		n, err := os.Stdin.Read(buf[:])
		if n > 0 {
			_, err := conn.Write(buf[:n])
			if err != nil {
				log.Fatal(err, string(debug.Stack()))
			}
			var buf [1024]byte
			n, remoteAddr, err := conn.ReadFromUDP(buf[:])
			if err != nil {
				log.Fatal(err, string(debug.Stack()))
			}
			fmt.Println(remoteAddr.String())
			fmt.Println(string(buf[:n]))
		}
		if err != nil {
			log.Fatal(err, string(debug.Stack()))
		}
	}
}
