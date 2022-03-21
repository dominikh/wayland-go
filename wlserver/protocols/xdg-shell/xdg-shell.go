// Code generated by wayland-scanner; DO NOT EDIT.

// Package xdgShell contains generated definitions of the xdg_shell Wayland protocol.
package xdgShell

import (
	"honnef.co/go/wayland/wlclient/protocols/wayland"
	"honnef.co/go/wayland/wlproto"
	"honnef.co/go/wayland/wlserver"
	"honnef.co/go/wayland/wlshared"
	"reflect"
)

var _ wlshared.Fixed

var interfaceNames = map[string]string{
	"xdg_wm_base":    "WmBase",
	"xdg_positioner": "Positioner",
	"xdg_surface":    "Surface",
	"xdg_toplevel":   "Toplevel",
	"xdg_popup":      "Popup",
}

var Interfaces = map[string]*wlproto.Interface{
	"xdg_wm_base":    WmBaseInterface,
	"xdg_positioner": PositionerInterface,
	"xdg_surface":    SurfaceInterface,
	"xdg_toplevel":   ToplevelInterface,
	"xdg_popup":      PopupInterface,
}

var Requests = map[string]*wlproto.Request{
	"xdg_wm_base_destroy":                      &WmBaseInterface.Requests[0],
	"xdg_wm_base_create_positioner":            &WmBaseInterface.Requests[1],
	"xdg_wm_base_get_xdg_surface":              &WmBaseInterface.Requests[2],
	"xdg_wm_base_pong":                         &WmBaseInterface.Requests[3],
	"xdg_positioner_destroy":                   &PositionerInterface.Requests[0],
	"xdg_positioner_set_size":                  &PositionerInterface.Requests[1],
	"xdg_positioner_set_anchor_rect":           &PositionerInterface.Requests[2],
	"xdg_positioner_set_anchor":                &PositionerInterface.Requests[3],
	"xdg_positioner_set_gravity":               &PositionerInterface.Requests[4],
	"xdg_positioner_set_constraint_adjustment": &PositionerInterface.Requests[5],
	"xdg_positioner_set_offset":                &PositionerInterface.Requests[6],
	"xdg_positioner_set_reactive":              &PositionerInterface.Requests[7],
	"xdg_positioner_set_parent_size":           &PositionerInterface.Requests[8],
	"xdg_positioner_set_parent_configure":      &PositionerInterface.Requests[9],
	"xdg_surface_destroy":                      &SurfaceInterface.Requests[0],
	"xdg_surface_get_toplevel":                 &SurfaceInterface.Requests[1],
	"xdg_surface_get_popup":                    &SurfaceInterface.Requests[2],
	"xdg_surface_set_window_geometry":          &SurfaceInterface.Requests[3],
	"xdg_surface_ack_configure":                &SurfaceInterface.Requests[4],
	"xdg_toplevel_destroy":                     &ToplevelInterface.Requests[0],
	"xdg_toplevel_set_parent":                  &ToplevelInterface.Requests[1],
	"xdg_toplevel_set_title":                   &ToplevelInterface.Requests[2],
	"xdg_toplevel_set_app_id":                  &ToplevelInterface.Requests[3],
	"xdg_toplevel_show_window_menu":            &ToplevelInterface.Requests[4],
	"xdg_toplevel_move":                        &ToplevelInterface.Requests[5],
	"xdg_toplevel_resize":                      &ToplevelInterface.Requests[6],
	"xdg_toplevel_set_max_size":                &ToplevelInterface.Requests[7],
	"xdg_toplevel_set_min_size":                &ToplevelInterface.Requests[8],
	"xdg_toplevel_set_maximized":               &ToplevelInterface.Requests[9],
	"xdg_toplevel_unset_maximized":             &ToplevelInterface.Requests[10],
	"xdg_toplevel_set_fullscreen":              &ToplevelInterface.Requests[11],
	"xdg_toplevel_unset_fullscreen":            &ToplevelInterface.Requests[12],
	"xdg_toplevel_set_minimized":               &ToplevelInterface.Requests[13],
	"xdg_popup_destroy":                        &PopupInterface.Requests[0],
	"xdg_popup_grab":                           &PopupInterface.Requests[1],
	"xdg_popup_reposition":                     &PopupInterface.Requests[2],
}

