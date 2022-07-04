package server

import (
	"net"
)

type packet struct {
	addr *net.UDPAddr
	data []byte
}
