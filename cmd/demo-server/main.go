package main

import (
	"log"
	"net"

	"honnef.co/go/wayland/wlserver"
	"honnef.co/go/wayland/wlserver/protocols/wayland"
)

type Keyboard struct {
	RepeatRate  int32
	RepeatDelay int32
}

func (k *Keyboard) Init(id wayland.Keyboard) {
	id.RepeatInfo(k.RepeatRate, k.RepeatDelay)
	id.Keymap(wayland.KeyboardKeymapFormatNoKeymap, 0, 0)
}

func (*Keyboard) Release(wayland.Keyboard)  {}
func (*Keyboard) OnDestroy(wlserver.Object) {}

type Seat struct {
	server   *Server
	resource wayland.Seat
}

func (s *Seat) GetPointer(obj wayland.Seat, id wayland.Pointer) wayland.PointerImplementation {
	panic("not implemented")
}
func (s *Seat) GetTouch(obj wayland.Seat, id wayland.Touch) wayland.TouchImplementation {
	panic("not implemented")
}

func (s *Seat) GetKeyboard(obj wayland.Seat, id wayland.Keyboard) wayland.KeyboardImplementation {
	s.server.keyboard.Init(id)
	return s.server.keyboard
}

func (s *Seat) Release(obj wayland.Seat) {
	// XXX let the server know we're done with the resource
}

func (s *Seat) OnDestroy(obj wlserver.Object) {}

type Server struct {
	keyboard *Keyboard
}

func (s *Server) bindSeat(sres wayland.Seat) wayland.SeatImplementation {
	sres.Capabilities(wayland.SeatCapabilityKeyboard)
	sres.Name("a nice seat")

	return &Seat{
		server:   s,
		resource: sres,
	}
}

func (s *Server) bindOutput(res wayland.Output) wayland.OutputImplementation {
	return nil
}

func main() {
	l, err := net.Listen("unix", "/run/user/1000/wayland-0")
	if err != nil {
		log.Fatal(err)
	}

	dsp := wlserver.NewDisplay(l.(*net.UnixListener))

	srv := &Server{
		keyboard: &Keyboard{4444, 8888},
	}

	wayland.AddSeatGlobal(dsp, 5, srv.bindSeat)
	wayland.AddOutputGlobal(dsp, 4, srv.bindOutput)

	go dsp.Run()
	for {
		select {
		case conn := <-dsp.NewConns():
			dsp.AddClient(conn)
		case msg := <-dsp.Messages():
			dsp.ProcessMessage(msg)
		case disco := <-dsp.Disconnects():
			log.Printf("client %d disconnected: %s", disco.Client.ID(), disco.Err)
			dsp.RemoveClient(disco.Client)
		}
	}
}
