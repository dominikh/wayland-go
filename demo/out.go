package demo

import "honnef.co/go/wayland"

// These errors are global and can be emitted in response to any
// server request.
const (
	// server couldn't find object
	DisplayErrorInvalidObject = 0
	// method doesn't exist on the specified interface or malformed request
	DisplayErrorInvalidMethod = 1
	// server is out of memory
	DisplayErrorNoMemory = 2
	// implementation error in compositor
	DisplayErrorImplementation = 3
)

var displayInterface = &wayland.Interface{
	Name:    "wl_display",
	Version: 1,
	Events:  []wayland.Event{(*DisplayEventError)(nil), (*DisplayEventDeleteID)(nil)},
}

type DisplayEventError struct {
	// object where the error occurred
	ObjectID wayland.Object
	// error code
	Code uint32
	// error description
	Message string
}

type DisplayEventDeleteID struct {
	// deleted object ID
	ID uint32
}

// The core global object.  This is a special singleton object.  It
// is used for internal Wayland protocol features.
type Display struct{ wayland.Proxy }

func (*Display) Interface() *wayland.Interface { return displayInterface }

// The sync request asks the server to emit the 'done' event
// on the returned wl_callback object.  Since requests are
// handled in-order and events are delivered in-order, this can
// be used as a barrier to ensure all previous requests and the
// resulting events have been handled.
//
// The object returned by this request will be destroyed by the
// compositor after the callback is fired and as such the client must not
// attempt to use it after that point.
//
// The callback_data passed in the callback is the event serial.
func (obj *Display) Sync() *Callback {
	const wl_display_sync = 0
	_ret := &Callback{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_display_sync, _ret)
	return _ret
}

// This request creates a registry object that allows the client
// to list and bind the global objects available from the
// compositor.
//
// It should be noted that the server side resources consumed in
// response to a get_registry request can only be released when the
// client disconnects, not when the client side proxy is destroyed.
// Therefore, clients should invoke get_registry as infrequently as
// possible to avoid wasting memory.
func (obj *Display) GetRegistry() *Registry {
	const wl_display_get_registry = 1
	_ret := &Registry{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_display_get_registry, _ret)
	return _ret
}

var registryInterface = &wayland.Interface{
	Name:    "wl_registry",
	Version: 1,
	Events:  []wayland.Event{(*RegistryEventGlobal)(nil), (*RegistryEventGlobalRemove)(nil)},
}

type RegistryEventGlobal struct {
	// numeric name of the global object
	Name uint32
	// interface implemented by the object
	Interface string
	// interface version
	Version uint32
}

type RegistryEventGlobalRemove struct {
	// numeric name of the global object
	Name uint32
}

// The singleton global registry object.  The server has a number of
// global objects that are available to all clients.  These objects
// typically represent an actual object in the server (for example,
// an input device) or they are singleton objects that provide
// extension functionality.
//
// When a client creates a registry object, the registry object
// will emit a global event for each global currently in the
// registry.  Globals come and go as a result of device or
// monitor hotplugs, reconfiguration or other events, and the
// registry will send out global and global_remove events to
// keep the client up to date with the changes.  To mark the end
// of the initial burst of events, the client can use the
// wl_display.sync request immediately after calling
// wl_display.get_registry.
//
// A client can bind to a global object by using the bind
// request.  This creates a client-side handle that lets the object
// emit events to the client and lets the client invoke requests on
// the object.
type Registry struct{ wayland.Proxy }

func (*Registry) Interface() *wayland.Interface { return registryInterface }

// Binds a new, client-created object to the server using the
// specified name as the identifier.
func (obj *Registry) Bind(name uint32, id wayland.Object, version uint32) {
	const wl_registry_bind = 0
	obj.Conn().SendRequest(obj, wl_registry_bind, name, id.Interface().Name, version, id)
}

var callbackInterface = &wayland.Interface{
	Name:    "wl_callback",
	Version: 1,
	Events:  []wayland.Event{(*CallbackEventDone)(nil)},
}

type CallbackEventDone struct {
	// request-specific data for the callback
	CallbackData uint32
}

// Clients can handle the 'done' event to get notified when
// the related request is done.
type Callback struct{ wayland.Proxy }

func (*Callback) Interface() *wayland.Interface { return callbackInterface }

var compositorInterface = &wayland.Interface{
	Name:    "wl_compositor",
	Version: 4,
	Events:  []wayland.Event{},
}

// A compositor.  This object is a singleton global.  The
// compositor is in charge of combining the contents of multiple
// surfaces into one displayable output.
type Compositor struct{ wayland.Proxy }

func (*Compositor) Interface() *wayland.Interface { return compositorInterface }

// Ask the compositor to create a new surface.
func (obj *Compositor) CreateSurface() *Surface {
	const wl_compositor_create_surface = 0
	_ret := &Surface{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_compositor_create_surface, _ret)
	return _ret
}

// Ask the compositor to create a new region.
func (obj *Compositor) CreateRegion() *Region {
	const wl_compositor_create_region = 1
	_ret := &Region{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_compositor_create_region, _ret)
	return _ret
}

var shmPoolInterface = &wayland.Interface{
	Name:    "wl_shm_pool",
	Version: 1,
	Events:  []wayland.Event{},
}

// The wl_shm_pool object encapsulates a piece of memory shared
// between the compositor and client.  Through the wl_shm_pool
// object, the client can allocate shared memory wl_buffer objects.
// All objects created through the same pool share the same
// underlying mapped memory. Reusing the mapped memory avoids the
// setup/teardown overhead and is useful when interactively resizing
// a surface or for many small buffers.
type ShmPool struct{ wayland.Proxy }

func (*ShmPool) Interface() *wayland.Interface { return shmPoolInterface }

// Create a wl_buffer object from the pool.
//
// The buffer is created offset bytes into the pool and has
// width and height as specified.  The stride argument specifies
// the number of bytes from the beginning of one row to the beginning
// of the next.  The format is the pixel format of the buffer and
// must be one of those advertised through the wl_shm.format event.
//
// A buffer will keep a reference to the pool it was created from
// so it is valid to destroy the pool immediately after creating
// a buffer from it.
func (obj *ShmPool) CreateBuffer(offset int32, width int32, height int32, stride int32, format uint32) *Buffer {
	const wl_shm_pool_create_buffer = 0
	_ret := &Buffer{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_shm_pool_create_buffer, _ret, offset, width, height, stride, format)
	return _ret
}

// Destroy the shared memory pool.
//
// The mmapped memory will be released when all
// buffers that have been created from this pool
// are gone.
func (obj *ShmPool) Destroy() {
	const wl_shm_pool_destroy = 1
	obj.Conn().SendRequest(obj, wl_shm_pool_destroy)
}

// This request will cause the server to remap the backing memory
// for the pool from the file descriptor passed when the pool was
// created, but using the new size.  This request can only be
// used to make the pool bigger.
func (obj *ShmPool) Resize(size int32) {
	const wl_shm_pool_resize = 2
	obj.Conn().SendRequest(obj, wl_shm_pool_resize, size)
}

// These errors can be emitted in response to wl_shm requests.
const (
	// buffer format is not known
	ShmErrorInvalidFormat = 0
	// invalid size or stride during pool or buffer creation
	ShmErrorInvalidStride = 1
	// mmapping the file descriptor failed
	ShmErrorInvalidFd = 2
)

// This describes the memory layout of an individual pixel.
//
// All renderers should support argb8888 and xrgb8888 but any other
// formats are optional and may not be supported by the particular
// renderer in use.
//
// The drm format codes match the macros defined in drm_fourcc.h, except
// argb8888 and xrgb8888. The formats actually supported by the compositor
// will be reported by the format event.
const (
	// 32-bit ARGB format, [31:0] A:R:G:B 8:8:8:8 little endian
	ShmFormatArgb8888 = 0
	// 32-bit RGB format, [31:0] x:R:G:B 8:8:8:8 little endian
	ShmFormatXrgb8888 = 1
	// 8-bit color index format, [7:0] C
	ShmFormatC8 = 0x20203843
	// 8-bit RGB format, [7:0] R:G:B 3:3:2
	ShmFormatRgb332 = 0x38424752
	// 8-bit BGR format, [7:0] B:G:R 2:3:3
	ShmFormatBgr233 = 0x38524742
	// 16-bit xRGB format, [15:0] x:R:G:B 4:4:4:4 little endian
	ShmFormatXrgb4444 = 0x32315258
	// 16-bit xBGR format, [15:0] x:B:G:R 4:4:4:4 little endian
	ShmFormatXbgr4444 = 0x32314258
	// 16-bit RGBx format, [15:0] R:G:B:x 4:4:4:4 little endian
	ShmFormatRgbx4444 = 0x32315852
	// 16-bit BGRx format, [15:0] B:G:R:x 4:4:4:4 little endian
	ShmFormatBgrx4444 = 0x32315842
	// 16-bit ARGB format, [15:0] A:R:G:B 4:4:4:4 little endian
	ShmFormatArgb4444 = 0x32315241
	// 16-bit ABGR format, [15:0] A:B:G:R 4:4:4:4 little endian
	ShmFormatAbgr4444 = 0x32314241
	// 16-bit RBGA format, [15:0] R:G:B:A 4:4:4:4 little endian
	ShmFormatRgba4444 = 0x32314152
	// 16-bit BGRA format, [15:0] B:G:R:A 4:4:4:4 little endian
	ShmFormatBgra4444 = 0x32314142
	// 16-bit xRGB format, [15:0] x:R:G:B 1:5:5:5 little endian
	ShmFormatXrgb1555 = 0x35315258
	// 16-bit xBGR 1555 format, [15:0] x:B:G:R 1:5:5:5 little endian
	ShmFormatXbgr1555 = 0x35314258
	// 16-bit RGBx 5551 format, [15:0] R:G:B:x 5:5:5:1 little endian
	ShmFormatRgbx5551 = 0x35315852
	// 16-bit BGRx 5551 format, [15:0] B:G:R:x 5:5:5:1 little endian
	ShmFormatBgrx5551 = 0x35315842
	// 16-bit ARGB 1555 format, [15:0] A:R:G:B 1:5:5:5 little endian
	ShmFormatArgb1555 = 0x35315241
	// 16-bit ABGR 1555 format, [15:0] A:B:G:R 1:5:5:5 little endian
	ShmFormatAbgr1555 = 0x35314241
	// 16-bit RGBA 5551 format, [15:0] R:G:B:A 5:5:5:1 little endian
	ShmFormatRgba5551 = 0x35314152
	// 16-bit BGRA 5551 format, [15:0] B:G:R:A 5:5:5:1 little endian
	ShmFormatBgra5551 = 0x35314142
	// 16-bit RGB 565 format, [15:0] R:G:B 5:6:5 little endian
	ShmFormatRgb565 = 0x36314752
	// 16-bit BGR 565 format, [15:0] B:G:R 5:6:5 little endian
	ShmFormatBgr565 = 0x36314742
	// 24-bit RGB format, [23:0] R:G:B little endian
	ShmFormatRgb888 = 0x34324752
	// 24-bit BGR format, [23:0] B:G:R little endian
	ShmFormatBgr888 = 0x34324742
	// 32-bit xBGR format, [31:0] x:B:G:R 8:8:8:8 little endian
	ShmFormatXbgr8888 = 0x34324258
	// 32-bit RGBx format, [31:0] R:G:B:x 8:8:8:8 little endian
	ShmFormatRgbx8888 = 0x34325852
	// 32-bit BGRx format, [31:0] B:G:R:x 8:8:8:8 little endian
	ShmFormatBgrx8888 = 0x34325842
	// 32-bit ABGR format, [31:0] A:B:G:R 8:8:8:8 little endian
	ShmFormatAbgr8888 = 0x34324241
	// 32-bit RGBA format, [31:0] R:G:B:A 8:8:8:8 little endian
	ShmFormatRgba8888 = 0x34324152
	// 32-bit BGRA format, [31:0] B:G:R:A 8:8:8:8 little endian
	ShmFormatBgra8888 = 0x34324142
	// 32-bit xRGB format, [31:0] x:R:G:B 2:10:10:10 little endian
	ShmFormatXrgb2101010 = 0x30335258
	// 32-bit xBGR format, [31:0] x:B:G:R 2:10:10:10 little endian
	ShmFormatXbgr2101010 = 0x30334258
	// 32-bit RGBx format, [31:0] R:G:B:x 10:10:10:2 little endian
	ShmFormatRgbx1010102 = 0x30335852
	// 32-bit BGRx format, [31:0] B:G:R:x 10:10:10:2 little endian
	ShmFormatBgrx1010102 = 0x30335842
	// 32-bit ARGB format, [31:0] A:R:G:B 2:10:10:10 little endian
	ShmFormatArgb2101010 = 0x30335241
	// 32-bit ABGR format, [31:0] A:B:G:R 2:10:10:10 little endian
	ShmFormatAbgr2101010 = 0x30334241
	// 32-bit RGBA format, [31:0] R:G:B:A 10:10:10:2 little endian
	ShmFormatRgba1010102 = 0x30334152
	// 32-bit BGRA format, [31:0] B:G:R:A 10:10:10:2 little endian
	ShmFormatBgra1010102 = 0x30334142
	// packed YCbCr format, [31:0] Cr0:Y1:Cb0:Y0 8:8:8:8 little endian
	ShmFormatYuyv = 0x56595559
	// packed YCbCr format, [31:0] Cb0:Y1:Cr0:Y0 8:8:8:8 little endian
	ShmFormatYvyu = 0x55595659
	// packed YCbCr format, [31:0] Y1:Cr0:Y0:Cb0 8:8:8:8 little endian
	ShmFormatUyvy = 0x59565955
	// packed YCbCr format, [31:0] Y1:Cb0:Y0:Cr0 8:8:8:8 little endian
	ShmFormatVyuy = 0x59555956
	// packed AYCbCr format, [31:0] A:Y:Cb:Cr 8:8:8:8 little endian
	ShmFormatAyuv = 0x56555941
	// 2 plane YCbCr Cr:Cb format, 2x2 subsampled Cr:Cb plane
	ShmFormatNv12 = 0x3231564e
	// 2 plane YCbCr Cb:Cr format, 2x2 subsampled Cb:Cr plane
	ShmFormatNv21 = 0x3132564e
	// 2 plane YCbCr Cr:Cb format, 2x1 subsampled Cr:Cb plane
	ShmFormatNv16 = 0x3631564e
	// 2 plane YCbCr Cb:Cr format, 2x1 subsampled Cb:Cr plane
	ShmFormatNv61 = 0x3136564e
	// 3 plane YCbCr format, 4x4 subsampled Cb (1) and Cr (2) planes
	ShmFormatYuv410 = 0x39565559
	// 3 plane YCbCr format, 4x4 subsampled Cr (1) and Cb (2) planes
	ShmFormatYvu410 = 0x39555659
	// 3 plane YCbCr format, 4x1 subsampled Cb (1) and Cr (2) planes
	ShmFormatYuv411 = 0x31315559
	// 3 plane YCbCr format, 4x1 subsampled Cr (1) and Cb (2) planes
	ShmFormatYvu411 = 0x31315659
	// 3 plane YCbCr format, 2x2 subsampled Cb (1) and Cr (2) planes
	ShmFormatYuv420 = 0x32315559
	// 3 plane YCbCr format, 2x2 subsampled Cr (1) and Cb (2) planes
	ShmFormatYvu420 = 0x32315659
	// 3 plane YCbCr format, 2x1 subsampled Cb (1) and Cr (2) planes
	ShmFormatYuv422 = 0x36315559
	// 3 plane YCbCr format, 2x1 subsampled Cr (1) and Cb (2) planes
	ShmFormatYvu422 = 0x36315659
	// 3 plane YCbCr format, non-subsampled Cb (1) and Cr (2) planes
	ShmFormatYuv444 = 0x34325559
	// 3 plane YCbCr format, non-subsampled Cr (1) and Cb (2) planes
	ShmFormatYvu444 = 0x34325659
	// [7:0] R
	ShmFormatR8 = 0x20203852
	// [15:0] R little endian
	ShmFormatR16 = 0x20363152
	// [15:0] R:G 8:8 little endian
	ShmFormatRg88 = 0x38384752
	// [15:0] G:R 8:8 little endian
	ShmFormatGr88 = 0x38385247
	// [31:0] R:G 16:16 little endian
	ShmFormatRg1616 = 0x32334752
	// [31:0] G:R 16:16 little endian
	ShmFormatGr1616 = 0x32335247
	// [63:0] x:R:G:B 16:16:16:16 little endian
	ShmFormatXrgb16161616f = 0x48345258
	// [63:0] x:B:G:R 16:16:16:16 little endian
	ShmFormatXbgr16161616f = 0x48344258
	// [63:0] A:R:G:B 16:16:16:16 little endian
	ShmFormatArgb16161616f = 0x48345241
	// [63:0] A:B:G:R 16:16:16:16 little endian
	ShmFormatAbgr16161616f = 0x48344241
	// [31:0] X:Y:Cb:Cr 8:8:8:8 little endian
	ShmFormatXyuv8888 = 0x56555958
	// [23:0] Cr:Cb:Y 8:8:8 little endian
	ShmFormatVuy888 = 0x34325556
	// Y followed by U then V, 10:10:10. Non-linear modifier only
	ShmFormatVuy101010 = 0x30335556
	// [63:0] Cr0:0:Y1:0:Cb0:0:Y0:0 10:6:10:6:10:6:10:6 little endian per 2 Y pixels
	ShmFormatY210 = 0x30313259
	// [63:0] Cr0:0:Y1:0:Cb0:0:Y0:0 12:4:12:4:12:4:12:4 little endian per 2 Y pixels
	ShmFormatY212 = 0x32313259
	// [63:0] Cr0:Y1:Cb0:Y0 16:16:16:16 little endian per 2 Y pixels
	ShmFormatY216 = 0x36313259
	// [31:0] A:Cr:Y:Cb 2:10:10:10 little endian
	ShmFormatY410 = 0x30313459
	// [63:0] A:0:Cr:0:Y:0:Cb:0 12:4:12:4:12:4:12:4 little endian
	ShmFormatY412 = 0x32313459
	// [63:0] A:Cr:Y:Cb 16:16:16:16 little endian
	ShmFormatY416 = 0x36313459
	// [31:0] X:Cr:Y:Cb 2:10:10:10 little endian
	ShmFormatXvyu2101010 = 0x30335658
	// [63:0] X:0:Cr:0:Y:0:Cb:0 12:4:12:4:12:4:12:4 little endian
	ShmFormatXvyu1216161616 = 0x36335658
	// [63:0] X:Cr:Y:Cb 16:16:16:16 little endian
	ShmFormatXvyu16161616 = 0x38345658
	// [63:0]   A3:A2:Y3:0:Cr0:0:Y2:0:A1:A0:Y1:0:Cb0:0:Y0:0  1:1:8:2:8:2:8:2:1:1:8:2:8:2:8:2 little endian
	ShmFormatY0l0 = 0x304c3059
	// [63:0]   X3:X2:Y3:0:Cr0:0:Y2:0:X1:X0:Y1:0:Cb0:0:Y0:0  1:1:8:2:8:2:8:2:1:1:8:2:8:2:8:2 little endian
	ShmFormatX0l0 = 0x304c3058
	// [63:0]   A3:A2:Y3:Cr0:Y2:A1:A0:Y1:Cb0:Y0  1:1:10:10:10:1:1:10:10:10 little endian
	ShmFormatY0l2 = 0x324c3059
	// [63:0]   X3:X2:Y3:Cr0:Y2:X1:X0:Y1:Cb0:Y0  1:1:10:10:10:1:1:10:10:10 little endian
	ShmFormatX0l2        = 0x324c3058
	ShmFormatYuv4208bit  = 0x38305559
	ShmFormatYuv42010bit = 0x30315559
	ShmFormatXrgb8888A8  = 0x38415258
	ShmFormatXbgr8888A8  = 0x38414258
	ShmFormatRgbx8888A8  = 0x38415852
	ShmFormatBgrx8888A8  = 0x38415842
	ShmFormatRgb888A8    = 0x38413852
	ShmFormatBgr888A8    = 0x38413842
	ShmFormatRgb565A8    = 0x38413552
	ShmFormatBgr565A8    = 0x38413542
	// non-subsampled Cr:Cb plane
	ShmFormatNv24 = 0x3432564e
	// non-subsampled Cb:Cr plane
	ShmFormatNv42 = 0x3234564e
	// 2x1 subsampled Cr:Cb plane, 10 bit per channel
	ShmFormatP210 = 0x30313250
	// 2x2 subsampled Cr:Cb plane 10 bits per channel
	ShmFormatP010 = 0x30313050
	// 2x2 subsampled Cr:Cb plane 12 bits per channel
	ShmFormatP012 = 0x32313050
	// 2x2 subsampled Cr:Cb plane 16 bits per channel
	ShmFormatP016 = 0x36313050
)

