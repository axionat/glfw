package glfw

//#cgo CFLAGS: -D_cgo=1
//#define GLFW_EXPOSE_NATIVE_WIN32
//#define GLFW_EXPOSE_NATIVE_WGL
//#define GLFW_INCLUDE_NONE
//#include "glfw/include/GLFW/glfw3.h"
//#include "glfw/include/GLFW/glfw3native.h"
//float getDPIScale(HWND handle);
//BOOL appendSeparator(HMENU handle);
//BOOL appendMenu(HMENU handle, int code, const char *title);
//BOOL appendPopup(HMENU handle, HMENU submenu, const char *title);
//BOOL showAndDestroyContextualMenu(HMENU menuHandle, HWND windowHandle, long x, long y);
//BOOL destroyMenu(HMENU handle);
//void showMessageBox(const char *caption, const char *message);
//void showToolsWindow(HWND parent);
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
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

// setMenuBar sets menu and all subMenus as menu bar menus
// this is so we know to redraw after enabled or checked status is changed for items
func (menu *Menu) setMenuBar() {
	menu.menuBar = true
	for _, entry := range menu.entries {
		if subMenu, ok := entry.(*SubMenu); ok {
			subMenu.menuBar = true
			subMenu.Menu.setMenuBar()
		}
	}
}

// SetMainMenu for this window
func (w *Window) SetMainMenu(menu *Menu) {
	// mark this menu tree as a menu bar
	menu.setMenuBar()

	ret := C.glfwGetWin32Window(w.data)
	C.SetMenu(ret, menu.handle)

	// destroy any previous menu
	if w.menu != nil {
		w.menu.Destroy()
	}
	w.menu = menu
}

