package server

import (
	"encoding/binary"
)

const (
	messageTypeHello uint8 = iota + 1
	messageTypeAck
	messageTypeData
)

type Message interface {
	Type() uint8
	Marshal() ([]byte, error)
}

type helloMessage struct {
}

func (m *helloMessage) Type() uint8 {
	return messageTypeHello
}

func (m helloMessage) Marshal() ([]byte, error) {
	return []byte{m.Type()}, nil
}

type ackMessage struct {
	AckID uint32
}

func (m *ackMessage) Type() uint8 {
	return messageTypeAck
}

func (m *ackMessage) Marshal() ([]byte, error) {
	var b [5]byte
	b[0] = m.Type()
	binary.BigEndian.PutUint32(b[1:], m.AckID)
	return b[:], nil
}

type dataMessage struct {
	SeqID uint32
	MsgID uint32

	Data []byte
}

func (m *dataMessage) Type() uint8 {
	return messageTypeData
}

func (m *dataMessage) Marshal() ([]byte, error) {
	size := len(m.Data) + 9
	b := make([]byte, size)
	b[0] = m.Type()
	binary.BigEndian.PutUint32(b[1:], m.SeqID)
	binary.BigEndian.PutUint32(b[5:], m.MsgID)
	copy(b[9:], m.Data)

	return b, nil
}
