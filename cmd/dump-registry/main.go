package main

import (
	"log"
	"net"

	"honnef.co/go/wayland/wlclient"
	"honnef.co/go/wayland/wlclient/protocols/wayland"
)

func roundtrip(dsp *wayland.Display) {
	queue := wlclient.NewEventQueue()
	cb := dsp.WithQueue(queue).Sync()
	var done bool
	cb.AddListener(wayland.CallbackEvents{
		Done: func(obj *wayland.Callback, _ uint32) {
			log.Println("callback fired")
			done = true
			cb.Destroy()
		},
	})
	for !done {
		queue.Dispatch()
	}
}

func main() {
	uc, err := net.Dial("unix", "/run/user/1000/wayland-0")
	if err != nil {
		log.Fatal(err)
	}
	c := wlclient.NewConn(uc.(*net.UnixConn))

	dsp := wayland.GetDisplay(c)

	registry := dsp.GetRegistry()
	registry.AddListener(wayland.RegistryEvents{
		Global: func(obj *wayland.Registry, name uint32, interface_ string, version uint32) {
			log.Println(obj, name, interface_, version)
		},
	})

	// wait until we've received the initial batch of registry events
	roundtrip(dsp)
	// call our callbacks
	dsp.Queue().Dispatch()
}
