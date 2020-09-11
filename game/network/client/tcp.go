package client

import (
	"fmt"
	"libs/netm"
	"libs/threads"
	"myFirstProject/game/content/worms"
	"net"
	"time"
)

// TCPResponces is lit of callbacks that are used to react to
// server packages
var TCPResponces = []func(b *netm.Buffer){

	func(b *netm.Buffer) {
		fmt.Println(b.Data())
	},
	// Sync worms
	worms.Q.Read,
	//Shoot Bullet
	worms.Q.TurretShoot,
	// segment destroyed
	worms.Q.PopSeg,
	// worm destroyed
	worms.Q.Pop,
	// worm segment order change
	worms.Q.MoveSegment,
	// segment deployed
	worms.Q.Deploy,
	// segment undeploy
	worms.Q.Undeploy,
	// spawn worm
	worms.Q.AddWormC,
	// bind player to worm
	worms.Q.BindToWorm,
}

// TCP is net.TCPConn wrapper
type TCP struct {
	*net.TCPConn
}

// NewTCP returns new connected tcp or error if connection failed
func NewTCP(addres string) (*TCP, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", addres)
	if err != nil {
		return nil, err
	}

	tcp := TCP{}

	tcp.TCPConn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return &tcp, nil
}

// Checkout confirms UDP connection
func (t *TCP) Checkout(timeout time.Duration) bool {
	t.SetReadDeadline(time.Now().Add(timeout))
	_, err := t.TCPConn.Read(make([]byte, 10))
	return err == nil
}

func (t *TCP) Read() {
	var (
		leftover *netm.Buffer
		supposed uint32
		remain   []byte
	)
	for {
		var buf [10000]byte
		n, err := t.TCPConn.Read(buf[0:])
		if err != nil {
			return
		}
		b := netm.New(buf[:n])
		var packets []*netm.Buffer
		packets, leftover, supposed, remain = b.Split(leftover, supposed, remain)
		threads.Queue.Post(func() {
			for _, b := range packets {
				TCPResponces[b.Uint8()](b)
			}
		})
	}
}
