import "dll"

kernel32 := dll.new("kernel32")
user32 := dll.new("user32.dll")

GetConsoleWindow := kernel32.proc("GetConsoleWindow")
ShowWindowAsync := user32.proc("ShowWindowAsync")
SetWindowPos := user32.proc("SetWindowPos")
SetWindowTextW := user32.proc("SetWindowTextW")
GetWindowTextW := user32.proc("GetWindowTextW")
IsWindowVisible := user32.proc("IsWindowVisible")
CloseWindow := user32.proc("CloseWindow")
SetForegroundWindow := user32.proc("SetForegroundWindow")
FlashWindow := user32.proc("FlashWindow")

// ShowWindow Constants
SW_HIDE := 0
SW_SHOW := 5
SW_SHOWNOACTIVATE := 4
SW_MINIMIZE := 6
SW_MAXIMIZE := 3
SW_RESTORE := 9

// SetWindowPos Constants
HWND_TOPMOST := -1
HWND_NOTOPMOST := -2
SWP_NOMOVE := 0x0002
SWP_NOSIZE := 0x0001

// Helper
hwnd := fn() {
	return GetConsoleWindow.call()
}

export {
	// Visibility
	hide : fn() {
		return ShowWindowAsync.call(hwnd(), SW_HIDE)
	},
	show : fn() {
		return ShowWindowAsync.call(hwnd(), SW_SHOW)
	},
	show_no_activate : fn() {
		return ShowWindowAsync.call(hwnd(), SW_SHOWNOACTIVATE)
	},
	toggle : fn() {
		if IsWindowVisible.call(hwnd()) != 0 {
			return ShowWindowAsync.call(hwnd(), SW_HIDE)
		}
		return ShowWindowAsync.call(hwnd(), SW_SHOW)
	},

	// Window States
	minimize : fn() {
		return ShowWindowAsync.call(hwnd(), SW_MINIMIZE)
	},
	maximize : fn() {
		return ShowWindowAsync.call(hwnd(), SW_MAXIMIZE)
	},
	restore : fn() {
		return ShowWindowAsync.call(hwnd(), SW_RESTORE)
	},

	// Window Title
	set_title : fn(title) {
		return SetWindowTextW.call(hwnd(), title)
	},
	get_title : fn() {
		buffer := bytes(256)
		GetWindowTextW.call(hwnd(), buffer, 256)
		return string(buffer)
	},

	// Always on Top
	set_topmost : fn() {
		return SetWindowPos.call(hwnd(), HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE)
	},
	set_not_topmost : fn() {
		return SetWindowPos.call(hwnd(), HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE)
	},

	// Focus
	set_focus : fn() {
		return SetForegroundWindow.call(hwnd())
	},
	flash : fn() {
		return FlashWindow.call(hwnd(), 1)
	},

	// Status
	is_visible : fn() {
		return IsWindowVisible.call(hwnd()) != 0
	},

	// Close
	close : fn() {
		return CloseWindow.call(hwnd())
	}
}