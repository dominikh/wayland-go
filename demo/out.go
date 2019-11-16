package demo

import "honnef.co/go/wayland"

const (
	DisplayErrorInvalidObject  = 0
	DisplayErrorInvalidMethod  = 1
	DisplayErrorNoMemory       = 2
	DisplayErrorImplementation = 3
)

var DisplayInterface = &wayland.Interface{
	Name:    "wl_display",
	Version: 1,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "error",
			Types: []interface{}{nil, uint32(0), ""},
		},
		wayland.MessageEvent{
			Name:  "delete_id",
			Types: []interface{}{uint32(0)},
		}},
}

type Display struct{ wayland.Proxy }

func (obj *Display) Sync() *Callback {
	const wl_display_sync = 0
	_ret := &Callback{}
	obj.Conn().NewProxy(0, _ret, CallbackInterface)
	obj.Conn().SendRequest(obj, wl_display_sync, _ret)
	return _ret
}
func (obj *Display) GetRegistry() *Registry {
	const wl_display_get_registry = 1
	_ret := &Registry{}
	obj.Conn().NewProxy(0, _ret, RegistryInterface)
	obj.Conn().SendRequest(obj, wl_display_get_registry, _ret)
	return _ret
}

var RegistryInterface = &wayland.Interface{
	Name:    "wl_registry",
	Version: 1,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "global",
			Types: []interface{}{uint32(0), "", uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "global_remove",
			Types: []interface{}{uint32(0)},
		}},
}

type Registry struct{ wayland.Proxy }

func (obj *Registry) Bind(name uint32, id wayland.Object) {
	const wl_registry_bind = 0
	obj.Conn().SendRequest(obj, wl_registry_bind, name, id)
}

var CallbackInterface = &wayland.Interface{
	Name:    "wl_callback",
	Version: 1,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "done",
			Types: []interface{}{uint32(0)},
		}},
}

type Callback struct{ wayland.Proxy }

var CompositorInterface = &wayland.Interface{
	Name:    "wl_compositor",
	Version: 4,
	Events:  []wayland.MessageEvent{},
}

type Compositor struct{ wayland.Proxy }

func (obj *Compositor) CreateSurface() *Surface {
	const wl_compositor_create_surface = 0
	_ret := &Surface{}
	obj.Conn().NewProxy(0, _ret, SurfaceInterface)
	obj.Conn().SendRequest(obj, wl_compositor_create_surface, _ret)
	return _ret
}
func (obj *Compositor) CreateRegion() *Region {
	const wl_compositor_create_region = 1
	_ret := &Region{}
	obj.Conn().NewProxy(0, _ret, RegionInterface)
	obj.Conn().SendRequest(obj, wl_compositor_create_region, _ret)
	return _ret
}

var ShmPoolInterface = &wayland.Interface{
	Name:    "wl_shm_pool",
	Version: 1,
	Events:  []wayland.MessageEvent{},
}

type ShmPool struct{ wayland.Proxy }

func (obj *ShmPool) CreateBuffer(offset int32, width int32, height int32, stride int32, format uint32) *Buffer {
	const wl_shm_pool_create_buffer = 0
	_ret := &Buffer{}
	obj.Conn().NewProxy(0, _ret, BufferInterface)
	obj.Conn().SendRequest(obj, wl_shm_pool_create_buffer, _ret, offset, width, height, stride, format)
	return _ret
}
func (obj *ShmPool) Destroy() {
	const wl_shm_pool_destroy = 1
	obj.Conn().SendRequest(obj, wl_shm_pool_destroy)
}
func (obj *ShmPool) Resize(size int32) {
	const wl_shm_pool_resize = 2
	obj.Conn().SendRequest(obj, wl_shm_pool_resize, size)
}

