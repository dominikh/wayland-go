// +build ignore

package main

import (
	"log"
	"net"

	"honnef.co/go/wayland"
	"honnef.co/go/wayland/demo"
)

func roundtrip(dsp *demo.Display) {
	queue := wayland.NewEventQueue()
	cb := dsp.WithQueue(queue).Sync()
	var done bool
	cb.AddListener(demo.CallbackEvents{
		Done: func(obj *demo.Callback, _ uint32) {
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
	c := wayland.NewConn(uc.(*net.UnixConn))
	dsp := &demo.Display{}
	c.NewProxy(1, dsp, nil)
	registry := dsp.GetRegistry()
	registry.AddListener(demo.RegistryEvents{
		Global: func(obj *demo.Registry, name uint32, interface_ string, version uint32) {
			log.Println(obj, name, interface_, version)
		},
	})

	roundtrip(dsp)
	log.Println("done")
	dsp.Queue().Dispatch()

	select {}
}
