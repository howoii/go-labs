package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/labs/rudp/rudp"
)

type State struct {
	sendSeq uint32

	sendChan chan rudp.Message
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8848")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	state := &State{
		sendSeq:  rand.Uint32(),
		sendChan: make(chan rudp.Message, 1024),
	}
	hello := &rudp.HelloMessage{
		SeqID: state.sendSeq,
	}
	state.sendChan <- hello

	go writeLoop(conn, state)
	go readLoop(conn, state)
	go ioLoop(state)

	select {}
}

func readLoop(conn *net.UDPConn, state *State) {
	buf := make([]byte, 1420)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		m, err := rudp.ParsePacket(buf[:n])
		if err != nil {
			log.Printf("parse packet err: %v\n", err)
		}
		fmt.Printf("Recv: %s\n", m)
		if sm, ok := m.(rudp.SequenceMsg); ok {
			state.sendChan <- &rudp.AckMessage{
				AckID: sm.GetSeqID(),
			}
		}
	}
}

func ioLoop(state *State) {
	var buf [1024]byte
	for {
		n, err := os.Stdin.Read(buf[:])
		if n > 0 {
			data := bytes.TrimSuffix(buf[:n], []byte{'\n'})
			if len(data) == 0 {
				continue
			}
			msg := &rudp.DataMessage{
				SeqID: atomic.AddUint32(&state.sendSeq, 1),
				Data:  data,
			}
			state.sendChan <- msg
		}
		if err != nil {
			log.Fatal(err, string(debug.Stack()))
		}
	}
}

func writeLoop(conn *net.UDPConn, state *State) {
	for {
		msg := <-state.sendChan
		_, err := conn.Write(msg.Marshal())
		if err != nil {
			log.Fatal(err, string(debug.Stack()))
		}
		fmt.Printf("Send: %s\n", msg)
	}
}