// GetMods manually queries and returns the current key modifiers for the given window
func (w *Window) GetMods() (mods ModifierKey) {
	if w == nil {
		return
	}
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

// GetDPIScale returns the DPI scaling in use for the window
func (w *Window) GetDPIScale() float32 {
	return float32(C.getDPIScale(w.GetWin32Window()))
}

//export goMenuCallback
func goMenuCallback(w *C.GLFWwindow, code C.int) {
	window := windows.get(w)
	if code == 0 || // default no callback code
		window == nil { // no window to look for callbacks in
		return
	}
	err := window.callbacks.execute(window, code)
	if err != nil {
		fmt.Printf("\n**************\nCallback not found: %3d  %s\n*****************\n",
			int(code), err.Error())
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
	handle  C.HMENU
	window  *Window
	entries []interface{}
	// set when menu is drawn as menu bar, used to determine whether to redraw when checked or enabled is changed on items
	menuBar bool
}

// GetEntries returns a slice of menu items, each entry will be one of:
//     *MenuItem
//     *SubMenu
//     nil (for separator)
func (menu *Menu) GetEntries() (entries []interface{}) {
	entries = make([]interface{}, len(menu.entries))
	copy(entries, menu.entries)
	return entries
}

// Destroy this menu and clean up resources associated with any callbacks.
func (menu *Menu) Destroy() {
	if menu == nil {
		return
	}

	if menu.window != nil {
		for _, entry := range menu.entries {
			switch entry := entry.(type) {
			case *MenuItem:
				menu.window.callbacks.Lock()
				delete(menu.window.callbacks.callbackMap, entry.code)
				menu.window.callbacks.Unlock()
				C.destroyMenu(entry.menu.handle)
			case *SubMenu:
				entry.Menu.Destroy()
				C.destroyMenu(entry.handle)
			}
		}
	}

	// after entries are destroyed, remove them from list
	// this makes a Destroy call safe to repeat
	menu.entries = nil

	C.destroyMenu(menu.handle)
}

// NewMenu returns a new menu ready for appending items and submenus
func NewMenu(w *Window) *Menu {
	if w == nil {
		return nil
	}
	return &Menu{
		handle: C.CreateMenu(),
		window: w,
	}
}

// NewContextualMenu returns a new popup menu ready for items and submenus to be added.
// After menu is assembled, it should be added as a general window contextual menu
// via SetcontextualCallback to be automatically shown on any right-click in the window
// or directly executed via Popup method
func NewContextualMenu(w *Window) *Menu {
	if w == nil {
		return nil
	}
	return &Menu{
		handle: C.CreatePopupMenu(),
		window: w,
	}
}

func (menu *Menu) showAndDestroy(x, y C.long) {
	C.showAndDestroyContextualMenu(menu.handle, menu.window.GetWin32Window(), x, y)
	// delay the destruction since callbacks may still be arriving based on this menu
	go func() {
		time.Sleep(time.Second)
		menu.Destroy()
	}()
}

// Popup displays the menu as a popup menu from the current cursor screen position.
// Menu must have been created via NewContextualMenu call.
// Menu will be destroyed after any action or click.
func (menu *Menu) Popup() {
	w := menu.window
	xpos, ypos := w.GetCursorPos()
	x, y := w.GetPos()
	x += int(xpos)
	y += int(ypos)
	menu.showAndDestroy(C.long(x), C.long(y))
}

// MenuItem is an item to be added to a menu, create via NewMenuItem or CoupledMenuItem
type MenuItem struct {
	title    string
	callback interface{} // func(), func(*Window), or func(*Window, ModifierKey)
	menu     *Menu
	code     C.int
	checked  bool // used to hold checked set prior to item being added to menu
	enabled  bool // used to hold enabled state prior to item being added to menuy
}

// Title returns the title of the menu item
func (mi *MenuItem) Title() string {
	if mi == nil {
		return "nil"
	}
	return mi.title
}

// Execute directly performs any callback action set for the menu item.
//
// If mods is provided, it will be used in place of current window modifiers.
//
// Retured error will be:
//    nil if callback performed,
//    ErrNoCallback if no callback is defined for this item
//    an error describing type mismatch error if callback is not a supported functiontype
func (mi *MenuItem) Execute(mods ...ModifierKey) error {
	if mi == nil || mi.callback == nil {
		return nil
	}
	var window *Window
	if mi.menu != nil {
		window = mi.menu.window
	}
	return doCallback(mi.callback, window, mods...)
}

// CoupledMenuItem returns a menu item coupled to a bool at a provided location.
//
// The menu items checked status will be set by and follow the boolean.  If the
// boolean value is changed other than by menu action, the checked state will
// be otu of sync until the next menu action.
func CoupledMenuItem(title string, target *bool) (item *MenuItem) {
	item = NewMenuItem(title, func() {
		*target = !*target
		item.SetChecked(*target)
	})
	item.checked = *target
	return item
}

// NewMenuItem returns a new menu item with an optional callback function.
// Callback, ir provied, must be one of the following types:
//     func()
//     func(*Window)
//     func(*Window, ModifierKey)
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
		title:    title,
		callback: callback,
		enabled:  true,
	}
}

// SetChecked adjusts the menu items checked status
func (mi *MenuItem) SetChecked(chk bool) {
	if mi.checked == chk {
		// no change, so no action required
		return
	}
	// update checked state
	mi.checked = chk
	if mi.menu == nil {
		// return, can not update not drawn item from C calls,
		// local state will update state when drawn
		return
	}
	// update status in C domain
	var status uint32 // C.MF_UNCKECKED == 0x0
	if chk {
		status = 0x8 // C.MF_CHECKED == 0x8
	}
	C.CheckMenuItem(mi.menu.handle, C.uint(mi.code), C.uint(status))
	// redraw if this is part of an active menu bar
	if mi.menu.menuBar {
		C.DrawMenuBar(mi.menu.window.GetWin32Window())
	}
}

// SetEnabled adjusts the menu items enabled vs disabled / grayed out status
func (mi *MenuItem) SetEnabled(enabled bool) {
	if enabled == mi.enabled {
		return
	}
	mi.enabled = enabled
	if mi.menu == nil { // item not added to menu yet
		// return, can not update not drawn item from C calls,
		// local state will update state when drawn
		return
	}
	// update status in C domain
	var status uint32 // C.MF_ENABLED == 0x0
	if !enabled {
		status = 0x1 // C.MF_GRAYED == 0x1
	}

	C.EnableMenuItem(mi.menu.handle, C.uint(mi.code), C.uint(status))
	// redraw if this is part of an active menu bar
	if mi.menu.menuBar {
		C.DrawMenuBar(mi.menu.window.GetWin32Window())
	}
}

