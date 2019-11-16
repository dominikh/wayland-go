package wayland

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

type Fixed uint32

func (f Fixed) Float64() float64 {
	panic("XXX")
}

func FromFloat64(f float64) Fixed {
	panic("XXX")
}

type ObjectID uint32
type NewID uint32

type Object interface {
	ID() ObjectID
	Conn() *Conn
	Interface() *Interface
	GetProxy() *Proxy
	NextEvent() (Event, bool)
	NextEventPoll() Event
}

var byteOrder binary.ByteOrder

func init() {
	var x uint32 = 0x01020304
	if *(*byte)(unsafe.Pointer(&x)) == 0x01 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
}

type Interface struct {
	Name    string
	Version uint32
	Events  []Event
}

type Proxy struct {
	id   ObjectID
	conn *Conn

	mu     sync.Mutex
	events []Event
	signal chan struct{}
}

type Event interface{}

func (p *Proxy) GetProxy() *Proxy { return p }
func (p *Proxy) ID() ObjectID     { return p.id }
func (p *Proxy) Conn() *Conn      { return p.conn }

func (p *Proxy) pushEvent(ev Event) {
	p.mu.Lock()
	p.events = append(p.events, ev)
	select {
	case p.signal <- struct{}{}:
	default:
	}
	p.mu.Unlock()
}

func (p *Proxy) NextEvent() (Event, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.events) > 0 {
		ev := p.events[0]
		copy(p.events, p.events[1:])
		p.events = p.events[:len(p.events)-1]
		return ev, true
	}
	return nil, false
}

func (p *Proxy) NextEventPoll() Event {
	for {
		p.mu.Lock()
		if len(p.events) > 0 {
			ev := p.events[0]
			copy(p.events, p.events[1:])
			p.events = p.events[:len(p.events)-1]
			p.mu.Unlock()
			return ev
		}
		p.mu.Unlock()
		<-p.signal
	}
}

type Conn struct {
	rw      *net.UnixConn
	objects map[ObjectID]Object
	debug   bool

	maxID ObjectID

	data []byte
	fds  []uintptr
}

func NewConn(rw *net.UnixConn) *Conn {
	c := &Conn{
		rw:      rw,
		objects: map[ObjectID]Object{},
		debug:   true,
		maxID:   1,
	}
	go func() {
		if err := c.readLoop(); err != nil {
			log.Println("error in read loop:", err)
		}
	}()
	return c
}

func (c *Conn) NewProxy(id ObjectID, obj Object) {
	if id == 0 {
		c.maxID++
		id = c.maxID
	}
	p := obj.GetProxy()
	*p = Proxy{
		id:     id,
		conn:   c,
		signal: make(chan struct{}),
	}
	c.objects[id] = obj
}

func (c *Conn) SendRequest(source Object, request int, args ...interface{}) {
	// XXX we need to be aware of destructors and stop tracking
	// objects that were destroyed

	// OPT(dh): cache buf in Conn
	buf := make([]byte, 8, 8+len(args)*4)
	var scratch [4]byte
	var fds []int
	for _, arg := range args {
		switch arg := arg.(type) {
		case int32:
			byteOrder.PutUint32(scratch[:], uint32(arg))
			buf = append(buf, scratch[:]...)
		case uint32:
			byteOrder.PutUint32(scratch[:], arg)
			buf = append(buf, scratch[:]...)
		case Fixed:
			byteOrder.PutUint32(scratch[:], uint32(arg))
			buf = append(buf, scratch[:]...)
		case string:
			byteOrder.PutUint32(scratch[:], uint32(len(arg)+1))
			buf = append(buf, scratch[:]...)
			buf = append(buf, arg...)
			buf = append(buf, 0)
			m := len(arg) + 1
			n := (m + 3) & ^3
			for i := n - m; i > 0; i-- {
				buf = append(buf, 0)
			}
			// XXX array
			// XXX fd
		case uintptr:
			fds = append(fds, int(arg))
		case Object:
			id := arg.ID()
			byteOrder.PutUint32(scratch[:], uint32(id))
			buf = append(buf, scratch[:]...)
		default:
			panic(fmt.Sprintf("internal error: unhandled type %T", arg))
		}
	}
	byteOrder.PutUint32(buf[0:4], uint32(source.ID()))
	byteOrder.PutUint16(buf[4:6], uint16(request))
	byteOrder.PutUint16(buf[6:8], uint16(len(buf)))

	var oob []byte
	if len(fds) > 0 {
		oob = syscall.UnixRights(fds...)
	}
	c.rw.WriteMsgUnix(buf, oob, nil)
}