var shmInterface = &wayland.Interface{
	Name:    "wl_shm",
	Version: 1,
	Events:  []wayland.Event{(*ShmEventFormat)(nil)},
}

type ShmEventFormat struct {
	// buffer pixel format
	Format uint32
}

// A singleton global object that provides support for shared
// memory.
//
// Clients can create wl_shm_pool objects using the create_pool
// request.
//
// At connection setup time, the wl_shm object emits one or more
// format events to inform clients about the valid pixel formats
// that can be used for buffers.
type Shm struct{ wayland.Proxy }

func (*Shm) Interface() *wayland.Interface { return shmInterface }

// Create a new wl_shm_pool object.
//
// The pool can be used to create shared memory based buffer
// objects.  The server will mmap size bytes of the passed file
// descriptor, to use as backing memory for the pool.
func (obj *Shm) CreatePool(fd uintptr, size int32) *ShmPool {
	const wl_shm_create_pool = 0
	_ret := &ShmPool{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_shm_create_pool, _ret, fd, size)
	return _ret
}

var bufferInterface = &wayland.Interface{
	Name:    "wl_buffer",
	Version: 1,
	Events:  []wayland.Event{(*BufferEventRelease)(nil)},
}

type BufferEventRelease struct {
}

// A buffer provides the content for a wl_surface. Buffers are
// created through factory interfaces such as wl_drm, wl_shm or
// similar. It has a width and a height and can be attached to a
// wl_surface, but the mechanism by which a client provides and
// updates the contents is defined by the buffer factory interface.
type Buffer struct{ wayland.Proxy }

func (*Buffer) Interface() *wayland.Interface { return bufferInterface }

// Destroy a buffer. If and how you need to release the backing
// storage is defined by the buffer factory interface.
//
// For possible side-effects to a surface, see wl_surface.attach.
func (obj *Buffer) Destroy() {
	const wl_buffer_destroy = 0
	obj.Conn().SendRequest(obj, wl_buffer_destroy)
}

const (
	// finish request was called untimely
	DataOfferErrorInvalidFinish = 0
	// action mask contains invalid values
	DataOfferErrorInvalidActionMask = 1
	// action argument has an invalid value
	DataOfferErrorInvalidAction = 2
	// offer doesn't accept this request
	DataOfferErrorInvalidOffer = 3
)

var dataOfferInterface = &wayland.Interface{
	Name:    "wl_data_offer",
	Version: 3,
	Events:  []wayland.Event{(*DataOfferEventOffer)(nil), (*DataOfferEventSourceActions)(nil), (*DataOfferEventAction)(nil)},
}

type DataOfferEventOffer struct {
	// offered mime type
	MimeType string
}

type DataOfferEventSourceActions struct {
	// actions offered by the data source
	SourceActions uint32
}

type DataOfferEventAction struct {
	// action selected by the compositor
	DndAction uint32
}

// A wl_data_offer represents a piece of data offered for transfer
// by another client (the source client).  It is used by the
// copy-and-paste and drag-and-drop mechanisms.  The offer
// describes the different mime types that the data can be
// converted to and provides the mechanism for transferring the
// data directly from the source client.
type DataOffer struct{ wayland.Proxy }

func (*DataOffer) Interface() *wayland.Interface { return dataOfferInterface }

// Indicate that the client can accept the given mime type, or
// NULL for not accepted.
//
// For objects of version 2 or older, this request is used by the
// client to give feedback whether the client can receive the given
// mime type, or NULL if none is accepted; the feedback does not
// determine whether the drag-and-drop operation succeeds or not.
//
// For objects of version 3 or newer, this request determines the
// final result of the drag-and-drop operation. If the end result
// is that no mime types were accepted, the drag-and-drop operation
// will be cancelled and the corresponding drag source will receive
// wl_data_source.cancelled. Clients may still use this event in
// conjunction with wl_data_source.action for feedback.
func (obj *DataOffer) Accept(serial uint32, mimeType string) {
	const wl_data_offer_accept = 0
	obj.Conn().SendRequest(obj, wl_data_offer_accept, serial, mimeType)
}

// To transfer the offered data, the client issues this request
// and indicates the mime type it wants to receive.  The transfer
// happens through the passed file descriptor (typically created
// with the pipe system call).  The source client writes the data
// in the mime type representation requested and then closes the
// file descriptor.
//
// The receiving client reads from the read end of the pipe until
// EOF and then closes its end, at which point the transfer is
// complete.
//
// This request may happen multiple times for different mime types,
// both before and after wl_data_device.drop. Drag-and-drop destination
// clients may preemptively fetch data or examine it more closely to
// determine acceptance.
func (obj *DataOffer) Receive(mimeType string, fd uintptr) {
	const wl_data_offer_receive = 1
	obj.Conn().SendRequest(obj, wl_data_offer_receive, mimeType, fd)
}

// Destroy the data offer.
func (obj *DataOffer) Destroy() {
	const wl_data_offer_destroy = 2
	obj.Conn().SendRequest(obj, wl_data_offer_destroy)
}

// Notifies the compositor that the drag destination successfully
// finished the drag-and-drop operation.
//
// Upon receiving this request, the compositor will emit
// wl_data_source.dnd_finished on the drag source client.
//
// It is a client error to perform other requests than
// wl_data_offer.destroy after this one. It is also an error to perform
// this request after a NULL mime type has been set in
// wl_data_offer.accept or no action was received through
// wl_data_offer.action.
//
// If wl_data_offer.finish request is received for a non drag and drop
// operation, the invalid_finish protocol error is raised.
func (obj *DataOffer) Finish() {
	const wl_data_offer_finish = 3
	obj.Conn().SendRequest(obj, wl_data_offer_finish)
}

// Sets the actions that the destination side client supports for
// this operation. This request may trigger the emission of
// wl_data_source.action and wl_data_offer.action events if the compositor
// needs to change the selected action.
//
// This request can be called multiple times throughout the
// drag-and-drop operation, typically in response to wl_data_device.enter
// or wl_data_device.motion events.
//
// This request determines the final result of the drag-and-drop
// operation. If the end result is that no action is accepted,
// the drag source will receive wl_drag_source.cancelled.
//
// The dnd_actions argument must contain only values expressed in the
// wl_data_device_manager.dnd_actions enum, and the preferred_action
// argument must only contain one of those values set, otherwise it
// will result in a protocol error.
//
// While managing an "ask" action, the destination drag-and-drop client
// may perform further wl_data_offer.receive requests, and is expected
// to perform one last wl_data_offer.set_actions request with a preferred
// action other than "ask" (and optionally wl_data_offer.accept) before
// requesting wl_data_offer.finish, in order to convey the action selected
// by the user. If the preferred action is not in the
// wl_data_offer.source_actions mask, an error will be raised.
//
// If the "ask" action is dismissed (e.g. user cancellation), the client
// is expected to perform wl_data_offer.destroy right away.
//
// This request can only be made on drag-and-drop offers, a protocol error
// will be raised otherwise.
func (obj *DataOffer) SetActions(dndActions uint32, preferredAction uint32) {
	const wl_data_offer_set_actions = 4
	obj.Conn().SendRequest(obj, wl_data_offer_set_actions, dndActions, preferredAction)
}

const (
	// action mask contains invalid values
	DataSourceErrorInvalidActionMask = 0
	// source doesn't accept this request
	DataSourceErrorInvalidSource = 1
)

var dataSourceInterface = &wayland.Interface{
	Name:    "wl_data_source",
	Version: 3,
	Events:  []wayland.Event{(*DataSourceEventTarget)(nil), (*DataSourceEventSend)(nil), (*DataSourceEventCancelled)(nil), (*DataSourceEventDndDropPerformed)(nil), (*DataSourceEventDndFinished)(nil), (*DataSourceEventAction)(nil)},
}

type DataSourceEventTarget struct {
	// mime type accepted by the target
	MimeType string
}

type DataSourceEventSend struct {
	// mime type for the data
	MimeType string
	// file descriptor for the data
	Fd uintptr
}

type DataSourceEventCancelled struct {
}

type DataSourceEventDndDropPerformed struct {
}

type DataSourceEventDndFinished struct {
}

type DataSourceEventAction struct {
	// action selected by the compositor
	DndAction uint32
}

// The wl_data_source object is the source side of a wl_data_offer.
// It is created by the source client in a data transfer and
// provides a way to describe the offered data and a way to respond
// to requests to transfer the data.
type DataSource struct{ wayland.Proxy }

func (*DataSource) Interface() *wayland.Interface { return dataSourceInterface }

// This request adds a mime type to the set of mime types
// advertised to targets.  Can be called several times to offer
// multiple types.
func (obj *DataSource) Offer(mimeType string) {
	const wl_data_source_offer = 0
	obj.Conn().SendRequest(obj, wl_data_source_offer, mimeType)
}

