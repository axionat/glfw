package glfw

//#cgo CFLAGS: -D_cgo=1
//#define GLFW_EXPOSE_NATIVE_WIN32
//#define GLFW_EXPOSE_NATIVE_WGL
//#define GLFW_INCLUDE_NONE
//#include "glfw/include/GLFW/glfw3.h"
//#include "glfw/include/GLFW/glfw3native.h"
//getMonitorPixelCount(HWND handle);
//BOOL appendSeparator(HMENU handle);
//BOOL appendMenu(HMENU handle, int code, const char *title);
//BOOL appendPopup(HMENU handle, HMENU submenu, const char *title);
//BOOL showAndDestroyContextualMenu(HMENU menuHandle, HWND windowHandle, long x, long y);
//BOOL destroyMenu(HMENU handle);
import "C"
import (
	"fmt"
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

// GetMods manually queries and returns the current key modifiers for the given window
func (w *Window) GetMods() (mods ModifierKey) {
	if w.GetKey(KeyLeftShift) == Press ||
		w.GetKey(KeyRightShift) == Press {
		mods += ModShift
	}
	if w.GetKey(KeyLeftControl) == Press ||
		w.GetKey(KeyRightControl) == Press {
		mods += ModControl
	}
	if w.GetKey(KeyLeftAlt) == Press ||
		w.GetKey(KeyRightAlt) == Press {
		mods += ModAlt
	}
	if w.GetKey(KeyLeftSuper) == Press ||
		w.GetKey(KeyRightSuper) == Press {
		mods += ModSuper
	}
	return mods
}

// GetMonitorScale receiver method
func (w *Window) GetMonitorScale() int {
	pixels := int(C.getMonitorPixelCount(w.GetWin32Window()))
	width, height := w.GetMonitor().GetPhysicalSize()
	return pixels / (width * height)
}

//export goMenuCallback
func goMenuCallback(w *C.GLFWwindow, code C.int) {
	window := windows.get(w)
	if callback := registry.menuCallbackMap[int(code)]; callback != nil {
		switch callback := callback.(type) {
		case func():
			callback()
		case func(*Window):
			callback(window)
		case func(*Window, ModifierKey):
			callback(window, window.GetMods())
		}
	}
}

//export goContextualMenuCallback
func goContextualMenuCallback(w *C.GLFWwindow, x, y C.long) bool {
	// returns true if a contextual menu was created
	window := windows.get(w)
	contextual := window.fContextualHolder
	if contextual == nil {
		// no contextual handler, pass false to click
		// will fall through to general mouse button handler
		return false
	}
	if menu := contextual(window); menu != nil {
		// an actual menu was created
		menu.showAndDestroy(x, y)
	}
	// true since we called a dedicated contexual handler
	return true
}

// Menu struct
type Menu struct {
	handle C.HMENU
	window *Window
}

// NewMenu constructor
func NewMenu(w *Window) *Menu {
	return &Menu{
		handle: C.CreateMenu(),
		window: w,
	}
}

// NewContextualMenu constructor
func NewContextualMenu(w *Window) *Menu {
	return &Menu{
		handle: C.CreatePopupMenu(),
		window: w,
	}
}

func (menu *Menu) showAndDestroy(x, y C.long) {
	C.showAndDestroyContextualMenu(menu.handle, menu.window.GetWin32Window(), x, y)
}

// ShowAndDestroy displays the menu then destroys it after any action or click
// outside the menu
func (menu *Menu) ShowAndDestroy() {
	w := menu.window
	xpos, ypos := w.GetCursorPos()
	x, y := w.GetPos()
	x += int(xpos)
	y += int(ypos)
	C.showAndDestroyContextualMenu(menu.handle, menu.window.GetWin32Window(),
		C.long(x), C.long(y))
}

// MenuItem struct
type MenuItem struct {
	Title    string
	Callback interface{} // func() func(*Window), or func(*Window, ModifierKey)
	menu     *Menu
	code     C.int
	checked  bool // used to hold checked set prior to item being added to menu
	enabled  bool // used to hold enabled state prior to item being added to menuy
}

// CoupledMenuItem returns a menu item coupled to a bool at a provided location
// checked state will follow the boolean; however, if value is changed separate
// from menu action the checked state can get out of sync
func CoupledMenuItem(title string, target *bool) (item *MenuItem) {
	item = NewMenuItem(title, func() {
		*target = !*target
		item.SetChecked(*target)
	})
	item.checked = *target
	return item
}

// NewMenuItem constructor
func NewMenuItem(title string, callback interface{}) *MenuItem {
	if callback != nil {
		// verify callback is a supported type
		switch callback := callback.(type) {
		case func():
		case func(*Window):
		case func(*Window, ModifierKey):
		default:
			panic(fmt.Sprintf("Unable to create menu item with unsupported type for callback: %T", callback))
		}
	}
	return &MenuItem{
		Title:    title,
		Callback: callback,
		enabled:  true,
	}
}

// SetChecked adjusts the menu items checked status
func (mi *MenuItem) SetChecked(chk bool) {
	if mi.menu == nil {
		// item not added to a menu yet, keep this state which will be applied
		// once item is added to menu
		mi.checked = chk
		return
	}
	if mi.checked == chk {
		return
	}
	mi.checked = chk
	if chk {
		C.CheckMenuItem(mi.menu.handle, C.uint(mi.code), 0x8) // MF_CHECKED == 0x8
	} else {
		C.CheckMenuItem(mi.menu.handle, C.uint(mi.code), 0x0) // C.MF_UNCKECKED == 0x0
	}
	C.DrawMenuBar(mi.menu.window.GetWin32Window())
}

// SetEnabled adjusts the menu items enabled vs disabled / grayed out status
func (mi *MenuItem) SetEnabled(enabled bool) {
	if mi.menu == nil {
		// item not added to a menu yet, keep this state which will be applied
		// once item is added to menu
		mi.enabled = enabled
		return
	}
	if enabled == mi.enabled {
		return
	}
	mi.enabled = enabled
	if enabled {
		C.EnableMenuItem(mi.menu.handle, C.uint(mi.code), 0x0) // C.MF_ENABLED == 0
	} else {
		C.EnableMenuItem(mi.menu.handle, C.uint(mi.code), 0x1) // C.MF_GRAYED == 0x1
	}
	C.DrawMenuBar(mi.menu.window.GetWin32Window())
}

// SubMenu struct
type SubMenu struct {
	*Menu
	Title string
}

// NewSubMenu constructor
func NewSubMenu(w *Window, title string) *SubMenu {
	return &SubMenu{
		Menu:  NewMenu(w),
		Title: title,
	}
}

// Destroy this menu. Not needed when menu is attached to a window via `SetMenu`.
func (menu *Menu) Destroy() {
	C.destroyMenu(menu.handle)
}

// AppendSeparator to this menu
func (menu *Menu) AppendSeparator() {
	C.appendSeparator(menu.handle)
}

// AppendMenuItem to this menu
func (menu *Menu) AppendMenuItem(menuItem *MenuItem) {
	code := registry.register(menuItem.Callback)
	menuItem.code = code
	menuItem.menu = menu

	title := C.CString(menuItem.Title)
	defer C.free(unsafe.Pointer(title))

	C.appendMenu(menu.handle, code, title)

	// use any pre-append item states that have been set
	// if checked or disabled then update
	if menuItem.checked {
		C.CheckMenuItem(menu.handle, C.uint(code), 0x8) // MF_CHECKED == 0x8
	}
	if !menuItem.enabled {
		C.EnableMenuItem(menu.handle, C.uint(code), 0x1) // C.MF_GRAYED == 0x1
	}
}

// AppendSubMenu to this menu
func (menu *Menu) AppendSubMenu(subMenu *SubMenu) {
	if subMenu.window != menu.window {
		panic("window does not match for menu and submenu: " + subMenu.Title)
	}

	title := C.CString(subMenu.Title)
	defer C.free(unsafe.Pointer(title))

	C.appendPopup(menu.handle, subMenu.handle, title)
}

var registry = newCallbackRegistry()

type callbackRegistry struct {
	sync.Mutex
	nextCode        int
	menuCallbackMap map[int]interface{}
}

func newCallbackRegistry() *callbackRegistry {
	return &callbackRegistry{
		nextCode:        13,
		menuCallbackMap: make(map[int]interface{}),
	}
}

func (registry *callbackRegistry) register(callback interface{}) C.int {
	registry.Lock()
	defer registry.Unlock()

	code := registry.nextCode
	registry.nextCode++
	registry.menuCallbackMap[code] = callback
	return C.int(code)
}
