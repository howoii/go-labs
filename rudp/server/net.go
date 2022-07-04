package server

import (
	"log"
	"net"
	"time"
)

const (
	MTU int = 1420
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
		conn:       conn,
		errChan:    make(chan error),
		packetChan: make(chan *packet),
	}
	l.startListen()
	return l, nil
}

// Listener server listener
type Listener struct {
	conn *net.UDPConn

	errChan    chan error
	packetChan chan *packet
}

// Accept implements the net.Listener Accept method.
func (l *Listener) Accept() (net.Conn, error) {
	return nil, nil
}

// Close implements the net.Listener Close method.
func (l *Listener) Close() error {
	return nil
}

// Addr implements the net.Listener Addr method.
func (l *Listener) Addr() net.Addr {
	return nil
}

func (l *Listener) startListen() {

}

func (l *Listener) readLoop() {
	for {
		var buf [MTU]byte
		n, addr, err := l.conn.ReadFromUDP(buf[:])
		if n > 0 {
			l.packetChan <- &packet{
				addr: addr,
				data: buf[:n],
			}
		}
		if err != nil {
			log.Printf("ReadFromUDP error: %v\n", err)
		}
	}
}

func (l *Listener) writeLoop() {

}

// Conn a udp connection
type Conn struct {
}

// Read implements the net.Conn Read method.
func (c *Conn) Read(b []byte) (int, error) {
	return len(b), nil
}

// Write implements the net.Conn Write method.
func (c *Conn) Write(b []byte) (int, error) {
	return len(b), nil
}

// Close implements the net.Conn Close method.
func (c *Conn) Close() error {
	return nil
}

// LocalAddr implements the net.Conn LocalAddr method.
func (c *Conn) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr implements the net.Conn RemoteAddr method.
func (c *Conn) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline implements the net.Conn SetDeadline method.
func (c Conn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements the net.Conn SetReadDeadline method.
func (c Conn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implements the net.Conn SetWriteDeadline method.
func (c Conn) SetWriteDeadline(t time.Time) error {
	return nil
}
