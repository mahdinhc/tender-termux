# WUI Standard Library for Tender

The `wui` module provides native Windows desktop application development capabilities for Tender, offering a comprehensive set of GUI controls, window management, and dialog functionality.

## Overview

The WUI library is a Windows-only module that provides:
- Native Windows GUI controls
- Event-driven programming model
- Custom drawing with canvas API
- File and folder dialogs
- Menu system
- Resource management (fonts, icons, colors, cursors)

## Module Functions

### Window & Dialog Creation

#### `new_window()`
Creates a new window object.

**Returns:** `window` object

**Example:**
```go
win := wui.new_window()
win.set_title("My Application")
win.set_size(400, 400)
win.show()
```

#### `new_file_open_dialog()`
Creates a new file open dialog object.

**Returns:** `fileopendialog` object

**Example:**
```go
dlg := wui.new_file_open_dialog()
dlg.add_filter("Text Files", "txt", "csv")
ok, file := dlg.execute_single_selection(win)
if ok {
    // Process file
}
```

#### `new_file_save_dialog()`
Creates a new file save dialog object.

**Returns:** `filesavedialog` object

#### `new_folder_select_dialog()`
Creates a new folder selection dialog object.

**Returns:** `folderselectdialog` object

### Control Creation

All control creation functions return a control object with methods for manipulation and event handling.

| Function | Description | Returns |
|----------|-------------|---------|
| `new_button()` | Creates a button control | `button` |
| `new_checkbox()` | Creates a checkbox | `checkbox` |
| `new_label()` | Creates a text label | `label` |
| `new_editline()` | Creates a single-line text input | `editline` |
| `new_textedit()` | Creates a multi-line text input | `textedit` |
| `new_combobox()` | Creates a dropdown list | `combobox` |
| `new_stringlist()` | Creates a list box for strings | `stringlist` |
| `new_stringtable(header1, ...)` | Creates a table control | `stringtable` |
| `new_slider()` | Creates a slider control | `slider` |
| `new_progressbar()` | Creates a progress bar | `progressbar` |
| `new_radiobutton()` | Creates a radio button | `radiobutton` |
| `new_intupdown()` | Creates integer up-down control | `intupdown` |
| `new_floatupdown()` | Creates float up-down control | `floatupdown` |
| `new_panel()` | Creates a panel container | `panel` |
| `new_paintbox()` | Creates a custom drawing area | `paintbox` |

### Menu Creation

#### `new_menu(name)`
Creates a new menu or submenu.

**Parameters:**
- `name`: The menu header text

**Returns:** `menu` object

#### `new_main_menu()`
Creates a main menu bar for a window.

**Returns:** `menu` object

#### `new_menu_string(text)`
Creates a clickable menu item.

**Parameters:**
- `text`: The menu item text

**Returns:** `menustring` object

#### `new_menu_separator()`
Creates a menu separator line.

**Returns:** `menuitem` object

### Resource Creation

#### `rgb_color(r, g, b)`
Creates a color object from RGB components.

**Parameters:**
- `r`: Red component (0-255)
- `g`: Green component (0-255)  
- `b`: Blue component (0-255)

**Returns:** `color` object

**Example:**
```go
red = wui.rgb_color(255, 0, 0)
```

#### `new_font(desc)`
Creates a font object from a description map.

**Parameters:**
- `desc`: A map containing font properties

**Properties:**
- `name`: Font family name
- `height`: Font height in pixels
- `bold`: Boolean
- `italic`: Boolean
- `underlined`: Boolean
- `striked_out`: Boolean

**Returns:** `font` object or `error`

**Example:**
```go
font := wui.new_font({
    name: "Arial",
    height: 12,
    bold: true
})
```

#### `new_cursor_from_image(image_bytes, x, y)`
Creates a cursor from image data.

**Parameters:**
- `image_bytes`: Byte array of PNG/JPEG/other image data
- `x`: Hotspot X coordinate
- `y`: Hotspot Y coordinate

**Returns:** `cursor` object or `error`

#### `new_icon_from_image(image_bytes)`
Creates an icon from image data.

**Parameters:**
- `image_bytes`: Byte array of PNG/JPEG/other image data

**Returns:** `icon` object or `error`

#### `new_image(data)`
Creates an image object from image data (for canvas drawing).

**Parameters:**
- `data`: Byte array of PNG/JPEG/other image data

**Returns:** `image` object or `error`

#### `new_image_from_hbitmap(bitmap, width, height)`
Creates an image from a Windows HBITMAP handle.

**Parameters:**
- `bitmap`: Integer handle
- `width`: Image width
- `height`: Image height

**Returns:** `image` object

