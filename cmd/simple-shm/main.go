package main

import (
	"log"
	"math/rand"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
	"honnef.co/go/wayland/wlclient"
	"honnef.co/go/wayland/wlcore"
)

type Display struct {
	display    *wlcore.Display
	registry   *wlcore.Registry
	compositor *wlcore.Compositor
	shm        *wlcore.Shm
	wmBase     *wlcore.XdgWmBase
	hasXRGB    bool
}

type Window struct {
	display          *Display
	width            int32
	height           int32
	surface          *wlcore.Surface
	xdgSurface       *wlcore.XdgSurface
	xdgToplevel      *wlcore.XdgToplevel
	waitForConfigure bool
	buffers          [2]Buffer
	callback         *wlcore.Callback
}

type Buffer struct {
	buffer  *wlcore.Buffer
	shmData []byte
	busy    bool
}

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

func createDisplay(c *wlclient.Conn) *Display {
	dsp := &Display{
		display: &wlcore.Display{},
	}
	c.NewProxy(1, dsp.display, nil)
	dsp.registry = dsp.display.GetRegistry()
	dsp.registry.AddListener(wlcore.RegistryEvents{
		Global: func(_ *wlcore.Registry, name uint32, iface string, version uint32) {
			switch iface {
			case "wl_compositor":
				dsp.compositor = &wlcore.Compositor{}
				c.NewProxy(0, dsp.compositor, nil)
				dsp.registry.Bind(name, dsp.compositor, 1)
			case "xdg_wm_base":
				dsp.wmBase = &wlcore.XdgWmBase{}
				c.NewProxy(0, dsp.wmBase, nil)
				dsp.registry.Bind(name, dsp.wmBase, 1)
			case "wl_shm":
				dsp.shm = &wlcore.Shm{}
				c.NewProxy(0, dsp.shm, nil)
				dsp.registry.Bind(name, dsp.shm, 1)
				dsp.shm.AddListener(wlcore.ShmEvents{
					Format: func(obj *wlcore.Shm, format uint32) {
						if format == wlcore.ShmFormatXrgb8888 {
							dsp.hasXRGB = true
						}
					},
				})
			}
		},
	})

	// make sure the server has processed all requests and sent out
	// all events, so that we have the full initial state of the
	// registry
	roundtrip(dsp.display)
	dsp.display.Queue().Dispatch()

	if dsp.shm == nil {
		log.Fatal("no SHM")
	}
	if dsp.wmBase == nil {
		log.Fatal("no XDG")
	}

	// this time make sure that we've processed all initial Shm events
	roundtrip(dsp.display)
	dsp.display.Queue().Dispatch()

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
	win.xdgSurface.AddListener(wlcore.XdgSurfaceEvents{
		Configure: func(_ *wlcore.XdgSurface, serial uint32) {
			win.xdgSurface.AckConfigure(serial)
			if win.waitForConfigure {
				redraw(win, nil, 0)
				win.waitForConfigure = false
			}
		},
	})

	win.xdgToplevel = win.xdgSurface.GetToplevel()
	win.xdgToplevel.SetTitle("simple-shm")
	win.surface.Commit()
	win.waitForConfigure = true
	return win
}

func redraw(win *Window, callback *wlcore.Callback, time uint32) {
	buf := windowNextBuffer(win)

	for i := range buf.shmData {
		buf.shmData[i] = byte(rand.Int())
	}

	win.surface.Attach(buf.buffer, 0, 0)
	win.surface.Damage(0, 0, win.width, win.height)
	if callback != nil {
		callback.Destroy()
	}
	win.callback = win.surface.Frame()
	win.callback.AddListener(wlcore.CallbackEvents{
		Done: func(_ *wlcore.Callback, data uint32) {
			redraw(win, win.callback, data)
		},
	})
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
		createShmBuffer(win.display, buf, win.width, win.height, wlcore.ShmFormatXrgb8888)
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
	buf.buffer.AddListener(wlcore.BufferEvents{
		Release: func(_ *wlcore.Buffer) {
			buf.busy = false
		},
	})
	pool.Destroy()
	unix.Close(fd)

	buf.shmData = data
}

func main() {
	uc, err := net.Dial("unix", "/run/user/1000/wayland-0")
	if err != nil {
		log.Fatal(err)
	}
	c := wlclient.NewConn(uc.(*net.UnixConn))

	dsp := createDisplay(c)
	createWindow(dsp, 250, 250)
	for {
		dsp.display.Queue().Dispatch()
	}
}
