//go:build windows

package stdlib

import (
	"image"
	"bytes"
	"fmt"

	"github.com/2dprototype/tender"
	"github.com/gonutz/wui/v2"
)

var wuiModule = map[string]tender.Object{
	"new_window": &tender.NativeFunction{
		Name:  "new_window",
		Value: wuiNewWindow,
	}, // new_window() => Window
	"message_box": &tender.NativeFunction{
		Name:  "message_box",
		Value: wuiMessageBox,
	}, // message_box(caption, text) => nil
	"message_box_error": &tender.NativeFunction{
		Name:  "message_box_error",
		Value: wuiMessageBoxError,
	}, // message_box_error(caption, text) => nil
	"message_box_info": &tender.NativeFunction{
		Name:  "message_box_info",
		Value: wuiMessageBoxInfo,
	}, // message_box_info(caption, text) => nil
	"message_box_warning": &tender.NativeFunction{
		Name:  "message_box_warning",
		Value: wuiMessageBoxWarning,
	}, // message_box_warning(caption, text) => nil
	"message_box_question": &tender.NativeFunction{
		Name:  "message_box_question",
		Value: wuiMessageBoxQuestion,
	}, // message_box_question(caption, text) => nil
	"message_box_ok_cancel": &tender.NativeFunction{
		Name:  "message_box_ok_cancel",
		Value: wuiMessageBoxOKCancel,
	}, // message_box_ok_cancel(caption, text) => bool
	"message_box_yes_no": &tender.NativeFunction{
		Name:  "message_box_yes_no",
		Value: wuiMessageBoxYesNo,
	}, // message_box_yes_no(caption, text) => bool
	"message_box_custom": &tender.NativeFunction{
		Name:  "message_box_custom",
		Value: wuiMessageBoxCustom,
	}, // message_box_custom(caption, text, flags) => int
	"new_button": &tender.NativeFunction{
		Name:  "new_button",
		Value: wuiNewButton,
	}, // new_button() => Button
	"new_checkbox": &tender.NativeFunction{
		Name:  "new_checkbox",
		Value: wuiNewCheckBox,
	}, // new_checkbox() => CheckBox
	"new_label": &tender.NativeFunction{
		Name:  "new_label",
		Value: wuiNewLabel,
	}, // new_label() => Label
	"new_editline": &tender.NativeFunction{
		Name:  "new_editline",
		Value: wuiNewEditLine,
	}, // new_edit_line() => EditLine
	"new_textedit": &tender.NativeFunction{
		Name:  "new_textedit",
		Value: wuiNewTextEdit,
	}, // new_text_edit() => TextEdit
	"new_combobox": &tender.NativeFunction{
		Name:  "new_combo_box",
		Value: wuiNewComboBox,
	}, // new_combo_box() => ComboBox
	"new_stringlist": &tender.NativeFunction{
		Name:  "new_stringlist",
		Value: wuiNewStringList,
	}, // new_string_list() => StringList
	"new_stringtable": &tender.NativeFunction{
		Name:  "new_stringtable",
		Value: wuiNewStringTable,
	}, // new_string_table(header1, ...) => StringTable
	"new_slider": &tender.NativeFunction{
		Name:  "new_slider",
		Value: wuiNewSlider,
	}, // new_slider() => Slider
	"new_progressbar": &tender.NativeFunction{
		Name:  "new_progressbar",
		Value: wuiNewProgressBar,
	}, // new_progress_bar() => ProgressBar
	"new_radiobutton": &tender.NativeFunction{
		Name:  "new_radiobutton",
		Value: wuiNewRadioButton,
	}, // new_radio_button() => RadioButton
	"new_intupdown": &tender.NativeFunction{
		Name:  "new_intupdown",
		Value: wuiNewIntUpDown,
	}, // new_int_up_down() => IntUpDown
	"new_floatupdown": &tender.NativeFunction{
		Name:  "new_floatupdown",
		Value: wuiNewFloatUpDown,
	}, // new_float_up_down() => FloatUpDown
	"new_panel": &tender.NativeFunction{
		Name:  "new_panel",
		Value: wuiNewPanel,
	}, // new_panel() => Panel
	"new_paintbox": &tender.NativeFunction{
		Name:  "new_paintbox",
		Value: wuiNewPaintBox,
	}, // new_paint_box() => PaintBox
	"new_file_open_dialog": &tender.NativeFunction{
		Name:  "new_file_open_dialog",
		Value: wuiNewFileOpenDialog,
	}, // new_file_open_dialog() => FileOpenDialog
	"new_file_save_dialog": &tender.NativeFunction{
		Name:  "new_file_save_dialog",
		Value: wuiNewFileSaveDialog,
	}, // new_file_save_dialog() => FileSaveDialog
	"new_folder_select_dialog": &tender.NativeFunction{
		Name:  "new_folder_select_dialog",
		Value: wuiNewFolderSelectDialog,
	}, // new_folder_select_dialog() => FolderSelectDialog
	"rgb_color": &tender.NativeFunction{
		Name:  "rgb_color",
		Value: wuiRGBColor,
	}, // rgb_color(r, g, b) => Color
	"new_font": &tender.NativeFunction{
		Name:  "new_font",
		Value: wuiNewFont,
	}, // new_font(desc) => Font/error
	"new_cursor_from_image": &tender.NativeFunction{
		Name:  "new_cursor_from_image",
		Value: wuiNewCursorFromImage,
	}, // new_cursor_from_image(image_bytes, x, y) => Cursor/error
	"new_icon_from_image": &tender.NativeFunction{
		Name:  "new_icon_from_image",
		Value: wuiNewIconFromImage,
	}, // new_icon_from_image(image_bytes) => Icon/error
	"new_menu": &tender.NativeFunction{
		Name:  "new_menu",
		Value: wuiNewMenu,
	}, // new_menu(name) => Menu
	"new_menu_string": &tender.NativeFunction{
		Name:  "new_menu_string",
		Value: wuiNewMenuString,
	}, // new_menu_string(text) => MenuString
	"new_menu_separator": &tender.NativeFunction{
		Name:  "new_menu_separator",
		Value: wuiNewMenuSeparator,
	}, // new_menu_separator() => MenuItem
	"new_main_menu": &tender.NativeFunction{
		Name:  "new_main_menu",
		Value: wuiNewMainMenu,
	}, // new_main_menu() => Menu
	"enabled": &tender.NativeFunction{
		Name:  "enabled",
		Value: wuiEnabled,
	}, // enabled(control) => bool
	"new_image": &tender.NativeFunction{
		Name:  "new_image",
		Value: wuiNewImage,
	},
	"new_image_from_hbitmap": &tender.NativeFunction{
		Name:  "new_image_from_hbitmap",
		Value: wuiNewImageFromHBITMAP,
	},
	"new_icon_from_exe_resource": &tender.NativeFunction{
		Name:  "new_icon_from_exe_resource",
		Value: wuiNewIconFromExeResource,
	},
	"new_icon_from_file": &tender.NativeFunction{
		Name:  "new_icon_from_file",
		Value: wuiNewIconFromFile,
	},
	"new_icon_from_reader": &tender.NativeFunction{
		Name:  "new_icon_from_reader",
		Value: wuiNewIconFromReader,
	},
	"rect": &tender.NativeFunction{
		Name:  "rect",
		Value: wuiRect, // already defined in wui_control.go
	},
	"visible": &tender.NativeFunction{
		Name:  "visible",
		Value: wuiVisible,
	}, // visible(control) => bool
}

