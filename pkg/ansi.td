export {
	// ====================================================
	// General ASCII Codes
	// ====================================================
	BEL: "\x07",                      // Bell
	BS: "\x08",                       // Backspace
	HT: "\x09",                       // Horizontal Tab
	LF: "\x0A",                       // Line Feed
	VT: "\x0B",                       // Vertical Tab
	FF: "\x0C",                       // Form Feed
	CR: "\x0D",                       // Carriage Return
	ESC: "\x1B",                      // Escape
	DEL: "\x7F",                      // Delete

	// ====================================================
	// Cursor Controls
	// ====================================================
	HOME: "\x1B[H",
	MOVE_TO: fn(line, column) { 
		return "\x1B[" + line + ";" + column + "H" 
	},
	CURSOR_UP: fn(lines) { 
		return "\x1B[" + lines + "A" 
	},
	CURSOR_DOWN: fn(lines) { 
		return "\x1B[" + lines + "B" 
	},
	CURSOR_RIGHT: fn(columns) { 
		return "\x1B[" + columns + "C" 
	},
	CURSOR_LEFT: fn(columns) { 
		return "\x1B[" + columns + "D" 
	},
	NEXT_LINE: fn(lines) { 
		return "\x1B[" + lines + "E" 
	},
	PREVIOUS_LINE: fn(lines) { 
		return "\x1B[" + lines + "F"
	},
	MOVE_TO_COLUMN: fn(column) { 
		return "\x1B[" + column + "G" 
	},
	REPORT_CURSOR_POSITION: "\x1B[6n",
	MOVE_UP: "\x1BM",
	SAVE_CURSOR_POSITION: "\x1B[s",
	RESTORE_CURSOR_POSITION: "\x1B[u",

	// ====================================================
	// Erase Functions
	// ====================================================
	ERASE_DISPLAY: "\x1B[J",
	ERASE_TO_END: "\x1B[0J",
	ERASE_TO_BEGINNING: "\x1B[1J",
	ERASE_ENTIRE_SCREEN: "\x1B[2J",
	ERASE_SAVED_LINES: "\x1B[3J",
	ERASE_LINE: "\x1B[K",
	ERASE_TO_END_OF_LINE: "\x1B[0K",
	ERASE_TO_BEGINNING_OF_LINE: "\x1B[1K",
	ERASE_ENTIRE_LINE: "\x1B[2K",

	// ====================================================
	// Scrolling Region
	// ====================================================
	SET_SCROLL_REGION: fn(top, bottom) { return "\x1B[" + top + ";" + bottom + "r" },
	RESET_SCROLL_REGION: "\x1B[r",

	// ====================================================
	// Tab Stops
	// ====================================================
	SET_TAB_STOP: "\x1BH",
	CLEAR_TAB_STOP: "\x1B[0g",
	CLEAR_ALL_TAB_STOPS: "\x1B[3g",

	// ====================================================
	// Keypad Modes
	// ====================================================
	APPLICATION_KEYPAD: "\x1B=",
	NORMAL_KEYPAD: "\x1B>",

	// ====================================================
	// Graphics Mode / SGR (Text Attributes)
	// ====================================================
	RESET_GRAPHICS_MODE: "\x1B[0m",

	// Style On
	SET_BOLD_MODE: "\x1B[1m",
	SET_DIM_MODE: "\x1B[2m",
	SET_ITALIC_MODE: "\x1B[3m",
	SET_UNDERLINE_MODE: "\x1B[4m",
	SET_BLINKING_MODE: "\x1B[5m",
	SET_INVERSE_MODE: "\x1B[7m",
	SET_HIDDEN_MODE: "\x1B[8m",
	SET_STRIKETHROUGH_MODE: "\x1B[9m",

	// Style Off
	BOLD_OFF: "\x1B[22m",              // Also resets DIM
	ITALIC_OFF: "\x1B[23m",
	UNDERLINE_OFF: "\x1B[24m",
	BLINK_OFF: "\x1B[25m",
	INVERSE_OFF: "\x1B[27m",
	HIDDEN_OFF: "\x1B[28m",
	STRIKE_OFF: "\x1B[29m",

	// Modern Underline Styles
	UNDERLINE_SINGLE: "\x1B[4m",
	UNDERLINE_DOUBLE: "\x1B[21m",
	UNDERLINE_CURLY: "\x1B[4:3m",
	UNDERLINE_DOTTED: "\x1B[4:4m",
	UNDERLINE_DASHED: "\x1B[4:5m",

	// Extra Attributes
	OVERLINE: "\x1B[53m",
	OVERLINE_OFF: "\x1B[55m",
	FRAME: "\x1B[51m",
	ENCIRCLE: "\x1B[52m",

	// ====================================================
	// Standard 16 Colors (Foreground)
	// ====================================================
	FG_BLACK: "\x1B[30m",
	FG_RED: "\x1B[31m",
	FG_GREEN: "\x1B[32m",
	FG_YELLOW: "\x1B[33m",
	FG_BLUE: "\x1B[34m",
	FG_MAGENTA: "\x1B[35m",
	FG_CYAN: "\x1B[36m",
	FG_WHITE: "\x1B[37m",
	FG_DEFAULT: "\x1B[39m",

	FG_BRIGHT_BLACK: "\x1B[90m",
	FG_BRIGHT_RED: "\x1B[91m",
	FG_BRIGHT_GREEN: "\x1B[92m",
	FG_BRIGHT_YELLOW: "\x1B[93m",
	FG_BRIGHT_BLUE: "\x1B[94m",
	FG_BRIGHT_MAGENTA: "\x1B[95m",
	FG_BRIGHT_CYAN: "\x1B[96m",
	FG_BRIGHT_WHITE: "\x1B[97m",

	// ====================================================
	// Standard 16 Colors (Background)
	// ====================================================
	BG_BLACK: "\x1B[40m",
	BG_RED: "\x1B[41m",
	BG_GREEN: "\x1B[42m",
	BG_YELLOW: "\x1B[43m",
	BG_BLUE: "\x1B[44m",
	BG_MAGENTA: "\x1B[45m",
	BG_CYAN: "\x1B[46m",
	BG_WHITE: "\x1B[47m",
	BG_DEFAULT: "\x1B[49m",

	BG_BRIGHT_BLACK: "\x1B[100m",
	BG_BRIGHT_RED: "\x1B[101m",
	BG_BRIGHT_GREEN: "\x1B[102m",
	BG_BRIGHT_YELLOW: "\x1B[103m",
	BG_BRIGHT_BLUE: "\x1B[104m",
	BG_BRIGHT_MAGENTA: "\x1B[105m",
	BG_BRIGHT_CYAN: "\x1B[106m",
	BG_BRIGHT_WHITE: "\x1B[107m",

	// ====================================================
	// 256 Colors
	// ====================================================
	SET_COLOR: fn(foreground, background) {
		return "\x1B[38;5;" + foreground + "m\x1B[48;5;" + background + "m"
	},
	SET_FG_256: fn(color) { 
		return "\x1B[38;5;" + color + "m" 
	},
	SET_BG_256: fn(color) { 
		return "\x1B[48;5;" + color + "m" 
	},

	// ====================================================
	// 24-bit Truecolor
	// ====================================================
	SET_TRUE_COLOR_FG: fn(r, g, b) {
		return "\x1B[38;2;" + r + ";" + g + ";" + b + "m"
	},
	SET_TRUE_COLOR_BG: fn(r, g, b) {
		return "\x1B[48;2;" + r + ";" + g + ";" + b + "m"
	},
	SET_TRUE_COLOR: fn(fr, fg, fb, br, bg, bb) {
		return "\x1B[38;2;" + fr + ";" + fg + ";" + fb +
		       "m\x1B[48;2;" + br + ";" + bg + ";" + bb + "m"
	},
	RESET_COLOR: "\x1B[39;49m",         // Reset FG & BG to default

	// ====================================================
	// Cursor Style (Shape)
	// ====================================================
	SET_CURSOR_STYLE: fn(style) { 
		return "\x1B[" + style + " q" 
	},
	CURSOR_BLINKING_BLOCK: "\x1B[1 q",
	CURSOR_STEADY_BLOCK: "\x1B[2 q",
	CURSOR_BLINKING_UNDERLINE: "\x1B[3 q",
	CURSOR_STEADY_UNDERLINE: "\x1B[4 q",
	CURSOR_BLINKING_BAR: "\x1B[5 q",
	CURSOR_STEADY_BAR: "\x1B[6 q",

	// ====================================================
	// Screen Modes (Private)
	// ====================================================
	SET_SCREEN_MODE: fn(value) { 
		return "\x1B[=" + value + "h" 
	},
	RESET_SCREEN_MODE: fn(value) { 
		return "\x1B[=" + value + "l" 
	},

	CURSOR_INVISIBLE: "\x1B[?25l",
	CURSOR_VISIBLE: "\x1B[?25h",
	RESTORE_SCREEN: "\x1B[?47l",
	SAVE_SCREEN: "\x1B[?47h",
	ENABLE_ALTERNATIVE_BUFFER: "\x1B[?1049h",
	DISABLE_ALTERNATIVE_BUFFER: "\x1B[?1049l",

	// ====================================================
	// Mouse Tracking
	// ====================================================
	MOUSE_OFF: "\x1B[?1000l\x1B[?1002l\x1B[?1003l\x1B[?1006l",
	MOUSE_ON_X10: "\x1B[?1000h\x1B[?1006h",       // Basic click only
	MOUSE_ON_BUTTON: "\x1B[?1002h\x1B[?1006h",    // Click + drag
	MOUSE_ON_MOTION: "\x1B[?1003h\x1B[?1006h",    // All mouse movement

	// ====================================================
	// Device Status Reports (DSR)
	// ====================================================
	REPORT_DEVICE_STATUS: "\x1B[5n",
	REPORT_DEVICE_ATTRIBUTES: "\x1B[c",

	// ====================================================
	// OSC (Operating System Commands)
	// ====================================================
	// Window / Tab Titles
	SET_WINDOW_TITLE: fn(title) { 
		return "\x1B]0;" + title + "\x07" // Icon + Title
	},
	SET_TAB_TITLE: fn(title) { 
		return "\x1B]1;" + title + "\x07" // Tab only
	},       
	SET_TITLE: fn(title) { 
		return "\x1B]2;" + title + "\x07" // Window title
	},           

	// Clipboard (OSC 52)
	// selection: 'c' = clipboard, 'p' = primary, 's' = secondary
	// data must be base64 encoded
	SET_CLIPBOARD: fn(selection, base64_data) {
		return "\x1B]52;" + selection + ";" + base64_data + "\x07"
	},

	// Hyperlinks (OSC 8)
	HYPERLINK: fn(id, url) {
		return "\x1B]8;id=" + id + ";" + url + "\x07"
	},
	HYPERLINK_END: "\x1B]8;;\x07",

	// ====================================================
	// Keyboard Strings
	// ====================================================
	REDEFINE_KEY: fn(code, string) { return "\x1B[" + code + ";" + string + "p" }
}