#### `new_icon_from_exe_resource(res_id)`
Creates an icon from an EXE resource.

**Parameters:**
- `res_id`: Resource ID

**Returns:** `icon` object or `error`

#### `new_icon_from_file(path)`
Creates an icon from a file.

**Parameters:**
- `path`: File path

**Returns:** `icon` object or `error`

#### `new_icon_from_reader(data)`
Creates an icon from byte data.

**Parameters:**
- `data`: Byte array

**Returns:** `icon` object or `error`

#### `rect(x, y, width, height)`
Creates a rectangle object.

**Parameters:**
- `x`, `y`, `width`, `height`: Integer values

**Returns:** `rectangle` object

### Message Boxes

| Function | Description | Return |
|----------|-------------|--------|
| `message_box(caption, text)` | Simple message box | `nil` |
| `message_box_error(caption, text)` | Error message box | `nil` |
| `message_box_info(caption, text)` | Information message box | `nil` |
| `message_box_warning(caption, text)` | Warning message box | `nil` |
| `message_box_question(caption, text)` | Question message box | `nil` |
| `message_box_ok_cancel(caption, text)` | OK/Cancel dialog | `bool` |
| `message_box_yes_no(caption, text)` | Yes/No dialog | `bool` |
| `message_box_custom(caption, text, flags)` | Custom button dialog | `int` |

**Example:**
```go
if wui.message_box_yes_no("Confirm", "Are you sure?") {
    // User clicked Yes
}
```

### Utility Functions

#### `enabled(control)`
Checks if a control is enabled.

**Parameters:**
- `control`: A control object

**Returns:** `bool`

#### `visible(control)`
Checks if a control is visible.

**Parameters:**
- `control`: A control object

**Returns:** `bool`

## Object Methods

### `window` Object

#### Methods
- `show()` - Shows the window
- `show_modal()` - Shows as modal dialog
- `close()` - Closes the window
- `destroy()` - Destroys the window
- `repaint()` - Forces redraw
- `scroll(dx, dy)` - Scrolls window contents
- `disable_alt_f4()` - Disables Alt+F4
- `enable_alt_f4()` - Enables Alt+F4
- `hide_console_on_start()` - Hides console window on startup
- `show_console_on_start()` - Shows console window on startup
- `add(control)` - Adds a control to the window
- `remove(control)` - Removes a control from the window
- `set_shortcut(callback, keys...)` - Sets global keyboard shortcut

#### Setters
- `set_title(title)`
- `set_size(width, height)`
- `set_position(x, y)`
- `set_bounds(x, y, width, height)`
- `set_alpha(alpha)` - 0-255
- `set_background(color)`
- `set_class_name(name)`
- `set_cursor(cursor)`
- `set_font(font)`
- `set_has_border(bool)`
- `set_has_close_button(bool)`
- `set_has_max_button(bool)`
- `set_has_min_button(bool)`
- `set_height(height)`
- `set_icon(icon)`
- `set_inner_bounds(x, y, width, height)`
- `set_inner_height(height)`
- `set_inner_position(x, y)`
- `set_inner_size(width, height)`
- `set_inner_width(width)`
- `set_inner_x(x)`
- `set_inner_y(y)`
- `set_menu(menu)`
- `set_resizable(bool)`
- `set_state(state)` - Window state constants
- `set_width(width)`
- `set_x(x)`
- `set_y(y)`

#### Event Callbacks
- `set_on_can_close(callback)` - Returns `bool` to allow close
- `set_on_char(callback)` - Receives `rune` argument
- `set_on_close(callback)`
- `set_on_key_down(callback)` - Receives `key` argument
- `set_on_key_up(callback)` - Receives `key` argument
- `set_on_mouse_down(callback)` - Receives `button, x, y`
- `set_on_mouse_move(callback)` - Receives `x, y`
- `set_on_mouse_up(callback)` - Receives `button, x, y`
- `set_on_mouse_wheel(callback)` - Receives `x, y, delta`
- `set_on_resize(callback)`
- `set_on_show(callback)`