const (
	ShmErrorInvalidFormat = 0
	ShmErrorInvalidStride = 1
	ShmErrorInvalidFd     = 2
)
const (
	ShmFormatArgb8888       = 0
	ShmFormatXrgb8888       = 1
	ShmFormatC8             = 0x20203843
	ShmFormatRgb332         = 0x38424752
	ShmFormatBgr233         = 0x38524742
	ShmFormatXrgb4444       = 0x32315258
	ShmFormatXbgr4444       = 0x32314258
	ShmFormatRgbx4444       = 0x32315852
	ShmFormatBgrx4444       = 0x32315842
	ShmFormatArgb4444       = 0x32315241
	ShmFormatAbgr4444       = 0x32314241
	ShmFormatRgba4444       = 0x32314152
	ShmFormatBgra4444       = 0x32314142
	ShmFormatXrgb1555       = 0x35315258
	ShmFormatXbgr1555       = 0x35314258
	ShmFormatRgbx5551       = 0x35315852
	ShmFormatBgrx5551       = 0x35315842
	ShmFormatArgb1555       = 0x35315241
	ShmFormatAbgr1555       = 0x35314241
	ShmFormatRgba5551       = 0x35314152
	ShmFormatBgra5551       = 0x35314142
	ShmFormatRgb565         = 0x36314752
	ShmFormatBgr565         = 0x36314742
	ShmFormatRgb888         = 0x34324752
	ShmFormatBgr888         = 0x34324742
	ShmFormatXbgr8888       = 0x34324258
	ShmFormatRgbx8888       = 0x34325852
	ShmFormatBgrx8888       = 0x34325842
	ShmFormatAbgr8888       = 0x34324241
	ShmFormatRgba8888       = 0x34324152
	ShmFormatBgra8888       = 0x34324142
	ShmFormatXrgb2101010    = 0x30335258
	ShmFormatXbgr2101010    = 0x30334258
	ShmFormatRgbx1010102    = 0x30335852
	ShmFormatBgrx1010102    = 0x30335842
	ShmFormatArgb2101010    = 0x30335241
	ShmFormatAbgr2101010    = 0x30334241
	ShmFormatRgba1010102    = 0x30334152
	ShmFormatBgra1010102    = 0x30334142
	ShmFormatYuyv           = 0x56595559
	ShmFormatYvyu           = 0x55595659
	ShmFormatUyvy           = 0x59565955
	ShmFormatVyuy           = 0x59555956
	ShmFormatAyuv           = 0x56555941
	ShmFormatNv12           = 0x3231564e
	ShmFormatNv21           = 0x3132564e
	ShmFormatNv16           = 0x3631564e
	ShmFormatNv61           = 0x3136564e
	ShmFormatYuv410         = 0x39565559
	ShmFormatYvu410         = 0x39555659
	ShmFormatYuv411         = 0x31315559
	ShmFormatYvu411         = 0x31315659
	ShmFormatYuv420         = 0x32315559
	ShmFormatYvu420         = 0x32315659
	ShmFormatYuv422         = 0x36315559
	ShmFormatYvu422         = 0x36315659
	ShmFormatYuv444         = 0x34325559
	ShmFormatYvu444         = 0x34325659
	ShmFormatR8             = 0x20203852
	ShmFormatR16            = 0x20363152
	ShmFormatRg88           = 0x38384752
	ShmFormatGr88           = 0x38385247
	ShmFormatRg1616         = 0x32334752
	ShmFormatGr1616         = 0x32335247
	ShmFormatXrgb16161616f  = 0x48345258
	ShmFormatXbgr16161616f  = 0x48344258
	ShmFormatArgb16161616f  = 0x48345241
	ShmFormatAbgr16161616f  = 0x48344241
	ShmFormatXyuv8888       = 0x56555958
	ShmFormatVuy888         = 0x34325556
	ShmFormatVuy101010      = 0x30335556
	ShmFormatY210           = 0x30313259
	ShmFormatY212           = 0x32313259
	ShmFormatY216           = 0x36313259
	ShmFormatY410           = 0x30313459
	ShmFormatY412           = 0x32313459
	ShmFormatY416           = 0x36313459
	ShmFormatXvyu2101010    = 0x30335658
	ShmFormatXvyu1216161616 = 0x36335658
	ShmFormatXvyu16161616   = 0x38345658
	ShmFormatY0l0           = 0x304c3059
	ShmFormatX0l0           = 0x304c3058
	ShmFormatY0l2           = 0x324c3059
	ShmFormatX0l2           = 0x324c3058
	ShmFormatYuv4208bit     = 0x38305559
	ShmFormatYuv42010bit    = 0x30315559
	ShmFormatXrgb8888A8     = 0x38415258
	ShmFormatXbgr8888A8     = 0x38414258
	ShmFormatRgbx8888A8     = 0x38415852
	ShmFormatBgrx8888A8     = 0x38415842
	ShmFormatRgb888A8       = 0x38413852
	ShmFormatBgr888A8       = 0x38413842
	ShmFormatRgb565A8       = 0x38413552
	ShmFormatBgr565A8       = 0x38413542
	ShmFormatNv24           = 0x3432564e
	ShmFormatNv42           = 0x3234564e
	ShmFormatP210           = 0x30313250
	ShmFormatP010           = 0x30313050
	ShmFormatP012           = 0x32313050
	ShmFormatP016           = 0x36313050
)

