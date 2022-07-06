package rudp

import (
	"errors"
	"fmt"
	"net"
)

type packet struct {
	addr *net.UDPAddr // remote address
	data []byte
}

func (p *packet) parse() (Message, error) {
	if len(p.data) < 1 {
		return nil, errors.New("packet: empty packet is not valid")
	}
	var m Message
	switch p.data[0] {
	case messageTypeAck:
		m = &AckMessage{}
	case messageTypeHello:
		m = &HelloMessage{}
	case messageTypeData:
		m = &DataMessage{}
	case messageTypeBye:
		m = &ByeMessage{}
	}
	if m != nil {
		err := m.Unmarshal(p.data[1:])
		return m, err
	}
	return nil, fmt.Errorf("packet: unexpected message type %v", p.data[0])
}