// Wrapper types for WUI controls
type WUIWindow struct {
	tender.ObjectImpl
	Value *wui.Window
}

func (w *WUIWindow) TypeName() string { return "wui_window" }
func (w *WUIWindow) String() string   { return "<window>" }
func (w *WUIWindow) Copy() tender.Object {
	return &WUIWindow{Value: w.Value}
}


type WUIFileOpenDialog struct {
	tender.ObjectImpl
	Value *wui.FileOpenDialog
}

func (f *WUIFileOpenDialog) TypeName() string { return "wui_fileopendialog" }
func (f *WUIFileOpenDialog) String() string   { return "<fileopendialog>" }
func (f *WUIFileOpenDialog) Copy() tender.Object {
	return &WUIFileOpenDialog{Value: f.Value}
}
func (f *WUIFileOpenDialog) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "add_filter":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) < 2 {
					return nil, tender.ErrWrongNumArguments
				}
				text, ok := tender.ToString(args[0])
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "text",
						Expected: "string(compatible)",
						Found:    args[0].TypeName(),
					}
				}
				ext1, ok := tender.ToString(args[1])
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "ext1",
						Expected: "string(compatible)",
						Found:    args[1].TypeName(),
					}
				}
				exts := make([]string, len(args)-2)
				for i := 2; i < len(args); i++ {
					e, ok := tender.ToString(args[i])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     fmt.Sprintf("ext%d", i-1),
							Expected: "string(compatible)",
							Found:    args[i].TypeName(),
						}
					}
					exts[i-2] = e
				}
				f.Value.AddFilter(text, ext1, exts...)
				return tender.NullValue, nil
			},
		}
	case "execute_single_selection":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				parent, ok := args[0].(*WUIWindow)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "parent",
						Expected: "window",
						Found:    args[0].TypeName(),
					}
				}
				ok, file := f.Value.ExecuteSingleSelection(parent.Value)
				return &tender.Array{
					Value: []tender.Object{
						tender.FromBool(ok),
						&tender.String{Value: file},
					},
				}, nil
			},
		}
	case "execute_multi_selection":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				parent, ok := args[0].(*WUIWindow)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "parent",
						Expected: "window",
						Found:    args[0].TypeName(),
					}
				}
				ok, files := f.Value.ExecuteMultiSelection(parent.Value)
				tenderFiles := make([]tender.Object, len(files))
				for i, fname := range files {
					tenderFiles[i] = &tender.String{Value: fname}
				}
				return &tender.Array{
					Value: []tender.Object{
						tender.FromBool(ok),
						&tender.Array{Value: tenderFiles},
					},
				}, nil
			},
		}
	case "set_filter_index":
		res = &tender.NativeFunction{
			Value: FuncAIR(f.Value.SetFilterIndex),
		}
	case "set_initial_path":
		res = &tender.NativeFunction{
			Value: FuncASR(f.Value.SetInitialPath),
		}
	case "set_title":
		res = &tender.NativeFunction{
			Value: FuncASR(f.Value.SetTitle),
		}
	}
	return
}


