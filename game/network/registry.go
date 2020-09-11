package network

import (
	"libs/netm"
	"net"
	"sync"
)

// Registry is thread safe map of connections
type Registry struct {
	sync.Mutex
	TCP, UDP netm.Buffer
	Conns    map[string]*ClientData
}

// Add adds connection
func (r *Registry) Add(conn *ClientData) {
	r.Lock()
	r.Conns[conn.RemoteAddr().String()] = conn
	r.Unlock()
}

// Remove removes the connection
func (r *Registry) Remove(addr string) {
	r.Lock()
	delete(r.Conns, addr)
	r.Unlock()
}

// Get gets the connection
func (r *Registry) Get(addr string) net.Conn {
	r.Lock()
	defer r.Unlock()
	return r.Conns[addr]
}

// Values returns slice of map values
func (r *Registry) Values() []*ClientData {
	values := []*ClientData{}

	r.Lock()
	for _, v := range r.Conns {
		values = append(values, v)
	}
	r.Unlock()

	return values
}

// Keys returns slice of map keys
func (r *Registry) Keys() []string {
	values := []string{}

	r.Lock()
	for k := range r.Conns {
		values = append(values, k)
	}
	r.Unlock()

	return values
}

// TCPAppendExcept appends packet to queue of all client except one with the address
func (r *Registry) TCPAppendExcept(b netm.Buffer, address string) {
	r.Lock()
	for ip, c := range r.Conns {
		if address != ip {
			c.TCP.Append(b)
		}
	}
	r.Unlock()
}

// TCPAppendOnly appends packet only to client with given address
func (r *Registry) TCPAppendOnly(b netm.Buffer, address string) {
	r.Lock()
	c, ok := r.Conns[address]
	r.Unlock()
	if ok {
		c.TCP.Append(b)
	}
}

// UDPAppendExcept appends packet to queue of all client except one with the address
func (r *Registry) UDPAppendExcept(b netm.Buffer, address string) {
	r.Lock()
	for ip, c := range r.Conns {
		if address != ip {
			c.UDP.Append(b)
		}
	}
	r.Unlock()
}

// UDPAppendOnly appends packet only to client with given address
func (r *Registry) UDPAppendOnly(b netm.Buffer, address string) {
	r.Lock()
	c, ok := r.Conns[address]
	r.Unlock()
	if ok {
		c.UDP.Append(b)
	}
}

// SendPackets sends all packets to clients
func (r *Registry) SendPackets(udp *net.UDPConn) {
	for _, c := range r.Values() {
		c.SendPackets(udp)
	}
	r.Clear()
}

// Clear clears buffers
func (r *Registry) Clear() {
	r.TCP.Clear()
	r.UDP.Clear()
}
