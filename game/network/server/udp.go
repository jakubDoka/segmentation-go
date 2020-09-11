package server

import (
	"net"
)

// UDP is wrapper for UDPConn
type UDP struct {
	*net.UDPConn
}

// NewUDP is UDP constructor
func NewUDP(addres string) (*UDP, error) {

	udpAddr, err := net.ResolveUDPAddr("udp4", addres)
	if err != nil {
		return nil, err
	}

	udp := UDP{}

	udp.UDPConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	return &udp, nil
}
