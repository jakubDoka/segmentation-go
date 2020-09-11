package network

import "net"

// PlayerConnect is event emmyted when player connects
type PlayerConnect struct {
	Conn net.Conn
}