type WUIFileSaveDialog struct {
	tender.ObjectImpl
	Value *wui.FileSaveDialog
}

func (f *WUIFileSaveDialog) TypeName() string { return "wui_filesavedialog" }
func (f *WUIFileSaveDialog) String() string   { return "<filesavedialog>" }
func (f *WUIFileSaveDialog) Copy() tender.Object {
	return &WUIFileSaveDialog{Value: f.Value}
}
func (f *WUIFileSaveDialog) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "add_filter":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) < 2 {
					return nil, tender.ErrWrongNumArguments
				}
				text, ok := tender.ToString(args[0])
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "text",
						Expected: "string(compatible)",
						Found:    args[0].TypeName(),
					}
				}
				ext1, ok := tender.ToString(args[1])
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "ext1",
						Expected: "string(compatible)",
						Found:    args[1].TypeName(),
					}
				}
				exts := make([]string, len(args)-2)
				for i := 2; i < len(args); i++ {
					e, ok := tender.ToString(args[i])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     fmt.Sprintf("ext%d", i-1),
							Expected: "string(compatible)",
							Found:    args[i].TypeName(),
						}
					}
					exts[i-2] = e
				}
				f.Value.AddFilter(text, ext1, exts...)
				return tender.NullValue, nil
			},
		}
	case "execute":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				parent, ok := args[0].(*WUIWindow)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "parent",
						Expected: "window",
						Found:    args[0].TypeName(),
					}
				}
				ok, file := f.Value.Execute(parent.Value)
				return &tender.Array{
					Value: []tender.Object{
						tender.FromBool(ok),
						&tender.String{Value: file},
					},
				}, nil
			},
		}
	case "set_append_ext":
		res = &tender.NativeFunction{
			Value: FuncABR(f.Value.SetAppendExt),
		}
	case "set_filter_index":
		res = &tender.NativeFunction{
			Value: FuncAIR(f.Value.SetFilterIndex),
		}
	case "set_initial_path":
		res = &tender.NativeFunction{
			Value: FuncASR(f.Value.SetInitialPath),
		}
	case "set_title":
		res = &tender.NativeFunction{
			Value: FuncASR(f.Value.SetTitle),
		}
	}
	return
}


type WUIFolderSelectDialog struct {
	tender.ObjectImpl
	Value *wui.FolderSelectDialog
}

func (f *WUIFolderSelectDialog) TypeName() string { return "wui_folderselectdialog" }
func (f *WUIFolderSelectDialog) String() string   { return "<folderselectdialog>" }
func (f *WUIFolderSelectDialog) Copy() tender.Object {
	return &WUIFolderSelectDialog{Value: f.Value}
}
func (f *WUIFolderSelectDialog) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "execute":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				parent, ok := args[0].(*WUIWindow)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "parent",
						Expected: "window",
						Found:    args[0].TypeName(),
					}
				}
				ok, folder := f.Value.Execute(parent.Value)
				return &tender.Array{
					Value: []tender.Object{
						tender.FromBool(ok),
						&tender.String{Value: folder},
					},
				}, nil
			},
		}
	case "set_title":
		res = &tender.NativeFunction{
			Value: FuncASR(f.Value.SetTitle),
		}
	}
	return
}


type WUIFont struct {
	tender.ObjectImpl
	Value *wui.Font
}

func (f *WUIFont) TypeName() string { return "wui_font" }
func (f *WUIFont) String() string   { return "<font>" }
func (f *WUIFont) Copy() tender.Object {
	return &WUIFont{Value: f.Value}
}

type WUICursor struct {
	tender.ObjectImpl
	Value *wui.Cursor
}

func (c *WUICursor) TypeName() string { return "wui_cursor" }
func (c *WUICursor) String() string   { return "<cursor>" }
func (c *WUICursor) Copy() tender.Object {
	return &WUICursor{Value: c.Value}
}

type WUIIcon struct {
	tender.ObjectImpl
	Value *wui.Icon
}

func (i *WUIIcon) TypeName() string { return "wui_icon" }
func (i *WUIIcon) String() string   { return "<icon>" }
func (i *WUIIcon) Copy() tender.Object {
	return &WUIIcon{Value: i.Value}
}

