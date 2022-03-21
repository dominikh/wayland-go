package wlserver

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"honnef.co/go/wayland/wlproto"
	"honnef.co/go/wayland/wlshared"
)

// TODO implement removal of globals, take the race documented in
// https://gitlab.freedesktop.org/wayland/wayland/-/issues/10 into consideration

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

	clients   map[*Client]struct{}
	globalsID uint32
	globals   map[uint32]global

	newConns    chan net.Conn
	messages    chan Message
	disconnects chan Disconnect
}

type Disconnect struct {
	Client *Client
	Err    error
}

func NewDisplay(l *net.UnixListener) *Display {
	return &Display{
		l:           l,
		clients:     make(map[*Client]struct{}),
		globals:     map[uint32]global{},
		newConns:    make(chan net.Conn),
		messages:    make(chan Message),
		disconnects: make(chan Disconnect),
	}
}

func (dsp *Display) NewConns() <-chan net.Conn {
	return dsp.newConns
}

func (dsp *Display) Messages() <-chan Message {
	return dsp.messages
}

func (dsp *Display) Disconnects() <-chan Disconnect {
	return dsp.disconnects
}

type global struct {
	iface   *wlproto.Interface
	version int
	bind    func(Object) ResourceImplementation
}

func (dsp *Display) AddGlobal(iface *wlproto.Interface, version int, bind func(Object) ResourceImplementation) uint32 {
	dsp.globalsID++
	name := dsp.globalsID
	if dsp.globalsID == math.MaxUint32 {
		// XXX reclaim names used by deleted globals
		panic("global counter overflow")
	}
	dsp.globals[name] = global{iface, version, bind}

	for c := range dsp.clients {
		for reg := range c.registries {
			c.sendEvent(reg, 0, uint32(name), iface.Name, uint32(version))
		}
	}

	return name
}

func (dsp *Display) RemoveGlobal(name uint32) {
	delete(dsp.globals, name)

	for c := range dsp.clients {
		for reg := range c.registries {
			c.sendEvent(reg, 1, uint32(name))
		}
	}
}

type Message struct {
	Client *Client
	Sender wlshared.ObjectID
	Opcode uint32
	Data   []byte
}

func (dsp *Display) AddClient(conn net.Conn) *Client {
	client := &Client{
		dsp:             dsp,
		id:              dsp.clientID,
		rw:              conn.(*net.UnixConn),
		objects:         map[wlshared.ObjectID]Object{},
		implementations: map[wlshared.ObjectID]ResourceImplementation{},
		registries:      map[wlshared.ObjectID]struct{}{},
	}
	dsp.clientID++
	dsp.clients[client] = struct{}{}
	go func() {
		err := client.readLoop(dsp.messages)
		client.rw.Close()
		if werr, ok := client.err.Load().(error); ok {
			// favour the write error over the read error
			err = werr
		}
		// XXX destroy all resources before officially disconnecting the client
		dsp.disconnects <- Disconnect{client, err}
	}()
	return client
}

func (dsp *Display) RemoveClient(client *Client) {
	// XXX properly disconnect the client if it isn't already disconnected
	delete(dsp.clients, client)
}

type buf []byte

func (b *buf) uint32() uint32 {
	n := byteOrder.Uint32(*b)
	*b = (*b)[4:]
	return n
}

func (b *buf) string() string {
	n := byteOrder.Uint32(*b)
	// -1 to skip terminating null byte
	data := (*b)[4 : 4+n-1]
	// strings are padded to 32-bit boundary
	n = (n + 3) &^ 0x03
	*b = (*b)[4+n:]
	return string(data)
}