var ShmInterface = &wayland.Interface{
	Name:    "wl_shm",
	Version: 1,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "format",
			Types: []interface{}{uint32(0)},
		}},
}

type Shm struct{ wayland.Proxy }

func (obj *Shm) CreatePool(fd int32, size int32) *ShmPool {
	const wl_shm_create_pool = 0
	_ret := &ShmPool{}
	obj.Conn().NewProxy(0, _ret, ShmPoolInterface)
	obj.Conn().SendRequest(obj, wl_shm_create_pool, _ret, fd, size)
	return _ret
}

var BufferInterface = &wayland.Interface{
	Name:    "wl_buffer",
	Version: 1,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "release",
			Types: []interface{}{},
		}},
}

type Buffer struct{ wayland.Proxy }

func (obj *Buffer) Destroy() {
	const wl_buffer_destroy = 0
	obj.Conn().SendRequest(obj, wl_buffer_destroy)
}

const (
	DataOfferErrorInvalidFinish     = 0
	DataOfferErrorInvalidActionMask = 1
	DataOfferErrorInvalidAction     = 2
	DataOfferErrorInvalidOffer      = 3
)

var DataOfferInterface = &wayland.Interface{
	Name:    "wl_data_offer",
	Version: 3,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "offer",
			Types: []interface{}{""},
		},
		wayland.MessageEvent{
			Name:  "source_actions",
			Types: []interface{}{uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "action",
			Types: []interface{}{uint32(0)},
		}},
}

type DataOffer struct{ wayland.Proxy }

func (obj *DataOffer) Accept(serial uint32, mimeType string) {
	const wl_data_offer_accept = 0
	obj.Conn().SendRequest(obj, wl_data_offer_accept, serial, mimeType)
}
func (obj *DataOffer) Receive(mimeType string, fd int32) {
	const wl_data_offer_receive = 1
	obj.Conn().SendRequest(obj, wl_data_offer_receive, mimeType, fd)
}
func (obj *DataOffer) Destroy() {
	const wl_data_offer_destroy = 2
	obj.Conn().SendRequest(obj, wl_data_offer_destroy)
}
func (obj *DataOffer) Finish() {
	const wl_data_offer_finish = 3
	obj.Conn().SendRequest(obj, wl_data_offer_finish)
}
func (obj *DataOffer) SetActions(dndActions uint32, preferredAction uint32) {
	const wl_data_offer_set_actions = 4
	obj.Conn().SendRequest(obj, wl_data_offer_set_actions, dndActions, preferredAction)
}

const (
	DataSourceErrorInvalidActionMask = 0
	DataSourceErrorInvalidSource     = 1
)

var DataSourceInterface = &wayland.Interface{
	Name:    "wl_data_source",
	Version: 3,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "target",
			Types: []interface{}{""},
		},
		wayland.MessageEvent{
			Name:  "send",
			Types: []interface{}{"", "XXX"},
		},
		wayland.MessageEvent{
			Name:  "cancelled",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "dnd_drop_performed",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "dnd_finished",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "action",
			Types: []interface{}{uint32(0)},
		}},
}

