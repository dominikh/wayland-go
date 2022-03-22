package wlserver

// This file contains an unexported copy of ./protocols/wayland, to break the circular dependency between wlserver and
// ./protocols/wayland.

import (
	"reflect"

	"honnef.co/go/wayland/wlproto"
	"honnef.co/go/wayland/wlshared"
)

type displayError uint32

const (
	displayErrorInvalidObject  displayError = 0
	displayErrorInvalidMethod  displayError = 1
	displayErrorNoMemory       displayError = 2
	displayErrorImplementation displayError = 3
)

var displayInterface = &wlproto.Interface{
	Name:    "wl_display",
	Version: 1,
	Type:    reflect.TypeOf(displayResource{}),
	Requests: []wlproto.Request{
		{
			Name:   "sync",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(displayImplementation.Sync),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeNewID, Aux: reflect.TypeOf(callbackResource{})},
			},
		},
		{
			Name:   "get_registry",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(displayImplementation.GetRegistry),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeNewID, Aux: reflect.TypeOf(registryResource{})},
			},
		},
	},
	Events: []wlproto.Event{
		{
			Name:  "error",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeObject},
				{Type: wlproto.ArgTypeUint},
				{Type: wlproto.ArgTypeString},
			},
		},
		{
			Name:  "delete_id",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
}

type displayResource struct{ Resource }

func (displayResource) Interface() *wlproto.Interface { return displayInterface }

type displayImplementation interface {
	Sync(obj displayResource, callback callbackResource) callbackImplementation
	GetRegistry(obj displayResource, registry registryResource) registryImplementation
}

func (obj displayResource) Error(objectId Object, code uint32, message string) {
	obj.Conn().SendEvent(obj, 0, objectId, code, message)
}

func (obj displayResource) DeleteID(id uint32) {
	obj.Conn().SendEvent(obj, 1, id)
}

var registryInterface = &wlproto.Interface{
	Name:    "wl_registry",
	Version: 1,
	Type:    reflect.TypeOf(registryResource{}),
	Requests: []wlproto.Request{
		{
			Name:   "bind",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(registryImplementation.Bind),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
				{Type: wlproto.ArgTypeString}, {Type: wlproto.ArgTypeUint}, {Type: wlproto.ArgTypeNewID},
			},
		},
	},
	Events: []wlproto.Event{
		{
			Name:  "global",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
				{Type: wlproto.ArgTypeString},
				{Type: wlproto.ArgTypeUint},
			},
		},
		{
			Name:  "global_remove",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
}

type registryResource struct{ Resource }

func (registryResource) Interface() *wlproto.Interface { return registryInterface }

type registryImplementation interface {
	Bind(obj registryResource, name uint32, idName string, idVersion uint32, id wlshared.ObjectID) ResourceImplementation
}

func (obj registryResource) Global(name uint32, interface_ string, version uint32) {
	obj.Conn().SendEvent(obj, 0, name, interface_, version)
}

func (obj registryResource) GlobalRemove(name uint32) {
	obj.Conn().SendEvent(obj, 1, name)
}

var callbackInterface = &wlproto.Interface{
	Name:     "wl_callback",
	Version:  1,
	Type:     reflect.TypeOf(callbackResource{}),
	Requests: []wlproto.Request{},
	Events: []wlproto.Event{
		{
			Name:  "done",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
}

type callbackResource struct{ Resource }

func (callbackResource) Interface() *wlproto.Interface { return callbackInterface }

type callbackImplementation interface{}

func (obj callbackResource) Done(callbackData uint32) {
	obj.Conn().SendEvent(obj, 0, callbackData)
}
