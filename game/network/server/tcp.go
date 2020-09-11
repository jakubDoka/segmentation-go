package server

import (
	"net"
)

// TCP is net.TCPListener wrapper
type TCP struct {
	*net.TCPListener
}

// NewTCP returns listening TCP or error if given port is not free or addres does not exist
func NewTCP(addres string) (*TCP, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addres)
	if err != nil {
		return nil, err
	}

	tcp := TCP{}

	tcp.TCPListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	return &tcp, nil
}

func (t *TCP) Read() {

}