type DataSource struct{ wayland.Proxy }

func (obj *DataSource) Offer(mimeType string) {
	const wl_data_source_offer = 0
	obj.Conn().SendRequest(obj, wl_data_source_offer, mimeType)
}
func (obj *DataSource) Destroy() {
	const wl_data_source_destroy = 1
	obj.Conn().SendRequest(obj, wl_data_source_destroy)
}
func (obj *DataSource) SetActions(dndActions uint32) {
	const wl_data_source_set_actions = 2
	obj.Conn().SendRequest(obj, wl_data_source_set_actions, dndActions)
}

const (
	DataDeviceErrorRole = 0
)

var DataDeviceInterface = &wayland.Interface{
	Name:    "wl_data_device",
	Version: 3,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "data_offer",
			Types: []interface{}{DataOfferInterface},
		},
		wayland.MessageEvent{
			Name:  "enter",
			Types: []interface{}{uint32(0), nil, wayland.Fixed(0), wayland.Fixed(0), nil},
		},
		wayland.MessageEvent{
			Name:  "leave",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "motion",
			Types: []interface{}{uint32(0), wayland.Fixed(0), wayland.Fixed(0)},
		},
		wayland.MessageEvent{
			Name:  "drop",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "selection",
			Types: []interface{}{nil},
		}},
}

type DataDevice struct{ wayland.Proxy }

func (obj *DataDevice) StartDrag(source *DataSource, origin *Surface, icon *Surface, serial uint32) {
	const wl_data_device_start_drag = 0
	obj.Conn().SendRequest(obj, wl_data_device_start_drag, source, origin, icon, serial)
}
func (obj *DataDevice) SetSelection(source *DataSource, serial uint32) {
	const wl_data_device_set_selection = 1
	obj.Conn().SendRequest(obj, wl_data_device_set_selection, source, serial)
}
func (obj *DataDevice) Release() {
	const wl_data_device_release = 2
	obj.Conn().SendRequest(obj, wl_data_device_release)
}

const (
	DataDeviceManagerDndActionNone = 0
	DataDeviceManagerDndActionCopy = 1
	DataDeviceManagerDndActionMove = 2
	DataDeviceManagerDndActionAsk  = 4
)

var DataDeviceManagerInterface = &wayland.Interface{
	Name:    "wl_data_device_manager",
	Version: 3,
	Events:  []wayland.MessageEvent{},
}

type DataDeviceManager struct{ wayland.Proxy }

func (obj *DataDeviceManager) CreateDataSource() *DataSource {
	const wl_data_device_manager_create_data_source = 0
	_ret := &DataSource{}
	obj.Conn().NewProxy(0, _ret, DataSourceInterface)
	obj.Conn().SendRequest(obj, wl_data_device_manager_create_data_source, _ret)
	return _ret
}
func (obj *DataDeviceManager) GetDataDevice(seat *Seat) *DataDevice {
	const wl_data_device_manager_get_data_device = 1
	_ret := &DataDevice{}
	obj.Conn().NewProxy(0, _ret, DataDeviceInterface)
	obj.Conn().SendRequest(obj, wl_data_device_manager_get_data_device, _ret, seat)
	return _ret
}

const (
	ShellErrorRole = 0
)

var ShellInterface = &wayland.Interface{
	Name:    "wl_shell",
	Version: 1,
	Events:  []wayland.MessageEvent{},
}

type Shell struct{ wayland.Proxy }

func (obj *Shell) GetShellSurface(surface *Surface) *ShellSurface {
	const wl_shell_get_shell_surface = 0
	_ret := &ShellSurface{}
	obj.Conn().NewProxy(0, _ret, ShellSurfaceInterface)
	obj.Conn().SendRequest(obj, wl_shell_get_shell_surface, _ret, surface)
	return _ret
}