func (dsp *Display) ProcessMessage(msg Message) {
	c := msg.Client
	d := buf(msg.Data)
	opcode := msg.Opcode
	sender := msg.Sender
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
			id := d.uint32()
			c.sendEvent(wlshared.ObjectID(id), evCallbackDone, uint32(0))
			c.sendEvent(idDisplay, evDisplayDeleteID, id)
		case reqDisplayGetRegistry:
			// XXX verify the ID isn't already in use
			id := d.uint32()

			c.registries[wlshared.ObjectID(id)] = struct{}{}
			for name, g := range c.dsp.globals {
				c.sendEvent(wlshared.ObjectID(id), 0, uint32(name), g.iface.Name, uint32(g.version))
			}
		default:
			// XXX invalid opcode
		}
	} else if _, ok := c.registries[sender]; ok {
		// request to a registry
		if opcode == reqRegistryBind {
			name := d.uint32()
			ifaceName := d.string()
			version := d.uint32()
			id := wlshared.ObjectID(d.uint32())
			_ = ifaceName

			g, ok := c.dsp.globals[name]
			if !ok {
				if name > c.dsp.globalsID {
					// XXX the client tried to bind to a global that has never existed and we should kill the client
				} else {
					// XXX the global has been removed and the client's bind raced with the removal. we should set up a
					// tombstone that later gets destroyed by the client.
				}
			}

			// XXX guard against in-use id
			// XXX verify that ifaceName matches the global's interface

			iface := g.iface
			res := Resource{
				conn:    c,
				id:      wlshared.ObjectID(id),
				version: version,
			}
			robj := reflect.New(iface.Type).Elem()
			robj.Field(0).Set(reflect.ValueOf(res))
			obj := robj.Interface().(Object)
			c.objects[id] = obj
			c.implementations[id] = c.dsp.globals[name].bind(obj)
			// TODO(dh): we should verify that bind returned the correct implementation, e.g. a global with
			// wayland.SeatInterface returns an implementation that implements wayland.SeatImplementation
		} else {
			// XXX invalid opcode
		}
	} else {
		obj, ok := c.objects[sender]
		if !ok {
			// TODO(dh): is it okay for objects to be unknown when we're the server, or should we kill the client?

			// unknown object
			return
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
				// XXX handle new_id with no specified interface, and thus no Aux
				if arg.Aux == nil {
					panic("not implemented")
				}
				// XXX set the resource's version. it will be tied to the version of the resource to which the request is being sent
				rv := reflect.New(arg.Aux).Elem()
				rv.Field(0).Set(reflect.ValueOf(res))
				v := rv.Interface().(Object)
				c.objects[wlshared.ObjectID(num)] = v
				args[i] = rv
			default:
				args[i] = reflect.ValueOf(argv)
			}
		}

		// XXX guard against opcodes that don't exist in our version of the protocol
		meth := obj.Interface().Requests[opcode].Method
		results := meth.Call(allArgs)

		n := 0
		for _, arg := range sig {
			if arg.Type == wlproto.ArgTypeNewID {
				obj.GetResource().SetImplementation(results[n].Interface().(ResourceImplementation))
				n++
			}
		}

		if obj.Interface().Requests[opcode].Type == "destructor" {
			delete(c.objects, obj.ID())
			delete(c.implementations, obj.ID())
			c.sendEvent(1, evDisplayDeleteID, obj.ID())
		}
	}
}

func objectString(obj Object) string {
	return fmt.Sprintf("%s@%d", obj.Interface().Name, obj.ID())
}

func (dsp *Display) Run() {
	for {
		conn, err := dsp.l.Accept()
		if err != nil {
			// XXX
			panic(err)
			return
		}
		dsp.newConns <- conn
	}
}

type Client struct {
	dsp *Display
	id  uint64
	rw  *net.UnixConn

	// TODO merge objects and implementations maps
	objects         map[wlshared.ObjectID]Object
	implementations map[wlshared.ObjectID]ResourceImplementation

	err atomic.Value

	// we track instances of wl_registry separately of other
	// resources, because we can't import the generated wayland
	// package and cannot use the generic dispatch code
	registries map[wlshared.ObjectID]struct{}

	fds []uintptr

	sendMu  sync.RWMutex
	sendBuf []byte
}

func (c *Client) ID() uint64 { return c.id }

func (c *Client) read(b []byte) (int, error) {
	// XXX can there be more than one SCM per message?
	// XXX in general, be more robust in handling SCM
	oob := make([]byte, 24)
	n, oobn, _, _, err := c.rw.ReadMsgUnix(b, oob)
	if err != nil {
		return n, err
	}
	if oobn == 24 {
		scm, err := syscall.ParseSocketControlMessage(oob[:oobn])
		if err != nil {
			return n, err
		}
		fds, err := syscall.ParseUnixRights(&scm[0])
		if err != nil {
			return n, err
		}
		c.fds = append(c.fds, uintptr(fds[0]))
	}
	return n, nil
}

func (c *Client) readFull(buf []byte) (n int, err error) {
	for n < len(buf) && err == nil {
		var nn int
		nn, err = c.read(buf[n:])
		n += nn
	}
	if n == len(buf) {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

func (c *Client) readLoop(msgs chan<- Message) error {
	// We are the server, thus we are reading requests
	hdr := make([]byte, 8)
	for {
		if _, err := c.readFull(hdr); err != nil {
			return err
		}
		sender := wlshared.ObjectID(byteOrder.Uint32(hdr[0:4]))
		h := byteOrder.Uint32(hdr[4:8])
		size := (h & 0xFFFF0000) >> 16
		if size < 8 {
			// XXX invalid size
		}
		size -= 8
		opcode := h & 0x0000FFFF

		buf := make([]byte, int(size))
		if _, err := c.readFull(buf); err != nil {
			return err
		}

		msgs <- Message{
			Client: c,
			Sender: sender,
			Opcode: opcode,
			Data:   buf,
		}
	}
}

type Resource struct {
	id      wlshared.ObjectID
	conn    *Client
	version uint32
}

type ResourceImplementation interface{}

func (p Resource) SetImplementation(impl ResourceImplementation) {
	p.conn.implementations[p.id] = impl
}

func (p Resource) GetResource() Resource { return p }
func (p Resource) Conn() *Client         { return p.conn }
func (p Resource) ID() wlshared.ObjectID { return p.id }
func (p Resource) Version() uint32       { return p.version }

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

	buf := c.sendBuf[:0]
	var oob []byte
	buf, oob = wlshared.EncodeRequest(buf, source, event, args)
	_, _, err := c.rw.WriteMsgUnix(buf, oob, nil)
	if err != nil {
		// Set c.err if it hasn't been set yet
		c.err.CompareAndSwap(nil, err)
		c.rw.Close()
	}

	c.sendBuf = buf[:0]
}
