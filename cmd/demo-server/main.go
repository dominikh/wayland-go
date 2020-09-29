/*
$ , weston-info
interface: 'wl_seat', version: 5, name: 1
	name: a nice seat
	capabilities: keyboard
	keyboard repeat rate: 1234
	keyboard repeat delay: 5678
*/

package main

import (
	"log"
	"net"

	"honnef.co/go/wayland/wlproto"
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
	keyboard *Keyboard

	resources map[wlserver.Object]struct{}
}

func (*Seat) Interface() (*wlproto.Interface, int) { return wayland.Seat{}.Interface(), 5 }

func (s *Seat) OnBind(res wlserver.Object) wlserver.ResourceImplementation {
	s.resources[res] = struct{}{}

	sres := res.(wayland.Seat)
	sres.Capabilities(wayland.SeatCapabilityKeyboard)
	sres.Name("a nice seat")

	return s
}

func (s *Seat) GetPointer(obj wayland.Seat, id wayland.Pointer) wayland.PointerImplementation {
	panic("not implemented")
}
func (s *Seat) GetTouch(obj wayland.Seat, id wayland.Touch) wayland.TouchImplementation {
	panic("not implemented")
}

func (s *Seat) GetKeyboard(obj wayland.Seat, id wayland.Keyboard) wayland.KeyboardImplementation {
	s.keyboard.Init(id)
	return s.keyboard
}

func (s *Seat) Release(obj wayland.Seat) {
	// XXX let the server know we're done with the resource
}

func (s *Seat) OnDestroy(obj wlserver.Object) {
	delete(s.resources, obj)
}

func main() {
	l, err := net.Listen("unix", "/run/user/1000/wayland-0")
	if err != nil {
		log.Fatal(err)
	}

	dsp := wlserver.NewDisplay(l.(*net.UnixListener))

	seat := &Seat{
		resources: map[wlserver.Object]struct{}{},
		keyboard:  &Keyboard{4444, 8888},
	}
	dsp.AddGlobal(seat)

	dsp.Run()
}
