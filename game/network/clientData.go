package network

import (
	"libs/netm"
	"net"
)

// ClientData is just data container
type ClientData struct {
	net.Conn
	*net.UDPAddr
	TCP, UDP netm.Buffer
	name, ip string
	ID       uint16
}

// SendPackets sends all packets to client tis data belongs to
func (c *ClientData) SendPackets(udp *net.UDPConn) {
	c.Write(append(c.TCP.Data(), Clients.TCP.Data()...))
	udp.WriteToUDP(append(c.UDP.Data(), Clients.UDP.Data()...), c.UDPAddr)
	c.Clear()
}

// Clear clears buffers
func (c *ClientData) Clear() {
	c.TCP.Clear()
	c.UDP.Clear()
}
