package main

import (
	"log"
	"net"

	"honnef.co/go/wayland/wlclient"
	"honnef.co/go/wayland/wlcore"
)

func roundtrip(dsp *wlcore.Display) {
	queue := wlclient.NewEventQueue()
	cb := dsp.WithQueue(queue).Sync()
	var done bool
	cb.AddListener(wlcore.CallbackEvents{
		Done: func(obj *wlcore.Callback, _ uint32) {
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

	// ID 1 is "special" and refers to the Display. no nice API for it
	// yet, hence explicit call to NewProxy.
	dsp := &wlcore.Display{}
	c.NewProxy(1, dsp, nil)

	registry := dsp.GetRegistry()
	registry.AddListener(wlcore.RegistryEvents{
		Global: func(obj *wlcore.Registry, name uint32, interface_ string, version uint32) {
			log.Println(obj, name, interface_, version)
		},
	})

	// wait until we've received the initial batch of registry events
	roundtrip(dsp)
	// call our callbacks
	dsp.Queue().Dispatch()
}
