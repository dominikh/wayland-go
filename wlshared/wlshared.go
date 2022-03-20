package wlshared

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"syscall"
	"unsafe"

	"honnef.co/go/wayland/wlproto"
)

var byteOrder binary.ByteOrder

func init() {
	var x uint32 = 0x01020304
	if *(*byte)(unsafe.Pointer(&x)) == 0x01 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
}

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
}

func EncodeRequest(buf []byte, source ObjectID, request int, args []interface{}) (data []byte, oob []byte) {
	buf = append(buf, 0, 0, 0, 0, 0, 0, 0, 0)
	var scratch [4]byte

	var fds []int
	for _, arg := range args {
		if v, ok := arg.(Object); ok {
			id := v.ID()
			byteOrder.PutUint32(scratch[:], uint32(id))
			buf = append(buf, scratch[:]...)
		} else {
			v := reflect.ValueOf(arg)
			switch v.Type().Kind() {
			case reflect.Int32:
				byteOrder.PutUint32(scratch[:], uint32(v.Int()))
				buf = append(buf, scratch[:]...)
			case reflect.Uint32:
				byteOrder.PutUint32(scratch[:], uint32(v.Uint()))
				buf = append(buf, scratch[:]...)
			case reflect.String:
				str := v.String()
				byteOrder.PutUint32(scratch[:], uint32(len(str)+1))
				buf = append(buf, scratch[:]...)
				buf = append(buf, str...)
				buf = append(buf, 0)
				m := len(str) + 1
				n := (m + 3) & ^3
				for i := n - m; i > 0; i-- {
					buf = append(buf, 0)
				}
			// XXX array
			case reflect.Uintptr:
				fds = append(fds, int(v.Uint()))
			default:
				panic(fmt.Sprintf("internal error: unhandled type %T", arg))
			}
		}
	}
	byteOrder.PutUint32(buf[0:4], uint32(source))
	byteOrder.PutUint16(buf[4:6], uint16(request))
	byteOrder.PutUint16(buf[6:8], uint16(len(buf)))

	if len(fds) > 0 {
		// OPT(dh): we send file descriptors so rarely that allocating
		// here isn't an issue.
		oob = syscall.UnixRights(fds...)
	}
	return buf, oob
}

func ParseArgument(arg wlproto.Arg, d []byte, off int) (newOff int, v interface{}) {
	var num uint32
	if arg.Type != wlproto.ArgTypeFd {
		num = byteOrder.Uint32(d[off:])
		off += 4
	}
	var out interface{}
	switch arg.Type {
	case wlproto.ArgTypeInt:
		out = int32(num)
	case wlproto.ArgTypeUint:
		if arg.Aux != nil {
			out = reflect.ValueOf(uint32(num)).Convert(arg.Aux).Interface()
		} else {
			out = uint32(num)
		}
	case wlproto.ArgTypeFixed:
		out = Fixed(num)
	case wlproto.ArgTypeString:
		out = string(d[off : off+int(num)-1])
		off += int(num)
		off = (off + 3) & ^3
	case wlproto.ArgTypeObject:
		out = ObjectID(num)
	case wlproto.ArgTypeArray:
		// XXX copy out the data, probably
		out = d[off : off+int(num)]
		off += int(num)
		off = (off + 3) &^ 3
	case wlproto.ArgTypeFd:
		out = nil
	case wlproto.ArgTypeNewID:
		// XXX a new_id might be just an integer, or also a string name and a version…
		out = ObjectID(num)
	default:
		panic("unreachable")
	}

	return off, out
}
