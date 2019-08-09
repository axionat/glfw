package glfw

//#cgo CFLAGS: -D_cgo=1
//#define GLFW_EXPOSE_NATIVE_WIN32
//#define GLFW_EXPOSE_NATIVE_WGL
//#define GLFW_INCLUDE_NONE
//#include "glfw/include/GLFW/glfw3.h"
//#include "glfw/include/GLFW/glfw3native.h"
//BOOL appendSeparator(HMENU handle);
//BOOL appendMenu(HMENU handle, int code, const char *title);
//BOOL appendPopup(HMENU handle, HMENU submenu, const char *title);
import "C"
import (
	"sync"
	"unsafe"
)

// GetWin32Adapter returns the adapter device name of the monitor.
func (m *Monitor) GetWin32Adapter() string {
	ret := C.glfwGetWin32Adapter(m.data)
	panicError()
	return C.GoString(ret)
}

// GetWin32Monitor returns the display device name of the monitor.
func (m *Monitor) GetWin32Monitor() string {
	ret := C.glfwGetWin32Monitor(m.data)
	panicError()
	return C.GoString(ret)
}

// GetWin32Window returns the HWND of the window.
func (w *Window) GetWin32Window() C.HWND {
	ret := C.glfwGetWin32Window(w.data)
	panicError()
	return ret
}

// GetWGLContext returns the HGLRC of the window.
func (w *Window) GetWGLContext() C.HGLRC {
	ret := C.glfwGetWGLContext(w.data)
	panicError()
	return ret
}

// SetMainMenu for this window
func (w *Window) SetMainMenu(menu *Menu) {
	ret := C.glfwGetWin32Window(w.data)
	C.SetMenu(ret, menu.handle)
}

//export goMenuCallback
func goMenuCallback(code C.int) {
	if callback := registry.menuCallbackMap[int(code)]; callback != nil {
		callback()
	}
}

// Menu struct
type Menu struct {
	handle C.HMENU
}

// NewMenu constructor
func NewMenu() *Menu {
	return &Menu{
		handle: C.CreateMenu(),
	}
}

// MenuItem struct
type MenuItem struct {
	Title    string
	Callback func()
}

// NewMenuItem constructor
func NewMenuItem(title string, callback func()) *MenuItem {
	return &MenuItem{
		Title:    title,
		Callback: callback,
	}
}

// SubMenu struct
type SubMenu struct {
	*Menu
	Title string
}

// NewSubMenu constructor
func NewSubMenu(title string) *SubMenu {
	return &SubMenu{
		Menu:  NewMenu(),
		Title: title,
	}
}

// AppendSeparator to this menu
func (menu *Menu) AppendSeparator() {
	C.appendSeparator(menu.handle)
}

// AppendMenuItem to this menu
func (menu *Menu) AppendMenuItem(menuItem *MenuItem) {
	title := C.CString(menuItem.Title)
	defer C.free(unsafe.Pointer(title))

	code := registry.register(menuItem.Callback)
	C.appendMenu(menu.handle, code, title)
}

// AppendSubMenu to this menu
func (menu *Menu) AppendSubMenu(subMenu *SubMenu) {
	title := C.CString(subMenu.Title)
	defer C.free(unsafe.Pointer(title))

	C.appendPopup(menu.handle, subMenu.handle, title)
}

var registry = newCalbackRegistry()

type callbackRegistry struct {
	sync.Mutex
	nextCode        int
	menuCallbackMap map[int]func()
}

func newCalbackRegistry() *callbackRegistry {
	return &callbackRegistry{
		nextCode:        13,
		menuCallbackMap: make(map[int]func()),
	}
}

func (registry *callbackRegistry) register(callback func()) C.int {
	registry.Lock()
	defer registry.Unlock()

	code := registry.nextCode
	registry.nextCode++
	registry.menuCallbackMap[code] = callback
	return C.int(code)
}