var Events = map[string]*wlproto.Event{
	"xdg_wm_base_ping":              &WmBaseInterface.Events[0],
	"xdg_surface_configure":         &SurfaceInterface.Events[0],
	"xdg_toplevel_configure":        &ToplevelInterface.Events[0],
	"xdg_toplevel_close":            &ToplevelInterface.Events[1],
	"xdg_toplevel_configure_bounds": &ToplevelInterface.Events[2],
	"xdg_popup_configure":           &PopupInterface.Events[0],
	"xdg_popup_popup_done":          &PopupInterface.Events[1],
	"xdg_popup_repositioned":        &PopupInterface.Events[2],
}

type WmBaseError uint32

const (
	// given wl_surface has another role
	WmBaseErrorRole WmBaseError = 0
	// xdg_wm_base was destroyed before children
	WmBaseErrorDefunctSurfaces WmBaseError = 1
	// the client tried to map or destroy a non-topmost popup
	WmBaseErrorNotTheTopmostPopup WmBaseError = 2
	// the client specified an invalid popup parent surface
	WmBaseErrorInvalidPopupParent WmBaseError = 3
	// the client provided an invalid surface state
	WmBaseErrorInvalidSurfaceState WmBaseError = 4
	// the client provided an invalid positioner
	WmBaseErrorInvalidPositioner WmBaseError = 5
)

