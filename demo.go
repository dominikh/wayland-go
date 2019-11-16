// +build ignore

package main

import (
	"fmt"
	"log"
	"net"

	"honnef.co/go/wayland"
	"honnef.co/go/wayland/demo"
)

func main() {
	uc, err := net.Dial("unix", "/run/user/1000/wayland-0")
	if err != nil {
		log.Fatal(err)
	}
	c := wayland.NewConn(uc.(*net.UnixConn))
	dsp := &demo.Display{}
	c.NewProxy(1, dsp)
	reg := dsp.GetRegistry()
	go func() {
		for {
			ev := reg.NextEvent()
			_ = ev
		}
	}()

	seat := &demo.Seat{}
	c.NewProxy(0, seat)
	reg.Bind(12, seat, 5)
	kb := seat.GetKeyboard()
	for {
		fmt.Printf("%#v\n",kb.NextEvent())
	}

	select {}
}