// Destroy the data source.
func (obj *DataSource) Destroy() {
	const wl_data_source_destroy = 1
	obj.Conn().SendRequest(obj, wl_data_source_destroy)
}

// Sets the actions that the source side client supports for this
// operation. This request may trigger wl_data_source.action and
// wl_data_offer.action events if the compositor needs to change the
// selected action.
//
// The dnd_actions argument must contain only values expressed in the
// wl_data_device_manager.dnd_actions enum, otherwise it will result
// in a protocol error.
//
// This request must be made once only, and can only be made on sources
// used in drag-and-drop, so it must be performed before
// wl_data_device.start_drag. Attempting to use the source other than
// for drag-and-drop will raise a protocol error.
func (obj *DataSource) SetActions(dndActions uint32) {
	const wl_data_source_set_actions = 2
	obj.Conn().SendRequest(obj, wl_data_source_set_actions, dndActions)
}

const (
	// given wl_surface has another role
	DataDeviceErrorRole = 0
)

var dataDeviceInterface = &wayland.Interface{
	Name:    "wl_data_device",
	Version: 3,
	Events:  []wayland.Event{(*DataDeviceEventDataOffer)(nil), (*DataDeviceEventEnter)(nil), (*DataDeviceEventLeave)(nil), (*DataDeviceEventMotion)(nil), (*DataDeviceEventDrop)(nil), (*DataDeviceEventSelection)(nil)},
}

type DataDeviceEventDataOffer struct {
	// the new data_offer object
	ID *DataOffer `wl:"new_id"`
}

type DataDeviceEventEnter struct {
	// serial number of the enter event
	Serial uint32
	// client surface entered
	Surface *Surface
	// surface-local x coordinate
	X wayland.Fixed
	// surface-local y coordinate
	Y wayland.Fixed
	// source data_offer object
	ID *DataOffer
}

type DataDeviceEventLeave struct {
}

type DataDeviceEventMotion struct {
	// timestamp with millisecond granularity
	Time uint32
	// surface-local x coordinate
	X wayland.Fixed
	// surface-local y coordinate
	Y wayland.Fixed
}

type DataDeviceEventDrop struct {
}

type DataDeviceEventSelection struct {
	// selection data_offer object
	ID *DataOffer
}

// There is one wl_data_device per seat which can be obtained
// from the global wl_data_device_manager singleton.
//
// A wl_data_device provides access to inter-client data transfer
// mechanisms such as copy-and-paste and drag-and-drop.
type DataDevice struct{ wayland.Proxy }

func (*DataDevice) Interface() *wayland.Interface { return dataDeviceInterface }

// This request asks the compositor to start a drag-and-drop
// operation on behalf of the client.
//
// The source argument is the data source that provides the data
// for the eventual data transfer. If source is NULL, enter, leave
// and motion events are sent only to the client that initiated the
// drag and the client is expected to handle the data passing
// internally.
//
// The origin surface is the surface where the drag originates and
// the client must have an active implicit grab that matches the
// serial.
//
// The icon surface is an optional (can be NULL) surface that
// provides an icon to be moved around with the cursor.  Initially,
// the top-left corner of the icon surface is placed at the cursor
// hotspot, but subsequent wl_surface.attach request can move the
// relative position. Attach requests must be confirmed with
// wl_surface.commit as usual. The icon surface is given the role of
// a drag-and-drop icon. If the icon surface already has another role,
// it raises a protocol error.
//
// The current and pending input regions of the icon wl_surface are
// cleared, and wl_surface.set_input_region is ignored until the
// wl_surface is no longer used as the icon surface. When the use
// as an icon ends, the current and pending input regions become
// undefined, and the wl_surface is unmapped.
func (obj *DataDevice) StartDrag(source *DataSource, origin *Surface, icon *Surface, serial uint32) {
	const wl_data_device_start_drag = 0
	obj.Conn().SendRequest(obj, wl_data_device_start_drag, source, origin, icon, serial)
}

// This request asks the compositor to set the selection
// to the data from the source on behalf of the client.
//
// To unset the selection, set the source to NULL.
func (obj *DataDevice) SetSelection(source *DataSource, serial uint32) {
	const wl_data_device_set_selection = 1
	obj.Conn().SendRequest(obj, wl_data_device_set_selection, source, serial)
}

// This request destroys the data device.
func (obj *DataDevice) Release() {
	const wl_data_device_release = 2
	obj.Conn().SendRequest(obj, wl_data_device_release)
}

// This is a bitmask of the available/preferred actions in a
// drag-and-drop operation.
//
// In the compositor, the selected action is a result of matching the
// actions offered by the source and destination sides.  "action" events
// with a "none" action will be sent to both source and destination if
// there is no match. All further checks will effectively happen on
// (source actions ∩ destination actions).
//
// In addition, compositors may also pick different actions in
// reaction to key modifiers being pressed. One common design that
// is used in major toolkits (and the behavior recommended for
// compositors) is:
//
// - If no modifiers are pressed, the first match (in bit order)
// will be used.
// - Pressing Shift selects "move", if enabled in the mask.
// - Pressing Control selects "copy", if enabled in the mask.
//
// Behavior beyond that is considered implementation-dependent.
// Compositors may for example bind other modifiers (like Alt/Meta)
// or drags initiated with other buttons than BTN_LEFT to specific
// actions (e.g. "ask").
const (
	// no action
	DataDeviceManagerDndActionNone = 0
	// copy action
	DataDeviceManagerDndActionCopy = 1
	// move action
	DataDeviceManagerDndActionMove = 2
	// ask action
	DataDeviceManagerDndActionAsk = 4
)

var dataDeviceManagerInterface = &wayland.Interface{
	Name:    "wl_data_device_manager",
	Version: 3,
	Events:  []wayland.Event{},
}

// The wl_data_device_manager is a singleton global object that
// provides access to inter-client data transfer mechanisms such as
// copy-and-paste and drag-and-drop.  These mechanisms are tied to
// a wl_seat and this interface lets a client get a wl_data_device
// corresponding to a wl_seat.
//
// Depending on the version bound, the objects created from the bound
// wl_data_device_manager object will have different requirements for
// functioning properly. See wl_data_source.set_actions,
// wl_data_offer.accept and wl_data_offer.finish for details.
type DataDeviceManager struct{ wayland.Proxy }

func (*DataDeviceManager) Interface() *wayland.Interface { return dataDeviceManagerInterface }

// Create a new data source.
func (obj *DataDeviceManager) CreateDataSource() *DataSource {
	const wl_data_device_manager_create_data_source = 0
	_ret := &DataSource{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_data_device_manager_create_data_source, _ret)
	return _ret
}

// Create a new data device for a given seat.
func (obj *DataDeviceManager) GetDataDevice(seat *Seat) *DataDevice {
	const wl_data_device_manager_get_data_device = 1
	_ret := &DataDevice{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_data_device_manager_get_data_device, _ret, seat)
	return _ret
}

const (
	// given wl_surface has another role
	ShellErrorRole = 0
)

var shellInterface = &wayland.Interface{
	Name:    "wl_shell",
	Version: 1,
	Events:  []wayland.Event{},
}

// This interface is implemented by servers that provide
// desktop-style user interfaces.
//
// It allows clients to associate a wl_shell_surface with
// a basic surface.
//
// Note! This protocol is deprecated and not intended for production use.
// For desktop-style user interfaces, use xdg_shell.
type Shell struct{ wayland.Proxy }

func (*Shell) Interface() *wayland.Interface { return shellInterface }

// Create a shell surface for an existing surface. This gives
// the wl_surface the role of a shell surface. If the wl_surface
// already has another role, it raises a protocol error.
//
// Only one shell surface can be associated with a given surface.
func (obj *Shell) GetShellSurface(surface *Surface) *ShellSurface {
	const wl_shell_get_shell_surface = 0
	_ret := &ShellSurface{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_shell_get_shell_surface, _ret, surface)
	return _ret
}

// These values are used to indicate which edge of a surface
// is being dragged in a resize operation. The server may
// use this information to adapt its behavior, e.g. choose
// an appropriate cursor image.
const (
	// no edge
	ShellSurfaceResizeNone = 0
	// top edge
	ShellSurfaceResizeTop = 1
	// bottom edge
	ShellSurfaceResizeBottom = 2
	// left edge
	ShellSurfaceResizeLeft = 4
	// top and left edges
	ShellSurfaceResizeTopLeft = 5
	// bottom and left edges
	ShellSurfaceResizeBottomLeft = 6
	// right edge
	ShellSurfaceResizeRight = 8
	// top and right edges
	ShellSurfaceResizeTopRight = 9
	// bottom and right edges
	ShellSurfaceResizeBottomRight = 10
)

// These flags specify details of the expected behaviour
// of transient surfaces. Used in the set_transient request.
const (
	// do not set keyboard focus
	ShellSurfaceTransientInactive = 0x1
)

// Hints to indicate to the compositor how to deal with a conflict
// between the dimensions of the surface and the dimensions of the
// output. The compositor is free to ignore this parameter.
const (
	// no preference, apply default policy
	ShellSurfaceFullscreenMethodDefault = 0
	// scale, preserve the surface's aspect ratio and center on output
	ShellSurfaceFullscreenMethodScale = 1
	// switch output mode to the smallest mode that can fit the surface, add black borders to compensate size mismatch
	ShellSurfaceFullscreenMethodDriver = 2
	// no upscaling, center on output and add black borders to compensate size mismatch
	ShellSurfaceFullscreenMethodFill = 3
)

var shellSurfaceInterface = &wayland.Interface{
	Name:    "wl_shell_surface",
	Version: 1,
	Events:  []wayland.Event{(*ShellSurfaceEventPing)(nil), (*ShellSurfaceEventConfigure)(nil), (*ShellSurfaceEventPopupDone)(nil)},
}

type ShellSurfaceEventPing struct {
	// serial number of the ping
	Serial uint32
}

type ShellSurfaceEventConfigure struct {
	// how the surface was resized
	Edges uint32
	// new width of the surface
	Width int32
	// new height of the surface
	Height int32
}

type ShellSurfaceEventPopupDone struct {
}

// An interface that may be implemented by a wl_surface, for
// implementations that provide a desktop-style user interface.
//
// It provides requests to treat surfaces like toplevel, fullscreen
// or popup windows, move, resize or maximize them, associate
// metadata like title and class, etc.
//
// On the server side the object is automatically destroyed when
// the related wl_surface is destroyed. On the client side,
// wl_shell_surface_destroy() must be called before destroying
// the wl_surface object.
type ShellSurface struct{ wayland.Proxy }

func (*ShellSurface) Interface() *wayland.Interface { return shellSurfaceInterface }

// A client must respond to a ping event with a pong request or
// the client may be deemed unresponsive.
func (obj *ShellSurface) Pong(serial uint32) {
	const wl_shell_surface_pong = 0
	obj.Conn().SendRequest(obj, wl_shell_surface_pong, serial)
}

// Start a pointer-driven move of the surface.
//
// This request must be used in response to a button press event.
// The server may ignore move requests depending on the state of
// the surface (e.g. fullscreen or maximized).
func (obj *ShellSurface) Move(seat *Seat, serial uint32) {
	const wl_shell_surface_move = 1
	obj.Conn().SendRequest(obj, wl_shell_surface_move, seat, serial)
}

// Start a pointer-driven resizing of the surface.
//
// This request must be used in response to a button press event.
// The server may ignore resize requests depending on the state of
// the surface (e.g. fullscreen or maximized).
func (obj *ShellSurface) Resize(seat *Seat, serial uint32, edges uint32) {
	const wl_shell_surface_resize = 2
	obj.Conn().SendRequest(obj, wl_shell_surface_resize, seat, serial, edges)
}

// Map the surface as a toplevel surface.
//
// A toplevel surface is not fullscreen, maximized or transient.
func (obj *ShellSurface) SetToplevel() {
	const wl_shell_surface_set_toplevel = 3
	obj.Conn().SendRequest(obj, wl_shell_surface_set_toplevel)
}

// Map the surface relative to an existing surface.
//
// The x and y arguments specify the location of the upper left
// corner of the surface relative to the upper left corner of the
// parent surface, in surface-local coordinates.
//
// The flags argument controls details of the transient behaviour.
func (obj *ShellSurface) SetTransient(parent *Surface, x int32, y int32, flags uint32) {
	const wl_shell_surface_set_transient = 4
	obj.Conn().SendRequest(obj, wl_shell_surface_set_transient, parent, x, y, flags)
}

// Map the surface as a fullscreen surface.
//
// If an output parameter is given then the surface will be made
// fullscreen on that output. If the client does not specify the
// output then the compositor will apply its policy - usually
// choosing the output on which the surface has the biggest surface
// area.
//
// The client may specify a method to resolve a size conflict
// between the output size and the surface size - this is provided
// through the method parameter.
//
// The framerate parameter is used only when the method is set
// to "driver", to indicate the preferred framerate. A value of 0
// indicates that the client does not care about framerate.  The
// framerate is specified in mHz, that is framerate of 60000 is 60Hz.
//
// A method of "scale" or "driver" implies a scaling operation of
// the surface, either via a direct scaling operation or a change of
// the output mode. This will override any kind of output scaling, so
// that mapping a surface with a buffer size equal to the mode can
// fill the screen independent of buffer_scale.
//
// A method of "fill" means we don't scale up the buffer, however
// any output scale is applied. This means that you may run into
// an edge case where the application maps a buffer with the same
// size of the output mode but buffer_scale 1 (thus making a
// surface larger than the output). In this case it is allowed to
// downscale the results to fit the screen.
//
// The compositor must reply to this request with a configure event
// with the dimensions for the output on which the surface will
// be made fullscreen.
func (obj *ShellSurface) SetFullscreen(method uint32, framerate uint32, output *Output) {
	const wl_shell_surface_set_fullscreen = 5
	obj.Conn().SendRequest(obj, wl_shell_surface_set_fullscreen, method, framerate, output)
}

// Map the surface as a popup.
//
// A popup surface is a transient surface with an added pointer
// grab.
//
// An existing implicit grab will be changed to owner-events mode,
// and the popup grab will continue after the implicit grab ends
// (i.e. releasing the mouse button does not cause the popup to
// be unmapped).
//
// The popup grab continues until the window is destroyed or a
// mouse button is pressed in any other client's window. A click
// in any of the client's surfaces is reported as normal, however,
// clicks in other clients' surfaces will be discarded and trigger
// the callback.
//
// The x and y arguments specify the location of the upper left
// corner of the surface relative to the upper left corner of the
// parent surface, in surface-local coordinates.
func (obj *ShellSurface) SetPopup(seat *Seat, serial uint32, parent *Surface, x int32, y int32, flags uint32) {
	const wl_shell_surface_set_popup = 6
	obj.Conn().SendRequest(obj, wl_shell_surface_set_popup, seat, serial, parent, x, y, flags)
}

