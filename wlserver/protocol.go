package wlserver

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math"
	"net"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	"honnef.co/go/wayland/wlproto"
	"honnef.co/go/wayland/wlshared"
)

// XXX handle removing globals, consider the race with bind

var byteOrder binary.ByteOrder

func init() {
	var x uint32 = 0x01020304
	if *(*byte)(unsafe.Pointer(&x)) == 0x01 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
}

type Display struct {
	l        *net.UnixListener
	clientID uint64

	globalsMu sync.RWMutex
	clients   []*Client
	globalsID uint32
	globals   map[uint32]registeredGlobal
}

func NewDisplay(l *net.UnixListener) *Display {
	return &Display{
		l:       l,
		globals: map[uint32]registeredGlobal{},
	}
}

type registeredGlobal struct {
	g       Global
	iface   *wlproto.Interface
	version int
}

type Global interface {
	OnBind(res Object) ResourceImplementation
}

func (dsp *Display) AddGlobal(g Global, iface *wlproto.Interface, version int) {
	dsp.globalsMu.Lock()
	defer dsp.globalsMu.Unlock()

	dsp.globalsID++
	name := dsp.globalsID
	if dsp.globalsID == math.MaxUint32 {
		panic("global counter overflow")
	}
	dsp.globals[name] = registeredGlobal{g, iface, version}
	for _, c := range dsp.clients {
		for reg := range c.registries {
			c.sendEvent(reg, 0, uint32(name), iface.Name, uint32(version))
		}
	}
}

func (dsp *Display) Run() {
	for {
		conn, err := dsp.l.Accept()
		if err != nil {
			// XXX
			panic(err)
		}

		client := &Client{
			dsp:             dsp,
			id:              dsp.clientID,
			rw:              conn.(*net.UnixConn),
			objects:         map[wlshared.ObjectID]Object{},
			implementations: map[wlshared.ObjectID]ResourceImplementation{},
			registries:      map[wlshared.ObjectID]struct{}{},
		}
		dsp.clientID++
		dsp.clients = append(dsp.clients, client)
		go func() {
			err := client.readLoop()
			client.rw.Close()

			// XXX destroy all the globals and objects

			idx := -1
			for i, oc := range dsp.clients {
				if client == oc {
					idx = i
					break
				}
			}
			dsp.globalsMu.Lock()
			copy(dsp.clients[idx:], dsp.clients[idx+1:])
			dsp.clients = dsp.clients[:len(dsp.clients)-1]
			dsp.globalsMu.Unlock()

			if err != nil {
				client.sendMu.Lock()
				if client.err == nil {
					client.err = err
				}
				client.sendMu.Unlock()
			}
			if errors.Is(err, io.EOF) {
				log.Printf("client %d disconnected", client.id)
			} else {
				log.Printf("fatal client error for client %d: %s", client.id, err)
			}
		}()
	}
}

type Client struct {
	dsp     *Display
	id      uint64
	rw      *net.UnixConn
	objects map[wlshared.ObjectID]Object

	implementations map[wlshared.ObjectID]ResourceImplementation

	// we track instances of wl_registry separately of other
	// resources, because we can't import the generated wayland
	// package and cannot use the generic dispatch code
	registries map[wlshared.ObjectID]struct{}

	data []byte
	fds  []uintptr

	sendMu  sync.RWMutex
	sendBuf []byte
	err     error
}

func (c *Client) read() error {
	b := make([]byte, 1<<16)
	// XXX can there be more than one SCM per message?
	// XXX in general, be more robust in handling SCM
	oob := make([]byte, 24)
	n, oobn, _, _, err := c.rw.ReadMsgUnix(b, oob)
	if err != nil {
		return err
	}
	c.data = append(c.data, b[:n]...)
	if oobn == 24 {
		scm, err := syscall.ParseSocketControlMessage(oob[:oobn])
		if err != nil {
			return err
		}
		fds, err := syscall.ParseUnixRights(&scm[0])
		if err != nil {
			return err
		}
		c.fds = append(c.fds, uintptr(fds[0]))
	}
	return nil
}