var WmBaseInterface = &wlproto.Interface{
	Name:    "xdg_wm_base",
	Version: 4,
	Type:    reflect.TypeOf(WmBase{}),
	Requests: []wlproto.Request{
		{
			Name:   "destroy",
			Type:   "destructor",
			Since:  1,
			Method: reflect.ValueOf(WmBaseImplementation.Destroy),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "create_positioner",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(WmBaseImplementation.CreatePositioner),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeNewID, Aux: reflect.TypeOf(Positioner{})},
			},
		},
		{
			Name:   "get_xdg_surface",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(WmBaseImplementation.GetXdgSurface),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeNewID, Aux: reflect.TypeOf(Surface{})},
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(wayland.Surface{})},
			},
		},
		{
			Name:   "pong",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(WmBaseImplementation.Pong),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
	Events: []wlproto.Event{
		{
			Name:  "ping",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
}

// The xdg_wm_base interface is exposed as a global object enabling clients
// to turn their wl_surfaces into windows in a desktop environment. It
// defines the basic functionality needed for clients and the compositor to
// create windows that can be dragged, resized, maximized, etc, as well as
// creating transient windows such as popup menus.
type WmBase struct{ wlserver.Resource }

func (WmBase) Interface() *wlproto.Interface { return WmBaseInterface }

type WmBaseImplementation interface {
	Destroy(obj WmBase)
	CreatePositioner(obj WmBase, id Positioner) PositionerImplementation
	GetXdgSurface(obj WmBase, id Surface, surface wayland.Surface) SurfaceImplementation
	Pong(obj WmBase, serial uint32)
}

func AddWmBaseGlobal(dsp *wlserver.Display, version int, bind func(res WmBase) WmBaseImplementation) {
	dsp.AddGlobal(WmBaseInterface, version, func(res wlserver.Object) wlserver.ResourceImplementation { return bind(res.(WmBase)) })
}

// The ping event asks the client if it's still alive. Pass the
// serial specified in the event back to the compositor by sending
// a "pong" request back with the specified serial. See xdg_wm_base.pong.
//
// Compositors can use this to determine if the client is still
// alive. It's unspecified what will happen if the client doesn't
// respond to the ping request, or in what timeframe. Clients should
// try to respond in a reasonable amount of time.
//
// A compositor is free to ping in any way it wants, but a client must
// always respond to any xdg_wm_base object it created.
func (obj WmBase) Ping(serial uint32) {
	obj.Conn().SendEvent(obj, 0, serial)
}

type PositionerError uint32

const (
	// invalid input provided
	PositionerErrorInvalidInput PositionerError = 0
)

type PositionerAnchor uint32

const (
	PositionerAnchorNone        PositionerAnchor = 0
	PositionerAnchorTop         PositionerAnchor = 1
	PositionerAnchorBottom      PositionerAnchor = 2
	PositionerAnchorLeft        PositionerAnchor = 3
	PositionerAnchorRight       PositionerAnchor = 4
	PositionerAnchorTopLeft     PositionerAnchor = 5
	PositionerAnchorBottomLeft  PositionerAnchor = 6
	PositionerAnchorTopRight    PositionerAnchor = 7
	PositionerAnchorBottomRight PositionerAnchor = 8
)

type PositionerGravity uint32

const (
	PositionerGravityNone        PositionerGravity = 0
	PositionerGravityTop         PositionerGravity = 1
	PositionerGravityBottom      PositionerGravity = 2
	PositionerGravityLeft        PositionerGravity = 3
	PositionerGravityRight       PositionerGravity = 4
	PositionerGravityTopLeft     PositionerGravity = 5
	PositionerGravityBottomLeft  PositionerGravity = 6
	PositionerGravityTopRight    PositionerGravity = 7
	PositionerGravityBottomRight PositionerGravity = 8
)

// The constraint adjustment value define ways the compositor will adjust
// the position of the surface, if the unadjusted position would result
// in the surface being partly constrained.
//
// Whether a surface is considered 'constrained' is left to the compositor
// to determine. For example, the surface may be partly outside the
// compositor's defined 'work area', thus necessitating the child surface's
// position be adjusted until it is entirely inside the work area.
//
// The adjustments can be combined, according to a defined precedence: 1)
// Flip, 2) Slide, 3) Resize.
type PositionerConstraintAdjustment uint32

const (
	// Don't alter the surface position even if it is constrained on some
	// axis, for example partially outside the edge of an output.
	PositionerConstraintAdjustmentNone PositionerConstraintAdjustment = 0
	// Slide the surface along the x axis until it is no longer constrained.
	//
	// First try to slide towards the direction of the gravity on the x axis
	// until either the edge in the opposite direction of the gravity is
	// unconstrained or the edge in the direction of the gravity is
	// constrained.
	//
	// Then try to slide towards the opposite direction of the gravity on the
	// x axis until either the edge in the direction of the gravity is
	// unconstrained or the edge in the opposite direction of the gravity is
	// constrained.
	PositionerConstraintAdjustmentSlideX PositionerConstraintAdjustment = 1
	// Slide the surface along the y axis until it is no longer constrained.
	//
	// First try to slide towards the direction of the gravity on the y axis
	// until either the edge in the opposite direction of the gravity is
	// unconstrained or the edge in the direction of the gravity is
	// constrained.
	//
	// Then try to slide towards the opposite direction of the gravity on the
	// y axis until either the edge in the direction of the gravity is
	// unconstrained or the edge in the opposite direction of the gravity is
	// constrained.
	PositionerConstraintAdjustmentSlideY PositionerConstraintAdjustment = 2
	// Invert the anchor and gravity on the x axis if the surface is
	// constrained on the x axis. For example, if the left edge of the
	// surface is constrained, the gravity is 'left' and the anchor is
	// 'left', change the gravity to 'right' and the anchor to 'right'.
	//
	// If the adjusted position also ends up being constrained, the resulting
	// position of the flip_x adjustment will be the one before the
	// adjustment.
	PositionerConstraintAdjustmentFlipX PositionerConstraintAdjustment = 4
	// Invert the anchor and gravity on the y axis if the surface is
	// constrained on the y axis. For example, if the bottom edge of the
	// surface is constrained, the gravity is 'bottom' and the anchor is
	// 'bottom', change the gravity to 'top' and the anchor to 'top'.
	//
	// The adjusted position is calculated given the original anchor
	// rectangle and offset, but with the new flipped anchor and gravity
	// values.
	//
	// If the adjusted position also ends up being constrained, the resulting
	// position of the flip_y adjustment will be the one before the
	// adjustment.
	PositionerConstraintAdjustmentFlipY PositionerConstraintAdjustment = 8
	// Resize the surface horizontally so that it is completely
	// unconstrained.
	PositionerConstraintAdjustmentResizeX PositionerConstraintAdjustment = 16
	// Resize the surface vertically so that it is completely unconstrained.
	PositionerConstraintAdjustmentResizeY PositionerConstraintAdjustment = 32
)

var PositionerInterface = &wlproto.Interface{
	Name:    "xdg_positioner",
	Version: 4,
	Type:    reflect.TypeOf(Positioner{}),
	Requests: []wlproto.Request{
		{
			Name:   "destroy",
			Type:   "destructor",
			Since:  1,
			Method: reflect.ValueOf(PositionerImplementation.Destroy),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "set_size",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(PositionerImplementation.SetSize),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:   "set_anchor_rect",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(PositionerImplementation.SetAnchorRect),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:   "set_anchor",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(PositionerImplementation.SetAnchor),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint, Aux: reflect.TypeOf(PositionerAnchor(0))},
			},
		},
		{
			Name:   "set_gravity",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(PositionerImplementation.SetGravity),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint, Aux: reflect.TypeOf(PositionerGravity(0))},
			},
		},
		{
			Name:   "set_constraint_adjustment",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(PositionerImplementation.SetConstraintAdjustment),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
		{
			Name:   "set_offset",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(PositionerImplementation.SetOffset),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:   "set_reactive",
			Type:   "",
			Since:  3,
			Method: reflect.ValueOf(PositionerImplementation.SetReactive),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "set_parent_size",
			Type:   "",
			Since:  3,
			Method: reflect.ValueOf(PositionerImplementation.SetParentSize),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:   "set_parent_configure",
			Type:   "",
			Since:  3,
			Method: reflect.ValueOf(PositionerImplementation.SetParentConfigure),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
	Events: []wlproto.Event{},
}

// The xdg_positioner provides a collection of rules for the placement of a
// child surface relative to a parent surface. Rules can be defined to ensure
// the child surface remains within the visible area's borders, and to
// specify how the child surface changes its position, such as sliding along
// an axis, or flipping around a rectangle. These positioner-created rules are
// constrained by the requirement that a child surface must intersect with or
// be at least partially adjacent to its parent surface.
//
// See the various requests for details about possible rules.
//
// At the time of the request, the compositor makes a copy of the rules
// specified by the xdg_positioner. Thus, after the request is complete the
// xdg_positioner object can be destroyed or reused; further changes to the
// object will have no effect on previous usages.
//
// For an xdg_positioner object to be considered complete, it must have a
// non-zero size set by set_size, and a non-zero anchor rectangle set by
// set_anchor_rect. Passing an incomplete xdg_positioner object when
// positioning a surface raises an error.
type Positioner struct{ wlserver.Resource }

func (Positioner) Interface() *wlproto.Interface { return PositionerInterface }

type PositionerImplementation interface {
	Destroy(obj Positioner)
	SetSize(obj Positioner, width int32, height int32)
	SetAnchorRect(obj Positioner, x int32, y int32, width int32, height int32)
	SetAnchor(obj Positioner, anchor PositionerAnchor)
	SetGravity(obj Positioner, gravity PositionerGravity)
	SetConstraintAdjustment(obj Positioner, constraintAdjustment uint32)
	SetOffset(obj Positioner, x int32, y int32)
	SetReactive(obj Positioner)
	SetParentSize(obj Positioner, parentWidth int32, parentHeight int32)
	SetParentConfigure(obj Positioner, serial uint32)
}

func AddPositionerGlobal(dsp *wlserver.Display, version int, bind func(res Positioner) PositionerImplementation) {
	dsp.AddGlobal(PositionerInterface, version, func(res wlserver.Object) wlserver.ResourceImplementation { return bind(res.(Positioner)) })
}

type SurfaceError uint32

const (
	SurfaceErrorNotConstructed     SurfaceError = 1
	SurfaceErrorAlreadyConstructed SurfaceError = 2
	SurfaceErrorUnconfiguredBuffer SurfaceError = 3
)

var SurfaceInterface = &wlproto.Interface{
	Name:    "xdg_surface",
	Version: 4,
	Type:    reflect.TypeOf(Surface{}),
	Requests: []wlproto.Request{
		{
			Name:   "destroy",
			Type:   "destructor",
			Since:  1,
			Method: reflect.ValueOf(SurfaceImplementation.Destroy),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "get_toplevel",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(SurfaceImplementation.GetToplevel),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeNewID, Aux: reflect.TypeOf(Toplevel{})},
			},
		},
		{
			Name:   "get_popup",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(SurfaceImplementation.GetPopup),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeNewID, Aux: reflect.TypeOf(Popup{})},
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(Surface{})},
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(Positioner{})},
			},
		},
		{
			Name:   "set_window_geometry",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(SurfaceImplementation.SetWindowGeometry),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:   "ack_configure",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(SurfaceImplementation.AckConfigure),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
	Events: []wlproto.Event{
		{
			Name:  "configure",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
}

// An interface that may be implemented by a wl_surface, for
// implementations that provide a desktop-style user interface.
//
// It provides a base set of functionality required to construct user
// interface elements requiring management by the compositor, such as
// toplevel windows, menus, etc. The types of functionality are split into
// xdg_surface roles.
//
// Creating an xdg_surface does not set the role for a wl_surface. In order
// to map an xdg_surface, the client must create a role-specific object
// using, e.g., get_toplevel, get_popup. The wl_surface for any given
// xdg_surface can have at most one role, and may not be assigned any role
// not based on xdg_surface.
//
// A role must be assigned before any other requests are made to the
// xdg_surface object.
//
// The client must call wl_surface.commit on the corresponding wl_surface
// for the xdg_surface state to take effect.
//
// Creating an xdg_surface from a wl_surface which has a buffer attached or
// committed is a client error, and any attempts by a client to attach or
// manipulate a buffer prior to the first xdg_surface.configure call must
// also be treated as errors.
//
// After creating a role-specific object and setting it up, the client must
// perform an initial commit without any buffer attached. The compositor
// will reply with an xdg_surface.configure event. The client must
// acknowledge it and is then allowed to attach a buffer to map the surface.
//
// Mapping an xdg_surface-based role surface is defined as making it
// possible for the surface to be shown by the compositor. Note that
// a mapped surface is not guaranteed to be visible once it is mapped.
//
// For an xdg_surface to be mapped by the compositor, the following
// conditions must be met:
// (1) the client has assigned an xdg_surface-based role to the surface
// (2) the client has set and committed the xdg_surface state and the
// role-dependent state to the surface
// (3) the client has committed a buffer to the surface
//
// A newly-unmapped surface is considered to have met condition (1) out
// of the 3 required conditions for mapping a surface if its role surface
// has not been destroyed, i.e. the client must perform the initial commit
// again before attaching a buffer.
type Surface struct{ wlserver.Resource }

func (Surface) Interface() *wlproto.Interface { return SurfaceInterface }

type SurfaceImplementation interface {
	Destroy(obj Surface)
	GetToplevel(obj Surface, id Toplevel) ToplevelImplementation
	GetPopup(obj Surface, id Popup, parent Surface, positioner Positioner) PopupImplementation
	SetWindowGeometry(obj Surface, x int32, y int32, width int32, height int32)
	AckConfigure(obj Surface, serial uint32)
}

func AddSurfaceGlobal(dsp *wlserver.Display, version int, bind func(res Surface) SurfaceImplementation) {
	dsp.AddGlobal(SurfaceInterface, version, func(res wlserver.Object) wlserver.ResourceImplementation { return bind(res.(Surface)) })
}

// The configure event marks the end of a configure sequence. A configure
// sequence is a set of one or more events configuring the state of the
// xdg_surface, including the final xdg_surface.configure event.
//
// Where applicable, xdg_surface surface roles will during a configure
// sequence extend this event as a latched state sent as events before the
// xdg_surface.configure event. Such events should be considered to make up
// a set of atomically applied configuration states, where the
// xdg_surface.configure commits the accumulated state.
//
// Clients should arrange their surface for the new states, and then send
// an ack_configure request with the serial sent in this configure event at
// some point before committing the new surface.
//
// If the client receives multiple configure events before it can respond
// to one, it is free to discard all but the last event it received.
func (obj Surface) Configure(serial uint32) {
	obj.Conn().SendEvent(obj, 0, serial)
}

type ToplevelError uint32

const (
	// provided value is not a valid variant of the resize_edge enum
	ToplevelErrorInvalidResizeEdge ToplevelError = 0
)

// These values are used to indicate which edge of a surface
// is being dragged in a resize operation.
type ToplevelResizeEdge uint32

const (
	ToplevelResizeEdgeNone        ToplevelResizeEdge = 0
	ToplevelResizeEdgeTop         ToplevelResizeEdge = 1
	ToplevelResizeEdgeBottom      ToplevelResizeEdge = 2
	ToplevelResizeEdgeLeft        ToplevelResizeEdge = 4
	ToplevelResizeEdgeTopLeft     ToplevelResizeEdge = 5
	ToplevelResizeEdgeBottomLeft  ToplevelResizeEdge = 6
	ToplevelResizeEdgeRight       ToplevelResizeEdge = 8
	ToplevelResizeEdgeTopRight    ToplevelResizeEdge = 9
	ToplevelResizeEdgeBottomRight ToplevelResizeEdge = 10
)

// The different state values used on the surface. This is designed for
// state values like maximized, fullscreen. It is paired with the
// configure event to ensure that both the client and the compositor
// setting the state can be synchronized.
//
// States set in this way are double-buffered. They will get applied on
// the next commit.
type ToplevelState uint32

const (
	// The surface is maximized. The window geometry specified in the configure
	// event must be obeyed by the client.
	//
	// The client should draw without shadow or other
	// decoration outside of the window geometry.
	ToplevelStateMaximized ToplevelState = 1
	// The surface is fullscreen. The window geometry specified in the
	// configure event is a maximum; the client cannot resize beyond it. For
	// a surface to cover the whole fullscreened area, the geometry
	// dimensions must be obeyed by the client. For more details, see
	// xdg_toplevel.set_fullscreen.
	ToplevelStateFullscreen ToplevelState = 2
	// The surface is being resized. The window geometry specified in the
	// configure event is a maximum; the client cannot resize beyond it.
	// Clients that have aspect ratio or cell sizing configuration can use
	// a smaller size, however.
	ToplevelStateResizing ToplevelState = 3
	// Client window decorations should be painted as if the window is
	// active. Do not assume this means that the window actually has
	// keyboard or pointer focus.
	ToplevelStateActivated ToplevelState = 4
	// The window is currently in a tiled layout and the left edge is
	// considered to be adjacent to another part of the tiling grid.
	ToplevelStateTiledLeft ToplevelState = 5
	// The window is currently in a tiled layout and the right edge is
	// considered to be adjacent to another part of the tiling grid.
	ToplevelStateTiledRight ToplevelState = 6
	// The window is currently in a tiled layout and the top edge is
	// considered to be adjacent to another part of the tiling grid.
	ToplevelStateTiledTop ToplevelState = 7
	// The window is currently in a tiled layout and the bottom edge is
	// considered to be adjacent to another part of the tiling grid.
	ToplevelStateTiledBottom ToplevelState = 8
)

var ToplevelInterface = &wlproto.Interface{
	Name:    "xdg_toplevel",
	Version: 4,
	Type:    reflect.TypeOf(Toplevel{}),
	Requests: []wlproto.Request{
		{
			Name:   "destroy",
			Type:   "destructor",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.Destroy),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "set_parent",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.SetParent),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(Toplevel{})},
			},
		},
		{
			Name:   "set_title",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.SetTitle),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeString},
			},
		},
		{
			Name:   "set_app_id",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.SetAppID),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeString},
			},
		},
		{
			Name:   "show_window_menu",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.ShowWindowMenu),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(wayland.Seat{})},
				{Type: wlproto.ArgTypeUint},
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:   "move",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.Move),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(wayland.Seat{})},
				{Type: wlproto.ArgTypeUint},
			},
		},
		{
			Name:   "resize",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.Resize),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(wayland.Seat{})},
				{Type: wlproto.ArgTypeUint},
				{Type: wlproto.ArgTypeUint, Aux: reflect.TypeOf(ToplevelResizeEdge(0))},
			},
		},
		{
			Name:   "set_max_size",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.SetMaxSize),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:   "set_min_size",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.SetMinSize),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:   "set_maximized",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.SetMaximized),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "unset_maximized",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.UnsetMaximized),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "set_fullscreen",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.SetFullscreen),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(wayland.Output{})},
			},
		},
		{
			Name:   "unset_fullscreen",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.UnsetFullscreen),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "set_minimized",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(ToplevelImplementation.SetMinimized),
			Args:   []wlproto.Arg{},
		},
	},
	Events: []wlproto.Event{
		{
			Name:  "configure",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeArray},
			},
		},
		{
			Name:  "close",
			Since: 1,
			Args:  []wlproto.Arg{},
		},
		{
			Name:  "configure_bounds",
			Since: 4,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
	},
}