// Map the surface as a maximized surface.
//
// If an output parameter is given then the surface will be
// maximized on that output. If the client does not specify the
// output then the compositor will apply its policy - usually
// choosing the output on which the surface has the biggest surface
// area.
//
// The compositor will reply with a configure event telling
// the expected new surface size. The operation is completed
// on the next buffer attach to this surface.
//
// A maximized surface typically fills the entire output it is
// bound to, except for desktop elements such as panels. This is
// the main difference between a maximized shell surface and a
// fullscreen shell surface.
//
// The details depend on the compositor implementation.
func (obj *ShellSurface) SetMaximized(output *Output) {
	const wl_shell_surface_set_maximized = 7
	obj.Conn().SendRequest(obj, wl_shell_surface_set_maximized, output)
}

// Set a short title for the surface.
//
// This string may be used to identify the surface in a task bar,
// window list, or other user interface elements provided by the
// compositor.
//
// The string must be encoded in UTF-8.
func (obj *ShellSurface) SetTitle(title string) {
	const wl_shell_surface_set_title = 8
	obj.Conn().SendRequest(obj, wl_shell_surface_set_title, title)
}

// Set a class for the surface.
//
// The surface class identifies the general class of applications
// to which the surface belongs. A common convention is to use the
// file name (or the full path if it is a non-standard location) of
// the application's .desktop file as the class.
func (obj *ShellSurface) SetClass(class string) {
	const wl_shell_surface_set_class = 9
	obj.Conn().SendRequest(obj, wl_shell_surface_set_class, class)
}

// These errors can be emitted in response to wl_surface requests.
const (
	// buffer scale value is invalid
	SurfaceErrorInvalidScale = 0
	// buffer transform value is invalid
	SurfaceErrorInvalidTransform = 1
)

var surfaceInterface = &wayland.Interface{
	Name:    "wl_surface",
	Version: 4,
	Events:  []wayland.Event{(*SurfaceEventEnter)(nil), (*SurfaceEventLeave)(nil)},
}

type SurfaceEventEnter struct {
	// output entered by the surface
	Output *Output
}

type SurfaceEventLeave struct {
	// output left by the surface
	Output *Output
}

// A surface is a rectangular area that may be displayed on zero
// or more outputs, and shown any number of times at the compositor's
// discretion. They can present wl_buffers, receive user input, and
// define a local coordinate system.
//
// The size of a surface (and relative positions on it) is described
// in surface-local coordinates, which may differ from the buffer
// coordinates of the pixel content, in case a buffer_transform
// or a buffer_scale is used.
//
// A surface without a "role" is fairly useless: a compositor does
// not know where, when or how to present it. The role is the
// purpose of a wl_surface. Examples of roles are a cursor for a
// pointer (as set by wl_pointer.set_cursor), a drag icon
// (wl_data_device.start_drag), a sub-surface
// (wl_subcompositor.get_subsurface), and a window as defined by a
// shell protocol (e.g. wl_shell.get_shell_surface).
//
// A surface can have only one role at a time. Initially a
// wl_surface does not have a role. Once a wl_surface is given a
// role, it is set permanently for the whole lifetime of the
// wl_surface object. Giving the current role again is allowed,
// unless explicitly forbidden by the relevant interface
// specification.
//
// Surface roles are given by requests in other interfaces such as
// wl_pointer.set_cursor. The request should explicitly mention
// that this request gives a role to a wl_surface. Often, this
// request also creates a new protocol object that represents the
// role and adds additional functionality to wl_surface. When a
// client wants to destroy a wl_surface, they must destroy this 'role
// object' before the wl_surface.
//
// Destroying the role object does not remove the role from the
// wl_surface, but it may stop the wl_surface from "playing the role".
// For instance, if a wl_subsurface object is destroyed, the wl_surface
// it was created for will be unmapped and forget its position and
// z-order. It is allowed to create a wl_subsurface for the same
// wl_surface again, but it is not allowed to use the wl_surface as
// a cursor (cursor is a different role than sub-surface, and role
// switching is not allowed).
type Surface struct{ wayland.Proxy }

func (*Surface) Interface() *wayland.Interface { return surfaceInterface }

// Deletes the surface and invalidates its object ID.
func (obj *Surface) Destroy() {
	const wl_surface_destroy = 0
	obj.Conn().SendRequest(obj, wl_surface_destroy)
}

// Set a buffer as the content of this surface.
//
// The new size of the surface is calculated based on the buffer
// size transformed by the inverse buffer_transform and the
// inverse buffer_scale. This means that the supplied buffer
// must be an integer multiple of the buffer_scale.
//
// The x and y arguments specify the location of the new pending
// buffer's upper left corner, relative to the current buffer's upper
// left corner, in surface-local coordinates. In other words, the
// x and y, combined with the new surface size define in which
// directions the surface's size changes.
//
// Surface contents are double-buffered state, see wl_surface.commit.
//
// The initial surface contents are void; there is no content.
// wl_surface.attach assigns the given wl_buffer as the pending
// wl_buffer. wl_surface.commit makes the pending wl_buffer the new
// surface contents, and the size of the surface becomes the size
// calculated from the wl_buffer, as described above. After commit,
// there is no pending buffer until the next attach.
//
// Committing a pending wl_buffer allows the compositor to read the
// pixels in the wl_buffer. The compositor may access the pixels at
// any time after the wl_surface.commit request. When the compositor
// will not access the pixels anymore, it will send the
// wl_buffer.release event. Only after receiving wl_buffer.release,
// the client may reuse the wl_buffer. A wl_buffer that has been
// attached and then replaced by another attach instead of committed
// will not receive a release event, and is not used by the
// compositor.
//
// If a pending wl_buffer has been committed to more than one wl_surface,
// the delivery of wl_buffer.release events becomes undefined. A well
// behaved client should not rely on wl_buffer.release events in this
// case. Alternatively, a client could create multiple wl_buffer objects
// from the same backing storage or use wp_linux_buffer_release.
//
// Destroying the wl_buffer after wl_buffer.release does not change
// the surface contents. However, if the client destroys the
// wl_buffer before receiving the wl_buffer.release event, the surface
// contents become undefined immediately.
//
// If wl_surface.attach is sent with a NULL wl_buffer, the
// following wl_surface.commit will remove the surface content.
func (obj *Surface) Attach(buffer *Buffer, x int32, y int32) {
	const wl_surface_attach = 1
	obj.Conn().SendRequest(obj, wl_surface_attach, buffer, x, y)
}

// This request is used to describe the regions where the pending
// buffer is different from the current surface contents, and where
// the surface therefore needs to be repainted. The compositor
// ignores the parts of the damage that fall outside of the surface.
//
// Damage is double-buffered state, see wl_surface.commit.
//
// The damage rectangle is specified in surface-local coordinates,
// where x and y specify the upper left corner of the damage rectangle.
//
// The initial value for pending damage is empty: no damage.
// wl_surface.damage adds pending damage: the new pending damage
// is the union of old pending damage and the given rectangle.
//
// wl_surface.commit assigns pending damage as the current damage,
// and clears pending damage. The server will clear the current
// damage as it repaints the surface.
//
// Note! New clients should not use this request. Instead damage can be
// posted with wl_surface.damage_buffer which uses buffer coordinates
// instead of surface coordinates.
func (obj *Surface) Damage(x int32, y int32, width int32, height int32) {
	const wl_surface_damage = 2
	obj.Conn().SendRequest(obj, wl_surface_damage, x, y, width, height)
}

// Request a notification when it is a good time to start drawing a new
// frame, by creating a frame callback. This is useful for throttling
// redrawing operations, and driving animations.
//
// When a client is animating on a wl_surface, it can use the 'frame'
// request to get notified when it is a good time to draw and commit the
// next frame of animation. If the client commits an update earlier than
// that, it is likely that some updates will not make it to the display,
// and the client is wasting resources by drawing too often.
//
// The frame request will take effect on the next wl_surface.commit.
// The notification will only be posted for one frame unless
// requested again. For a wl_surface, the notifications are posted in
// the order the frame requests were committed.
//
// The server must send the notifications so that a client
// will not send excessive updates, while still allowing
// the highest possible update rate for clients that wait for the reply
// before drawing again. The server should give some time for the client
// to draw and commit after sending the frame callback events to let it
// hit the next output refresh.
//
// A server should avoid signaling the frame callbacks if the
// surface is not visible in any way, e.g. the surface is off-screen,
// or completely obscured by other opaque surfaces.
//
// The object returned by this request will be destroyed by the
// compositor after the callback is fired and as such the client must not
// attempt to use it after that point.
//
// The callback_data passed in the callback is the current time, in
// milliseconds, with an undefined base.
func (obj *Surface) Frame() *Callback {
	const wl_surface_frame = 3
	_ret := &Callback{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_surface_frame, _ret)
	return _ret
}

// This request sets the region of the surface that contains
// opaque content.
//
// The opaque region is an optimization hint for the compositor
// that lets it optimize the redrawing of content behind opaque
// regions.  Setting an opaque region is not required for correct
// behaviour, but marking transparent content as opaque will result
// in repaint artifacts.
//
// The opaque region is specified in surface-local coordinates.
//
// The compositor ignores the parts of the opaque region that fall
// outside of the surface.
//
// Opaque region is double-buffered state, see wl_surface.commit.
//
// wl_surface.set_opaque_region changes the pending opaque region.
// wl_surface.commit copies the pending region to the current region.
// Otherwise, the pending and current regions are never changed.
//
// The initial value for an opaque region is empty. Setting the pending
// opaque region has copy semantics, and the wl_region object can be
// destroyed immediately. A NULL wl_region causes the pending opaque
// region to be set to empty.
func (obj *Surface) SetOpaqueRegion(region *Region) {
	const wl_surface_set_opaque_region = 4
	obj.Conn().SendRequest(obj, wl_surface_set_opaque_region, region)
}

// This request sets the region of the surface that can receive
// pointer and touch events.
//
// Input events happening outside of this region will try the next
// surface in the server surface stack. The compositor ignores the
// parts of the input region that fall outside of the surface.
//
// The input region is specified in surface-local coordinates.
//
// Input region is double-buffered state, see wl_surface.commit.
//
// wl_surface.set_input_region changes the pending input region.
// wl_surface.commit copies the pending region to the current region.
// Otherwise the pending and current regions are never changed,
// except cursor and icon surfaces are special cases, see
// wl_pointer.set_cursor and wl_data_device.start_drag.
//
// The initial value for an input region is infinite. That means the
// whole surface will accept input. Setting the pending input region
// has copy semantics, and the wl_region object can be destroyed
// immediately. A NULL wl_region causes the input region to be set
// to infinite.
func (obj *Surface) SetInputRegion(region *Region) {
	const wl_surface_set_input_region = 5
	obj.Conn().SendRequest(obj, wl_surface_set_input_region, region)
}

// Surface state (input, opaque, and damage regions, attached buffers,
// etc.) is double-buffered. Protocol requests modify the pending state,
// as opposed to the current state in use by the compositor. A commit
// request atomically applies all pending state, replacing the current
// state. After commit, the new pending state is as documented for each
// related request.
//
// On commit, a pending wl_buffer is applied first, and all other state
// second. This means that all coordinates in double-buffered state are
// relative to the new wl_buffer coming into use, except for
// wl_surface.attach itself. If there is no pending wl_buffer, the
// coordinates are relative to the current surface contents.
//
// All requests that need a commit to become effective are documented
// to affect double-buffered state.
//
// Other interfaces may add further double-buffered surface state.
func (obj *Surface) Commit() {
	const wl_surface_commit = 6
	obj.Conn().SendRequest(obj, wl_surface_commit)
}

// This request sets an optional transformation on how the compositor
// interprets the contents of the buffer attached to the surface. The
// accepted values for the transform parameter are the values for
// wl_output.transform.
//
// Buffer transform is double-buffered state, see wl_surface.commit.
//
// A newly created surface has its buffer transformation set to normal.
//
// wl_surface.set_buffer_transform changes the pending buffer
// transformation. wl_surface.commit copies the pending buffer
// transformation to the current one. Otherwise, the pending and current
// values are never changed.
//
// The purpose of this request is to allow clients to render content
// according to the output transform, thus permitting the compositor to
// use certain optimizations even if the display is rotated. Using
// hardware overlays and scanning out a client buffer for fullscreen
// surfaces are examples of such optimizations. Those optimizations are
// highly dependent on the compositor implementation, so the use of this
// request should be considered on a case-by-case basis.
//
// Note that if the transform value includes 90 or 270 degree rotation,
// the width of the buffer will become the surface height and the height
// of the buffer will become the surface width.
//
// If transform is not one of the values from the
// wl_output.transform enum the invalid_transform protocol error
// is raised.
func (obj *Surface) SetBufferTransform(transform int32) {
	const wl_surface_set_buffer_transform = 7
	obj.Conn().SendRequest(obj, wl_surface_set_buffer_transform, transform)
}

// This request sets an optional scaling factor on how the compositor
// interprets the contents of the buffer attached to the window.
//
// Buffer scale is double-buffered state, see wl_surface.commit.
//
// A newly created surface has its buffer scale set to 1.
//
// wl_surface.set_buffer_scale changes the pending buffer scale.
// wl_surface.commit copies the pending buffer scale to the current one.
// Otherwise, the pending and current values are never changed.
//
// The purpose of this request is to allow clients to supply higher
// resolution buffer data for use on high resolution outputs. It is
// intended that you pick the same buffer scale as the scale of the
// output that the surface is displayed on. This means the compositor
// can avoid scaling when rendering the surface on that output.
//
// Note that if the scale is larger than 1, then you have to attach
// a buffer that is larger (by a factor of scale in each dimension)
// than the desired surface size.
//
// If scale is not positive the invalid_scale protocol error is
// raised.
func (obj *Surface) SetBufferScale(scale int32) {
	const wl_surface_set_buffer_scale = 8
	obj.Conn().SendRequest(obj, wl_surface_set_buffer_scale, scale)
}