type WUIMenu struct {
	tender.ObjectImpl
	Value *wui.Menu
}

func (m *WUIMenu) TypeName() string { return "wui_menu" }
func (m *WUIMenu) String() string   { return "<menu>" }
func (m *WUIMenu) Copy() tender.Object {
	return &WUIMenu{Value: m.Value}
}

type WUIMenuString struct {
	tender.ObjectImpl
	Value *wui.MenuString
}

func (m *WUIMenuString) TypeName() string { return "wui_menustring" }
func (m *WUIMenuString) String() string   { return "<menustring>" }
func (m *WUIMenuString) Copy() tender.Object {
	return &WUIMenuString{Value: m.Value}
}

type WUIMenuItem struct {
	tender.ObjectImpl
	Value wui.MenuItem
}

func (m *WUIMenuItem) TypeName() string { return "wui_menuitem" }
func (m *WUIMenuItem) String() string   { return "<menuitem>" }
func (m *WUIMenuItem) Copy() tender.Object {
	return &WUIMenuItem{Value: m.Value}
}


// Helper function to extract font from wrapper
func extractFont(obj tender.Object) (*wui.Font, bool) {
	if font, ok := obj.(*WUIFont); ok {
		return font.Value, true
	}
	return nil, false
}

func wuiNewWindow(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	window := wui.NewWindow()
	return &WUIWindow{Value: window}, nil
}


func wuiNewFileOpenDialog(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	fileOpenDialog := wui.NewFileOpenDialog()
	return &WUIFileOpenDialog{Value: fileOpenDialog}, nil
}

func wuiNewFileSaveDialog(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	fileSaveDialog := wui.NewFileSaveDialog()
	return &WUIFileSaveDialog{Value: fileSaveDialog}, nil
}

func wuiNewFolderSelectDialog(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	folderSelectDialog := wui.NewFolderSelectDialog()
	return &WUIFolderSelectDialog{Value: folderSelectDialog}, nil
}

func wuiRGBColor(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 3 {
		err = tender.ErrWrongNumArguments
		return
	}

	r, ok := tender.ToInt(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "red",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	g, ok := tender.ToInt(args[1])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "green",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	b, ok := tender.ToInt(args[2])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "blue",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
		return
	}

	color := wui.RGB(uint8(r), uint8(g), uint8(b))
	return &WUIColor{Value: color}, nil
}

func wuiNewFont(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}

	descMap, ok := args[0].(*tender.Map)
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "font_desc",
			Expected: "map",
			Found:    args[0].TypeName(),
		}
		return
	}

	desc := wui.FontDesc{}

	if val, ok := descMap.Value["name"]; ok {
		desc.Name, _ = tender.ToString(val)
	}
	if val, ok := descMap.Value["height"]; ok {
		desc.Height, _ = tender.ToInt(val)
	}
	if val, ok := descMap.Value["bold"]; ok {
		desc.Bold, _ = tender.ToBool(val)
	}
	if val, ok := descMap.Value["italic"]; ok {
		desc.Italic, _ = tender.ToBool(val)
	}
	if val, ok := descMap.Value["underlined"]; ok {
		desc.Underlined, _ = tender.ToBool(val)
	}
	if val, ok := descMap.Value["striked_out"]; ok {
		desc.StrikedOut, _ = tender.ToBool(val)
	}

	font, err := wui.NewFont(desc)
	if err != nil {
		ret = wrapError(err)
		return
	}

	return &WUIFont{Value: font}, nil
}

func wuiNewCursorFromImage(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 3 {
		err = tender.ErrWrongNumArguments
		return
	}

	imageBytes, err := ToFileData(args[0])
	if err != nil {
		return nil, err
	}

	x, ok := tender.ToInt(args[1])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "x",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
		return
	}

	y, ok := tender.ToInt(args[2])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "y",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		ret = wrapError(err)
		return
	}

	cursor, err := wui.NewCursorFromImage(img, x, y)
	if err != nil {
		ret = wrapError(err)
		return
	}

	return &WUICursor{Value: cursor}, nil
}

func wuiNewImage(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}

	data, err := ToFileData(args[0])
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return wrapError(err), nil
	}

	wuiImg := wui.NewImage(img)
	return &WUIImage{Value: wuiImg}, nil
}

func wuiNewImageFromHBITMAP(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 3 {
		err = tender.ErrWrongNumArguments
		return
	}
	bitmap, _ := tender.ToInt(args[0])
	width, _ := tender.ToInt(args[1])
	height, _ := tender.ToInt(args[2])
	wuiImg := wui.NewImageFromHBITMAP(uintptr(bitmap), width, height)
	return &WUIImage{Value: wuiImg}, nil
}

func wuiNewIconFromExeResource(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}
	resID, _ := tender.ToInt(args[0])
	icon, err := wui.NewIconFromExeResource(resID)
	if err != nil {
		return wrapError(err), nil
	}
	return &WUIIcon{Value: icon}, nil
}