// This interface defines an xdg_surface role which allows a surface to,
// among other things, set window-like properties such as maximize,
// fullscreen, and minimize, set application-specific metadata like title and
// id, and well as trigger user interactive operations such as interactive
// resize and move.
//
// Unmapping an xdg_toplevel means that the surface cannot be shown
// by the compositor until it is explicitly mapped again.
// All active operations (e.g., move, resize) are canceled and all
// attributes (e.g. title, state, stacking, ...) are discarded for
// an xdg_toplevel surface when it is unmapped. The xdg_toplevel returns to
// the state it had right after xdg_surface.get_toplevel. The client
// can re-map the toplevel by perfoming a commit without any buffer
// attached, waiting for a configure event and handling it as usual (see
// xdg_surface description).
//
// Attaching a null buffer to a toplevel unmaps the surface.
type Toplevel struct{ wlserver.Resource }

func (Toplevel) Interface() *wlproto.Interface { return ToplevelInterface }

type ToplevelImplementation interface {
	Destroy(obj Toplevel)
	SetParent(obj Toplevel, parent Toplevel)
	SetTitle(obj Toplevel, title string)
	SetAppID(obj Toplevel, appId string)
	ShowWindowMenu(obj Toplevel, seat wayland.Seat, serial uint32, x int32, y int32)
	Move(obj Toplevel, seat wayland.Seat, serial uint32)
	Resize(obj Toplevel, seat wayland.Seat, serial uint32, edges ToplevelResizeEdge)
	SetMaxSize(obj Toplevel, width int32, height int32)
	SetMinSize(obj Toplevel, width int32, height int32)
	SetMaximized(obj Toplevel)
	UnsetMaximized(obj Toplevel)
	SetFullscreen(obj Toplevel, output wayland.Output)
	UnsetFullscreen(obj Toplevel)
	SetMinimized(obj Toplevel)
}

