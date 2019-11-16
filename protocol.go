package wayland

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
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
	Name     string
	Version  int
	Requests []MessageRequest
	Events   []MessageEvent
}

type MessageRequest struct {
	Name  string
	Types []string
}

type MessageEvent struct {
	Name  string
	Types []interface{}
}

type Proxy struct {
	id    ObjectID
	conn  *Conn
}

func (p *Proxy) GetProxy() *Proxy      { return p }
func (p *Proxy) ID() ObjectID          { return p.id }
func (p *Proxy) Conn() *Conn           { return p.conn }

type Conn struct {
	rw      io.ReadWriter
	objects map[ObjectID]Object
	debug   bool

	maxID ObjectID
}

func NewConn(rw io.ReadWriter) *Conn {
	c := &Conn{
		rw:      rw,
		objects: map[ObjectID]Object{},
		debug:   true,
		maxID:   1,
	}
	go func() {
		if err := c.read(); err != nil {
			log.Println("error in read loop:", err)
		}
	}()

	// dsp := &Display{
	// 	Proxy:Proxy{
	// 		id: 1,
	// 		conn: c,
	// 		iface: DisplayInterface,
	// 	},
	// }
	// dsp.GetRegistry()
	return c
}

func (c *Conn) Test() {
	b := make([]byte, 12)
	b[3] = 1
	byteOrder.PutUint32(b[4:], 8<<16|1)
	b[11] = 2
	c.rw.Write(b)
}

func (c *Conn) NewProxy(id ObjectID, obj Object)  {
	if id == 0 {
		c.maxID++
		id = c.maxID
	}
	p := obj.GetProxy()
	*p = Proxy{
		id:    id,
		conn:  c,
	}
	c.objects[id] = obj
}

func (c *Conn) SendRequest(source Object, request int, args ...interface{}) {
	// OPT(dh): cache buf in Conn
	buf := make([]byte, 8, 8+len(args)*4)
	var scratch [4]byte
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
			panic("string not implemented")
			// XXX array
			// XXX fd
		case Object:
			id := arg.ID()
			byteOrder.PutUint32(scratch[:], uint32(id))
			buf = append(buf, scratch[:]...)
		}
	}
	byteOrder.PutUint32(buf[0:4], uint32(source.ID()))
	byteOrder.PutUint16(buf[4:6], uint16(request))
	byteOrder.PutUint16(buf[6:8], uint16(len(buf)))
	log.Println(len(buf), request)
	c.rw.Write(buf)
}

func (c *Conn) read() error {
	// XXX it's a stream protocol, not a message protocol, so don't expect to read full messages

	b := make([]byte, 1<<16)
	for {
		n, err := c.rw.Read(b)
		log.Printf("read %d bytes, err = %v", n, err)
		if err != nil {
			return err
		}
		d := b[:n]

		for len(d) > 0 {
			if len(d) < 8 {
				// XXX invalid header
			}

			sender := ObjectID(byteOrder.Uint32(d[0:4]))
			h := byteOrder.Uint32(d[4:8])
			size := (h & 0xFFFF0000) >> 16
			opcode := h & 0x0000FFFF

			if c.debug {
				log.Printf("event: sender = %d, opcode = %d, size = %d", sender, opcode, size)
			}

			if size < 8 {
				// XXX invalid size
			}
			obj, ok := c.objects[sender]
			if !ok {
				// XXX unknown object
			}
			m := obj.Interface().Events[opcode]
			off := 8 // skip the header
			for _, typ := range m.Types {
				switch typ.(type) {
				case int32:
					fmt.Println(int32(byteOrder.Uint32(d[off:])))
					off += 4
				case uint32:
					fmt.Println(byteOrder.Uint32(d[off:]))
					off += 4
				case Fixed:
					fmt.Println(Fixed(byteOrder.Uint32(d[off:])))
					off += 4
				case string:
					n := int(byteOrder.Uint32(d[off:]))
					s := string(d[off+4 : off+4+n-1])
					fmt.Println(s)
					off += 4 + n
					off = (off + 3) & ^3
				case *Interface:
					// new_id
					fmt.Println(byteOrder.Uint32(d[off:]))
					off += 4
				default:
					// object
					obj := c.objects[ObjectID(byteOrder.Uint32(d[off:]))]
					fmt.Println(obj)
					off += 4
				}
			}

			d = d[size:]
		}
	}
}
