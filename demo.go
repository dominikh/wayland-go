// +build ignore

package main

import (
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
	c := wayland.NewConn(uc)
	dsp := &demo.Display{}
	c.NewProxy(1, dsp, demo.DisplayInterface)
	dsp.GetRegistry()
	select {}
}