func AddToplevelGlobal(dsp *wlserver.Display, version int, bind func(res Toplevel) ToplevelImplementation) {
	dsp.AddGlobal(ToplevelInterface, version, func(res wlserver.Object) wlserver.ResourceImplementation { return bind(res.(Toplevel)) })
}

// This configure event asks the client to resize its toplevel surface or
// to change its state. The configured state should not be applied
// immediately. See xdg_surface.configure for details.
//
// The width and height arguments specify a hint to the window
// about how its surface should be resized in window geometry
// coordinates. See set_window_geometry.
//
// If the width or height arguments are zero, it means the client
// should decide its own window dimension. This may happen when the
// compositor needs to configure the state of the surface but doesn't
// have any information about any previous or expected dimension.
//
// The states listed in the event specify how the width/height
// arguments should be interpreted, and possibly how it should be
// drawn.
//
// Clients must send an ack_configure in response to this event. See
// xdg_surface.configure and xdg_surface.ack_configure for details.
func (obj Toplevel) Configure(width int32, height int32, states []byte) {
	obj.Conn().SendEvent(obj, 0, width, height, states)
}

// The close event is sent by the compositor when the user
// wants the surface to be closed. This should be equivalent to
// the user clicking the close button in client-side decorations,
// if your application has any.
//
// This is only a request that the user intends to close the
// window. The client may choose to ignore this request, or show
// a dialog to ask the user to save their data, etc.
func (obj Toplevel) Close() {
	obj.Conn().SendEvent(obj, 1)
}

