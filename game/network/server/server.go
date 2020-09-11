package server

import (
	"fmt"
	"libs/events"
	"libs/netm"
	"libs/threads"
	"myFirstProject/game"
	"myFirstProject/game/content/worms"
	"myFirstProject/game/network"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

// TCPResponces is lit of callbacks that are used to react to
// client packages packages
var TCPResponces = []func(b *netm.Buffer){
	func(b *netm.Buffer) { fmt.Println("placeholder in server") },

	// updating input state of clients worm
	worms.Q.Goto,
	worms.Q.Undeploy,
}

// UDPResponces is lit of callbacks that are used to react to
// client UDP packets
var UDPResponces = []func(b *netm.Buffer){
	// Sync worms
	func(b *netm.Buffer) {
		panic("this should not happen")
	},
	worms.Q.ReadInput,
	worms.Q.SetVel,
}

// Server manages client connections and sends the packets
type Server struct {
	*UDP
	*TCP
	UDPAddr
	Opened bool
}

// Open launches the server but don't accept any connections yet
func (s *Server) Open(addres string) error {
	udp, err := NewUDP(addres)
	if err != nil {
		return err
	}

	tcp, err := NewTCP(addres)
	if err != nil {
		return err
	}

	s.UDP = udp
	s.TCP = tcp
	s.Opened = true

	return nil
}

// AcceptCons starts accept loop so player can connect now
func (s *Server) AcceptCons() {
	go s.listenUDP()
	for s.Opened {
		conn, err := s.Accept()
		if err != nil {
			continue
		}

		go s.AcceptClient(conn)

	}
}

// Sync sync the state of game of client
func (s *Server) Sync(cd network.ClientData) {
	threads.Queue.Post(func() {
		b := netm.Buffer{}
		worms.Q.Write(&b)
		cd.Write(b.LenAndData())

		network.Clients.Add(&cd)

		atomic.AddInt32(&game.Halted, -1)

		events.Handler.Fire(network.PlayerConnect{Conn: cd.Conn})
	})
}

// AcceptClient ensures udp connection or closes tcp connection if attempt fails
func (s *Server) AcceptClient(conn net.Conn) {
	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]

	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 30)
		if !s.Empty() {
			addr := s.Take()
			if ip != addr.IP.String() { // If two clients connect at the same time tis might happen
				continue
			}
			cd := network.ClientData{
				Conn:    conn,
				UDPAddr: addr,
			}
			conn.Write([]byte{0}) // Confirming connection with client

			atomic.AddInt32(&game.Halted, 1)
			s.Sync(cd)

			go s.listen(conn)
			return
		}
	}
	conn.Close()

}

func (s *Server) listen(conn net.Conn) {
	var (
		buf      [10000]byte
		leftover *netm.Buffer
		supposed uint32
		remain   []byte
	)
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			network.Clients.Remove(conn.RemoteAddr().String())
			conn.Close()
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

// Read starts listening loop
func (s *Server) listenUDP() {
	var (
		buf      [10000]byte
		leftover *netm.Buffer
		supposed uint32
		remain   []byte
	)

	for {
		n, addr, err := s.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}

		b := netm.New(buf[:n])

		var packets []*netm.Buffer
		packets, leftover, supposed, remain = b.Split(leftover, supposed, remain)

		if len(packets) == 0 {
			continue
		}

		threads.Queue.Post(func() {
			for _, p := range packets {
				id := p.Uint8()
				if id == 0 {
					s.Put(addr)
					continue
				}
				UDPResponces[id](p)
			}
		})

	}
}
