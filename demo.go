// +build ignore

package main

import (
	"log"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
	"honnef.co/go/wayland"
	"honnef.co/go/wayland/demo"
)

type Display struct {
	display    *demo.Display
	registry   *demo.Registry
	compositor *demo.Compositor
	shm        *demo.Shm
	wmBase     *demo.XdgWmBase
	hasXRGB    bool
}

type Window struct {
	display          *Display
	width            int32
	height           int32
	surface          *demo.Surface
	xdgSurface       *demo.XdgSurface
	xdgToplevel      *demo.XdgToplevel
	waitForConfigure bool
	buffers          [2]Buffer
	callback         *demo.Callback
}

type Buffer struct {
	buffer  *demo.Buffer
	shmData []byte
	busy    bool
}

func createDisplay(c *wayland.Conn) *Display {
	dsp := &Display{
		display: &demo.Display{},
	}
	c.NewProxy(1, dsp.display)
	dsp.registry = dsp.display.GetRegistry()

	// make sure the server has processed all requests and sent out
	// all events, so that we have the full initial state of the
	// registry
	dsp.display.Sync().NextEventPoll()

	for {
		ev, ok := dsp.registry.NextEvent()
		if !ok {
			break
		}
		g, ok := ev.(*demo.RegistryEventGlobal)
		if !ok {
			continue
		}
		switch g.Interface {
		case "wl_compositor":
			dsp.compositor = &demo.Compositor{}
			c.NewProxy(0, dsp.compositor)
			dsp.registry.Bind(g.Name, dsp.compositor, 1)
		case "xdg_wm_base":
			dsp.wmBase = &demo.XdgWmBase{}
			c.NewProxy(0, dsp.wmBase)
			dsp.registry.Bind(g.Name, dsp.wmBase, 1)
		case "wl_shm":
			dsp.shm = &demo.Shm{}
			c.NewProxy(0, dsp.shm)
			dsp.registry.Bind(g.Name, dsp.shm, 1)
		}
	}

	if dsp.shm == nil {
		log.Fatal("no SHM")
	}
	if dsp.wmBase == nil {
		log.Fatal("no XDG")
	}

	// this time we wait for all initial events from Shm
	dsp.display.Sync().NextEventPoll()

	for {
		ev, ok := dsp.shm.NextEvent()
		if !ok {
			break
		}
		f := ev.(*demo.ShmEventFormat)
		if f.Format == demo.ShmFormatXrgb8888 {
			dsp.hasXRGB = true
		}
	}

	if !dsp.hasXRGB {
		log.Fatal("no XRGB8888")
	}

	return dsp
}

func createWindow(dsp *Display, width, height int32) *Window {
	win := &Window{
		display: dsp,
		width:   width,
		height:  height,
		surface: dsp.compositor.CreateSurface(),
	}

	win.xdgSurface = dsp.wmBase.GetXdgSurface(win.surface)
	go func() {
		for {
			ev := win.xdgSurface.NextEventPoll()
			if ev, ok := ev.(*demo.XdgSurfaceEventConfigure); ok {
				win.xdgSurface.AckConfigure(ev.Serial)
				if win.waitForConfigure {
					redraw(win, nil, 0)
					// XXX race condition
					win.waitForConfigure = false
				}
			}
		}
	}()
	win.xdgToplevel = win.xdgSurface.GetToplevel()
	win.xdgToplevel.SetTitle("simple-shm")
	win.surface.Commit()
	win.waitForConfigure = true
	return win
}

func redraw(win *Window, callback *demo.Callback, time uint32) {
	buf := windowNextBuffer(win)

	// for i := range buf.shmData {
	// 	buf.shmData[i] = byte(rand.Int())
	// }

	win.surface.Attach(buf.buffer, 0, 0)
	win.surface.Damage(0, 0, win.width, win.height)
	if callback != nil {
		// XXX destroy callback
		// callback.Destroy()
	}
	win.callback = win.surface.Frame()
	go func() {
		for {
			ev := win.callback.NextEventPoll().(*demo.CallbackEventDone)
			redraw(win, win.callback, ev.CallbackData)
		}
	}()
	buf.busy = true
	win.surface.Commit()
}

func paintPixels() {

}

func windowNextBuffer(win *Window) *Buffer {
	var buf *Buffer
	if !win.buffers[0].busy {
		buf = &win.buffers[0]
	} else {
		buf = &win.buffers[1]
	}

	if buf.buffer == nil {
		createShmBuffer(win.display, buf, win.width, win.height, demo.ShmFormatXrgb8888)

		// XXX memset
	}

	return buf
}

func createShmBuffer(dsp *Display, buf *Buffer, width, height int32, format uint32) {
	stride := width * 4
	size := stride * height
	fd, err := unix.MemfdCreate("", 0)
	if err != nil {
		log.Fatal(err)
	}
	unix.Ftruncate(fd, int64(size))
	data, err := syscall.Mmap(fd, 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}

	pool := dsp.shm.CreatePool(uintptr(fd), size)
	buf.buffer = pool.CreateBuffer(0, width, height, stride, format)
	go func() {
		for {
			ev := buf.buffer.NextEventPoll()
			if _, ok := ev.(*demo.BufferEventRelease); ok {
				// XXX race
				buf.busy = false
			}
		}
	}()
	pool.Destroy()
	unix.Close(fd)

	buf.shmData = data
}

func main() {
	uc, err := net.Dial("unix", "/run/user/1000/wayland-0")
	if err != nil {
		log.Fatal(err)
	}
	c := wayland.NewConn(uc.(*net.UnixConn))

	dsp := createDisplay(c)
	window := createWindow(dsp, 250, 250)
	_ = window
	select {}
}