// The configure_bounds event may be sent prior to a xdg_toplevel.configure
// event to communicate the bounds a window geometry size is recommended
// to constrain to.
//
// The passed width and height are in surface coordinate space. If width
// and height are 0, it means bounds is unknown and equivalent to as if no
// configure_bounds event was ever sent for this surface.
//
// The bounds can for example correspond to the size of a monitor excluding
// any panels or other shell components, so that a surface isn't created in
// a way that it cannot fit.
//
// The bounds may change at any point, and in such a case, a new
// xdg_toplevel.configure_bounds will be sent, followed by
// xdg_toplevel.configure and xdg_surface.configure.
func (obj Toplevel) ConfigureBounds(width int32, height int32) {
	obj.Conn().SendEvent(obj, 2, width, height)
}

type PopupError uint32

const (
	// tried to grab after being mapped
	PopupErrorInvalidGrab PopupError = 0
)

var PopupInterface = &wlproto.Interface{
	Name:    "xdg_popup",
	Version: 4,
	Type:    reflect.TypeOf(Popup{}),
	Requests: []wlproto.Request{
		{
			Name:   "destroy",
			Type:   "destructor",
			Since:  1,
			Method: reflect.ValueOf(PopupImplementation.Destroy),
			Args:   []wlproto.Arg{},
		},
		{
			Name:   "grab",
			Type:   "",
			Since:  1,
			Method: reflect.ValueOf(PopupImplementation.Grab),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(wayland.Seat{})},
				{Type: wlproto.ArgTypeUint},
			},
		},
		{
			Name:   "reposition",
			Type:   "",
			Since:  3,
			Method: reflect.ValueOf(PopupImplementation.Reposition),
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeObject, Aux: reflect.TypeOf(Positioner{})},
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
	Events: []wlproto.Event{
		{
			Name:  "configure",
			Since: 1,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
				{Type: wlproto.ArgTypeInt},
			},
		},
		{
			Name:  "popup_done",
			Since: 1,
			Args:  []wlproto.Arg{},
		},
		{
			Name:  "repositioned",
			Since: 3,
			Args: []wlproto.Arg{
				{Type: wlproto.ArgTypeUint},
			},
		},
	},
}