const (
	ShellSurfaceResizeNone        = 0
	ShellSurfaceResizeTop         = 1
	ShellSurfaceResizeBottom      = 2
	ShellSurfaceResizeLeft        = 4
	ShellSurfaceResizeTopLeft     = 5
	ShellSurfaceResizeBottomLeft  = 6
	ShellSurfaceResizeRight       = 8
	ShellSurfaceResizeTopRight    = 9
	ShellSurfaceResizeBottomRight = 10
)
const (
	ShellSurfaceTransientInactive = 0x1
)
const (
	ShellSurfaceFullscreenMethodDefault = 0
	ShellSurfaceFullscreenMethodScale   = 1
	ShellSurfaceFullscreenMethodDriver  = 2
	ShellSurfaceFullscreenMethodFill    = 3
)

var ShellSurfaceInterface = &wayland.Interface{
	Name:    "wl_shell_surface",
	Version: 1,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "ping",
			Types: []interface{}{uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "configure",
			Types: []interface{}{uint32(0), int32(0), int32(0)},
		},
		wayland.MessageEvent{
			Name:  "popup_done",
			Types: []interface{}{},
		}},
}

type ShellSurface struct{ wayland.Proxy }

func (obj *ShellSurface) Pong(serial uint32) {
	const wl_shell_surface_pong = 0
	obj.Conn().SendRequest(obj, wl_shell_surface_pong, serial)
}
func (obj *ShellSurface) Move(seat *Seat, serial uint32) {
	const wl_shell_surface_move = 1
	obj.Conn().SendRequest(obj, wl_shell_surface_move, seat, serial)
}
func (obj *ShellSurface) Resize(seat *Seat, serial uint32, edges uint32) {
	const wl_shell_surface_resize = 2
	obj.Conn().SendRequest(obj, wl_shell_surface_resize, seat, serial, edges)
}
func (obj *ShellSurface) SetToplevel() {
	const wl_shell_surface_set_toplevel = 3
	obj.Conn().SendRequest(obj, wl_shell_surface_set_toplevel)
}
func (obj *ShellSurface) SetTransient(parent *Surface, x int32, y int32, flags uint32) {
	const wl_shell_surface_set_transient = 4
	obj.Conn().SendRequest(obj, wl_shell_surface_set_transient, parent, x, y, flags)
}
func (obj *ShellSurface) SetFullscreen(method uint32, framerate uint32, output *Output) {
	const wl_shell_surface_set_fullscreen = 5
	obj.Conn().SendRequest(obj, wl_shell_surface_set_fullscreen, method, framerate, output)
}
func (obj *ShellSurface) SetPopup(seat *Seat, serial uint32, parent *Surface, x int32, y int32, flags uint32) {
	const wl_shell_surface_set_popup = 6
	obj.Conn().SendRequest(obj, wl_shell_surface_set_popup, seat, serial, parent, x, y, flags)
}
func (obj *ShellSurface) SetMaximized(output *Output) {
	const wl_shell_surface_set_maximized = 7
	obj.Conn().SendRequest(obj, wl_shell_surface_set_maximized, output)
}
func (obj *ShellSurface) SetTitle(title string) {
	const wl_shell_surface_set_title = 8
	obj.Conn().SendRequest(obj, wl_shell_surface_set_title, title)
}
func (obj *ShellSurface) SetClass(class string) {
	const wl_shell_surface_set_class = 9
	obj.Conn().SendRequest(obj, wl_shell_surface_set_class, class)
}

const (
	SurfaceErrorInvalidScale     = 0
	SurfaceErrorInvalidTransform = 1
)

var SurfaceInterface = &wayland.Interface{
	Name:    "wl_surface",
	Version: 4,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "enter",
			Types: []interface{}{nil},
		},
		wayland.MessageEvent{
			Name:  "leave",
			Types: []interface{}{nil},
		}},
}

type Surface struct{ wayland.Proxy }

