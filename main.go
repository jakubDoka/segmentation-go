package main

import (
	"libs/debugm"
	"libs/events"
	"libs/graphics/animations"
	"libs/graphics/g2d/colors"
	"libs/netm"
	"libs/textbox"
	"libs/threads"
	"math"
	"math/rand"
	"myFirstProject/game"

	"myFirstProject/game/content"
	"myFirstProject/game/content/bullets"
	"myFirstProject/game/content/effects"
	"myFirstProject/game/content/modules/turrets"
	"myFirstProject/game/content/segments"
	"myFirstProject/game/content/worms"
	"myFirstProject/game/gui"
	"myFirstProject/game/network"
	"myFirstProject/game/network/client"
	"myFirstProject/game/network/server"
	"myFirstProject/game/world/ent/ai"

	"myFirstProject/game/world/graphics"
	"myFirstProject/game/world/tilemap"
	"time"

	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

//launcher
func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Segmentation",
		Bounds: pixel.R(0, 0, 1000, 600),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	game.Win = win
	if err != nil {
		panic(err)
	}

	content.Load()
	graphics.Cam = graphics.Camera{
		Scale:       pixel.V(1, 1),
		Vec:         pixel.ZV,
		MinZ:        .001,
		MaxZ:        2,
		Sensitivity: .2,
	}
	text := textbox.NewDrawer(pixel.ZV, text.Atlas7x13)
	ui := gui.New(content.UiSheet.Picture, text)
	s := server.Server{}
	events.Handler.On(network.PlayerConnect{}, func(i interface{}) {
		threads.Queue.Post(func() {
			e := i.(network.PlayerConnect)
			worm := worms.Standard.New(worms.Q.Take(), 0)
			for ii := 0; ii < 10; ii++ {
				seg := segments.Standard.New(pixel.ZV)
				turrets.Standard.New(seg)
				worm.Add(seg)
			}
			worms.Q.AddWorm(worm)
			worm.Client = true
			worm.IP = e.Conn.RemoteAddr().String()
			threads.Queue.Post(func() {
				b := netm.Buffer{}
				b.PutUint8(9)
				b.PutUint16(worm.ID)
				e.Conn.Write(b.LenAndData())
			})
		})
	})

	tBox := textbox.New(text, 100)
	tBox.Pos = game.Win.Bounds().Center()
	tBox.Scl = pixel.V(4, 4)
	tBox.Content = "127.0.0.1:6000"
	ui.AddText(tBox)
	tilemap.World.AddTeam(1)
	tilemap.World.AddTeam(0)

	txt := gui.ButtonText{
		Text:  "HOST",
		Scale: pixel.V(4, 4),
	}

	host := gui.NewButton(
		pixel.V(600, 200),
		pixel.V(1.5, 1.5),
		content.UiSheet.Regs[0],
		content.UiSheet.Regs[1],
		content.UiSheet.Regs[0],
		pixel.ZR,
	)
	host.Txt = txt
	host.Call = func() {

		game.IsServer = true
		err := s.Open(tBox.Content)
		if err != nil {
			panic(err)
		}
		go s.AcceptCons()
		worm := worms.Standard.New(worms.Q.Take(), 1)
		worm.Control.Set(true)
		graphics.Cam.Existence = worm
		for ii := 0; ii < 10; ii++ {
			seg := segments.Standard.New(pixel.ZV)
			turrets.Standard.New(seg)
			worm.Add(seg)
		}

		worm = worms.Standard.New(worms.Q.Take(), 0)
		for ii := 0; ii < 1; ii++ {
			seg := segments.Standard.New(pixel.ZV)
			turrets.Standard.New(seg)
			worm.Add(seg)
		}

		threads.Queue.AddCycle(func() {
			b := netm.Buffer{}
			worms.Q.WriteUpdate(&b)
			network.Clients.UDP.Append(b)
		}, time.Millisecond*100)
		ui.Hidden = true
	}
	ui.Add(host)
	join := gui.NewButton(
		pixel.V(400, 200),
		pixel.V(1.5, 1.5),
		content.UiSheet.Regs[0],
		content.UiSheet.Regs[1],
		content.UiSheet.Regs[0],
		pixel.ZR,
	)
	c := client.Client{}
	txt.Text = "JOIN"
	join.Txt = txt
	join.Call = func() {
		content.Mixer = rand.Float64()
		content.Load()
		game.IsNetworking = true

		c.Connect(tBox.Content)
		go c.ConnectUDP(3, time.Second*3)
		ui.Hidden = true
	}
	ui.Add(join)

	//sd := debugm.New()

	fp := debugm.FpsPrinter{Freq: 1}
	fp.Start()
	now := time.Now()
	a, b := color.RGBA{255, 0, 0, 100}, color.RGBA{0, 255, 0, 100}

	for !win.Closed() {
		game.Delta = time.Since(now).Seconds()

		now = time.Now()

		threads.Queue.Run()

		if game.IsServer && game.Halted == 0 {
			network.Clients.SendPackets(s.UDPConn)
		}

		if game.IsNetworking {
			client.UpdateQueue.Run()
			game.Alfa = math.Min(time.Since(game.LastUpdate).Seconds()/game.UpdateSpacing, 2)

			network.Server.SendPackets(c.TCPConn, c.UDPConn)
		}

		ai.Scanner.Update()

		fp.Count++

		animations.Queue.Update(game.Delta)

		graphics.Cam.Update()

		effects.Update(game.Delta)

		bullets.Q.Update()

		worms.Q.Update()

		win.Clear(colors.BlendWithT(a, b, .28))
		ui.Draw(win, graphics.Matrix)
		text.Clear()
		content.Draw()
		win.Update()
	}
}

func main() {
	//starter
	pixelgl.Run(run)
}