func (c *Client) readAtLeast(n int) error {
	for len(c.data) < n {
		if err := c.read(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) readLoop() error {
	// We are the server, thus we are reading requests
	for {
		if err := c.readAtLeast(8); err != nil {
			return err
		}
		sender := wlshared.ObjectID(byteOrder.Uint32(c.data[0:4]))
		h := byteOrder.Uint32(c.data[4:8])
		size := (h & 0xFFFF0000) >> 16
		if size < 8 {
			// XXX invalid size
		}
		size -= 8
		opcode := h & 0x0000FFFF
		c.data = c.data[8:]
		if err := c.readAtLeast(int(size)); err != nil {
			return err
		}

		d := c.data[:size]
		c.data = c.data[size:]

		// special-case requests to the display or the registry, to
		// avoid a circular dependency between this package and the
		// generated wayland package.
		const (
			idDisplay             = 1
			reqDisplaySync        = 0
			reqDisplayGetRegistry = 1
			reqRegistryBind       = 0
			evCallbackDone        = 0
			evDisplayDeleteID     = 1
		)
		if sender == idDisplay {
			// request to the display
			switch opcode {
			case reqDisplaySync:
				id := byteOrder.Uint32(d)
				c.sendEvent(wlshared.ObjectID(id), evCallbackDone, uint32(0))
				c.sendEvent(idDisplay, evDisplayDeleteID, id)
			case reqDisplayGetRegistry:
				// XXX verify the ID isn't already in use
				id := byteOrder.Uint32(d)

				c.dsp.globalsMu.RLock()
				c.registries[wlshared.ObjectID(id)] = struct{}{}
				for name, g := range c.dsp.globals {
					c.sendEvent(wlshared.ObjectID(id), 0, uint32(name), g.iface.Name, uint32(g.version))
				}
				c.dsp.globalsMu.RUnlock()
			default:
				// XXX invalid opcode
			}
			continue
		} else if _, ok := c.registries[sender]; ok {
			// request to a registry
			if opcode == reqRegistryBind {
				// bind(name uint32, id new_id)

				name := byteOrder.Uint32(d)
				ifaceLen := byteOrder.Uint32(d[4:])
				ifaceName := string(d[8 : 8+ifaceLen])
				id := wlshared.ObjectID(byteOrder.Uint32(d[8+ifaceLen:]))
				_ = ifaceName

				// XXX guard against invalid name
				// XXX guard against in-use id
				// XXX verify that ifaceName matches the global's interface

				iface := c.dsp.globals[name].iface
				res := Resource{
					conn: c,
					id:   wlshared.ObjectID(id),
				}
				robj := reflect.New(iface.Type).Elem()
				robj.Field(0).Set(reflect.ValueOf(res))
				obj := robj.Interface().(Object)
				c.objects[id] = obj
				c.implementations[id] = c.dsp.globals[name].g.OnBind(obj)
			} else {
				// XXX invalid opcode
			}
			continue
		}

		obj, ok := c.objects[sender]
		if !ok {
			// TODO(dh): is it okay for objects to be unknown when we're the server, or should we kill the client?

			// unknown object
			continue
		}
		off := 0
		// XXX guard against invalid opcodes
		sig := obj.Interface().Requests[opcode].Args
		allArgs := make([]reflect.Value, len(sig)+2)
		impl := c.implementations[sender]
		allArgs[0] = reflect.ValueOf(impl)
		allArgs[1] = reflect.ValueOf(obj)
		args := allArgs[2:]
		for i, arg := range sig {
			newOff, argv := wlshared.ParseArgument(arg, d, off)
			off = newOff

			switch arg.Type {
			case wlproto.ArgTypeObject:
				// XXX guard against invalid object id
				args[i] = reflect.ValueOf(c.objects[argv.(wlshared.ObjectID)])
			case wlproto.ArgTypeFd:
				fd := c.fds[0]
				copy(c.fds, c.fds[1:])
				c.fds = c.fds[:len(c.fds)-1]
				args[i] = reflect.ValueOf(uintptr(fd))
			case wlproto.ArgTypeNewID:
				// XXX verify that the new ID isn't already in use
				// XXX verify that the new ID is in the client's ID space
				num := argv.(wlshared.ObjectID)
				res := Resource{
					conn: c,
					id:   wlshared.ObjectID(num),
				}
				rv := reflect.New(arg.Aux).Elem()
				rv.Field(0).Set(reflect.ValueOf(res))
				v := rv.Interface().(Object)
				c.objects[wlshared.ObjectID(num)] = v
				args[i] = rv
			default:
				args[i] = reflect.ValueOf(argv)
			}
		}

		meth := obj.Interface().Requests[opcode].Method
		results := meth.Call(allArgs)

		n := 0
		for _, arg := range sig {
			if arg.Type == wlproto.ArgTypeNewID {
				obj.GetResource().SetImplementation(results[n].Interface().(ResourceImplementation))
				n++
			}
		}

		c.sendMu.RLock()
		err := c.err
		c.sendMu.RUnlock()
		if err != nil {
			return err
		}
	}
}

type Resource struct {
	id   wlshared.ObjectID
	conn *Client
}

type ResourceImplementation interface {
	OnDestroy(Object)
}

func (p Resource) SetImplementation(impl ResourceImplementation) {
	p.conn.implementations[p.id] = impl
}

func (p Resource) GetResource() Resource { return p }
func (p Resource) Conn() *Client         { return p.conn }
func (p Resource) ID() wlshared.ObjectID { return p.id }

type Object interface {
	ID() wlshared.ObjectID
	Conn() *Client
	Interface() *wlproto.Interface
	GetResource() Resource
}

func (c *Client) SendEvent(source wlshared.Object, event int, args ...interface{}) {
	c.sendEvent(source.ID(), event, args...)
}

func (c *Client) sendEvent(source wlshared.ObjectID, event int, args ...interface{}) {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	if c.err != nil {
		return
	}

	buf := c.sendBuf[:0]
	var oob []byte
	buf, oob = wlshared.EncodeRequest(buf, source, event, args)
	_, _, c.err = c.rw.WriteMsgUnix(buf, oob, nil)

	c.sendBuf = buf[:0]
}