// SubMenu struct
type SubMenu struct {
	*Menu
	Title string
}

// NewSubMenu constructor
func NewSubMenu(w *Window, title string) *SubMenu {
	if w == nil {
		// we have to have a window to create a SubMenu
		return nil
	}
	return &SubMenu{
		Menu:  NewMenu(w),
		Title: title,
	}
}

// AppendSeparator to this menu
func (menu *Menu) AppendSeparator() {
	menu.entries = append(menu.entries, nil)
	C.appendSeparator(menu.handle)
}

// AppendMenuItem to this menu
func (menu *Menu) AppendMenuItem(menuItem *MenuItem) {
	menu.entries = append(menu.entries, menuItem)
	var code C.int
	if menuItem.callback != nil {
		code = menu.window.callbacks.register(menuItem.callback)
		menuItem.code = code
	}
	menuItem.menu = menu

	title := C.CString(menuItem.title)
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
	menu.entries = append(menu.entries, subMenu)
	if subMenu.window != menu.window {
		panic("window does not match for menu and submenu: " + subMenu.Title)
	}

	title := C.CString(subMenu.Title)
	defer C.free(unsafe.Pointer(title))

	C.appendPopup(menu.handle, subMenu.handle, title)
}

type callbackRegistry struct {
	sync.Mutex
	callbackMap map[C.int]interface{}
}

// ErrNoCallback is error returned when trying to execute a callback on an item with no
// associated callback defined
var ErrNoCallback = errors.New("No callback defined for this item")

// last unique menu item handle
var lastCode = int32(13)

func (registry *callbackRegistry) execute(w *Window, code C.int) error {
	if registry == nil {
		return ErrNoCallback
	}
	registry.Lock()
	defer registry.Unlock()
	callback := registry.callbackMap[code]
	return doCallback(callback, w)
}

// doCallback executes a callback based on the type found.
// If mods is not provided and is required for the callback, current modifiers
// for the provided window are obtained and used.
func doCallback(callback interface{}, w *Window, mods ...ModifierKey) error {
	if callback == nil {
		return ErrNoCallback
	}
	switch callback := callback.(type) {
	case func():
		callback()
		return nil
	case func(*Window):
		callback(w)
		return nil
	case func(*Window, ModifierKey):
		if len(mods) > 0 {
			callback(w, mods[0])
		} else {
			callback(w, w.GetMods())
		}
		return nil
	default:
		return fmt.Errorf("Unable to perform unsupported callback function: %T   [allowed functions: func(), func(*glfw.Window), and func(*glfw.Window, glfw.ModifierKey) ]", callback)
	}
}

func (registry *callbackRegistry) register(callback interface{}) C.int {
	if callback == nil {
		return 0
	}
	// verify callback is a supported type
	switch callback := callback.(type) {
	case func():
	case func(*Window):
	case func(*Window, ModifierKey):
	default:
		panic(fmt.Sprintf("Unable to register menu callback with unsupported type (func(), func(*glfw.Window), and func(*glfw.Window, glfw.ModifierKey) allowed): %T", callback))
	}

	registry.Lock()
	defer registry.Unlock()

	// initialize if needed
	if registry.callbackMap == nil {
		registry.callbackMap = make(map[C.int]interface{})
	}

	code := C.int(atomic.AddInt32(&lastCode, 1))
	registry.callbackMap[code] = callback
	return code
}

// ShowMessageBox with a simple OK button
func ShowMessageBox(caption, message string) {
	c := C.CString(caption)
	defer C.free(unsafe.Pointer(c))

	m := C.CString(message)
	defer C.free(unsafe.Pointer(m))

	C.showMessageBox(c, m)
}

func (w *Window) ShowToolsWindow() {
	C.showToolsWindow(w.GetWin32Window())
}