func (obj *Surface) Destroy() {
	const wl_surface_destroy = 0
	obj.Conn().SendRequest(obj, wl_surface_destroy)
}
func (obj *Surface) Attach(buffer *Buffer, x int32, y int32) {
	const wl_surface_attach = 1
	obj.Conn().SendRequest(obj, wl_surface_attach, buffer, x, y)
}
func (obj *Surface) Damage(x int32, y int32, width int32, height int32) {
	const wl_surface_damage = 2
	obj.Conn().SendRequest(obj, wl_surface_damage, x, y, width, height)
}
func (obj *Surface) Frame() *Callback {
	const wl_surface_frame = 3
	_ret := &Callback{}
	obj.Conn().NewProxy(0, _ret, CallbackInterface)
	obj.Conn().SendRequest(obj, wl_surface_frame, _ret)
	return _ret
}
func (obj *Surface) SetOpaqueRegion(region *Region) {
	const wl_surface_set_opaque_region = 4
	obj.Conn().SendRequest(obj, wl_surface_set_opaque_region, region)
}
func (obj *Surface) SetInputRegion(region *Region) {
	const wl_surface_set_input_region = 5
	obj.Conn().SendRequest(obj, wl_surface_set_input_region, region)
}
func (obj *Surface) Commit() {
	const wl_surface_commit = 6
	obj.Conn().SendRequest(obj, wl_surface_commit)
}
func (obj *Surface) SetBufferTransform(transform int32) {
	const wl_surface_set_buffer_transform = 7
	obj.Conn().SendRequest(obj, wl_surface_set_buffer_transform, transform)
}
func (obj *Surface) SetBufferScale(scale int32) {
	const wl_surface_set_buffer_scale = 8
	obj.Conn().SendRequest(obj, wl_surface_set_buffer_scale, scale)
}
func (obj *Surface) DamageBuffer(x int32, y int32, width int32, height int32) {
	const wl_surface_damage_buffer = 9
	obj.Conn().SendRequest(obj, wl_surface_damage_buffer, x, y, width, height)
}

const (
	SeatCapabilityPointer  = 1
	SeatCapabilityKeyboard = 2
	SeatCapabilityTouch    = 4
)

var SeatInterface = &wayland.Interface{
	Name:    "wl_seat",
	Version: 7,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "capabilities",
			Types: []interface{}{uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "name",
			Types: []interface{}{""},
		}},
}

type Seat struct{ wayland.Proxy }

func (obj *Seat) GetPointer() *Pointer {
	const wl_seat_get_pointer = 0
	_ret := &Pointer{}
	obj.Conn().NewProxy(0, _ret, PointerInterface)
	obj.Conn().SendRequest(obj, wl_seat_get_pointer, _ret)
	return _ret
}
func (obj *Seat) GetKeyboard() *Keyboard {
	const wl_seat_get_keyboard = 1
	_ret := &Keyboard{}
	obj.Conn().NewProxy(0, _ret, KeyboardInterface)
	obj.Conn().SendRequest(obj, wl_seat_get_keyboard, _ret)
	return _ret
}
func (obj *Seat) GetTouch() *Touch {
	const wl_seat_get_touch = 2
	_ret := &Touch{}
	obj.Conn().NewProxy(0, _ret, TouchInterface)
	obj.Conn().SendRequest(obj, wl_seat_get_touch, _ret)
	return _ret
}
func (obj *Seat) Release() {
	const wl_seat_release = 3
	obj.Conn().SendRequest(obj, wl_seat_release)
}

const (
	PointerErrorRole = 0
)
const (
	PointerButtonStateReleased = 0
	PointerButtonStatePressed  = 1
)
const (
	PointerAxisVerticalScroll   = 0
	PointerAxisHorizontalScroll = 1
)
const (
	PointerAxisSourceWheel      = 0
	PointerAxisSourceFinger     = 1
	PointerAxisSourceContinuous = 2
	PointerAxisSourceWheelTilt  = 3
)

