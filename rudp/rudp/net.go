package rudp

import (
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

const (
	MTU     int = 1420
	BACKLOG int = 64
)

func Listen(addr string) (net.Listener, error) {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", a)
	if err != nil {
		return nil, err
	}
	l := &Listener{
		addr:       a,
		conn:       conn,
		errChan:    make(chan error),
		recvChan:   make(chan *packet),
		sendChan:   make(chan *packet),
		acceptChan: make(chan *Conn, BACKLOG),

		connections: make(map[string]*Conn),
	}
	l.startListen()
	return l, nil
}

// Listener server listener
type Listener struct {
	addr *net.UDPAddr
	conn *net.UDPConn

	cmu         sync.Mutex // protect connections
	connections map[string]*Conn

	errChan    chan error
	recvChan   chan *packet
	sendChan   chan *packet
	acceptChan chan *Conn
}

// Accept implements the net.Listener Accept method.
func (l *Listener) Accept() (net.Conn, error) {
	c := <-l.acceptChan
	c.start()
	return c, nil
}

// Close implements the net.Listener Close method.
func (l *Listener) Close() error {
	return nil
}

// Addr implements the net.Listener Addr method.
func (l *Listener) Addr() net.Addr {
	return l.addr
}

func (l *Listener) startListen() {
	go l.readLoop()
	go l.writeLoop()
	go l.processLoop()
}

func (l *Listener) readLoop() {
	for {
		var buf [MTU]byte
		n, addr, err := l.conn.ReadFromUDP(buf[:])
		if n > 0 {
			l.recvChan <- &packet{
				addr: addr,
				data: buf[:n],
			}
		}
		if err != nil {
			log.Printf("readFromUDP error: %v\n", err)
		}
	}
}

func (l *Listener) writeLoop() {
	for {
		pack := <-l.sendChan
		_, err := l.conn.WriteToUDP(pack.data, pack.addr)
		if err != nil {
			log.Printf("writeToUDP error: %v\n", err)
		}
	}
}

func (l *Listener) processLoop() {
	for {
		pack := <-l.recvChan
		msg, err := pack.parse()
		if err != nil {
			log.Printf("packet parse error: %v\n", err)
		}
		switch m := msg.(type) {
		case *HelloMessage:
			err := l.handleNewConnection(pack.addr, m)
			if err != nil {
				log.Printf("process hello message error: %v\n", err)
			}
		case *AckMessage:
			err := l.handleAck(pack.addr, m)
			if err != nil {
				log.Printf("process ack message error: %v\n", err)
			}
		case *DataMessage:
			err := l.handleData(pack.addr, m)
			if err != nil {
				log.Printf("process data message error: %v\n", err)
			}
		case *ByeMessage:
			// TODO
		}
	}
}

func (l *Listener) getConn(addr string) *Conn {
	l.cmu.Lock()
	defer l.cmu.Unlock()
	return l.connections[addr]
}

func (l *Listener) deleteConn(addr string) *Conn {
	l.cmu.Lock()
	defer l.cmu.Unlock()
	c := l.connections[addr]
	delete(l.connections, addr)
	return c
}

func (l *Listener) handleNewConnection(addr *net.UDPAddr, msg *HelloMessage) error {
	l.cmu.Lock()
	defer l.cmu.Unlock()

	if conn := l.connections[addr.String()]; conn != nil {
		return errors.New("unexpect hello message")
	}

	conn := newConn(l, addr, &connParams{
		RecvSeq: msg.SeqID + 1,
	})
	l.connections[addr.String()] = conn

	m := &HelloMessage{
		SeqID: conn.nextSeqID.Next(),
		AckID: msg.SeqID,
	}
	conn.sendMessage(m)

	return nil
}

func (l *Listener) handleAck(addr *net.UDPAddr, msg *AckMessage) error {
	conn := l.getConn(addr.String())
	if conn == nil {
		return errors.New("unexpect ack message, connection not exist")
	}
	tm := conn.removeAckQueue(msg)
	if tm == nil {
		return errors.New("unexpect ack message")
	}
	switch tm.msg.(type) {
	case *HelloMessage:
		conn.state = connStateEstablish
		select {
		case l.acceptChan <- conn:
			return nil
		default:
			l.deleteConn(addr.String())
			conn.Close()
			return errors.New("discard connection since backlog is full")
		}
	default:
		tm.Finish(nil)
	}
	return nil
}

func (l *Listener) handleData(addr *net.UDPAddr, msg *DataMessage) error {
	conn := l.getConn(addr.String())
	if conn == nil {
		return errors.New("unexpect data message, connection not exist")
	}
	if atomic.LoadUint32((*uint32)(&conn.state)) != uint32(connStateEstablish) {
		return errors.New("unexpect data message, connection not established")
	}

	conn.ack(msg.GetSeqID())
	err := conn.receiveData(msg)
	if err != nil {
		conn.fillRecvBuff()
	}
	return err
}
