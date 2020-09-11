package network

import (
	"libs/netm"
	"net"
)

// Host holds buffers used for sending packets to server
type Host struct {
	TCP, UDP netm.Buffer
}

// SendPackets sends all packets to server
func (r *Host) SendPackets(tcp *net.TCPConn, udp *net.UDPConn) {
	tcp.Write(r.TCP.Data())
	udp.Write(r.UDP.Data())
	r.Clear()
}

// Clear clears buffers
func (r *Host) Clear() {
	r.TCP.Clear()
	r.UDP.Clear()
}