// A popup surface is a short-lived, temporary surface. It can be used to
// implement for example menus, popovers, tooltips and other similar user
// interface concepts.
//
// A popup can be made to take an explicit grab. See xdg_popup.grab for
// details.
//
// When the popup is dismissed, a popup_done event will be sent out, and at
// the same time the surface will be unmapped. See the xdg_popup.popup_done
// event for details.
//
// Explicitly destroying the xdg_popup object will also dismiss the popup and
// unmap the surface. Clients that want to dismiss the popup when another
// surface of their own is clicked should dismiss the popup using the destroy
// request.
//
// A newly created xdg_popup will be stacked on top of all previously created
// xdg_popup surfaces associated with the same xdg_toplevel.
//
// The parent of an xdg_popup must be mapped (see the xdg_surface
// description) before the xdg_popup itself.
//
// The client must call wl_surface.commit on the corresponding wl_surface
// for the xdg_popup state to take effect.
type Popup struct{ wlserver.Resource }

func (Popup) Interface() *wlproto.Interface { return PopupInterface }

type PopupImplementation interface {
	Destroy(obj Popup)
	Grab(obj Popup, seat wayland.Seat, serial uint32)
	Reposition(obj Popup, positioner Positioner, token uint32)
}

func AddPopupGlobal(dsp *wlserver.Display, version int, bind func(res Popup) PopupImplementation) {
	dsp.AddGlobal(PopupInterface, version, func(res wlserver.Object) wlserver.ResourceImplementation { return bind(res.(Popup)) })
}

