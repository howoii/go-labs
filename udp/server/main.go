package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("", ":8848")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("server listens on: %s\n", conn.LocalAddr().String())
	for {
		var buf [1024]byte
		n, remoteAddr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("receive from %v\n", remoteAddr.String())
		_, err = conn.WriteToUDP(buf[:n], remoteAddr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