var PointerInterface = &wayland.Interface{
	Name:    "wl_pointer",
	Version: 7,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "enter",
			Types: []interface{}{uint32(0), nil, wayland.Fixed(0), wayland.Fixed(0)},
		},
		wayland.MessageEvent{
			Name:  "leave",
			Types: []interface{}{uint32(0), nil},
		},
		wayland.MessageEvent{
			Name:  "motion",
			Types: []interface{}{uint32(0), wayland.Fixed(0), wayland.Fixed(0)},
		},
		wayland.MessageEvent{
			Name:  "button",
			Types: []interface{}{uint32(0), uint32(0), uint32(0), uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "axis",
			Types: []interface{}{uint32(0), uint32(0), wayland.Fixed(0)},
		},
		wayland.MessageEvent{
			Name:  "frame",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "axis_source",
			Types: []interface{}{uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "axis_stop",
			Types: []interface{}{uint32(0), uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "axis_discrete",
			Types: []interface{}{uint32(0), int32(0)},
		}},
}

type Pointer struct{ wayland.Proxy }

func (obj *Pointer) SetCursor(serial uint32, surface *Surface, hotspotX int32, hotspotY int32) {
	const wl_pointer_set_cursor = 0
	obj.Conn().SendRequest(obj, wl_pointer_set_cursor, serial, surface, hotspotX, hotspotY)
}
func (obj *Pointer) Release() {
	const wl_pointer_release = 1
	obj.Conn().SendRequest(obj, wl_pointer_release)
}

const (
	KeyboardKeymapFormatNoKeymap = 0
	KeyboardKeymapFormatXkbV1    = 1
)
const (
	KeyboardKeyStateReleased = 0
	KeyboardKeyStatePressed  = 1
)

var KeyboardInterface = &wayland.Interface{
	Name:    "wl_keyboard",
	Version: 7,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "keymap",
			Types: []interface{}{uint32(0), "XXX", uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "enter",
			Types: []interface{}{uint32(0), nil, "XXX"},
		},
		wayland.MessageEvent{
			Name:  "leave",
			Types: []interface{}{uint32(0), nil},
		},
		wayland.MessageEvent{
			Name:  "key",
			Types: []interface{}{uint32(0), uint32(0), uint32(0), uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "modifiers",
			Types: []interface{}{uint32(0), uint32(0), uint32(0), uint32(0), uint32(0)},
		},
		wayland.MessageEvent{
			Name:  "repeat_info",
			Types: []interface{}{int32(0), int32(0)},
		}},
}

type Keyboard struct{ wayland.Proxy }

func (obj *Keyboard) Release() {
	const wl_keyboard_release = 0
	obj.Conn().SendRequest(obj, wl_keyboard_release)
}

var TouchInterface = &wayland.Interface{
	Name:    "wl_touch",
	Version: 7,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "down",
			Types: []interface{}{uint32(0), uint32(0), nil, int32(0), wayland.Fixed(0), wayland.Fixed(0)},
		},
		wayland.MessageEvent{
			Name:  "up",
			Types: []interface{}{uint32(0), uint32(0), int32(0)},
		},
		wayland.MessageEvent{
			Name:  "motion",
			Types: []interface{}{uint32(0), int32(0), wayland.Fixed(0), wayland.Fixed(0)},
		},
		wayland.MessageEvent{
			Name:  "frame",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "cancel",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "shape",
			Types: []interface{}{int32(0), wayland.Fixed(0), wayland.Fixed(0)},
		},
		wayland.MessageEvent{
			Name:  "orientation",
			Types: []interface{}{int32(0), wayland.Fixed(0)},
		}},
}

type Touch struct{ wayland.Proxy }

func (obj *Touch) Release() {
	const wl_touch_release = 0
	obj.Conn().SendRequest(obj, wl_touch_release)
}

const (
	OutputSubpixelUnknown       = 0
	OutputSubpixelNone          = 1
	OutputSubpixelHorizontalRgb = 2
	OutputSubpixelHorizontalBgr = 3
	OutputSubpixelVerticalRgb   = 4
	OutputSubpixelVerticalBgr   = 5
)
const (
	OutputTransformNormal     = 0
	OutputTransform90         = 1
	OutputTransform180        = 2
	OutputTransform270        = 3
	OutputTransformFlipped    = 4
	OutputTransformFlipped90  = 5
	OutputTransformFlipped180 = 6
	OutputTransformFlipped270 = 7
)
const (
	OutputModeCurrent   = 0x1
	OutputModePreferred = 0x2
)