// This request is used to describe the regions where the pending
// buffer is different from the current surface contents, and where
// the surface therefore needs to be repainted. The compositor
// ignores the parts of the damage that fall outside of the surface.
//
// Damage is double-buffered state, see wl_surface.commit.
//
// The damage rectangle is specified in buffer coordinates,
// where x and y specify the upper left corner of the damage rectangle.
//
// The initial value for pending damage is empty: no damage.
// wl_surface.damage_buffer adds pending damage: the new pending
// damage is the union of old pending damage and the given rectangle.
//
// wl_surface.commit assigns pending damage as the current damage,
// and clears pending damage. The server will clear the current
// damage as it repaints the surface.
//
// This request differs from wl_surface.damage in only one way - it
// takes damage in buffer coordinates instead of surface-local
// coordinates. While this generally is more intuitive than surface
// coordinates, it is especially desirable when using wp_viewport
// or when a drawing library (like EGL) is unaware of buffer scale
// and buffer transform.
//
// Note: Because buffer transformation changes and damage requests may
// be interleaved in the protocol stream, it is impossible to determine
// the actual mapping between surface and buffer damage until
// wl_surface.commit time. Therefore, compositors wishing to take both
// kinds of damage into account will have to accumulate damage from the
// two requests separately and only transform from one to the other
// after receiving the wl_surface.commit.
func (obj *Surface) DamageBuffer(x int32, y int32, width int32, height int32) {
	const wl_surface_damage_buffer = 9
	obj.Conn().SendRequest(obj, wl_surface_damage_buffer, x, y, width, height)
}

// This is a bitmask of capabilities this seat has; if a member is
// set, then it is present on the seat.
const (
	// the seat has pointer devices
	SeatCapabilityPointer = 1
	// the seat has one or more keyboards
	SeatCapabilityKeyboard = 2
	// the seat has touch devices
	SeatCapabilityTouch = 4
)

var seatInterface = &wayland.Interface{
	Name:    "wl_seat",
	Version: 7,
	Events:  []wayland.Event{(*SeatEventCapabilities)(nil), (*SeatEventName)(nil)},
}

type SeatEventCapabilities struct {
	// capabilities of the seat
	Capabilities uint32
}

type SeatEventName struct {
	// seat identifier
	Name string
}

// A seat is a group of keyboards, pointer and touch devices. This
// object is published as a global during start up, or when such a
// device is hot plugged.  A seat typically has a pointer and
// maintains a keyboard focus and a pointer focus.
type Seat struct{ wayland.Proxy }

func (*Seat) Interface() *wayland.Interface { return seatInterface }

// The ID provided will be initialized to the wl_pointer interface
// for this seat.
//
// This request only takes effect if the seat has the pointer
// capability, or has had the pointer capability in the past.
// It is a protocol violation to issue this request on a seat that has
// never had the pointer capability.
func (obj *Seat) GetPointer() *Pointer {
	const wl_seat_get_pointer = 0
	_ret := &Pointer{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_seat_get_pointer, _ret)
	return _ret
}

// The ID provided will be initialized to the wl_keyboard interface
// for this seat.
//
// This request only takes effect if the seat has the keyboard
// capability, or has had the keyboard capability in the past.
// It is a protocol violation to issue this request on a seat that has
// never had the keyboard capability.
func (obj *Seat) GetKeyboard() *Keyboard {
	const wl_seat_get_keyboard = 1
	_ret := &Keyboard{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_seat_get_keyboard, _ret)
	return _ret
}

// The ID provided will be initialized to the wl_touch interface
// for this seat.
//
// This request only takes effect if the seat has the touch
// capability, or has had the touch capability in the past.
// It is a protocol violation to issue this request on a seat that has
// never had the touch capability.
func (obj *Seat) GetTouch() *Touch {
	const wl_seat_get_touch = 2
	_ret := &Touch{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_seat_get_touch, _ret)
	return _ret
}

// Using this request a client can tell the server that it is not going to
// use the seat object anymore.
func (obj *Seat) Release() {
	const wl_seat_release = 3
	obj.Conn().SendRequest(obj, wl_seat_release)
}

const (
	// given wl_surface has another role
	PointerErrorRole = 0
)

// Describes the physical state of a button that produced the button
// event.
const (
	// the button is not pressed
	PointerButtonStateReleased = 0
	// the button is pressed
	PointerButtonStatePressed = 1
)

// Describes the axis types of scroll events.
const (
	// vertical axis
	PointerAxisVerticalScroll = 0
	// horizontal axis
	PointerAxisHorizontalScroll = 1
)

// Describes the source types for axis events. This indicates to the
// client how an axis event was physically generated; a client may
// adjust the user interface accordingly. For example, scroll events
// from a "finger" source may be in a smooth coordinate space with
// kinetic scrolling whereas a "wheel" source may be in discrete steps
// of a number of lines.
//
// The "continuous" axis source is a device generating events in a
// continuous coordinate space, but using something other than a
// finger. One example for this source is button-based scrolling where
// the vertical motion of a device is converted to scroll events while
// a button is held down.
//
// The "wheel tilt" axis source indicates that the actual device is a
// wheel but the scroll event is not caused by a rotation but a
// (usually sideways) tilt of the wheel.
const (
	// a physical wheel rotation
	PointerAxisSourceWheel = 0
	// finger on a touch surface
	PointerAxisSourceFinger = 1
	// continuous coordinate space
	PointerAxisSourceContinuous = 2
	// a physical wheel tilt
	PointerAxisSourceWheelTilt = 3
)

var pointerInterface = &wayland.Interface{
	Name:    "wl_pointer",
	Version: 7,
	Events:  []wayland.Event{(*PointerEventEnter)(nil), (*PointerEventLeave)(nil), (*PointerEventMotion)(nil), (*PointerEventButton)(nil), (*PointerEventAxis)(nil), (*PointerEventFrame)(nil), (*PointerEventAxisSource)(nil), (*PointerEventAxisStop)(nil), (*PointerEventAxisDiscrete)(nil)},
}

type PointerEventEnter struct {
	// serial number of the enter event
	Serial uint32
	// surface entered by the pointer
	Surface *Surface
	// surface-local x coordinate
	SurfaceX wayland.Fixed
	// surface-local y coordinate
	SurfaceY wayland.Fixed
}

type PointerEventLeave struct {
	// serial number of the leave event
	Serial uint32
	// surface left by the pointer
	Surface *Surface
}

type PointerEventMotion struct {
	// timestamp with millisecond granularity
	Time uint32
	// surface-local x coordinate
	SurfaceX wayland.Fixed
	// surface-local y coordinate
	SurfaceY wayland.Fixed
}

type PointerEventButton struct {
	// serial number of the button event
	Serial uint32
	// timestamp with millisecond granularity
	Time uint32
	// button that produced the event
	Button uint32
	// physical state of the button
	State uint32
}

type PointerEventAxis struct {
	// timestamp with millisecond granularity
	Time uint32
	// axis type
	Axis uint32
	// length of vector in surface-local coordinate space
	Value wayland.Fixed
}

type PointerEventFrame struct {
}

type PointerEventAxisSource struct {
	// source of the axis event
	AxisSource uint32
}

type PointerEventAxisStop struct {
	// timestamp with millisecond granularity
	Time uint32
	// the axis stopped with this event
	Axis uint32
}

type PointerEventAxisDiscrete struct {
	// axis type
	Axis uint32
	// number of steps
	Discrete int32
}

// The wl_pointer interface represents one or more input devices,
// such as mice, which control the pointer location and pointer_focus
// of a seat.
//
// The wl_pointer interface generates motion, enter and leave
// events for the surfaces that the pointer is located over,
// and button and axis events for button presses, button releases
// and scrolling.
type Pointer struct{ wayland.Proxy }

func (*Pointer) Interface() *wayland.Interface { return pointerInterface }

// Set the pointer surface, i.e., the surface that contains the
// pointer image (cursor). This request gives the surface the role
// of a cursor. If the surface already has another role, it raises
// a protocol error.
//
// The cursor actually changes only if the pointer
// focus for this device is one of the requesting client's surfaces
// or the surface parameter is the current pointer surface. If
// there was a previous surface set with this request it is
// replaced. If surface is NULL, the pointer image is hidden.
//
// The parameters hotspot_x and hotspot_y define the position of
// the pointer surface relative to the pointer location. Its
// top-left corner is always at (x, y) - (hotspot_x, hotspot_y),
// where (x, y) are the coordinates of the pointer location, in
// surface-local coordinates.
//
// On surface.attach requests to the pointer surface, hotspot_x
// and hotspot_y are decremented by the x and y parameters
// passed to the request. Attach must be confirmed by
// wl_surface.commit as usual.
//
// The hotspot can also be updated by passing the currently set
// pointer surface to this request with new values for hotspot_x
// and hotspot_y.
//
// The current and pending input regions of the wl_surface are
// cleared, and wl_surface.set_input_region is ignored until the
// wl_surface is no longer used as the cursor. When the use as a
// cursor ends, the current and pending input regions become
// undefined, and the wl_surface is unmapped.
func (obj *Pointer) SetCursor(serial uint32, surface *Surface, hotspotX int32, hotspotY int32) {
	const wl_pointer_set_cursor = 0
	obj.Conn().SendRequest(obj, wl_pointer_set_cursor, serial, surface, hotspotX, hotspotY)
}

// Using this request a client can tell the server that it is not going to
// use the pointer object anymore.
//
// This request destroys the pointer proxy object, so clients must not call
// wl_pointer_destroy() after using this request.
func (obj *Pointer) Release() {
	const wl_pointer_release = 1
	obj.Conn().SendRequest(obj, wl_pointer_release)
}

// This specifies the format of the keymap provided to the
// client with the wl_keyboard.keymap event.
const (
	// no keymap; client must understand how to interpret the raw keycode
	KeyboardKeymapFormatNoKeymap = 0
	// libxkbcommon compatible; to determine the xkb keycode, clients must add 8 to the key event keycode
	KeyboardKeymapFormatXkbV1 = 1
)

// Describes the physical state of a key that produced the key event.
const (
	// key is not pressed
	KeyboardKeyStateReleased = 0
	// key is pressed
	KeyboardKeyStatePressed = 1
)

var keyboardInterface = &wayland.Interface{
	Name:    "wl_keyboard",
	Version: 7,
	Events:  []wayland.Event{(*KeyboardEventKeymap)(nil), (*KeyboardEventEnter)(nil), (*KeyboardEventLeave)(nil), (*KeyboardEventKey)(nil), (*KeyboardEventModifiers)(nil), (*KeyboardEventRepeatInfo)(nil)},
}

type KeyboardEventKeymap struct {
	// keymap format
	Format uint32
	// keymap file descriptor
	Fd uintptr
	// keymap size, in bytes
	Size uint32
}

type KeyboardEventEnter struct {
	// serial number of the enter event
	Serial uint32
	// surface gaining keyboard focus
	Surface *Surface
	// the currently pressed keys
	Keys []byte
}

type KeyboardEventLeave struct {
	// serial number of the leave event
	Serial uint32
	// surface that lost keyboard focus
	Surface *Surface
}

type KeyboardEventKey struct {
	// serial number of the key event
	Serial uint32
	// timestamp with millisecond granularity
	Time uint32
	// key that produced the event
	Key uint32
	// physical state of the key
	State uint32
}

type KeyboardEventModifiers struct {
	// serial number of the modifiers event
	Serial uint32
	// depressed modifiers
	ModsDepressed uint32
	// latched modifiers
	ModsLatched uint32
	// locked modifiers
	ModsLocked uint32
	// keyboard layout
	Group uint32
}

type KeyboardEventRepeatInfo struct {
	// the rate of repeating keys in characters per second
	Rate int32
	// delay in milliseconds since key down until repeating starts
	Delay int32
}

// The wl_keyboard interface represents one or more keyboards
// associated with a seat.
type Keyboard struct{ wayland.Proxy }

func (*Keyboard) Interface() *wayland.Interface { return keyboardInterface }

// release the keyboard object
func (obj *Keyboard) Release() {
	const wl_keyboard_release = 0
	obj.Conn().SendRequest(obj, wl_keyboard_release)
}

var touchInterface = &wayland.Interface{
	Name:    "wl_touch",
	Version: 7,
	Events:  []wayland.Event{(*TouchEventDown)(nil), (*TouchEventUp)(nil), (*TouchEventMotion)(nil), (*TouchEventFrame)(nil), (*TouchEventCancel)(nil), (*TouchEventShape)(nil), (*TouchEventOrientation)(nil)},
}

type TouchEventDown struct {
	// serial number of the touch down event
	Serial uint32
	// timestamp with millisecond granularity
	Time uint32
	// surface touched
	Surface *Surface
	// the unique ID of this touch point
	ID int32
	// surface-local x coordinate
	X wayland.Fixed
	// surface-local y coordinate
	Y wayland.Fixed
}

type TouchEventUp struct {
	// serial number of the touch up event
	Serial uint32
	// timestamp with millisecond granularity
	Time uint32
	// the unique ID of this touch point
	ID int32
}

type TouchEventMotion struct {
	// timestamp with millisecond granularity
	Time uint32
	// the unique ID of this touch point
	ID int32
	// surface-local x coordinate
	X wayland.Fixed
	// surface-local y coordinate
	Y wayland.Fixed
}

type TouchEventFrame struct {
}

type TouchEventCancel struct {
}

type TouchEventShape struct {
	// the unique ID of this touch point
	ID int32
	// length of the major axis in surface-local coordinates
	Major wayland.Fixed
	// length of the minor axis in surface-local coordinates
	Minor wayland.Fixed
}

type TouchEventOrientation struct {
	// the unique ID of this touch point
	ID int32
	// angle between major axis and positive surface y-axis in degrees
	Orientation wayland.Fixed
}

// The wl_touch interface represents a touchscreen
// associated with a seat.
//
// Touch interactions can consist of one or more contacts.
// For each contact, a series of events is generated, starting
// with a down event, followed by zero or more motion events,
// and ending with an up event. Events relating to the same
// contact point can be identified by the ID of the sequence.
type Touch struct{ wayland.Proxy }

func (*Touch) Interface() *wayland.Interface { return touchInterface }

// release the touch object
func (obj *Touch) Release() {
	const wl_touch_release = 0
	obj.Conn().SendRequest(obj, wl_touch_release)
}

// This enumeration describes how the physical
// pixels on an output are laid out.
const (
	// unknown geometry
	OutputSubpixelUnknown = 0
	// no geometry
	OutputSubpixelNone = 1
	// horizontal RGB
	OutputSubpixelHorizontalRgb = 2
	// horizontal BGR
	OutputSubpixelHorizontalBgr = 3
	// vertical RGB
	OutputSubpixelVerticalRgb = 4
	// vertical BGR
	OutputSubpixelVerticalBgr = 5
)

// This describes the transform that a compositor will apply to a
// surface to compensate for the rotation or mirroring of an
// output device.
//
// The flipped values correspond to an initial flip around a
// vertical axis followed by rotation.
//
// The purpose is mainly to allow clients to render accordingly and
// tell the compositor, so that for fullscreen surfaces, the
// compositor will still be able to scan out directly from client
// surfaces.
const (
	// no transform
	OutputTransformNormal = 0
	// 90 degrees counter-clockwise
	OutputTransform90 = 1
	// 180 degrees counter-clockwise
	OutputTransform180 = 2
	// 270 degrees counter-clockwise
	OutputTransform270 = 3
	// 180 degree flip around a vertical axis
	OutputTransformFlipped = 4
	// flip and rotate 90 degrees counter-clockwise
	OutputTransformFlipped90 = 5
	// flip and rotate 180 degrees counter-clockwise
	OutputTransformFlipped180 = 6
	// flip and rotate 270 degrees counter-clockwise
	OutputTransformFlipped270 = 7
)