func wuiNewIconFromFile(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}
	path, ok := tender.ToString(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "path",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}
	icon, err := wui.NewIconFromFile(path)
	if err != nil {
		return wrapError(err), nil
	}
	return &WUIIcon{Value: icon}, nil
}

func wuiNewIconFromReader(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}
	// We need an io.Reader, but tender doesn't have a built-in reader type.
	// We can accept bytes and use bytes.NewReader.
	bytesData, ok := tender.ToByteSlice(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "bytes",
			Found:    args[0].TypeName(),
		}
		return
	}
	reader := bytes.NewReader(bytesData)
	icon, err := wui.NewIconFromReader(reader)
	if err != nil {
		return wrapError(err), nil
	}
	return &WUIIcon{Value: icon}, nil
}

func wuiNewIconFromImage(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}

	imageBytes, ok := tender.ToByteSlice(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "image",
			Expected: "bytes",
			Found:    args[0].TypeName(),
		}
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		ret = wrapError(err)
		return
	}

	icon, err := wui.NewIconFromImage(img)
	if err != nil {
		ret = wrapError(err)
		return
	}

	return &WUIIcon{Value: icon}, nil
}

func wuiNewMenu(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}

	name, ok := tender.ToString(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	menu := wui.NewMenu(name)
	return &WUIMenu{Value: menu}, nil
}

func wuiNewMenuString(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}

	text, ok := tender.ToString(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "text",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	menuString := wui.NewMenuString(text)
	return &WUIMenuString{Value: menuString}, nil
}

func wuiNewMenuSeparator(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	menuSeparator := wui.NewMenuSeparator()
	return &WUIMenuItem{Value: menuSeparator}, nil
}

func wuiNewMainMenu(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	mainMenu := wui.NewMainMenu()
	return &WUIMenu{Value: mainMenu}, nil
}

func wuiEnabled(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}

	control, ok := extractControl(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "control",
			Expected: "control",
			Found:    args[0].TypeName(),
		}
		return
	}

	if control.Enabled() {
		ret = tender.TrueValue
	} else {
		ret = tender.FalseValue
	}
	return
}

func wuiVisible(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 1 {
		err = tender.ErrWrongNumArguments
		return
	}

	control, ok := extractControl(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "control",
			Expected: "control",
			Found:    args[0].TypeName(),
		}
		return
	}

	if control.Visible() {
		ret = tender.TrueValue
	} else {
		ret = tender.FalseValue
	}
	return
}