#### Getters
- `title()` - Returns `string`
- `size()` - Returns `[width, height]`
- `position()` - Returns `[x, y]`
- `bounds()` - Returns `[x, y, width, height]`
- `alpha()` - Returns `int`
- `children()` - Returns `array` of controls
- `class_name()` - Returns `string`
- `cursor()` - Returns `cursor` or `null`
- `enabled()` - Returns `bool`
- `font()` - Returns `font`
- `get_background()` - Returns `color`
- `handle()` - Returns `int` (HWND)
- `has_border()` - Returns `bool`
- `has_close_button()` - Returns `bool`
- `has_max_button()` - Returns `bool`
- `has_min_button()` - Returns `bool`
- `height()` - Returns `int`
- `icon()` - Returns `icon` or `null`
- `inner_bounds()` - Returns `[x, y, width, height]`
- `inner_height()` - Returns `int`
- `inner_position()` - Returns `[x, y]`
- `inner_size()` - Returns `[width, height]`
- `inner_width()` - Returns `int`
- `inner_x()` - Returns `int`
- `inner_y()` - Returns `int`
- `menu()` - Returns `menu` or `null`
- `monitor()` - Returns `int` (monitor handle)
- `parent()` - Returns `parent` or `null`
- `resizable()` - Returns `bool`
- `state()` - Returns `int`
- `visible()` - Returns `bool`
- `width()` - Returns `int`
- `x()` - Returns `int`
- `y()` - Returns `int`

### Common Control Methods

Most controls share these methods:

#### Setters
- `set_text(text)`
- `set_bounds(x, y, width, height)`
- `set_font(font)`
- `set_enabled(bool)`
- `set_visible(bool)`

#### Getters
- `text()` - Returns `string`
- `font()` - Returns `font`
- `enabled()` - Returns `bool`
- `visible()` - Returns `bool`
- `has_focus()` - Returns `bool`

#### Methods
- `focus()` - Sets keyboard focus

### `button` Object

- `set_onclick(callback)` - Called when button is clicked

### `checkbox` Object

- `set_checked(bool)`
- `checked()` - Returns `bool`
- `set_on_change(callback)` - Receives `checked` argument

### `combobox` Object

- `set_items(array)`
- `add_item(item)`
- `clear()`
- `set_selected_index(index)`
- `selected_index()` - Returns `int`
- `items()` - Returns `array`
- `set_on_change(callback)` - Receives `new_index` argument

### `editline` Object

- `set_text(text)`
- `text()` - Returns `string`
- `set_character_limit(limit)`
- `character_limit()` - Returns `int`
- `set_is_password(bool)`
- `is_password()` - Returns `bool`
- `set_read_only(bool)`
- `read_only()` - Returns `bool`
- `select_all()`
- `set_selection(start, end)`
- `cursor_position()` - Returns `[start, end]`
- `set_cursor_position(pos)`
- `set_on_tab_focus(callback)`
- `set_on_text_change(callback)`

### `intupdown` & `floatupdown` Objects

- `set_min(value)`
- `set_max(value)`
- `set_min_max(min, max)`
- `set_value(value)`
- `value()` - Returns `int` or `float`
- `min()` - Returns `int` or `float`
- `max()` - Returns `int` or `float`
- `min_max()` - Returns `[min, max]`
- `set_on_value_change(callback)` - Receives `new_value`

**FloatUpDown specific:**
- `precision()` - Returns `int`
- `set_precision(precision)`

### `label` Object

- `set_alignment(alignment)` - Use constants `wui.align_left`, `wui.align_center`, `wui.align_right`
- `alignment()` - Returns `int`

### `paintbox` Object

- `set_on_paint(callback)` - Receives `canvas` object
- `set_on_mouse_move(callback)` - Receives `x, y`
- `paint()` - Manually triggers repaint

### `canvas` Object (from `paintbox`)

**Drawing Methods:**
- `draw_rect(x, y, width, height, color)`
- `fill_rect(x, y, width, height, color)`
- `draw_ellipse(x, y, width, height, color)`
- `fill_ellipse(x, y, width, height, color)`
- `line(x1, y1, x2, y2, color)`
- `arc(x, y, width, height, fromAngle, dAngle, color)`
- `draw_pie(x, y, width, height, fromAngle, dAngle, color)`
- `fill_pie(x, y, width, height, fromAngle, dAngle, color)`
- `polygon(points, color)` - Points: `[[x,y], ...]`
- `polyline(points, color)`
- `text_out(x, y, text, color)`
- `text_rect(x, y, w, h, text, color)`
- `text_rect_format(x, y, w, h, text, format, color)`
- `draw_image(image, rect, destX, destY)`

**Font & Text:**
- `set_font(font)`
- `text_extent(text)` - Returns `[width, height]`
- `text_rect_extent(text, width)` - Returns `[width, height]`

**State:**
- `push_draw_region(x, y, width, height)`
- `pop_draw_region()`
- `clear_draw_regions()`

**Getters:**
- `size()` - Returns `[width, height]`
- `width()` - Returns `int`
- `height()` - Returns `int`
- `handle()` - Returns `int` (HDC)

### `stringlist` Object