// These flags describe properties of an output mode.
// They are used in the flags bitfield of the mode event.
const (
	// indicates this is the current mode
	OutputModeCurrent = 0x1
	// indicates this is the preferred mode
	OutputModePreferred = 0x2
)

var outputInterface = &wayland.Interface{
	Name:    "wl_output",
	Version: 3,
	Events:  []wayland.Event{(*OutputEventGeometry)(nil), (*OutputEventMode)(nil), (*OutputEventDone)(nil), (*OutputEventScale)(nil)},
}

type OutputEventGeometry struct {
	// x position within the global compositor space
	X int32
	// y position within the global compositor space
	Y int32
	// width in millimeters of the output
	PhysicalWidth int32
	// height in millimeters of the output
	PhysicalHeight int32
	// subpixel orientation of the output
	Subpixel int32
	// textual description of the manufacturer
	Make string
	// textual description of the model
	Model string
	// transform that maps framebuffer to output
	Transform int32
}

type OutputEventMode struct {
	// bitfield of mode flags
	Flags uint32
	// width of the mode in hardware units
	Width int32
	// height of the mode in hardware units
	Height int32
	// vertical refresh rate in mHz
	Refresh int32
}

type OutputEventDone struct {
}

type OutputEventScale struct {
	// scaling factor of output
	Factor int32
}

// An output describes part of the compositor geometry.  The
// compositor works in the 'compositor coordinate system' and an
// output corresponds to a rectangular area in that space that is
// actually visible.  This typically corresponds to a monitor that
// displays part of the compositor space.  This object is published
// as global during start up, or when a monitor is hotplugged.
type Output struct{ wayland.Proxy }

func (*Output) Interface() *wayland.Interface { return outputInterface }

// Using this request a client can tell the server that it is not going to
// use the output object anymore.
func (obj *Output) Release() {
	const wl_output_release = 0
	obj.Conn().SendRequest(obj, wl_output_release)
}

var regionInterface = &wayland.Interface{
	Name:    "wl_region",
	Version: 1,
	Events:  []wayland.Event{},
}

// A region object describes an area.
//
// Region objects are used to describe the opaque and input
// regions of a surface.
type Region struct{ wayland.Proxy }

func (*Region) Interface() *wayland.Interface { return regionInterface }

// Destroy the region.  This will invalidate the object ID.
func (obj *Region) Destroy() {
	const wl_region_destroy = 0
	obj.Conn().SendRequest(obj, wl_region_destroy)
}

// Add the specified rectangle to the region.
func (obj *Region) Add(x int32, y int32, width int32, height int32) {
	const wl_region_add = 1
	obj.Conn().SendRequest(obj, wl_region_add, x, y, width, height)
}

// Subtract the specified rectangle from the region.
func (obj *Region) Subtract(x int32, y int32, width int32, height int32) {
	const wl_region_subtract = 2
	obj.Conn().SendRequest(obj, wl_region_subtract, x, y, width, height)
}

const (
	// the to-be sub-surface is invalid
	SubcompositorErrorBadSurface = 0
)

var subcompositorInterface = &wayland.Interface{
	Name:    "wl_subcompositor",
	Version: 1,
	Events:  []wayland.Event{},
}

// The global interface exposing sub-surface compositing capabilities.
// A wl_surface, that has sub-surfaces associated, is called the
// parent surface. Sub-surfaces can be arbitrarily nested and create
// a tree of sub-surfaces.
//
// The root surface in a tree of sub-surfaces is the main
// surface. The main surface cannot be a sub-surface, because
// sub-surfaces must always have a parent.
//
// A main surface with its sub-surfaces forms a (compound) window.
// For window management purposes, this set of wl_surface objects is
// to be considered as a single window, and it should also behave as
// such.
//
// The aim of sub-surfaces is to offload some of the compositing work
// within a window from clients to the compositor. A prime example is
// a video player with decorations and video in separate wl_surface
// objects. This should allow the compositor to pass YUV video buffer
// processing to dedicated overlay hardware when possible.
type Subcompositor struct{ wayland.Proxy }

func (*Subcompositor) Interface() *wayland.Interface { return subcompositorInterface }

// Informs the server that the client will not be using this
// protocol object anymore. This does not affect any other
// objects, wl_subsurface objects included.
func (obj *Subcompositor) Destroy() {
	const wl_subcompositor_destroy = 0
	obj.Conn().SendRequest(obj, wl_subcompositor_destroy)
}

// Create a sub-surface interface for the given surface, and
// associate it with the given parent surface. This turns a
// plain wl_surface into a sub-surface.
//
// The to-be sub-surface must not already have another role, and it
// must not have an existing wl_subsurface object. Otherwise a protocol
// error is raised.
//
// Adding sub-surfaces to a parent is a double-buffered operation on the
// parent (see wl_surface.commit). The effect of adding a sub-surface
// becomes visible on the next time the state of the parent surface is
// applied.
//
// This request modifies the behaviour of wl_surface.commit request on
// the sub-surface, see the documentation on wl_subsurface interface.
func (obj *Subcompositor) GetSubsurface(surface *Surface, parent *Surface) *Subsurface {
	const wl_subcompositor_get_subsurface = 1
	_ret := &Subsurface{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, wl_subcompositor_get_subsurface, _ret, surface, parent)
	return _ret
}

const (
	// wl_surface is not a sibling or the parent
	SubsurfaceErrorBadSurface = 0
)

var subsurfaceInterface = &wayland.Interface{
	Name:    "wl_subsurface",
	Version: 1,
	Events:  []wayland.Event{},
}

// An additional interface to a wl_surface object, which has been
// made a sub-surface. A sub-surface has one parent surface. A
// sub-surface's size and position are not limited to that of the parent.
// Particularly, a sub-surface is not automatically clipped to its
// parent's area.
//
// A sub-surface becomes mapped, when a non-NULL wl_buffer is applied
// and the parent surface is mapped. The order of which one happens
// first is irrelevant. A sub-surface is hidden if the parent becomes
// hidden, or if a NULL wl_buffer is applied. These rules apply
// recursively through the tree of surfaces.
//
// The behaviour of a wl_surface.commit request on a sub-surface
// depends on the sub-surface's mode. The possible modes are
// synchronized and desynchronized, see methods
// wl_subsurface.set_sync and wl_subsurface.set_desync. Synchronized
// mode caches the wl_surface state to be applied when the parent's
// state gets applied, and desynchronized mode applies the pending
// wl_surface state directly. A sub-surface is initially in the
// synchronized mode.
//
// Sub-surfaces have also other kind of state, which is managed by
// wl_subsurface requests, as opposed to wl_surface requests. This
// state includes the sub-surface position relative to the parent
// surface (wl_subsurface.set_position), and the stacking order of
// the parent and its sub-surfaces (wl_subsurface.place_above and
// .place_below). This state is applied when the parent surface's
// wl_surface state is applied, regardless of the sub-surface's mode.
// As the exception, set_sync and set_desync are effective immediately.
//
// The main surface can be thought to be always in desynchronized mode,
// since it does not have a parent in the sub-surfaces sense.
//
// Even if a sub-surface is in desynchronized mode, it will behave as
// in synchronized mode, if its parent surface behaves as in
// synchronized mode. This rule is applied recursively throughout the
// tree of surfaces. This means, that one can set a sub-surface into
// synchronized mode, and then assume that all its child and grand-child
// sub-surfaces are synchronized, too, without explicitly setting them.
//
// If the wl_surface associated with the wl_subsurface is destroyed, the
// wl_subsurface object becomes inert. Note, that destroying either object
// takes effect immediately. If you need to synchronize the removal
// of a sub-surface to the parent surface update, unmap the sub-surface
// first by attaching a NULL wl_buffer, update parent, and then destroy
// the sub-surface.
//
// If the parent wl_surface object is destroyed, the sub-surface is
// unmapped.
type Subsurface struct{ wayland.Proxy }

func (*Subsurface) Interface() *wayland.Interface { return subsurfaceInterface }

// The sub-surface interface is removed from the wl_surface object
// that was turned into a sub-surface with a
// wl_subcompositor.get_subsurface request. The wl_surface's association
// to the parent is deleted, and the wl_surface loses its role as
// a sub-surface. The wl_surface is unmapped immediately.
func (obj *Subsurface) Destroy() {
	const wl_subsurface_destroy = 0
	obj.Conn().SendRequest(obj, wl_subsurface_destroy)
}

// This schedules a sub-surface position change.
// The sub-surface will be moved so that its origin (top left
// corner pixel) will be at the location x, y of the parent surface
// coordinate system. The coordinates are not restricted to the parent
// surface area. Negative values are allowed.
//
// The scheduled coordinates will take effect whenever the state of the
// parent surface is applied. When this happens depends on whether the
// parent surface is in synchronized mode or not. See
// wl_subsurface.set_sync and wl_subsurface.set_desync for details.
//
// If more than one set_position request is invoked by the client before
// the commit of the parent surface, the position of a new request always
// replaces the scheduled position from any previous request.
//
// The initial position is 0, 0.
func (obj *Subsurface) SetPosition(x int32, y int32) {
	const wl_subsurface_set_position = 1
	obj.Conn().SendRequest(obj, wl_subsurface_set_position, x, y)
}

// This sub-surface is taken from the stack, and put back just
// above the reference surface, changing the z-order of the sub-surfaces.
// The reference surface must be one of the sibling surfaces, or the
// parent surface. Using any other surface, including this sub-surface,
// will cause a protocol error.
//
// The z-order is double-buffered. Requests are handled in order and
// applied immediately to a pending state. The final pending state is
// copied to the active state the next time the state of the parent
// surface is applied. When this happens depends on whether the parent
// surface is in synchronized mode or not. See wl_subsurface.set_sync and
// wl_subsurface.set_desync for details.
//
// A new sub-surface is initially added as the top-most in the stack
// of its siblings and parent.
func (obj *Subsurface) PlaceAbove(sibling *Surface) {
	const wl_subsurface_place_above = 2
	obj.Conn().SendRequest(obj, wl_subsurface_place_above, sibling)
}

// The sub-surface is placed just below the reference surface.
// See wl_subsurface.place_above.
func (obj *Subsurface) PlaceBelow(sibling *Surface) {
	const wl_subsurface_place_below = 3
	obj.Conn().SendRequest(obj, wl_subsurface_place_below, sibling)
}

// Change the commit behaviour of the sub-surface to synchronized
// mode, also described as the parent dependent mode.
//
// In synchronized mode, wl_surface.commit on a sub-surface will
// accumulate the committed state in a cache, but the state will
// not be applied and hence will not change the compositor output.
// The cached state is applied to the sub-surface immediately after
// the parent surface's state is applied. This ensures atomic
// updates of the parent and all its synchronized sub-surfaces.
// Applying the cached state will invalidate the cache, so further
// parent surface commits do not (re-)apply old state.
//
// See wl_subsurface for the recursive effect of this mode.
func (obj *Subsurface) SetSync() {
	const wl_subsurface_set_sync = 4
	obj.Conn().SendRequest(obj, wl_subsurface_set_sync)
}

// Change the commit behaviour of the sub-surface to desynchronized
// mode, also described as independent or freely running mode.
//
// In desynchronized mode, wl_surface.commit on a sub-surface will
// apply the pending state directly, without caching, as happens
// normally with a wl_surface. Calling wl_surface.commit on the
// parent surface has no effect on the sub-surface's wl_surface
// state. This mode allows a sub-surface to be updated on its own.
//
// If cached state exists when wl_surface.commit is called in
// desynchronized mode, the pending state is added to the cached
// state, and applied as a whole. This invalidates the cache.
//
// Note: even if a sub-surface is set to desynchronized, a parent
// sub-surface may override it to behave as synchronized. For details,
// see wl_subsurface.
//
// If a surface's parent surface behaves as desynchronized, then
// the cached state is applied on set_desync.
func (obj *Subsurface) SetDesync() {
	const wl_subsurface_set_desync = 5
	obj.Conn().SendRequest(obj, wl_subsurface_set_desync)
}

const (
	// given wl_surface has another role
	XdgWmBaseErrorRole = 0
	// xdg_wm_base was destroyed before children
	XdgWmBaseErrorDefunctSurfaces = 1
	// the client tried to map or destroy a non-topmost popup
	XdgWmBaseErrorNotTheTopmostPopup = 2
	// the client specified an invalid popup parent surface
	XdgWmBaseErrorInvalidPopupParent = 3
	// the client provided an invalid surface state
	XdgWmBaseErrorInvalidSurfaceState = 4
	// the client provided an invalid positioner
	XdgWmBaseErrorInvalidPositioner = 5
)

var xdgWmBaseInterface = &wayland.Interface{
	Name:    "xdg_wm_base",
	Version: 2,
	Events:  []wayland.Event{(*XdgWmBaseEventPing)(nil)},
}

type XdgWmBaseEventPing struct {
	// pass this to the pong request
	Serial uint32
}

// The xdg_wm_base interface is exposed as a global object enabling clients
// to turn their wl_surfaces into windows in a desktop environment. It
// defines the basic functionality needed for clients and the compositor to
// create windows that can be dragged, resized, maximized, etc, as well as
// creating transient windows such as popup menus.
type XdgWmBase struct{ wayland.Proxy }

func (*XdgWmBase) Interface() *wayland.Interface { return xdgWmBaseInterface }

// Destroy this xdg_wm_base object.
//
// Destroying a bound xdg_wm_base object while there are surfaces
// still alive created by this xdg_wm_base object instance is illegal
// and will result in a protocol error.
func (obj *XdgWmBase) Destroy() {
	const xdg_wm_base_destroy = 0
	obj.Conn().SendRequest(obj, xdg_wm_base_destroy)
}

// Create a positioner object. A positioner object is used to position
// surfaces relative to some parent surface. See the interface description
// and xdg_surface.get_popup for details.
func (obj *XdgWmBase) CreatePositioner() *XdgPositioner {
	const xdg_wm_base_create_positioner = 1
	_ret := &XdgPositioner{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, xdg_wm_base_create_positioner, _ret)
	return _ret
}

// This creates an xdg_surface for the given surface. While xdg_surface
// itself is not a role, the corresponding surface may only be assigned
// a role extending xdg_surface, such as xdg_toplevel or xdg_popup.
//
// This creates an xdg_surface for the given surface. An xdg_surface is
// used as basis to define a role to a given surface, such as xdg_toplevel
// or xdg_popup. It also manages functionality shared between xdg_surface
// based surface roles.
//
// See the documentation of xdg_surface for more details about what an
// xdg_surface is and how it is used.
func (obj *XdgWmBase) GetXdgSurface(surface *Surface) *XdgSurface {
	const xdg_wm_base_get_xdg_surface = 2
	_ret := &XdgSurface{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, xdg_wm_base_get_xdg_surface, _ret, surface)
	return _ret
}