// This event asks the popup surface to configure itself given the
// configuration. The configured state should not be applied immediately.
// See xdg_surface.configure for details.
//
// The x and y arguments represent the position the popup was placed at
// given the xdg_positioner rule, relative to the upper left corner of the
// window geometry of the parent surface.
//
// For version 2 or older, the configure event for an xdg_popup is only
// ever sent once for the initial configuration. Starting with version 3,
// it may be sent again if the popup is setup with an xdg_positioner with
// set_reactive requested, or in response to xdg_popup.reposition requests.
func (obj Popup) Configure(x int32, y int32, width int32, height int32) {
	obj.Conn().SendEvent(obj, 0, x, y, width, height)
}

// The popup_done event is sent out when a popup is dismissed by the
// compositor. The client should destroy the xdg_popup object at this
// point.
func (obj Popup) PopupDone() {
	obj.Conn().SendEvent(obj, 1)
}

// The repositioned event is sent as part of a popup configuration
// sequence, together with xdg_popup.configure and lastly
// xdg_surface.configure to notify the completion of a reposition request.
//
// The repositioned event is to notify about the completion of a
// xdg_popup.reposition request. The token argument is the token passed
// in the xdg_popup.reposition request.
//
// Immediately after this event is emitted, xdg_popup.configure and
// xdg_surface.configure will be sent with the updated size and position,
// as well as a new configure serial.
//
// The client should optionally update the content of the popup, but must
// acknowledge the new popup configuration for the new position to take
// effect. See xdg_surface.ack_configure for details.
func (obj Popup) Repositioned(token uint32) {
	obj.Conn().SendEvent(obj, 2, token)
}
