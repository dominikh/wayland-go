package wlproto

import "reflect"

type Interface struct {
	Name     string
	Version  uint32
	Type     reflect.Type
	Requests []Request
	Events   []Event
}

type Arg struct {
	Type ArgType
	Aux  reflect.Type
}

type ArgType byte

const (
	ArgTypeInt ArgType = iota + 1
	ArgTypeUint
	ArgTypeFixed
	ArgTypeString
	ArgTypeObject
	ArgTypeNewID
	ArgTypeArray
	ArgTypeFd
)

type Request struct {
	Name   string
	Method reflect.Value
	Type   string
	Since  uint32
	Args   []Arg
}

type Event struct {
	Name  string
	Since uint32
	Args  []Arg
}