// A client must respond to a ping event with a pong request or
// the client may be deemed unresponsive. See xdg_wm_base.ping.
func (obj *XdgWmBase) Pong(serial uint32) {
	const xdg_wm_base_pong = 3
	obj.Conn().SendRequest(obj, xdg_wm_base_pong, serial)
}

const (
	// invalid input provided
	XdgPositionerErrorInvalidInput = 0
)
const (
	XdgPositionerAnchorNone        = 0
	XdgPositionerAnchorTop         = 1
	XdgPositionerAnchorBottom      = 2
	XdgPositionerAnchorLeft        = 3
	XdgPositionerAnchorRight       = 4
	XdgPositionerAnchorTopLeft     = 5
	XdgPositionerAnchorBottomLeft  = 6
	XdgPositionerAnchorTopRight    = 7
	XdgPositionerAnchorBottomRight = 8
)
const (
	XdgPositionerGravityNone        = 0
	XdgPositionerGravityTop         = 1
	XdgPositionerGravityBottom      = 2
	XdgPositionerGravityLeft        = 3
	XdgPositionerGravityRight       = 4
	XdgPositionerGravityTopLeft     = 5
	XdgPositionerGravityBottomLeft  = 6
	XdgPositionerGravityTopRight    = 7
	XdgPositionerGravityBottomRight = 8
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
const (
	// Don't alter the surface position even if it is constrained on some
	// axis, for example partially outside the edge of an output.
	XdgPositionerConstraintAdjustmentNone = 0
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
	XdgPositionerConstraintAdjustmentSlideX = 1
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
	XdgPositionerConstraintAdjustmentSlideY = 2
	// Invert the anchor and gravity on the x axis if the surface is
	// constrained on the x axis. For example, if the left edge of the
	// surface is constrained, the gravity is 'left' and the anchor is
	// 'left', change the gravity to 'right' and the anchor to 'right'.
	//
	// If the adjusted position also ends up being constrained, the resulting
	// position of the flip_x adjustment will be the one before the
	// adjustment.
	XdgPositionerConstraintAdjustmentFlipX = 4
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
	XdgPositionerConstraintAdjustmentFlipY = 8
	// Resize the surface horizontally so that it is completely
	// unconstrained.
	XdgPositionerConstraintAdjustmentResizeX = 16
	// Resize the surface vertically so that it is completely unconstrained.
	XdgPositionerConstraintAdjustmentResizeY = 32
)

var xdgPositionerInterface = &wayland.Interface{
	Name:    "xdg_positioner",
	Version: 2,
	Events:  []wayland.Event{},
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
type XdgPositioner struct{ wayland.Proxy }

func (*XdgPositioner) Interface() *wayland.Interface { return xdgPositionerInterface }

// Notify the compositor that the xdg_positioner will no longer be used.
func (obj *XdgPositioner) Destroy() {
	const xdg_positioner_destroy = 0
	obj.Conn().SendRequest(obj, xdg_positioner_destroy)
}

// Set the size of the surface that is to be positioned with the positioner
// object. The size is in surface-local coordinates and corresponds to the
// window geometry. See xdg_surface.set_window_geometry.
//
// If a zero or negative size is set the invalid_input error is raised.
func (obj *XdgPositioner) SetSize(width int32, height int32) {
	const xdg_positioner_set_size = 1
	obj.Conn().SendRequest(obj, xdg_positioner_set_size, width, height)
}

// Specify the anchor rectangle within the parent surface that the child
// surface will be placed relative to. The rectangle is relative to the
// window geometry as defined by xdg_surface.set_window_geometry of the
// parent surface.
//
// When the xdg_positioner object is used to position a child surface, the
// anchor rectangle may not extend outside the window geometry of the
// positioned child's parent surface.
//
// If a negative size is set the invalid_input error is raised.
func (obj *XdgPositioner) SetAnchorRect(x int32, y int32, width int32, height int32) {
	const xdg_positioner_set_anchor_rect = 2
	obj.Conn().SendRequest(obj, xdg_positioner_set_anchor_rect, x, y, width, height)
}

// Defines the anchor point for the anchor rectangle. The specified anchor
// is used derive an anchor point that the child surface will be
// positioned relative to. If a corner anchor is set (e.g. 'top_left' or
// 'bottom_right'), the anchor point will be at the specified corner;
// otherwise, the derived anchor point will be centered on the specified
// edge, or in the center of the anchor rectangle if no edge is specified.
func (obj *XdgPositioner) SetAnchor(anchor uint32) {
	const xdg_positioner_set_anchor = 3
	obj.Conn().SendRequest(obj, xdg_positioner_set_anchor, anchor)
}

// Defines in what direction a surface should be positioned, relative to
// the anchor point of the parent surface. If a corner gravity is
// specified (e.g. 'bottom_right' or 'top_left'), then the child surface
// will be placed towards the specified gravity; otherwise, the child
// surface will be centered over the anchor point on any axis that had no
// gravity specified.
func (obj *XdgPositioner) SetGravity(gravity uint32) {
	const xdg_positioner_set_gravity = 4
	obj.Conn().SendRequest(obj, xdg_positioner_set_gravity, gravity)
}

// Specify how the window should be positioned if the originally intended
// position caused the surface to be constrained, meaning at least
// partially outside positioning boundaries set by the compositor. The
// adjustment is set by constructing a bitmask describing the adjustment to
// be made when the surface is constrained on that axis.
//
// If no bit for one axis is set, the compositor will assume that the child
// surface should not change its position on that axis when constrained.
//
// If more than one bit for one axis is set, the order of how adjustments
// are applied is specified in the corresponding adjustment descriptions.
//
// The default adjustment is none.
func (obj *XdgPositioner) SetConstraintAdjustment(constraintAdjustment uint32) {
	const xdg_positioner_set_constraint_adjustment = 5
	obj.Conn().SendRequest(obj, xdg_positioner_set_constraint_adjustment, constraintAdjustment)
}

// Specify the surface position offset relative to the position of the
// anchor on the anchor rectangle and the anchor on the surface. For
// example if the anchor of the anchor rectangle is at (x, y), the surface
// has the gravity bottom|right, and the offset is (ox, oy), the calculated
// surface position will be (x + ox, y + oy). The offset position of the
// surface is the one used for constraint testing. See
// set_constraint_adjustment.
//
// An example use case is placing a popup menu on top of a user interface
// element, while aligning the user interface element of the parent surface
// with some user interface element placed somewhere in the popup surface.
func (obj *XdgPositioner) SetOffset(x int32, y int32) {
	const xdg_positioner_set_offset = 6
	obj.Conn().SendRequest(obj, xdg_positioner_set_offset, x, y)
}

const (
	XdgSurfaceErrorNotConstructed     = 1
	XdgSurfaceErrorAlreadyConstructed = 2
	XdgSurfaceErrorUnconfiguredBuffer = 3
)

var xdgSurfaceInterface = &wayland.Interface{
	Name:    "xdg_surface",
	Version: 2,
	Events:  []wayland.Event{(*XdgSurfaceEventConfigure)(nil)},
}

type XdgSurfaceEventConfigure struct {
	// serial of the configure event
	Serial uint32
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
// has not been destroyed.
type XdgSurface struct{ wayland.Proxy }

func (*XdgSurface) Interface() *wayland.Interface { return xdgSurfaceInterface }

// Destroy the xdg_surface object. An xdg_surface must only be destroyed
// after its role object has been destroyed.
func (obj *XdgSurface) Destroy() {
	const xdg_surface_destroy = 0
	obj.Conn().SendRequest(obj, xdg_surface_destroy)
}

// This creates an xdg_toplevel object for the given xdg_surface and gives
// the associated wl_surface the xdg_toplevel role.
//
// See the documentation of xdg_toplevel for more details about what an
// xdg_toplevel is and how it is used.
func (obj *XdgSurface) GetToplevel() *XdgToplevel {
	const xdg_surface_get_toplevel = 1
	_ret := &XdgToplevel{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, xdg_surface_get_toplevel, _ret)
	return _ret
}

// This creates an xdg_popup object for the given xdg_surface and gives
// the associated wl_surface the xdg_popup role.
//
// If null is passed as a parent, a parent surface must be specified using
// some other protocol, before committing the initial state.
//
// See the documentation of xdg_popup for more details about what an
// xdg_popup is and how it is used.
func (obj *XdgSurface) GetPopup(parent *XdgSurface, positioner *XdgPositioner) *XdgPopup {
	const xdg_surface_get_popup = 2
	_ret := &XdgPopup{}
	obj.Conn().NewProxy(0, _ret)
	obj.Conn().SendRequest(obj, xdg_surface_get_popup, _ret, parent, positioner)
	return _ret
}

// The window geometry of a surface is its "visible bounds" from the
// user's perspective. Client-side decorations often have invisible
// portions like drop-shadows which should be ignored for the
// purposes of aligning, placing and constraining windows.
//
// The window geometry is double buffered, and will be applied at the
// time wl_surface.commit of the corresponding wl_surface is called.
//
// When maintaining a position, the compositor should treat the (x, y)
// coordinate of the window geometry as the top left corner of the window.
// A client changing the (x, y) window geometry coordinate should in
// general not alter the position of the window.
//
// Once the window geometry of the surface is set, it is not possible to
// unset it, and it will remain the same until set_window_geometry is
// called again, even if a new subsurface or buffer is attached.
//
// If never set, the value is the full bounds of the surface,
// including any subsurfaces. This updates dynamically on every
// commit. This unset is meant for extremely simple clients.
//
// The arguments are given in the surface-local coordinate space of
// the wl_surface associated with this xdg_surface.
//
// The width and height must be greater than zero. Setting an invalid size
// will raise an error. When applied, the effective window geometry will be
// the set window geometry clamped to the bounding rectangle of the
// combined geometry of the surface of the xdg_surface and the associated
// subsurfaces.
func (obj *XdgSurface) SetWindowGeometry(x int32, y int32, width int32, height int32) {
	const xdg_surface_set_window_geometry = 3
	obj.Conn().SendRequest(obj, xdg_surface_set_window_geometry, x, y, width, height)
}

// When a configure event is received, if a client commits the
// surface in response to the configure event, then the client
// must make an ack_configure request sometime before the commit
// request, passing along the serial of the configure event.
//
// For instance, for toplevel surfaces the compositor might use this
// information to move a surface to the top left only when the client has
// drawn itself for the maximized or fullscreen state.
//
// If the client receives multiple configure events before it
// can respond to one, it only has to ack the last configure event.
//
// A client is not required to commit immediately after sending
// an ack_configure request - it may even ack_configure several times
// before its next surface commit.
//
// A client may send multiple ack_configure requests before committing, but
// only the last request sent before a commit indicates which configure
// event the client really is responding to.
func (obj *XdgSurface) AckConfigure(serial uint32) {
	const xdg_surface_ack_configure = 4
	obj.Conn().SendRequest(obj, xdg_surface_ack_configure, serial)
}

// These values are used to indicate which edge of a surface
// is being dragged in a resize operation.
const (
	XdgToplevelResizeEdgeNone        = 0
	XdgToplevelResizeEdgeTop         = 1
	XdgToplevelResizeEdgeBottom      = 2
	XdgToplevelResizeEdgeLeft        = 4
	XdgToplevelResizeEdgeTopLeft     = 5
	XdgToplevelResizeEdgeBottomLeft  = 6
	XdgToplevelResizeEdgeRight       = 8
	XdgToplevelResizeEdgeTopRight    = 9
	XdgToplevelResizeEdgeBottomRight = 10
)

// The different state values used on the surface. This is designed for
// state values like maximized, fullscreen. It is paired with the
// configure event to ensure that both the client and the compositor
// setting the state can be synchronized.
//
// States set in this way are double-buffered. They will get applied on
// the next commit.
const (
	// The surface is maximized. The window geometry specified in the configure
	// event must be obeyed by the client.
	//
	// The client should draw without shadow or other
	// decoration outside of the window geometry.
	XdgToplevelStateMaximized = 1
	// The surface is fullscreen. The window geometry specified in the
	// configure event is a maximum; the client cannot resize beyond it. For
	// a surface to cover the whole fullscreened area, the geometry
	// dimensions must be obeyed by the client. For more details, see
	// xdg_toplevel.set_fullscreen.
	XdgToplevelStateFullscreen = 2
	// The surface is being resized. The window geometry specified in the
	// configure event is a maximum; the client cannot resize beyond it.
	// Clients that have aspect ratio or cell sizing configuration can use
	// a smaller size, however.
	XdgToplevelStateResizing = 3
	// Client window decorations should be painted as if the window is
	// active. Do not assume this means that the window actually has
	// keyboard or pointer focus.
	XdgToplevelStateActivated = 4
	// The window is currently in a tiled layout and the left edge is
	// considered to be adjacent to another part of the tiling grid.
	XdgToplevelStateTiledLeft = 5
	// The window is currently in a tiled layout and the right edge is
	// considered to be adjacent to another part of the tiling grid.
	XdgToplevelStateTiledRight = 6
	// The window is currently in a tiled layout and the top edge is
	// considered to be adjacent to another part of the tiling grid.
	XdgToplevelStateTiledTop = 7
	// The window is currently in a tiled layout and the bottom edge is
	// considered to be adjacent to another part of the tiling grid.
	XdgToplevelStateTiledBottom = 8
)

var xdgToplevelInterface = &wayland.Interface{
	Name:    "xdg_toplevel",
	Version: 2,
	Events:  []wayland.Event{(*XdgToplevelEventConfigure)(nil), (*XdgToplevelEventClose)(nil)},
}

type XdgToplevelEventConfigure struct {
	Width  int32
	Height int32
	States []byte
}

type XdgToplevelEventClose struct {
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
// an xdg_toplevel surface when it is unmapped.
//
// Attaching a null buffer to a toplevel unmaps the surface.
type XdgToplevel struct{ wayland.Proxy }

func (*XdgToplevel) Interface() *wayland.Interface { return xdgToplevelInterface }

// This request destroys the role surface and unmaps the surface;
// see "Unmapping" behavior in interface section for details.
func (obj *XdgToplevel) Destroy() {
	const xdg_toplevel_destroy = 0
	obj.Conn().SendRequest(obj, xdg_toplevel_destroy)
}

// Set the "parent" of this surface. This surface should be stacked
// above the parent surface and all other ancestor surfaces.
//
// Parent windows should be set on dialogs, toolboxes, or other
// "auxiliary" surfaces, so that the parent is raised when the dialog
// is raised.
//
// Setting a null parent for a child window removes any parent-child
// relationship for the child. Setting a null parent for a window which
// currently has no parent is a no-op.
//
// If the parent is unmapped then its children are managed as
// though the parent of the now-unmapped parent has become the
// parent of this surface. If no parent exists for the now-unmapped
// parent then the children are managed as though they have no
// parent surface.
func (obj *XdgToplevel) SetParent(parent *XdgToplevel) {
	const xdg_toplevel_set_parent = 1
	obj.Conn().SendRequest(obj, xdg_toplevel_set_parent, parent)
}

// Set a short title for the surface.
//
// This string may be used to identify the surface in a task bar,
// window list, or other user interface elements provided by the
// compositor.
//
// The string must be encoded in UTF-8.
func (obj *XdgToplevel) SetTitle(title string) {
	const xdg_toplevel_set_title = 2
	obj.Conn().SendRequest(obj, xdg_toplevel_set_title, title)
}

// Set an application identifier for the surface.
//
// The app ID identifies the general class of applications to which
// the surface belongs. The compositor can use this to group multiple
// surfaces together, or to determine how to launch a new application.
//
// For D-Bus activatable applications, the app ID is used as the D-Bus
// service name.
//
// The compositor shell will try to group application surfaces together
// by their app ID. As a best practice, it is suggested to select app
// ID's that match the basename of the application's .desktop file.
// For example, "org.freedesktop.FooViewer" where the .desktop file is
// "org.freedesktop.FooViewer.desktop".
//
// Like other properties, a set_app_id request can be sent after the
// xdg_toplevel has been mapped to update the property.
//
// See the desktop-entry specification [0] for more details on
// application identifiers and how they relate to well-known D-Bus
// names and .desktop files.
//
// [0] http://standards.freedesktop.org/desktop-entry-spec/
func (obj *XdgToplevel) SetAppID(appId string) {
	const xdg_toplevel_set_app_id = 3
	obj.Conn().SendRequest(obj, xdg_toplevel_set_app_id, appId)
}

// Clients implementing client-side decorations might want to show
// a context menu when right-clicking on the decorations, giving the
// user a menu that they can use to maximize or minimize the window.
//
// This request asks the compositor to pop up such a window menu at
// the given position, relative to the local surface coordinates of
// the parent surface. There are no guarantees as to what menu items
// the window menu contains.
//
// This request must be used in response to some sort of user action
// like a button press, key press, or touch down event.
func (obj *XdgToplevel) ShowWindowMenu(seat *Seat, serial uint32, x int32, y int32) {
	const xdg_toplevel_show_window_menu = 4
	obj.Conn().SendRequest(obj, xdg_toplevel_show_window_menu, seat, serial, x, y)
}

// Start an interactive, user-driven move of the surface.
//
// This request must be used in response to some sort of user action
// like a button press, key press, or touch down event. The passed
// serial is used to determine the type of interactive move (touch,
// pointer, etc).
//
// The server may ignore move requests depending on the state of
// the surface (e.g. fullscreen or maximized), or if the passed serial
// is no longer valid.
//
// If triggered, the surface will lose the focus of the device
// (wl_pointer, wl_touch, etc) used for the move. It is up to the
// compositor to visually indicate that the move is taking place, such as
// updating a pointer cursor, during the move. There is no guarantee
// that the device focus will return when the move is completed.
func (obj *XdgToplevel) Move(seat *Seat, serial uint32) {
	const xdg_toplevel_move = 5
	obj.Conn().SendRequest(obj, xdg_toplevel_move, seat, serial)
}

// Start a user-driven, interactive resize of the surface.
//
// This request must be used in response to some sort of user action
// like a button press, key press, or touch down event. The passed
// serial is used to determine the type of interactive resize (touch,
// pointer, etc).
//
// The server may ignore resize requests depending on the state of
// the surface (e.g. fullscreen or maximized).
//
// If triggered, the client will receive configure events with the
// "resize" state enum value and the expected sizes. See the "resize"
// enum value for more details about what is required. The client
// must also acknowledge configure events using "ack_configure". After
// the resize is completed, the client will receive another "configure"
// event without the resize state.
//
// If triggered, the surface also will lose the focus of the device
// (wl_pointer, wl_touch, etc) used for the resize. It is up to the
// compositor to visually indicate that the resize is taking place,
// such as updating a pointer cursor, during the resize. There is no
// guarantee that the device focus will return when the resize is
// completed.
//
// The edges parameter specifies how the surface should be resized,
// and is one of the values of the resize_edge enum. The compositor
// may use this information to update the surface position for
// example when dragging the top left corner. The compositor may also
// use this information to adapt its behavior, e.g. choose an
// appropriate cursor image.
func (obj *XdgToplevel) Resize(seat *Seat, serial uint32, edges uint32) {
	const xdg_toplevel_resize = 6
	obj.Conn().SendRequest(obj, xdg_toplevel_resize, seat, serial, edges)
}

// Set a maximum size for the window.
//
// The client can specify a maximum size so that the compositor does
// not try to configure the window beyond this size.
//
// The width and height arguments are in window geometry coordinates.
// See xdg_surface.set_window_geometry.
//
// Values set in this way are double-buffered. They will get applied
// on the next commit.
//
// The compositor can use this information to allow or disallow
// different states like maximize or fullscreen and draw accurate
// animations.
//
// Similarly, a tiling window manager may use this information to
// place and resize client windows in a more effective way.
//
// The client should not rely on the compositor to obey the maximum
// size. The compositor may decide to ignore the values set by the
// client and request a larger size.
//
// If never set, or a value of zero in the request, means that the
// client has no expected maximum size in the given dimension.
// As a result, a client wishing to reset the maximum size
// to an unspecified state can use zero for width and height in the
// request.
//
// Requesting a maximum size to be smaller than the minimum size of
// a surface is illegal and will result in a protocol error.
//
// The width and height must be greater than or equal to zero. Using
// strictly negative values for width and height will result in a
// protocol error.
func (obj *XdgToplevel) SetMaxSize(width int32, height int32) {
	const xdg_toplevel_set_max_size = 7
	obj.Conn().SendRequest(obj, xdg_toplevel_set_max_size, width, height)
}

// Set a minimum size for the window.
//
// The client can specify a minimum size so that the compositor does
// not try to configure the window below this size.
//
// The width and height arguments are in window geometry coordinates.
// See xdg_surface.set_window_geometry.
//
// Values set in this way are double-buffered. They will get applied
// on the next commit.
//
// The compositor can use this information to allow or disallow
// different states like maximize or fullscreen and draw accurate
// animations.
//
// Similarly, a tiling window manager may use this information to
// place and resize client windows in a more effective way.
//
// The client should not rely on the compositor to obey the minimum
// size. The compositor may decide to ignore the values set by the
// client and request a smaller size.
//
// If never set, or a value of zero in the request, means that the
// client has no expected minimum size in the given dimension.
// As a result, a client wishing to reset the minimum size
// to an unspecified state can use zero for width and height in the
// request.
//
// Requesting a minimum size to be larger than the maximum size of
// a surface is illegal and will result in a protocol error.
//
// The width and height must be greater than or equal to zero. Using
// strictly negative values for width and height will result in a
// protocol error.
func (obj *XdgToplevel) SetMinSize(width int32, height int32) {
	const xdg_toplevel_set_min_size = 8
	obj.Conn().SendRequest(obj, xdg_toplevel_set_min_size, width, height)
}

// Maximize the surface.
//
// After requesting that the surface should be maximized, the compositor
// will respond by emitting a configure event. Whether this configure
// actually sets the window maximized is subject to compositor policies.
// The client must then update its content, drawing in the configured
// state. The client must also acknowledge the configure when committing
// the new content (see ack_configure).
//
// It is up to the compositor to decide how and where to maximize the
// surface, for example which output and what region of the screen should
// be used.
//
// If the surface was already maximized, the compositor will still emit
// a configure event with the "maximized" state.
//
// If the surface is in a fullscreen state, this request has no direct
// effect. It may alter the state the surface is returned to when
// unmaximized unless overridden by the compositor.
func (obj *XdgToplevel) SetMaximized() {
	const xdg_toplevel_set_maximized = 9
	obj.Conn().SendRequest(obj, xdg_toplevel_set_maximized)
}

// Unmaximize the surface.
//
// After requesting that the surface should be unmaximized, the compositor
// will respond by emitting a configure event. Whether this actually
// un-maximizes the window is subject to compositor policies.
// If available and applicable, the compositor will include the window
// geometry dimensions the window had prior to being maximized in the
// configure event. The client must then update its content, drawing it in
// the configured state. The client must also acknowledge the configure
// when committing the new content (see ack_configure).
//
// It is up to the compositor to position the surface after it was
// unmaximized; usually the position the surface had before maximizing, if
// applicable.
//
// If the surface was already not maximized, the compositor will still
// emit a configure event without the "maximized" state.
//
// If the surface is in a fullscreen state, this request has no direct
// effect. It may alter the state the surface is returned to when
// unmaximized unless overridden by the compositor.
func (obj *XdgToplevel) UnsetMaximized() {
	const xdg_toplevel_unset_maximized = 10
	obj.Conn().SendRequest(obj, xdg_toplevel_unset_maximized)
}

// Make the surface fullscreen.
//
// After requesting that the surface should be fullscreened, the
// compositor will respond by emitting a configure event. Whether the
// client is actually put into a fullscreen state is subject to compositor
// policies. The client must also acknowledge the configure when
// committing the new content (see ack_configure).
//
// The output passed by the request indicates the client's preference as
// to which display it should be set fullscreen on. If this value is NULL,
// it's up to the compositor to choose which display will be used to map
// this surface.
//
// If the surface doesn't cover the whole output, the compositor will
// position the surface in the center of the output and compensate with
// with border fill covering the rest of the output. The content of the
// border fill is undefined, but should be assumed to be in some way that
// attempts to blend into the surrounding area (e.g. solid black).
//
// If the fullscreened surface is not opaque, the compositor must make
// sure that other screen content not part of the same surface tree (made
// up of subsurfaces, popups or similarly coupled surfaces) are not
// visible below the fullscreened surface.
func (obj *XdgToplevel) SetFullscreen(output *Output) {
	const xdg_toplevel_set_fullscreen = 11
	obj.Conn().SendRequest(obj, xdg_toplevel_set_fullscreen, output)
}

// Make the surface no longer fullscreen.
//
// After requesting that the surface should be unfullscreened, the
// compositor will respond by emitting a configure event.
// Whether this actually removes the fullscreen state of the client is
// subject to compositor policies.
//
// Making a surface unfullscreen sets states for the surface based on the following:
// * the state(s) it may have had before becoming fullscreen
// * any state(s) decided by the compositor
// * any state(s) requested by the client while the surface was fullscreen
//
// The compositor may include the previous window geometry dimensions in
// the configure event, if applicable.
//
// The client must also acknowledge the configure when committing the new
// content (see ack_configure).
func (obj *XdgToplevel) UnsetFullscreen() {
	const xdg_toplevel_unset_fullscreen = 12
	obj.Conn().SendRequest(obj, xdg_toplevel_unset_fullscreen)
}

// Request that the compositor minimize your surface. There is no
// way to know if the surface is currently minimized, nor is there
// any way to unset minimization on this surface.
//
// If you are looking to throttle redrawing when minimized, please
// instead use the wl_surface.frame event for this, as this will
// also work with live previews on windows in Alt-Tab, Expose or
// similar compositor features.
func (obj *XdgToplevel) SetMinimized() {
	const xdg_toplevel_set_minimized = 13
	obj.Conn().SendRequest(obj, xdg_toplevel_set_minimized)
}

const (
	// tried to grab after being mapped
	XdgPopupErrorInvalidGrab = 0
)

var xdgPopupInterface = &wayland.Interface{
	Name:    "xdg_popup",
	Version: 2,
	Events:  []wayland.Event{(*XdgPopupEventConfigure)(nil), (*XdgPopupEventPopupDone)(nil)},
}

type XdgPopupEventConfigure struct {
	// x position relative to parent surface window geometry
	X int32
	// y position relative to parent surface window geometry
	Y int32
	// window geometry width
	Width int32
	// window geometry height
	Height int32
}

type XdgPopupEventPopupDone struct {
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
// The x and y arguments passed when creating the popup object specify
// where the top left of the popup should be placed, relative to the
// local surface coordinates of the parent surface. See
// xdg_surface.get_popup. An xdg_popup must intersect with or be at least
// partially adjacent to its parent surface.
//
// The client must call wl_surface.commit on the corresponding wl_surface
// for the xdg_popup state to take effect.
type XdgPopup struct{ wayland.Proxy }

func (*XdgPopup) Interface() *wayland.Interface { return xdgPopupInterface }

// This destroys the popup. Explicitly destroying the xdg_popup
// object will also dismiss the popup, and unmap the surface.
//
// If this xdg_popup is not the "topmost" popup, a protocol error
// will be sent.
func (obj *XdgPopup) Destroy() {
	const xdg_popup_destroy = 0
	obj.Conn().SendRequest(obj, xdg_popup_destroy)
}

// This request makes the created popup take an explicit grab. An explicit
// grab will be dismissed when the user dismisses the popup, or when the
// client destroys the xdg_popup. This can be done by the user clicking
// outside the surface, using the keyboard, or even locking the screen
// through closing the lid or a timeout.
//
// If the compositor denies the grab, the popup will be immediately
// dismissed.
//
// This request must be used in response to some sort of user action like a
// button press, key press, or touch down event. The serial number of the
// event should be passed as 'serial'.
//
// The parent of a grabbing popup must either be an xdg_toplevel surface or
// another xdg_popup with an explicit grab. If the parent is another
// xdg_popup it means that the popups are nested, with this popup now being
// the topmost popup.
//
// Nested popups must be destroyed in the reverse order they were created
// in, e.g. the only popup you are allowed to destroy at all times is the
// topmost one.
//
// When compositors choose to dismiss a popup, they may dismiss every
// nested grabbing popup as well. When a compositor dismisses popups, it
// will follow the same dismissing order as required from the client.
//
// The parent of a grabbing popup must either be another xdg_popup with an
// active explicit grab, or an xdg_popup or xdg_toplevel, if there are no
// explicit grabs already taken.
//
// If the topmost grabbing popup is destroyed, the grab will be returned to
// the parent of the popup, if that parent previously had an explicit grab.
//
// If the parent is a grabbing popup which has already been dismissed, this
// popup will be immediately dismissed. If the parent is a popup that did
// not take an explicit grab, an error will be raised.
//
// During a popup grab, the client owning the grab will receive pointer
// and touch events for all their surfaces as normal (similar to an
// "owner-events" grab in X11 parlance), while the top most grabbing popup
// will always have keyboard focus.
func (obj *XdgPopup) Grab(seat *Seat, serial uint32) {
	const xdg_popup_grab = 1
	obj.Conn().SendRequest(obj, xdg_popup_grab, seat, serial)
}