// Implement IndexGet for all wrapper types to expose their methods
func (w *WUIWindow) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	// --- Methods (actions) ---
	case "show":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				err := w.Value.Show()
				if err != nil {
					return wrapError(err), nil
				}
				return tender.NullValue, nil
			},
		}
	case "show_modal":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				err := w.Value.ShowModal()
				if err != nil {
					return wrapError(err), nil
				}
				return tender.NullValue, nil
			},
		}
	case "close":
		res = &tender.NativeFunction{
			Value: FuncAR(w.Value.Close),
		}
	case "destroy":
		res = &tender.NativeFunction{
			Value: FuncAR(w.Value.Destroy),
		}
	case "repaint":
		res = &tender.NativeFunction{
			Value: FuncAR(w.Value.Repaint),
		}
	case "scroll":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				dx, _ := tender.ToInt(args[0])
				dy, _ := tender.ToInt(args[1])
				w.Value.Scroll(dx, dy)
				return tender.NullValue, nil
			},
		}
	case "disable_alt_f4":
		res = &tender.NativeFunction{
			Value: FuncAR(w.Value.DisableAltF4),
		}
	case "enable_alt_f4":
		res = &tender.NativeFunction{
			Value: FuncAR(w.Value.EnableAltF4),
		}
	case "hide_console_on_start":
		res = &tender.NativeFunction{
			Value: FuncAR(w.Value.HideConsoleOnStart),
		}
	case "show_console_on_start":
		res = &tender.NativeFunction{
			Value: FuncAR(w.Value.ShowConsoleOnStart),
		}
	case "add":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				control, ok := extractControl(args[0])
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "control",
						Expected: "control",
						Found:    args[0].TypeName(),
					}
				}
				w.Value.Add(control)
				return tender.NullValue, nil
			},
		}
	case "remove":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				control, ok := extractControl(args[0])
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "control",
						Expected: "control",
						Found:    args[0].TypeName(),
					}
				}
				w.Value.Remove(control)
				return tender.NullValue, nil
			},
		}
	case "set_shortcut":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) < 2 {
					return nil, tender.ErrWrongNumArguments
				}
				// first arg is function, rest are keys
				callback := args[0]
				keys := make([]wui.Key, len(args)-1)
				for i := 1; i < len(args); i++ {
					k, ok := tender.ToInt(args[i])
					if !ok {
						return nil, tender.ErrInvalidArgumentType{
							Name:     fmt.Sprintf("key %d", i),
							Expected: "int(compatible)",
							Found:    args[i].TypeName(),
						}
					}
					keys[i-1] = wui.Key(k)
				}
				w.Value.SetShortcut(func() {
					tender.WrapFuncCall(vm, callback)
				}, keys...)
				return tender.NullValue, nil
			},
		}
	case "set_menu":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				menu, ok := args[0].(*WUIMenu)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "menu",
						Expected: "menu",
						Found:    args[0].TypeName(),
					}
				}
				w.Value.SetMenu(menu.Value)
				return tender.NullValue, nil
			},
		}

	// --- Setters ---
	case "set_title":
		res = &tender.NativeFunction{
			Value: FuncASR(w.Value.SetTitle),
		}
	case "set_size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				width, _ := tender.ToInt(args[0])
				height, _ := tender.ToInt(args[1])
				w.Value.SetSize(width, height)
				return tender.NullValue, nil
			},
		}
	case "set_position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				x, _ := tender.ToInt(args[0])
				y, _ := tender.ToInt(args[1])
				w.Value.SetPosition(x, y)
				return tender.NullValue, nil
			},
		}
	case "set_bounds":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 4 {
					return nil, tender.ErrWrongNumArguments
				}
				x, _ := tender.ToInt(args[0])
				y, _ := tender.ToInt(args[1])
				width, _ := tender.ToInt(args[2])
				height, _ := tender.ToInt(args[3])
				w.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	case "set_alpha":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				a, _ := tender.ToInt(args[0])
				w.Value.SetAlpha(uint8(a))
				return tender.NullValue, nil
			},
		}
	case "set_background":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				color, ok := extractColor(args[0])
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "color",
						Expected: "color",
						Found:    args[0].TypeName(),
					}
				}
				w.Value.SetBackground(color)
				return tender.NullValue, nil
			},
		}
	case "set_class_name":
		res = &tender.NativeFunction{
			Value: FuncASR(w.Value.SetClassName),
		}
	case "set_cursor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				cursor, ok := args[0].(*WUICursor)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "cursor",
						Expected: "cursor",
						Found:    args[0].TypeName(),
					}
				}
				w.Value.SetCursor(cursor.Value)
				return tender.NullValue, nil
			},
		}
	case "set_font":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				font, ok := extractFont(args[0])
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "font",
						Expected: "font",
						Found:    args[0].TypeName(),
					}
				}
				w.Value.SetFont(font)
				return tender.NullValue, nil
			},
		}
	case "set_has_border":
		res = &tender.NativeFunction{
			Value: FuncABR(w.Value.SetHasBorder),
		}
	case "set_has_close_button":
		res = &tender.NativeFunction{
			Value: FuncABR(w.Value.SetHasCloseButton),
		}
	case "set_has_max_button":
		res = &tender.NativeFunction{
			Value: FuncABR(w.Value.SetHasMaxButton),
		}
	case "set_has_min_button":
		res = &tender.NativeFunction{
			Value: FuncABR(w.Value.SetHasMinButton),
		}
	case "set_height":
		res = &tender.NativeFunction{
			Value: FuncAIR(w.Value.SetHeight),
		}
	case "set_icon":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				icon, ok := args[0].(*WUIIcon)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "icon",
						Expected: "icon",
						Found:    args[0].TypeName(),
					}
				}
				w.Value.SetIcon(icon.Value)
				return tender.NullValue, nil
			},
		}
	case "set_inner_bounds":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 4 {
					return nil, tender.ErrWrongNumArguments
				}
				x, _ := tender.ToInt(args[0])
				y, _ := tender.ToInt(args[1])
				width, _ := tender.ToInt(args[2])
				height, _ := tender.ToInt(args[3])
				w.Value.SetInnerBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	case "set_inner_height":
		res = &tender.NativeFunction{
			Value: FuncAIR(w.Value.SetInnerHeight),
		}
	case "set_inner_position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				x, _ := tender.ToInt(args[0])
				y, _ := tender.ToInt(args[1])
				w.Value.SetInnerPosition(x, y)
				return tender.NullValue, nil
			},
		}
	case "set_inner_size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				width, _ := tender.ToInt(args[0])
				height, _ := tender.ToInt(args[1])
				w.Value.SetInnerSize(width, height)
				return tender.NullValue, nil
			},
		}
	case "set_inner_width":
		res = &tender.NativeFunction{
			Value: FuncAIR(w.Value.SetInnerWidth),
		}
	case "set_inner_x":
		res = &tender.NativeFunction{
			Value: FuncAIR(w.Value.SetInnerX),
		}
	case "set_inner_y":
		res = &tender.NativeFunction{
			Value: FuncAIR(w.Value.SetInnerY),
		}
	case "set_resizable":
		res = &tender.NativeFunction{
			Value: FuncABR(w.Value.SetResizable),
		}
	case "set_state":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				state, _ := tender.ToInt(args[0])
				w.Value.SetState(wui.WindowState(state))
				return tender.NullValue, nil
			},
		}
	case "set_width":
		res = &tender.NativeFunction{
			Value: FuncAIR(w.Value.SetWidth),
		}
	case "set_x":
		res = &tender.NativeFunction{
			Value: FuncAIR(w.Value.SetX),
		}
	case "set_y":
		res = &tender.NativeFunction{
			Value: FuncAIR(w.Value.SetY),
		}

	// --- Setters for event callbacks (they take a function) ---
	case "set_on_can_close":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnCanClose(func() bool {
					ret, err := tender.WrapFuncCall(vm, args[0])
					if err != nil {
						return true // default to allow close
					}
					ok, _ := tender.ToBool(ret)
					return ok
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_char":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnChar(func(r rune) {
					tender.WrapFuncCall(vm, args[0], &tender.String{Value: string(r)})
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_close":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnClose(func() {
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_key_down":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnKeyDown(func(key int) {
					tender.WrapFuncCall(vm, args[0], &tender.Int{Value: int64(key)})
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_key_up":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnKeyUp(func(key int) {
					tender.WrapFuncCall(vm, args[0], &tender.Int{Value: int64(key)})
				})
				return tender.NullValue, nil
			},
		}
	// case "set_on_message":
		// res = &tender.NativeFunction{
			// NeedVMObj: true,
			// Value: func(args ...tender.Object) (tender.Object, error) {
				// vm := args[0].(*tender.VMObj).Value
				// args = args[1:]
				// if len(args) != 1 {
					// return nil, tender.ErrWrongNumArguments
				// }
				// w.Value.SetOnMessage(func(msg uint32, wParam, lParam uintptr) uintptr {
					// tender.WrapFuncCall(vm, args[0],
						// &tender.Int{Value: int64(msg)},
						// &tender.Int{Value: int64(wParam)},
						// &tender.Int{Value: int64(lParam)})
					// return 0
				// })
				// return tender.NullValue, nil
			// },
		// }
	case "set_on_mouse_down":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnMouseDown(func(button wui.MouseButton, x, y int) {
					tender.WrapFuncCall(vm, args[0],
						&tender.Int{Value: int64(button)},
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)})
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_mouse_move":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnMouseMove(func(x, y int) {
					tender.WrapFuncCall(vm, args[0],
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)})
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_mouse_up":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnMouseUp(func(button wui.MouseButton, x, y int) {
					tender.WrapFuncCall(vm, args[0],
						&tender.Int{Value: int64(button)},
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)})
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_mouse_wheel":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnMouseWheel(func(x, y int, delta float64) {
					tender.WrapFuncCall(vm, args[0],
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)},
						&tender.Float{Value: delta})
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_resize":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnResize(func() {
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_show":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				w.Value.SetOnShow(func() {
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}

	// --- Getters ---
	case "title":
		res = &tender.NativeFunction{
			Value: FuncARS(w.Value.Title),
		}
	case "size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				width, height := w.Value.Size()
				return &tender.Array{
					Value: []tender.Object{
						&tender.Int{Value: int64(width)},
						&tender.Int{Value: int64(height)},
					},
				}, nil
			},
		}
	case "position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				x, y := w.Value.Position()
				return &tender.Array{
					Value: []tender.Object{
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)},
					},
				}, nil
			},
		}
	case "bounds":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				x, y, width, height := w.Value.Bounds()
				return &tender.Array{
					Value: []tender.Object{
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)},
						&tender.Int{Value: int64(width)},
						&tender.Int{Value: int64(height)},
					},
				}, nil
			},
		}
	case "alpha":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.Alpha())}, nil
			},
		}
	case "children":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				children := w.Value.Children()
				tenderChildren := make([]tender.Object, len(children))
				for i, child := range children {
					tenderChildren[i] = wrapControl(child)
				}
				return &tender.Array{Value: tenderChildren}, nil
			},
		}
	case "class_name":
		res = &tender.NativeFunction{
			Value: FuncARS(w.Value.ClassName),
		}
	case "cursor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				cursor := w.Value.Cursor()
				if cursor == nil {
					return tender.NullValue, nil
				}
				return &WUICursor{Value: cursor}, nil
			},
		}
	case "enabled":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.FromBool(w.Value.Enabled()), nil
			},
		}
	case "font":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				font := w.Value.Font()
				return &WUIFont{Value: font}, nil
			},
		}
	case "get_background":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &WUIColor{Value: w.Value.GetBackground()}, nil
			},
		}
	case "handle":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.Handle())}, nil
			},
		}
	case "has_border":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.FromBool(w.Value.HasBorder()), nil
			},
		}
	case "has_close_button":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.FromBool(w.Value.HasCloseButton()), nil
			},
		}
	case "has_max_button":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.FromBool(w.Value.HasMaxButton()), nil
			},
		}
	case "has_min_button":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.FromBool(w.Value.HasMinButton()), nil
			},
		}
	case "height":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.Height())}, nil
			},
		}
	case "icon":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				icon := w.Value.Icon()
				if icon == nil {
					return tender.NullValue, nil
				}
				return &WUIIcon{Value: icon}, nil
			},
		}
	case "inner_bounds":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				x, y, width, height := w.Value.InnerBounds()
				return &tender.Array{
					Value: []tender.Object{
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)},
						&tender.Int{Value: int64(width)},
						&tender.Int{Value: int64(height)},
					},
				}, nil
			},
		}
	case "inner_height":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.InnerHeight())}, nil
			},
		}
	case "inner_position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				x, y := w.Value.InnerPosition()
				return &tender.Array{
					Value: []tender.Object{
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)},
					},
				}, nil
			},
		}
	case "inner_size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				width, height := w.Value.InnerSize()
				return &tender.Array{
					Value: []tender.Object{
						&tender.Int{Value: int64(width)},
						&tender.Int{Value: int64(height)},
					},
				}, nil
			},
		}
	case "inner_width":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.InnerWidth())}, nil
			},
		}
	case "inner_x":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.InnerX())}, nil
			},
		}
	case "inner_y":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.InnerY())}, nil
			},
		}
	case "menu":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				menu := w.Value.Menu()
				if menu == nil {
					return tender.NullValue, nil
				}
				return &WUIMenu{Value: menu}, nil
			},
		}
	case "monitor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.Monitor())}, nil
			},
		}
	case "on_can_close":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_char":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_close":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_key_down":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_key_up":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_message":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_mouse_down":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_mouse_move":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_mouse_up":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_mouse_wheel":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_resize":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_show":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "parent":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				parent := w.Value.Parent()
				if parent == nil {
					return tender.NullValue, nil
				}
				// Wrap parent container - for now return null as we don't have a generic container wrapper
				return tender.NullValue, nil
			},
		}
	case "resizable":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.FromBool(w.Value.Resizable()), nil
			},
		}
	case "state":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.State())}, nil
			},
		}
	case "visible":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.FromBool(w.Value.Visible()), nil
			},
		}
	case "width":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.Width())}, nil
			},
		}
	case "x":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.X())}, nil
			},
		}
	case "y":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(w.Value.Y())}, nil
			},
		}
	}
	return
}