- `add_item(item)`
- `clear()`
- `set_items(array)`
- `set_selected_index(index)`
- `selected_index()` - Returns `int`
- `items()` - Returns `array`
- `set_on_change(callback)` - Receives `new_index`

### `stringtable` Object

- `clear()`
- `delete_row(row)`
- `set_cell(col, row, text)`
- `col_count()` - Returns `int`
- `row_count()` - Returns `int`
- `selected_row()` - Returns `int`
- `set_on_selection_change(callback)`
- `set_on_tab_focus(callback)`

### `slider` Object

- `set_min(min)`
- `set_max(max)`
- `set_min_max(min, max)`
- `min()` - Returns `int`
- `max()` - Returns `int`
- `min_max()` - Returns `[min, max]`
- `set_value(value)`
- `value()` - Returns `int`
- `set_on_change(callback)` - Receives `new_value`
- `set_orientation(orientation)` - Constants available
- `orientation()` - Returns `int`
- `set_tick_frequency(freq)`
- `set_ticks_visible(bool)`
- `set_arrow_increment(increment)`
- `set_mouse_increment(increment)`

### `progressbar` Object

- `set_value(value)` - Float 0.0-1.0
- `value()` - Returns `float`
- `set_vertical(bool)`
- `vertical()` - Returns `bool`
- `set_moves_forever(bool)`

### `panel` Object

- `add(control)`
- `remove(control)`
- `children()` - Returns `array`
- `inner_bounds()` - Returns `[x, y, width, height]`
- `set_border_style(style)`
- `set_on_resize(callback)`

### `radiobutton` Object

- `set_checked(bool)`
- `checked()` - Returns `bool`
- `set_on_check(callback)` - Receives `checked` argument

### `color` Object

- `r` - Red component (read-only)
- `g` - Green component (read-only)
- `b` - Blue component (read-only)

### `menu` Object

- `add(menuitem)` - Add `menustring`, `menu`, or `menuitem` (separator)

### `menustring` Object

- `set_text(text)`
- `text()` - Returns `string`
- `set_checked(bool)`
- `checked()` - Returns `bool`
- `set_on_click(callback)`

### `fileopendialog` Object

- `add_filter(text, ext1, ...)`
- `execute_single_selection(parent)` - Returns `[bool, string]`
- `execute_multi_selection(parent)` - Returns `[bool, array]`
- `set_filter_index(index)`
- `set_initial_path(path)`
- `set_title(title)`

### `filesavedialog` Object

- `add_filter(text, ext1, ...)`
- `execute(parent)` - Returns `[bool, string]`
- `set_append_ext(bool)`
- `set_filter_index(index)`
- `set_initial_path(path)`
- `set_title(title)`

### `folderselectdialog` Object

- `execute(parent)` - Returns `[bool, string]`
- `set_title(title)`

## Complete Examples

### Basic Window

```go
import "wui"

win := wui.new_window()
win.set_title("Hello Tender")
win.set_size(400, 300)

btn := wui.new_button()
btn.set_text("Click Me")
btn.set_bounds(150, 120, 100, 30)
btn.set_onclick(fn() {
    wui.message_box("Hello", "Button was clicked!")
})

win.add(btn)
win.show()
```

### Custom Drawing with PaintBox

```go
import "wui"

win := wui.new_window()
win.set_title("Drawing Example")
win.set_size(500, 400)

pb := wui.new_paintbox()
pb.set_bounds(0, 0, 500, 400)
pb.set_on_paint(fn(canvas) {
    canvas.fill_rect(50, 50, 200, 100, wui.rgb_color(255, 0, 0))
    canvas.fill_ellipse(300, 50, 100, 100, wui.rgb_color(0, 0, 255))
    canvas.set_font(wui.new_font({name: "Arial", height: 24}))
    canvas.text_out(100, 200, "Hello Canvas!", wui.rgb_color(0, 0, 0))
    canvas.line(50, 300, 400, 350, wui.rgb_color(0, 255, 0))
})

pb.paint()
win.add(pb)
win.show()
```

### File Dialog Example

```go
import "wui"

win := wui.new_window()
win.set_title("File Example")
win.set_size(400, 200)

btn := wui.new_button()
btn.set_text("Open File")
btn.set_bounds(150, 80, 100, 30)
btn.set_onclick(fn() {
    dlg := wui.new_file_open_dialog()
    dlg.set_title("Select a file")
    dlg.add_filter("Text Files", "txt", "csv")
    dlg.add_filter("All Files", "*")
    
    ok, file := dlg.execute_single_selection(win)
    if ok {
        wui.message_box_info("File Selected", "You selected: " + file)
    }
})

win.add(btn)
win.show()
```