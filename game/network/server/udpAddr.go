package server

import (
	"net"
	"sync"
)

// UDPAddr is thread save udp address
type UDPAddr struct {
	sync.Mutex
	addr *net.UDPAddr
}

// Empty returns whether address is nil
func (u *UDPAddr) Empty() bool {
	u.Lock()
	defer u.Unlock()
	return u.addr == nil
}

// Take returns addres and sets it to nil
func (u *UDPAddr) Take() *net.UDPAddr {
	u.Lock()
	defer u.Unlock()
	addr := u.addr
	u.addr = nil
	return addr
}

// Put sets address or panics if its not nil
func (u *UDPAddr) Put(addr *net.UDPAddr) {
	u.Lock()
	if u.addr != nil {
		panic("setting addres that is not empty")
	}
	u.addr = addr
	u.Unlock()
}