func (c *Conn) read() {
	b := make([]byte, 1<<16)
	// XXX can there be more than one SCM per message?
	// XXX in general, be more robust in handling SCM
	oob := make([]byte, 24)
	n, oobn, _, _, err := c.rw.ReadMsgUnix(b, oob)
	if err != nil {
		panic(err)
	}
	c.data = append(c.data, b[:n]...)
	if oobn == 24 {
		scm, err := syscall.ParseSocketControlMessage(oob[:oobn])
		if err != nil {
			panic(err)
		}
		fds, err := syscall.ParseUnixRights(&scm[0])
		if err != nil {
			panic(err)
		}
		c.fds = append(c.fds, uintptr(fds[0]))
	}
}

func (c *Conn) readAtLeast(n int) {
	for len(c.data) < n {
		c.read()
	}
}

func (c *Conn) readLoop() error {
	var (
		tInt32   = reflect.TypeOf(int32(0))
		tUint32  = reflect.TypeOf(uint32(0))
		tFixed   = reflect.TypeOf(Fixed(0))
		tString  = reflect.TypeOf("")
		tObject  = reflect.TypeOf((*Object)(nil)).Elem()
		tUintptr = reflect.TypeOf(uintptr(0))
		tArray   = reflect.TypeOf([]byte{})
	)

	for {
		c.readAtLeast(8)
		sender := ObjectID(byteOrder.Uint32(c.data[0:4]))
		h := byteOrder.Uint32(c.data[4:8])
		size := (h & 0xFFFF0000) >> 16
		if size < 8 {
			// XXX invalid size
		}
		size -= 8
		opcode := h & 0x0000FFFF
		c.data = c.data[8:]
		c.readAtLeast(int(size))

		d := c.data[:size]
		c.data = c.data[size:]

		obj, ok := c.objects[sender]
		if !ok {
			// XXX unknown object
		}
		off := 0
		evT := reflect.TypeOf(obj.Interface().Events[opcode])
		ev := reflect.New(evT.Elem())
		elem := ev.Elem()
		for i := 0; i < elem.NumField(); i++ {
			f := elem.Field(i)
			fT := evT.Elem().Field(i)

			var num uint32
			if fT.Type != tUintptr {
				num = byteOrder.Uint32(d[off:])
				off += 4
			}
			switch fT.Type {
			case tInt32:
				f.SetInt(int64(num))
			case tUint32:
				f.SetUint(uint64(num))
			case tFixed:
				f.SetUint(uint64(num))
			case tString:
				s := string(d[off : off+int(num)-1])
				f.SetString(s)
				off += int(num)
				off = (off + 3) & ^3
			case tObject:
				if evT.Elem().Field(i).Tag.Get("wl") == "new_id" {
					v := reflect.New(f.Type().Elem()).Interface().(Object)
					c.NewProxy(ObjectID(num), v)
					c.objects[ObjectID(num)] = v
				} else {
					f.Set(reflect.ValueOf(c.objects[ObjectID(num)]))
				}
			case tUintptr:
				fd := c.fds[0]
				copy(c.fds, c.fds[1:])
				c.fds = c.fds[:len(c.fds)-1]
				f.SetUint(uint64(fd))
			case tArray:
				// XXX copy out the data, probably
				b := d[off : off+int(num)]
				f.Set(reflect.ValueOf(b))
				off += int(num)
				off = (off + 3) &^ 3
			default:
				// XXX support arrays and file descriptors
				panic(fmt.Sprintf("internal error: unexpected type %v", fT.Type))
			}
		}
		obj.GetProxy().pushEvent(ev.Interface().(Event))
	}
}