var OutputInterface = &wayland.Interface{
	Name:    "wl_output",
	Version: 3,
	Events: []wayland.MessageEvent{
		wayland.MessageEvent{
			Name:  "geometry",
			Types: []interface{}{int32(0), int32(0), int32(0), int32(0), int32(0), "", "", int32(0)},
		},
		wayland.MessageEvent{
			Name:  "mode",
			Types: []interface{}{uint32(0), int32(0), int32(0), int32(0)},
		},
		wayland.MessageEvent{
			Name:  "done",
			Types: []interface{}{},
		},
		wayland.MessageEvent{
			Name:  "scale",
			Types: []interface{}{int32(0)},
		}},
}

type Output struct{ wayland.Proxy }

func (obj *Output) Release() {
	const wl_output_release = 0
	obj.Conn().SendRequest(obj, wl_output_release)
}

var RegionInterface = &wayland.Interface{
	Name:    "wl_region",
	Version: 1,
	Events:  []wayland.MessageEvent{},
}

type Region struct{ wayland.Proxy }

func (obj *Region) Destroy() {
	const wl_region_destroy = 0
	obj.Conn().SendRequest(obj, wl_region_destroy)
}
func (obj *Region) Add(x int32, y int32, width int32, height int32) {
	const wl_region_add = 1
	obj.Conn().SendRequest(obj, wl_region_add, x, y, width, height)
}
func (obj *Region) Subtract(x int32, y int32, width int32, height int32) {
	const wl_region_subtract = 2
	obj.Conn().SendRequest(obj, wl_region_subtract, x, y, width, height)
}

const (
	SubcompositorErrorBadSurface = 0
)

var SubcompositorInterface = &wayland.Interface{
	Name:    "wl_subcompositor",
	Version: 1,
	Events:  []wayland.MessageEvent{},
}

type Subcompositor struct{ wayland.Proxy }

func (obj *Subcompositor) Destroy() {
	const wl_subcompositor_destroy = 0
	obj.Conn().SendRequest(obj, wl_subcompositor_destroy)
}
func (obj *Subcompositor) GetSubsurface(surface *Surface, parent *Surface) *Subsurface {
	const wl_subcompositor_get_subsurface = 1
	_ret := &Subsurface{}
	obj.Conn().NewProxy(0, _ret, SubsurfaceInterface)
	obj.Conn().SendRequest(obj, wl_subcompositor_get_subsurface, _ret, surface, parent)
	return _ret
}

const (
	SubsurfaceErrorBadSurface = 0
)

var SubsurfaceInterface = &wayland.Interface{
	Name:    "wl_subsurface",
	Version: 1,
	Events:  []wayland.MessageEvent{},
}

type Subsurface struct{ wayland.Proxy }

func (obj *Subsurface) Destroy() {
	const wl_subsurface_destroy = 0
	obj.Conn().SendRequest(obj, wl_subsurface_destroy)
}
func (obj *Subsurface) SetPosition(x int32, y int32) {
	const wl_subsurface_set_position = 1
	obj.Conn().SendRequest(obj, wl_subsurface_set_position, x, y)
}
func (obj *Subsurface) PlaceAbove(sibling *Surface) {
	const wl_subsurface_place_above = 2
	obj.Conn().SendRequest(obj, wl_subsurface_place_above, sibling)
}
func (obj *Subsurface) PlaceBelow(sibling *Surface) {
	const wl_subsurface_place_below = 3
	obj.Conn().SendRequest(obj, wl_subsurface_place_below, sibling)
}
func (obj *Subsurface) SetSync() {
	const wl_subsurface_set_sync = 4
	obj.Conn().SendRequest(obj, wl_subsurface_set_sync)
}
func (obj *Subsurface) SetDesync() {
	const wl_subsurface_set_desync = 5
	obj.Conn().SendRequest(obj, wl_subsurface_set_desync)
}
