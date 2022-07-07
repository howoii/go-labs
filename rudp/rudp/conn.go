package rudp

import (
	"errors"
	"math"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type connState uint32

const (
	connStatePending connState = iota + 1
	connStateEstablish
)

const (
	halfUint32     = math.MaxUint32 >> 1
	initialTimeout = time.Second * 3
	buffSize       = 1 << 13
)

type SeqID uint32

func (id *SeqID) Init(v uint32) {
	if v != 0 {
		*id = SeqID(v)
	} else {
		*id = SeqID(rand.Uint32())
	}
}

func (id *SeqID) Next() uint32 {
	return atomic.AddUint32((*uint32)(id), 1)
}

func (id *SeqID) Curr() uint32 {
	return atomic.LoadUint32((*uint32)(id))
}

func (id SeqID) Before(oid SeqID) bool {
	return (id < oid && (oid-id) < halfUint32) || (id > oid && (id-oid) > halfUint32)
}

// Conn a rudp connection
type Conn struct {
	l     *Listener
	rAddr *net.UDPAddr // remote address

	state     connState
	closeFlag uint32

	nextSeqID SeqID

	ackMu    sync.Mutex // protect ackQueue
	ackQueue map[uint32]*timeoutMessage

	sendBuf chan *timeoutMessage
	recvBuf chan *DataMessage

	sortMu     sync.Mutex // protect following fields
	nextRecvID SeqID
	sortQueue  *MHeap
}

type connParams struct {
	RecvSeq uint32
}

func newConn(l *Listener, addr *net.UDPAddr, params *connParams) *Conn {
	c := &Conn{
		l:     l,
		rAddr: addr,
		state: connStatePending,

		ackQueue: make(map[uint32]*timeoutMessage),
	}
	c.nextSeqID.Init(0)
	c.nextRecvID.Init(params.RecvSeq)
	go c.processTimeout()
	return c
}

// Read implements the net.Conn Read method.
func (c *Conn) Read(b []byte) (int, error) {
	if c.isClosing() {
		return 0, ErrReadClosed
	}

	var m *DataMessage
	select {
	case m = <-c.recvBuf:
		goto got
	default:
		c.fillRecvBuff()
	}
	m = <-c.recvBuf
got:
	if m == nil {
		return 0, ErrReadClosed
	}
	n := copy(b, m.Data)
	return n, nil
}

// Write implements the net.Conn Write method.
func (c *Conn) Write(b []byte) (int, error) {
	if len(b) > MTU {
		return 0, ErrMaxMTU
	}
	if c.isClosing() {
		return 0, ErrWriteClosed
	}
	msg := &DataMessage{
		SeqID: c.nextSeqID.Next(),
		Data:  append([]byte{}, b...),
	}
	tm := c.sendMessage(msg)
	err := <-tm.errChan
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// Close implements the net.Conn Close method.
func (c *Conn) Close() error {
	if atomic.CompareAndSwapUint32(&c.closeFlag, 0, 1) {
		// do some cleanup work
	}
	return nil
}

// LocalAddr implements the net.Conn LocalAddr method.
func (c *Conn) LocalAddr() net.Addr {
	return c.l.addr
}

// RemoteAddr implements the net.Conn RemoteAddr method.
func (c *Conn) RemoteAddr() net.Addr {
	return c.rAddr
}

// SetDeadline implements the net.Conn SetDeadline method.
func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements the net.Conn SetReadDeadline method.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implements the net.Conn SetWriteDeadline method.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *Conn) isClosing() bool {
	return atomic.LoadUint32(&c.closeFlag) == 1
}

func (c *Conn) start() {
	c.recvBuf = make(chan *DataMessage, buffSize)
	c.sortQueue = newMHeap()
	go c.processSendBuf()
}

func (c *Conn) removeAckQueue(msg *AckMessage) *timeoutMessage {
	ackID := msg.AckID
	c.ackMu.Lock()
	defer c.ackMu.Unlock()

	tm := c.ackQueue[ackID]
	delete(c.ackQueue, ackID)
	return tm
}

func (c *Conn) putAckQueue(tm *timeoutMessage) {
	c.ackMu.Lock()
	defer c.ackMu.Unlock()
	tm.deadline = time.Now().Add(tm.timeout)
	c.ackQueue[tm.msg.GetSeqID()] = tm
}

func (c *Conn) receiveData(msg *DataMessage) error {
	c.sortMu.Lock()
	defer c.sortMu.Unlock()
	if SeqID(msg.SeqID).Before(c.nextRecvID) {
		return errors.New("receive expired message")
	}
	if c.sortQueue.Has(msg) {
		return errors.New("receive repeat message")
	}
	c.sortQueue.Push(msg)
	return nil
}

func (c *Conn) fillRecvBuff() {
	c.sortMu.Lock()
	defer c.sortMu.Unlock()

	for c.sortQueue.Len() > 0 {
		msg := c.sortQueue.Top()
		if msg.SeqID != c.nextRecvID.Curr() {
			return
		}
		select {
		case c.recvBuf <- msg:
			c.sortQueue.Pop()
			c.nextRecvID.Next()
		default:
			return
		}
	}
}

func (c *Conn) ack(ackID uint32) {
	m := &AckMessage{
		AckID: ackID,
	}
	c.sendMessage(m)
}

// chanMessage put message into listener's send channel
func (c *Conn) chanMessage(msg Message) {
	c.l.sendChan <- &packet{
		addr: c.rAddr,
		data: msg.Marshal(),
	}
}

func (c *Conn) sendMessage(msg Message) *timeoutMessage {
	if sm, ok := msg.(SequenceMsg); ok {
		return c.putMessage(sm, initialTimeout)
	}
	c.chanMessage(msg)
	return nil
}

func (c *Conn) putMessage(m SequenceMsg, timeout time.Duration) *timeoutMessage {
	tm := &timeoutMessage{
		msg:     m,
		timeout: timeout,
		errChan: make(chan error),
	}
	if c.sendBuf != nil {
		c.sendBuf <- tm
		return tm
	} else {
		c.chanMessage(m)
		c.putAckQueue(tm)
		return tm
	}
}

func (c *Conn) reSendMessage(tm *timeoutMessage) {
	tm.timeout *= 2
	c.chanMessage(tm.msg)
	tm.deadline = time.Now().Add(tm.timeout)
}

func (c *Conn) processSendBuf() {
	c.sendBuf = make(chan *timeoutMessage, buffSize)
	for {
		tm := <-c.sendBuf
		c.chanMessage(tm.msg)
		c.putAckQueue(tm)
	}
}

func (c *Conn) processTimeout() {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		if c.isClosing() {
			return
		}
		<-tick.C

		c.ackMu.Lock()
		now := time.Now()
		for _, v := range c.ackQueue {
			if v.deadline.Before(now) {
				c.reSendMessage(v)
			}
		}
		c.ackMu.Unlock()
	}
}
