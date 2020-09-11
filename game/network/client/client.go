package client

import (
	"libs/netm"
	"libs/threads"
	"time"
)

// UpdateQueue for calling the client update.
var UpdateQueue = threads.Call{}

// Client is for handling server responces and sending player state and input
type Client struct {
	*UDP
	*TCP
}

// Connect connects client to network
func (c *Client) Connect(addres string) error {
	udp, err := NewUDP(addres)
	if err != nil {
		return err
	}

	tcp, err := NewTCP(addres)
	if err != nil {
		return err
	}

	c.UDP = udp
	c.TCP = tcp

	return nil
}

// Disconnect disconnects the client
func (c *Client) Disconnect() {
	c.UDP.Close()
	c.TCP.Close()
}

// ConnectUDP tries to make udp connection with server. If this fails client will disconnect.
func (c *Client) ConnectUDP(attempts int, timeout time.Duration) {
	failed := true
	b := netm.Buffer{}
	b.PutUint8(0)
	b.AddLen()
	for i := 0; i < attempts; i++ {
		c.UDP.Write(b.Data())
		if c.Checkout(timeout) {
			failed = false
			break
		}
	}

	if failed {
		c.Disconnect()
		return
	}

	c.TCP.SetReadDeadline(time.Time{})

	go c.TCP.Read()
	go c.UDP.Read()
}
