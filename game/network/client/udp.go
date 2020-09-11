package client

import (
	"libs/netm"
	"myFirstProject/game/content/worms"
	"net"
)

// UDPResponces is lit of callbacks that are used to react to
// server packages
var UDPResponces = []func(b *netm.Buffer){
	// Sync worms
	func(b *netm.Buffer) {
		panic("this should not happen")
	},
	worms.Q.ReadUpdate,
}

// UDP is net.UDPConn wrapper
type UDP struct {
	*net.UDPConn
	*net.UDPAddr
}

// NewUDP returns new connected udp or error if connection failed
func NewUDP(addres string) (*UDP, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", addres)
	if err != nil {
		return nil, err
	}

	udp := UDP{UDPAddr: udpAddr}
	udp.UDPConn, err = net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &udp, nil
}

func (u *UDP) Read() {
	var (
		buf      [10000]byte
		leftover *netm.Buffer
		supposed uint32
		remain   []byte
	)
	for {

		n, err := u.UDPConn.Read(buf[0:])
		if err != nil {
			return
		}
		b := netm.New(buf[:n])
		var packets []*netm.Buffer
		packets, leftover, supposed, remain = b.Split(leftover, supposed, remain)
		if len(packets) == 0 {
			continue
		}

		/*threads.Queue.Post(func() {
			game.UpdateSpacing = time.Since(game.LastUpdate).Seconds()
			game.LastUpdate = time.Now()
		})*/

		UpdateQueue.Set(func() { UDPResponces[packets[0].Uint8()](packets[0]) })
	}
}