// Color wrapper
type WUIColor struct {
	tender.ObjectImpl
	Value wui.Color
}

func (c *WUIColor) TypeName() string { return "wui_color" }
func (c *WUIColor) String() string   { return fmt.Sprintf("<color r:%d g:%d b:%d>", c.Value.R(), c.Value.G(), c.Value.B()) }
func (c *WUIColor) Copy() tender.Object {
	return &WUIColor{Value: c.Value}
}

func (c *WUIColor) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "r":
		res = &tender.Int{Value: int64(c.Value.R())}
	case "g":
		res = &tender.Int{Value: int64(c.Value.G())}
	case "b":
		res = &tender.Int{Value: int64(c.Value.B())}
	}
	return
}

func (m *WUIMenu) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "add":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				menuItem, ok := args[0].(*WUIMenuItem)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "menuitem",
						Expected: "menuitem",
						Found:    args[0].TypeName(),
					}
				}
				m.Value.Add(menuItem.Value)
				return tender.NullValue, nil
			},
		}
	}
	return
}

func (m *WUIMenuString) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "set_text":
		res = &tender.NativeFunction{
			Value: FuncASR(m.Value.SetText),
		}
	case "text":
		res = &tender.NativeFunction{
			Value: FuncARS(m.Value.Text),
		}
	case "checked":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.FromBool(m.Value.Checked()), nil
			},
		}
	case "set_checked":
		res = &tender.NativeFunction{
			Value: FuncABR(m.Value.SetChecked),
		}
	case "on_click":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "set_on_click":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				m.Value.SetOnClick(func() {
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	}
	return
}

// Common helper functions for WUI module
func wuiSetBounds(control interface{ SetBounds(x, y, width, height int) }, args []tender.Object) (tender.Object, error) {
	if len(args) != 4 {
		return nil, tender.ErrWrongNumArguments
	}
	x, _ := tender.ToInt(args[0])
	y, _ := tender.ToInt(args[1])
	width, _ := tender.ToInt(args[2])
	height, _ := tender.ToInt(args[3])
	control.SetBounds(x, y, width, height)
	return tender.NullValue, nil
}

func wuiSetPosition(control interface{ SetPosition(x, y int) }, args []tender.Object) (tender.Object, error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	x, _ := tender.ToInt(args[0])
	y, _ := tender.ToInt(args[1])
	control.SetPosition(x, y)
	return tender.NullValue, nil
}

func wuiSetSize(control interface{ SetSize(width, height int) }, args []tender.Object) (tender.Object, error) {
	if len(args) != 2 {
		return nil, tender.ErrWrongNumArguments
	}
	width, _ := tender.ToInt(args[0])
	height, _ := tender.ToInt(args[1])
	control.SetSize(width, height)
	return tender.NullValue, nil
}