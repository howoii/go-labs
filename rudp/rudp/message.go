package rudp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	messageTypeHello uint8 = iota + 1
	messageTypeAck
	messageTypeData
	messageTypeBye
)

type MessageType uint8

func (m MessageType) String() string {
	switch uint8(m) {
	case messageTypeHello:
		return "HELLO"
	case messageTypeAck:
		return "ACK"
	case messageTypeData:
		return "DATA"
	case messageTypeBye:
		return "BYE"
	}
	return ""
}

type Message interface {
	Type() uint8
	Marshal() []byte
	Unmarshal([]byte) error
	String() string
}

type SequenceMsg interface {
	Message
	GetSeqID() uint32
	SetSeqID(seqID uint32)
}

type timeoutMessage struct {
	msg SequenceMsg

	timeout  time.Duration
	deadline time.Time

	errChan chan error
}

func (tm *timeoutMessage) Finish(err error) {
	if tm.errChan == nil {
		return
	}
	if err != nil {
		select {
		case tm.errChan <- err:
		default:
		}
	} else {
		close(tm.errChan)
	}
}

type HelloMessage struct {
	SeqID uint32
	AckID uint32
}

func (m *HelloMessage) Type() uint8 {
	return messageTypeHello
}

func (m *HelloMessage) Marshal() []byte {
	var b [9]byte
	b[0] = m.Type()
	binary.BigEndian.PutUint32(b[1:], m.SeqID)
	binary.BigEndian.PutUint32(b[5:], m.AckID)

	return b[:]
}

func (m *HelloMessage) Unmarshal(data []byte) error {
	if len(data) != 8 {
		return errors.New("message: invalid length of AckMessage")
	}
	m.SeqID = binary.BigEndian.Uint32(data)
	m.AckID = binary.BigEndian.Uint32(data[4:])

	return nil
}

func (m *HelloMessage) String() string {
	return fmt.Sprintf("%s, Seq: %d, Ack: %d", MessageType(m.Type()), m.SeqID, m.AckID)
}

func (m *HelloMessage) GetSeqID() uint32 {
	return m.SeqID
}

func (m *HelloMessage) SetSeqID(seqID uint32) {
	m.SeqID = seqID
}

type AckMessage struct {
	AckID uint32
}

func (m *AckMessage) Type() uint8 {
	return messageTypeAck
}

func (m *AckMessage) Marshal() []byte {
	var b [5]byte
	b[0] = m.Type()
	binary.BigEndian.PutUint32(b[1:], m.AckID)
	return b[:]
}

func (m *AckMessage) Unmarshal(data []byte) error {
	if len(data) != 4 {
		return errors.New("message: invalid length of AckMessage")
	}
	m.AckID = binary.BigEndian.Uint32(data)
	return nil
}

func (m *AckMessage) String() string {
	return fmt.Sprintf("%s, Ack: %d", MessageType(m.Type()), m.AckID)
}

type DataMessage struct {
	SeqID uint32

	Data []byte
}

func (m *DataMessage) Type() uint8 {
	return messageTypeData
}

func (m *DataMessage) Marshal() []byte {
	size := len(m.Data) + 5
	b := make([]byte, size)
	b[0] = m.Type()
	binary.BigEndian.PutUint32(b[1:], m.SeqID)
	copy(b[5:], m.Data)

	return b
}

func (m *DataMessage) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return errors.New("message: invalid length of DataMessage")
	}
	m.SeqID = binary.BigEndian.Uint32(data)
	m.Data = make([]byte, len(data[4:]))
	copy(m.Data, data[4:])

	return nil
}

func (m *DataMessage) String() string {
	return fmt.Sprintf("%s, Seq: %d, Data: %v", MessageType(m.Type()), m.SeqID, string(m.Data))
}

func (m *DataMessage) GetSeqID() uint32 {
	return m.SeqID
}

func (m *DataMessage) SetSeqID(seqID uint32) {
	m.SeqID = seqID
}

type ByeMessage struct {
}

func (m *ByeMessage) Type() uint8 {
	return messageTypeBye
}

func (m *ByeMessage) Marshal() []byte {
	var b [1]byte
	b[0] = m.Type()
	return b[:]
}

func (m *ByeMessage) Unmarshal(data []byte) error {
	return nil
}

func (m *ByeMessage) String() string {
	return fmt.Sprintf("%s", MessageType(m.Type()))
}
