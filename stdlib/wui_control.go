//go:build windows

package stdlib

import (
	"fmt"
	
	"github.com/2dprototype/tender"
	"github.com/gonutz/wui/v2"
)

type WUIButton struct {
	tender.ObjectImpl
	Value *wui.Button
}

func (b *WUIButton) TypeName() string { return "button" }
func (b *WUIButton) String() string   { return "<button>" }
func (b *WUIButton) Copy() tender.Object {
	return &WUIButton{Value: b.Value}
}

type WUICheckBox struct {
	tender.ObjectImpl
	Value *wui.CheckBox
}

func (c *WUICheckBox) TypeName() string { return "checkbox" }
func (c *WUICheckBox) String() string   { return "<checkbox>" }
func (c *WUICheckBox) Copy() tender.Object {
	return &WUICheckBox{Value: c.Value}
}

type WUILabel struct {
	tender.ObjectImpl
	Value *wui.Label
}

func (l *WUILabel) TypeName() string { return "label" }
func (l *WUILabel) String() string   { return "<label>" }
func (l *WUILabel) Copy() tender.Object {
	return &WUILabel{Value: l.Value}
}

type WUIEditLine struct {
	tender.ObjectImpl
	Value *wui.EditLine
}

func (e *WUIEditLine) TypeName() string { return "editline" }
func (e *WUIEditLine) String() string   { return "<editline>" }
func (e *WUIEditLine) Copy() tender.Object {
	return &WUIEditLine{Value: e.Value}
}

type WUITextEdit struct {
	tender.ObjectImpl
	Value *wui.TextEdit
}

func (t *WUITextEdit) TypeName() string { return "textedit" }
func (t *WUITextEdit) String() string   { return "<textedit>" }
func (t *WUITextEdit) Copy() tender.Object {
	return &WUITextEdit{Value: t.Value}
}
func (t *WUITextEdit) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	// Methods
	case "select_all":
		res = &tender.NativeFunction{
			Value: FuncAR(t.Value.SelectAll),
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
				t.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	case "set_character_limit":
		res = &tender.NativeFunction{
			Value: FuncAIR(t.Value.SetCharacterLimit),
		}
	case "set_cursor_position":
		res = &tender.NativeFunction{
			Value: FuncAIR(t.Value.SetCursorPosition),
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
				t.Value.SetFont(font)
				return tender.NullValue, nil
			},
		}
	case "set_on_tab_focus":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				t.Value.SetOnTabFocus(func() {
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_text_change":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				t.Value.SetOnTextChange(func() {
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "set_read_only":
		res = &tender.NativeFunction{
			Value: FuncABR(t.Value.SetReadOnly),
		}
	case "set_selection":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				start, _ := tender.ToInt(args[0])
				end, _ := tender.ToInt(args[1])
				t.Value.SetSelection(start, end)
				return tender.NullValue, nil
			},
		}
	case "set_text":
		res = &tender.NativeFunction{
			Value: FuncASR(t.Value.SetText),
		}
	case "set_word_wrap":
		res = &tender.NativeFunction{
			Value: FuncABR(t.Value.SetWordWrap),
		}
	case "set_writes_tabs":
		res = &tender.NativeFunction{
			Value: FuncABR(t.Value.SetWritesTabs),
		}

	// Getters
	case "character_limit":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(t.Value.CharacterLimit())}, nil
			},
		}
	case "cursor_position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				start, end := t.Value.CursorPosition()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(start)},
					&tender.Int{Value: int64(end)},
				}}, nil
			},
		}
	case "on_tab_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_text_change":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "read_only":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if t.Value.ReadOnly() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "text":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.String{Value: t.Value.Text()}, nil
			},
		}
	case "word_wrap":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if t.Value.WordWrap() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "writes_tabs":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if t.Value.WritesTabs() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	}
	return
}

type WUIComboBox struct {
	tender.ObjectImpl
	Value *wui.ComboBox
}

func (c *WUIComboBox) TypeName() string { return "combobox" }
func (c *WUIComboBox) String() string   { return "<combobox>" }
func (c *WUIComboBox) Copy() tender.Object {
	return &WUIComboBox{Value: c.Value}
}



type WUISlider struct {
	tender.ObjectImpl
	Value *wui.Slider
}

func (s *WUISlider) TypeName() string { return "slider" }
func (s *WUISlider) String() string   { return "<slider>" }
func (s *WUISlider) Copy() tender.Object {
	return &WUISlider{Value: s.Value}
}
func (s *WUISlider) IndexGet(index tender.Object) (res tender.Object, err error) {
    strIdx, ok := index.(*tender.String)
    if !ok {
        return nil, tender.ErrInvalidIndexType
    }

    switch strIdx.Value {
    // Getters
    case "anchors":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                horizontal, vertical := s.Value.Anchors()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(horizontal)},
                    &tender.Int{Value: int64(vertical)},
                }}, nil
            },
        }
    case "arrow_increment":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.ArrowIncrement())}, nil
            },
        }
    case "bounds":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, y, width, height := s.Value.Bounds()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(x)},
                    &tender.Int{Value: int64(y)},
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "cursor_position":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.CursorPosition())}, nil
            },
        }
    case "enabled":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                if s.Value.Enabled() {
                    return tender.TrueValue, nil
                }
                return tender.FalseValue, nil
            },
        }
    case "handle":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.Handle())}, nil
            },
        }
    case "height":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.Height())}, nil
            },
        }
    case "horizontal_anchor":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.HorizontalAnchor())}, nil
            },
        }
    case "max":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.Max())}, nil
            },
        }
    case "min":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.Min())}, nil
            },
        }
    case "min_max":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                min, max := s.Value.MinMax()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(min)},
                    &tender.Int{Value: int64(max)},
                }}, nil
            },
        }
    case "mouse_increment":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.MouseIncrement())}, nil
            },
        }
    case "on_change":
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
    case "on_tab_focus":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return tender.NullValue, nil
            },
        }
    case "orientation":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.Orientation())}, nil
            },
        }
    case "parent":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                parent := s.Value.Parent()
                if parent == nil {
                    return tender.NullValue, nil
                }
                return tender.NullValue, nil
            },
        }
    case "position":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, y := s.Value.Position()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(x)},
                    &tender.Int{Value: int64(y)},
                }}, nil
            },
        }
    case "size":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                width, height := s.Value.Size()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "tick_frequency":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.TickFrequency())}, nil
            },
        }
    case "tick_position":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.TickPosition())}, nil
            },
        }
    case "ticks_visible":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                if s.Value.TicksVisible() {
                    return tender.TrueValue, nil
                }
                return tender.FalseValue, nil
            },
        }
    case "vertical_anchor":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.VerticalAnchor())}, nil
            },
        }
    case "visible":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                if s.Value.Visible() {
                    return tender.TrueValue, nil
                }
                return tender.FalseValue, nil
            },
        }
    case "width":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.Width())}, nil
            },
        }
    case "x":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.X())}, nil
            },
        }
    case "y":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.Y())}, nil
            },
        }

    // Setters
    case "set_anchors":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                horizontal, _ := tender.ToInt(args[0])
                vertical, _ := tender.ToInt(args[1])
                s.Value.SetAnchors(wui.Anchor(horizontal), wui.Anchor(vertical))
                return tender.NullValue, nil
            },
        }
    case "set_arrow_increment":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetArrowIncrement),
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
                s.Value.SetBounds(x, y, width, height)
                return tender.NullValue, nil
            },
        }
    case "set_cursor_position":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetCursorPosition),
        }
    case "set_enabled":
        res = &tender.NativeFunction{
            Value: FuncABR(s.Value.SetEnabled),
        }
    case "set_height":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetHeight),
        }
    case "set_horizontal_anchor":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                anchor, _ := tender.ToInt(args[0])
                s.Value.SetHorizontalAnchor(wui.Anchor(anchor))
                return tender.NullValue, nil
            },
        }
    case "set_max":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetMax),
        }
    case "set_min":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetMin),
        }
    case "set_min_max":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                min, _ := tender.ToInt(args[0])
                max, _ := tender.ToInt(args[1])
                s.Value.SetMinMax(min, max)
                return tender.NullValue, nil
            },
        }
    case "set_mouse_increment":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetMouseIncrement),
        }
    case "set_on_change":
        res = &tender.NativeFunction{
            NeedVMObj: true,
            Value: func(args ...tender.Object) (tender.Object, error) {
                vm := args[0].(*tender.VMObj).Value
                args = args[1:]
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                s.Value.SetOnChange(func(cursor int){
                    tender.WrapFuncCall(vm, args[0], &tender.Int{Value: int64(cursor)})
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
                s.Value.SetOnResize(func(){
                    tender.WrapFuncCall(vm, args[0])
                })
                return tender.NullValue, nil
            },
        }
    case "set_on_tab_focus":
        res = &tender.NativeFunction{
            NeedVMObj: true,
            Value: func(args ...tender.Object) (tender.Object, error) {
                vm := args[0].(*tender.VMObj).Value
                args = args[1:]
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                s.Value.SetOnTabFocus(func(){
                    tender.WrapFuncCall(vm, args[0])
                })
                return tender.NullValue, nil
            },
        }
    case "set_orientation":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                orientation, _ := tender.ToInt(args[0])
                s.Value.SetOrientation(wui.SliderOrientation(orientation))
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
                s.Value.SetPosition(x, y)
                return tender.NullValue, nil
            },
        }
    case "set_size":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                width, _ := tender.ToInt(args[0])
                height, _ := tender.ToInt(args[1])
                s.Value.SetSize(width, height)
                return tender.NullValue, nil
            },
        }
    case "set_tick_frequency":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetTickFrequency),
        }
    case "set_tick_position":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                position, _ := tender.ToInt(args[0])
                s.Value.SetTickPosition(wui.TickPosition(position))
                return tender.NullValue, nil
            },
        }
    case "set_ticks_visible":
        res = &tender.NativeFunction{
            Value: FuncABR(s.Value.SetTicksVisible),
        }
    case "set_vertical_anchor":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                anchor, _ := tender.ToInt(args[0])
                s.Value.SetVerticalAnchor(wui.Anchor(anchor))
                return tender.NullValue, nil
            },
        }
    case "set_visible":
        res = &tender.NativeFunction{
            Value: FuncABR(s.Value.SetVisible),
        }
    case "set_width":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetWidth),
        }
    case "set_x":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetX),
        }
    case "set_y":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.SetY),
        }
    }
    return
}


func wrapControl(control wui.Control) tender.Object {
    switch c := control.(type) {
    case *wui.Button:
        return &WUIButton{Value: c}
    case *wui.CheckBox:
        return &WUICheckBox{Value: c}
    case *wui.Label:
        return &WUILabel{Value: c}
    case *wui.EditLine:
        return &WUIEditLine{Value: c}
    case *wui.TextEdit:
        return &WUITextEdit{Value: c}
    case *wui.ComboBox:
        return &WUIComboBox{Value: c}
    case *wui.StringList:
        return &WUIStringList{Value: c}
    case *wui.StringTable:
        return &WUIStringTable{Value: c}
    case *wui.Slider:
        return &WUISlider{Value: c}
    case *wui.ProgressBar:
        return &WUIProgressBar{Value: c}
    case *wui.RadioButton:
        return &WUIRadioButton{Value: c}
    case *wui.IntUpDown:
        return &WUIIntUpDown{Value: c}
    case *wui.FloatUpDown:
        return &WUIFloatUpDown{Value: c}
    case *wui.Panel:
        return &WUIPanel{Value: c}
    case *wui.PaintBox:
        return &WUIPaintBox{Value: c}
    default:
        return tender.NullValue
    }
}

type WUIProgressBar struct {
	tender.ObjectImpl
	Value *wui.ProgressBar
}

func (p *WUIProgressBar) TypeName() string { return "progressbar" }
func (p *WUIProgressBar) String() string   { return "<progressbar>" }
func (p *WUIProgressBar) Copy() tender.Object {
	return &WUIProgressBar{Value: p.Value}
}
func (p *WUIProgressBar) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	// Getters
	case "anchors":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				horizontal, vertical := p.Value.Anchors()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(horizontal)},
					&tender.Int{Value: int64(vertical)},
				}}, nil
			},
		}
	case "bounds":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				x, y, width, height := p.Value.Bounds()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(x)},
					&tender.Int{Value: int64(y)},
					&tender.Int{Value: int64(width)},
					&tender.Int{Value: int64(height)},
				}}, nil
			},
		}
	case "enabled":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if p.Value.Enabled() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "handle":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.Handle())}, nil
			},
		}
	case "height":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.Height())}, nil
			},
		}
	case "horizontal_anchor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.HorizontalAnchor())}, nil
			},
		}
	case "moves_forever":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if p.Value.MovesForever() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
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
	case "parent":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				parent := p.Value.Parent()
				if parent == nil {
					return tender.NullValue, nil
				}
				return tender.NullValue, nil
			},
		}
	case "position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				x, y := p.Value.Position()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(x)},
					&tender.Int{Value: int64(y)},
				}}, nil
			},
		}
	case "size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				width, height := p.Value.Size()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(width)},
					&tender.Int{Value: int64(height)},
				}}, nil
			},
		}
	case "value":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Float{Value: p.Value.Value()}, nil
			},
		}
	case "vertical":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if p.Value.Vertical() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "vertical_anchor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.VerticalAnchor())}, nil
			},
		}
	case "visible":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if p.Value.Visible() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "width":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.Width())}, nil
			},
		}
	case "x":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.X())}, nil
			},
		}
	case "y":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.Y())}, nil
			},
		}

	// Setters
	case "set_anchors":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				horizontal, _ := tender.ToInt(args[0])
				vertical, _ := tender.ToInt(args[1])
				p.Value.SetAnchors(wui.Anchor(horizontal), wui.Anchor(vertical))
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
				p.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	case "set_enabled":
		res = &tender.NativeFunction{
			Value: FuncABR(p.Value.SetEnabled),
		}
	case "set_height":
		res = &tender.NativeFunction{
			Value: FuncAIR(p.Value.SetHeight),
		}
	case "set_horizontal_anchor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				anchor, _ := tender.ToInt(args[0])
				p.Value.SetHorizontalAnchor(wui.Anchor(anchor))
				return tender.NullValue, nil
			},
		}
	case "set_moves_forever":
		res = &tender.NativeFunction{
			Value: FuncABR(p.Value.SetMovesForever),
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
				p.Value.SetOnResize(func(){
					tender.WrapFuncCall(vm, args[0])
				})
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
				p.Value.SetPosition(x, y)
				return tender.NullValue, nil
			},
		}
	case "set_size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				width, _ := tender.ToInt(args[0])
				height, _ := tender.ToInt(args[1])
				p.Value.SetSize(width, height)
				return tender.NullValue, nil
			},
		}
	case "set_value":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				val, _ := tender.ToFloat64(args[0])
				p.Value.SetValue(val)
				return tender.NullValue, nil
			},
		}
	case "set_vertical":
		res = &tender.NativeFunction{
			Value: FuncABR(p.Value.SetVertical),
		}
	case "set_vertical_anchor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				anchor, _ := tender.ToInt(args[0])
				p.Value.SetVerticalAnchor(wui.Anchor(anchor))
				return tender.NullValue, nil
			},
		}
	case "set_visible":
		res = &tender.NativeFunction{
			Value: FuncABR(p.Value.SetVisible),
		}
	case "set_width":
		res = &tender.NativeFunction{
			Value: FuncAIR(p.Value.SetWidth),
		}
	case "set_x":
		res = &tender.NativeFunction{
			Value: FuncAIR(p.Value.SetX),
		}
	case "set_y":
		res = &tender.NativeFunction{
			Value: FuncAIR(p.Value.SetY),
		}
	}
	return
}

type WUIRadioButton struct {
	tender.ObjectImpl
	Value *wui.RadioButton
}

func (r *WUIRadioButton) TypeName() string { return "radiobutton" }
func (r *WUIRadioButton) String() string   { return "<radiobutton>" }
func (r *WUIRadioButton) Copy() tender.Object {
	return &WUIRadioButton{Value: r.Value}
}
func (r *WUIRadioButton) IndexGet(index tender.Object) (res tender.Object, err error) {
    strIdx, ok := index.(*tender.String)
    if !ok {
        return nil, tender.ErrInvalidIndexType
    }

    switch strIdx.Value {
    case "set_text":
        res = &tender.NativeFunction{
            Value: FuncASR(r.Value.SetText),
        }
    case "text":
        res = &tender.NativeFunction{
            Value: FuncARS(r.Value.Text),
        }
    case "set_checked":
        res = &tender.NativeFunction{
            Value: FuncABR(r.Value.SetChecked),
        }
    case "checked":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                if r.Value.Checked() {
                    return tender.TrueValue, nil
                }
                return tender.FalseValue, nil
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
                r.Value.SetFont(font)
                return tender.NullValue, nil
            },
        }
    case "font":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                font := r.Value.Font()
                return &WUIFont{Value: font}, nil
            },
        }
    case "focus":
        res = &tender.NativeFunction{
            Value: FuncAR(r.Value.Focus),
        }
    case "has_focus":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                if r.Value.HasFocus() {
                    return tender.TrueValue, nil
                }
                return tender.FalseValue, nil
            },
        }
    case "on_check":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return tender.NullValue, nil
            },
        }
    case "set_on_check":
        res = &tender.NativeFunction{
            NeedVMObj: true,
            Value: func(args ...tender.Object) (tender.Object, error) {
                vm := args[0].(*tender.VMObj).Value
                args = args[1:]
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                r.Value.SetOnCheck(func(checked bool){
                    var arg tender.Object
                    if checked {
                        arg = tender.TrueValue
                    } else {
                        arg = tender.FalseValue
                    }
                    tender.WrapFuncCall(vm, args[0], arg)
                })
                return tender.NullValue, nil
            },
        }
    case "set_on_tab_focus":
        res = &tender.NativeFunction{
            NeedVMObj: true,
            Value: func(args ...tender.Object) (tender.Object, error) {
                vm := args[0].(*tender.VMObj).Value
                args = args[1:]
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                r.Value.SetOnTabFocus(func(){
                    tender.WrapFuncCall(vm, args[0])
                })
                return tender.NullValue, nil
            },
        }
    case "on_tab_focus":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
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
                r.Value.SetBounds(x, y, width, height)
                return tender.NullValue, nil
            },
        }
    }
    return
}

type WUIIntUpDown struct {
	tender.ObjectImpl
	Value *wui.IntUpDown
}

func (i *WUIIntUpDown) TypeName() string { return "intupdown" }
func (i *WUIIntUpDown) String() string   { return "<intupdown>" }
func (i *WUIIntUpDown) Copy() tender.Object {
	return &WUIIntUpDown{Value: i.Value}
}

type WUIFloatUpDown struct {
	tender.ObjectImpl
	Value *wui.FloatUpDown
}

func (f *WUIFloatUpDown) TypeName() string { return "floatupdown" }
func (f *WUIFloatUpDown) String() string   { return "<floatupdown>" }
func (f *WUIFloatUpDown) Copy() tender.Object {
	return &WUIFloatUpDown{Value: f.Value}
}

type WUIPanel struct {
	tender.ObjectImpl
	Value *wui.Panel
}

func (p *WUIPanel) TypeName() string { return "panel" }
func (p *WUIPanel) String() string   { return "<panel>" }
func (p *WUIPanel) Copy() tender.Object {
	return &WUIPanel{Value: p.Value}
}
func (p *WUIPanel) IndexGet(index tender.Object) (res tender.Object, err error) {
    strIdx, ok := index.(*tender.String)
    if !ok {
        return nil, tender.ErrInvalidIndexType
    }

    switch strIdx.Value {
    // Container methods
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
                p.Value.Add(control)
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
                p.Value.Remove(control)
                return tender.NullValue, nil
            },
        }
    case "children":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                children := p.Value.Children()
                tenderChildren := make([]tender.Object, len(children))
                for i, child := range children {
                    // Wrap each child control
                    tenderChildren[i] = wrapControl(child)
                }
                return &tender.Array{Value: tenderChildren}, nil
            },
        }

    // Property getters
    case "anchors":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                horizontal, vertical := p.Value.Anchors()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(horizontal)},
                    &tender.Int{Value: int64(vertical)},
                }}, nil
            },
        }
    case "border_style":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.BorderStyle())}, nil
            },
        }
    case "bounds":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, y, width, height := p.Value.Bounds()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(x)},
                    &tender.Int{Value: int64(y)},
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "enabled":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                if p.Value.Enabled() {
                    return tender.TrueValue, nil
                }
                return tender.FalseValue, nil
            },
        }
    case "font":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                font := p.Value.Font()
                return &WUIFont{Value: font}, nil
            },
        }
    case "handle":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.Handle())}, nil
            },
        }
    case "height":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.Height())}, nil
            },
        }
    case "horizontal_anchor":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.HorizontalAnchor())}, nil
            },
        }
    case "inner_bounds":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, y, width, height := p.Value.InnerBounds()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(x)},
                    &tender.Int{Value: int64(y)},
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "inner_height":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.InnerHeight())}, nil
            },
        }
    case "inner_position":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, y := p.Value.InnerPosition()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(x)},
                    &tender.Int{Value: int64(y)},
                }}, nil
            },
        }
    case "inner_size":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                width, height := p.Value.InnerSize()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "inner_width":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.InnerWidth())}, nil
            },
        }
    case "inner_x":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.InnerX())}, nil
            },
        }
    case "inner_y":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.InnerY())}, nil
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
    case "parent":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                parent := p.Value.Parent()
                if parent == nil {
                    return tender.NullValue, nil
                }
                // Need to wrap parent container
                return tender.NullValue, nil
            },
        }
    case "position":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, y := p.Value.Position()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(x)},
                    &tender.Int{Value: int64(y)},
                }}, nil
            },
        }
    case "size":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                width, height := p.Value.Size()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "vertical_anchor":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.VerticalAnchor())}, nil
            },
        }
    case "visible":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                if p.Value.Visible() {
                    return tender.TrueValue, nil
                }
                return tender.FalseValue, nil
            },
        }
    case "width":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.Width())}, nil
            },
        }
    case "x":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.X())}, nil
            },
        }
    case "y":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(p.Value.Y())}, nil
            },
        }

    // Setters
    case "set_anchors":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                horizontal, _ := tender.ToInt(args[0])
                vertical, _ := tender.ToInt(args[1])
                p.Value.SetAnchors(wui.Anchor(horizontal), wui.Anchor(vertical))
                return tender.NullValue, nil
            },
        }
    case "set_border_style":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                style, _ := tender.ToInt(args[0])
                p.Value.SetBorderStyle(wui.PanelBorderStyle(style))
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
                p.Value.SetBounds(x, y, width, height)
                return tender.NullValue, nil
            },
        }
    case "set_enabled":
        res = &tender.NativeFunction{
            Value: FuncABR(p.Value.SetEnabled),
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
                p.Value.SetFont(font)
                return tender.NullValue, nil
            },
        }
    case "set_height":
        res = &tender.NativeFunction{
            Value: FuncAIR(p.Value.SetHeight),
        }
    case "set_horizontal_anchor":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                anchor, _ := tender.ToInt(args[0])
                p.Value.SetHorizontalAnchor(wui.Anchor(anchor))
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
                p.Value.SetInnerBounds(x, y, width, height)
                return tender.NullValue, nil
            },
        }
    case "set_inner_height":
        res = &tender.NativeFunction{
            Value: FuncAIR(p.Value.SetInnerHeight),
        }
    case "set_inner_position":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                p.Value.SetInnerPosition(x, y)
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
                p.Value.SetInnerSize(width, height)
                return tender.NullValue, nil
            },
        }
    case "set_inner_width":
        res = &tender.NativeFunction{
            Value: FuncAIR(p.Value.SetInnerWidth),
        }
    case "set_inner_x":
        res = &tender.NativeFunction{
            Value: FuncAIR(p.Value.SetInnerX),
        }
    case "set_inner_y":
        res = &tender.NativeFunction{
            Value: FuncAIR(p.Value.SetInnerY),
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
                p.Value.SetOnResize(func(){
                    tender.WrapFuncCall(vm, args[0])
                })
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
                p.Value.SetPosition(x, y)
                return tender.NullValue, nil
            },
        }
    case "set_size":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                width, _ := tender.ToInt(args[0])
                height, _ := tender.ToInt(args[1])
                p.Value.SetSize(width, height)
                return tender.NullValue, nil
            },
        }
    case "set_vertical_anchor":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                anchor, _ := tender.ToInt(args[0])
                p.Value.SetVerticalAnchor(wui.Anchor(anchor))
                return tender.NullValue, nil
            },
        }
    case "set_visible":
        res = &tender.NativeFunction{
            Value: FuncABR(p.Value.SetVisible),
        }
    case "set_width":
        res = &tender.NativeFunction{
            Value: FuncAIR(p.Value.SetWidth),
        }
    case "set_x":
        res = &tender.NativeFunction{
            Value: FuncAIR(p.Value.SetX),
        }
    case "set_y":
        res = &tender.NativeFunction{
            Value: FuncAIR(p.Value.SetY),
        }
    }
    return
}

type WUIPaintBox struct {
	tender.ObjectImpl
	Value *wui.PaintBox
}

func (p *WUIPaintBox) TypeName() string { return "paintbox" }
func (p *WUIPaintBox) String() string   { return "<paintbox>" }
func (p *WUIPaintBox) Copy() tender.Object {
	return &WUIPaintBox{Value: p.Value}
}


type WUIStringList struct {
	tender.ObjectImpl
	Value *wui.StringList
}

func (s *WUIStringList) TypeName() string { return "stringlist" }
func (s *WUIStringList) String() string   { return "<stringlist>" }
func (s *WUIStringList) Copy() tender.Object {
	return &WUIStringList{Value: s.Value}
}
func (s *WUIStringList) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	// Methods
	case "add_item":
		res = &tender.NativeFunction{
			Value: FuncASR(s.Value.AddItem),
		}
	case "clear":
		res = &tender.NativeFunction{
			Value: FuncAR(s.Value.Clear),
		}
	case "focus":
		res = &tender.NativeFunction{
			Value: FuncAR(s.Value.Focus),
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
				s.Value.SetFont(font)
				return tender.NullValue, nil
			},
		}
	case "set_items":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				arr, ok := args[0].(*tender.Array)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "items",
						Expected: "array",
						Found:    args[0].TypeName(),
					}
				}
				items := make([]string, len(arr.Value))
				for i, item := range arr.Value {
					items[i], _ = tender.ToString(item)
				}
				s.Value.SetItems(items)
				return tender.NullValue, nil
			},
		}
	case "set_on_change":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				s.Value.SetOnChange(func(newIndex int) {
					tender.WrapFuncCall(vm, args[0], &tender.Int{Value: int64(newIndex)})
				})
				return tender.NullValue, nil
			},
		}
	case "set_selected_index":
		res = &tender.NativeFunction{
			Value: FuncAIR(s.Value.SetSelectedIndex),
		}
	case "set_text":
		res = &tender.NativeFunction{
			Value: FuncASR(s.Value.SetText),
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
				s.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}

	// Getters
	case "font":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &WUIFont{Value: s.Value.Font()}, nil
			},
		}
	case "has_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if s.Value.HasFocus() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "items":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				items := s.Value.Items()
				tenderItems := make([]tender.Object, len(items))
				for i, item := range items {
					tenderItems[i] = &tender.String{Value: item}
				}
				return &tender.Array{Value: tenderItems}, nil
			},
		}
	case "on_change":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "selected_index":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(s.Value.SelectedIndex())}, nil
			},
		}
	case "text":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.String{Value: s.Value.Text()}, nil
			},
		}
	}
	return
}

type WUIStringTable struct {
	tender.ObjectImpl
	Value *wui.StringTable
}

func (s *WUIStringTable) TypeName() string { return "stringtable" }
func (s *WUIStringTable) String() string   { return "<stringtable>" }
func (s *WUIStringTable) Copy() tender.Object {
	return &WUIStringTable{Value: s.Value}
}
func (s *WUIStringTable) IndexGet(index tender.Object) (res tender.Object, err error) {
    strIdx, ok := index.(*tender.String)
    if !ok {
        return nil, tender.ErrInvalidIndexType
    }

    switch strIdx.Value {
    // Methods
    case "clear":
        res = &tender.NativeFunction{
            Value: FuncAR(s.Value.Clear),
        }
    case "delete_row":
        res = &tender.NativeFunction{
            Value: FuncAIR(s.Value.DeleteRow),
        }
    case "focus":
        res = &tender.NativeFunction{
            Value: FuncAR(s.Value.Focus),
        }
    case "set_cell":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 3 {
                    return nil, tender.ErrWrongNumArguments
                }
                col, _ := tender.ToInt(args[0])
                row, _ := tender.ToInt(args[1])
                str, ok := tender.ToString(args[2])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "text",
                        Expected: "string(compatible)",
                        Found:    args[2].TypeName(),
                    }
                }
                s.Value.SetCell(col, row, str)
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
                s.Value.SetFont(font)
                return tender.NullValue, nil
            },
        }
    case "set_on_selection_change":
        res = &tender.NativeFunction{
            NeedVMObj: true,
            Value: func(args ...tender.Object) (tender.Object, error) {
                vm := args[0].(*tender.VMObj).Value
                args = args[1:]
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                s.Value.SetOnSelectionChange(func(){
                    tender.WrapFuncCall(vm, args[0])
                })
                return tender.NullValue, nil
            },
        }
    case "set_on_tab_focus":
        res = &tender.NativeFunction{
            NeedVMObj: true,
            Value: func(args ...tender.Object) (tender.Object, error) {
                vm := args[0].(*tender.VMObj).Value
                args = args[1:]
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                s.Value.SetOnTabFocus(func(){
                    tender.WrapFuncCall(vm, args[0])
                })
                return tender.NullValue, nil
            },
        }
    case "set_text":
        res = &tender.NativeFunction{
            Value: FuncASR(s.Value.SetText),
        }

    // Getters
    case "col_count":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.ColCount())}, nil
            },
        }
    case "font":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                font := s.Value.Font()
                return &WUIFont{Value: font}, nil
            },
        }
    case "has_focus":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                if s.Value.HasFocus() {
                    return tender.TrueValue, nil
                }
                return tender.FalseValue, nil
            },
        }
    case "on_selection_change":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return tender.NullValue, nil
            },
        }
    case "on_tab_focus":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return tender.NullValue, nil
            },
        }
    case "row_count":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.RowCount())}, nil
            },
        }
    case "selected_row":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(s.Value.SelectedRow())}, nil
            },
        }
    case "text":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.String{Value: s.Value.Text()}, nil
            },
        }

    // Set bounds
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
                s.Value.SetBounds(x, y, width, height)
                return tender.NullValue, nil
            },
        }
    }
    return
}

func wuiNewButton(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	button := wui.NewButton()
	return &WUIButton{Value: button}, nil
}

func wuiNewCheckBox(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	checkbox := wui.NewCheckBox()
	return &WUICheckBox{Value: checkbox}, nil
}

func wuiNewLabel(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	label := wui.NewLabel()
	return &WUILabel{Value: label}, nil
}

func wuiNewEditLine(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	editLine := wui.NewEditLine()
	return &WUIEditLine{Value: editLine}, nil
}

func wuiNewTextEdit(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	textEdit := wui.NewTextEdit()
	return &WUITextEdit{Value: textEdit}, nil
}

func wuiNewComboBox(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	comboBox := wui.NewComboBox()
	return &WUIComboBox{Value: comboBox}, nil
}

func wuiNewStringList(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	stringList := wui.NewStringList()
	return &WUIStringList{Value: stringList}, nil
}

func wuiNewStringTable(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) < 1 {
		err = tender.ErrWrongNumArguments
		return
	}

	header1, ok := tender.ToString(args[0])
	if !ok {
		err = tender.ErrInvalidArgumentType{
			Name:     "first header",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
		return
	}

	headers := []string{header1}
	for i := 1; i < len(args); i++ {
		header, ok := tender.ToString(args[i])
		if !ok {
			err = tender.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("header %d", i+1),
				Expected: "string(compatible)",
				Found:    args[i].TypeName(),
			}
			return
		}
		headers = append(headers, header)
	}

	stringTable := wui.NewStringTable(header1, headers[1:]...)
	return &WUIStringTable{Value: stringTable}, nil
}

func wuiNewSlider(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	slider := wui.NewSlider()
	return &WUISlider{Value: slider}, nil
}

func wuiNewProgressBar(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	progressBar := wui.NewProgressBar()
	return &WUIProgressBar{Value: progressBar}, nil
}

func wuiNewRadioButton(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	radioButton := wui.NewRadioButton()
	return &WUIRadioButton{Value: radioButton}, nil
}

func wuiNewIntUpDown(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	intUpDown := wui.NewIntUpDown()
	return &WUIIntUpDown{Value: intUpDown}, nil
}

func wuiNewFloatUpDown(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	floatUpDown := wui.NewFloatUpDown()
	return &WUIFloatUpDown{Value: floatUpDown}, nil
}

func wuiNewPanel(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	panel := wui.NewPanel()
	return &WUIPanel{Value: panel}, nil
}

func wuiNewPaintBox(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 0 {
		err = tender.ErrWrongNumArguments
		return
	}

	paintBox := wui.NewPaintBox()
	return &WUIPaintBox{Value: paintBox}, nil
}


// Helper function to extract control from wrapper
func extractControl(obj tender.Object) (wui.Control, bool) {
	switch c := obj.(type) {
	case *WUIButton:
		return c.Value, true
	case *WUICheckBox:
		return c.Value, true
	case *WUILabel:
		return c.Value, true
	case *WUIEditLine:
		return c.Value, true
	case *WUITextEdit:
		return c.Value, true
	case *WUIComboBox:
		return c.Value, true
	case *WUIStringList:
		return c.Value, true
	case *WUIStringTable:
		return c.Value, true
	case *WUISlider:
		return c.Value, true
	case *WUIProgressBar:
		return c.Value, true
	case *WUIRadioButton:
		return c.Value, true
	case *WUIIntUpDown:
		return c.Value, true
	case *WUIFloatUpDown:
		return c.Value, true
	case *WUIPanel:
		return c.Value, true
	case *WUIPaintBox:
		return c.Value, true
	default:
		return nil, false
	}
}





func (b *WUIButton) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "set_text":
		res = &tender.NativeFunction{
			Value: FuncASR(b.Value.SetText),
		}
	case "text":
		res = &tender.NativeFunction{
			Value: FuncARS(b.Value.Text),
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
				b.Value.SetFont(font)
				return tender.NullValue, nil
			},
		}
	case "font":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				font := b.Value.Font()
				return &WUIFont{Value: font}, nil
			},
		}
	case "focus":
		res = &tender.NativeFunction{
			Value: FuncAR(b.Value.Focus),
		}
	case "has_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if b.Value.HasFocus() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "set_onclick":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:] // the first arg is VMObj inserted by VM
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				b.Value.SetOnClick(func(){
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "onclick":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				// Return nil as we can't return a Go function to script
				return tender.NullValue, nil
			},
		}
	case "set_on_tab_focus":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:] // the first arg is VMObj inserted by VM
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				b.Value.SetOnTabFocus(func(){
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "on_tab_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				// Return nil as we can't return a Go function to script
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
				b.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	case "set_enabled":
		res = &tender.NativeFunction{
			Value: FuncABR(b.Value.SetEnabled),
		}
	case "enabled":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if b.Value.Enabled() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "set_visible":
		res = &tender.NativeFunction{
			Value: FuncABR(b.Value.SetVisible),
		}
	case "visible":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if b.Value.Visible() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	}
	return
}

// WUICheckBox updates
func (c *WUICheckBox) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "set_text":
		res = &tender.NativeFunction{
			Value: FuncASR(c.Value.SetText),
		}
	case "text":
		res = &tender.NativeFunction{
			Value: FuncARS(c.Value.Text),
		}
	case "set_checked":
		res = &tender.NativeFunction{
			Value: FuncABR(c.Value.SetChecked),
		}
	case "checked":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if c.Value.Checked() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
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
				c.Value.SetFont(font)
				return tender.NullValue, nil
			},
		}
	case "font":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				font := c.Value.Font()
				return &WUIFont{Value: font}, nil
			},
		}
	case "focus":
		res = &tender.NativeFunction{
			Value: FuncAR(c.Value.Focus),
		}
	case "has_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if c.Value.HasFocus() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "set_on_change":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				c.Value.SetOnChange(func(checked bool){
					var arg tender.Object
					if checked {
						arg = tender.TrueValue
					} else {
						arg = tender.FalseValue
					}
					tender.WrapFuncCall(vm, args[0], arg)
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_tab_focus":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				c.Value.SetOnTabFocus(func(){
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "on_tab_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
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
				c.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	}
	return
}

// WUIComboBox updates
func (c *WUIComboBox) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "set_text":
		res = &tender.NativeFunction{
			Value: FuncASR(c.Value.SetText),
		}
	case "text":
		res = &tender.NativeFunction{
			Value: FuncARS(c.Value.Text),
		}
	case "set_items":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				arr, ok := args[0].(*tender.Array)
				if !ok {
					return nil, tender.ErrInvalidArgumentType{
						Name:     "items",
						Expected: "array",
						Found:    args[0].TypeName(),
					}
				}
				items := make([]string, len(arr.Value))
				for i, item := range arr.Value {
					items[i], _ = tender.ToString(item)
				}
				c.Value.SetItems(items)
				return tender.NullValue, nil
			},
		}
	case "add_item":
		res = &tender.NativeFunction{
			Value: FuncASR(c.Value.AddItem),
		}
	case "clear":
		res = &tender.NativeFunction{
			Value: FuncAR(c.Value.Clear),
		}
	case "set_selected_index":
		res = &tender.NativeFunction{
			Value: FuncAIR(c.Value.SetSelectedIndex),
		}
	case "selected_index":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(c.Value.SelectedIndex())}, nil
			},
		}
	case "items":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				items := c.Value.Items()
				tenderItems := make([]tender.Object, len(items))
				for i, item := range items {
					tenderItems[i] = &tender.String{Value: item}
				}
				return &tender.Array{Value: tenderItems}, nil
			},
		}
	case "focus":
		res = &tender.NativeFunction{
			Value: FuncAR(c.Value.Focus),
		}
	case "has_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if c.Value.HasFocus() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
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
				c.Value.SetFont(font)
				return tender.NullValue, nil
			},
		}
	case "font":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				font := c.Value.Font()
				return &WUIFont{Value: font}, nil
			},
		}
	case "set_on_change":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				c.Value.SetOnChange(func(newIndex int){
					tender.WrapFuncCall(vm, args[0], &tender.Int{Value: int64(newIndex)})
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_tab_focus":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				c.Value.SetOnTabFocus(func(){
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "on_tab_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
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
				c.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	}
	return
}

func (e *WUIEditLine) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "set_text":
		res = &tender.NativeFunction{
			Value: FuncASR(e.Value.SetText),
		}
	case "text":
		res = &tender.NativeFunction{
			Value: FuncARS(e.Value.Text),
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
				e.Value.SetBounds(x, y, width, height)
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
				e.Value.SetFont(font)
				return tender.NullValue, nil
			},
		}
	case "font":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				font := e.Value.Font()
				return &WUIFont{Value: font}, nil
			},
		}
	case "focus":
		res = &tender.NativeFunction{
			Value: FuncAR(e.Value.Focus),
		}
	case "has_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if e.Value.HasFocus() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "character_limit":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(e.Value.CharacterLimit())}, nil
			},
		}
	case "set_character_limit":
		res = &tender.NativeFunction{
			Value: FuncAIR(e.Value.SetCharacterLimit),
		}
	case "cursor_position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				start, end := e.Value.CursorPosition()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(start)},
					&tender.Int{Value: int64(end)},
				}}, nil
			},
		}
	case "set_cursor_position":
		res = &tender.NativeFunction{
			Value: FuncAIR(e.Value.SetCursorPosition),
		}
	case "is_password":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if e.Value.IsPassword() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "set_is_password":
		res = &tender.NativeFunction{
			Value: FuncABR(e.Value.SetIsPassword),
		}
	case "read_only":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if e.Value.ReadOnly() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "set_read_only":
		res = &tender.NativeFunction{
			Value: FuncABR(e.Value.SetReadOnly),
		}
	case "select_all":
		res = &tender.NativeFunction{
			Value: FuncAR(e.Value.SelectAll),
		}
	case "set_selection":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				start, _ := tender.ToInt(args[0])
				end, _ := tender.ToInt(args[1])
				e.Value.SetSelection(start, end)
				return tender.NullValue, nil
			},
		}
	case "set_on_tab_focus":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				e.Value.SetOnTabFocus(func(){
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_text_change":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				e.Value.SetOnTextChange(func(){
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "on_tab_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_text_change":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	}
	return
}

// WUIFloatUpDown updates
func (f *WUIFloatUpDown) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "set_min":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				min, _ := tender.ToFloat64(args[0])
				f.Value.SetMin(min)
				return tender.NullValue, nil
			},
		}
	case "set_max":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				max, _ := tender.ToFloat64(args[0])
				f.Value.SetMax(max)
				return tender.NullValue, nil
			},
		}
	case "set_value":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				val, _ := tender.ToFloat64(args[0])
				f.Value.SetValue(val)
				return tender.NullValue, nil
			},
		}
	case "value":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Float{Value: float64(f.Value.Value())}, nil
			},
		}
	case "min":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Float{Value: f.Value.Min()}, nil
			},
		}
	case "max":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Float{Value: f.Value.Max()}, nil
			},
		}
	case "min_max":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				min, max := f.Value.MinMax()
				return &tender.Array{Value: []tender.Object{
					&tender.Float{Value: min},
					&tender.Float{Value: max},
				}}, nil
			},
		}
	case "set_min_max":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				min, _ := tender.ToFloat64(args[0])
				max, _ := tender.ToFloat64(args[1])
				f.Value.SetMinMax(min, max)
				return tender.NullValue, nil
			},
		}
	case "precision":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(f.Value.Precision())}, nil
			},
		}
	case "set_precision":
		res = &tender.NativeFunction{
			Value: FuncAIR(f.Value.SetPrecision),
		}
	case "cursor_position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				start, end := f.Value.CursorPosition()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(start)},
					&tender.Int{Value: int64(end)},
				}}, nil
			},
		}
	case "set_cursor_position":
		res = &tender.NativeFunction{
			Value: FuncAIR(f.Value.SetCursorPosition),
		}
	case "select_all":
		res = &tender.NativeFunction{
			Value: FuncAR(f.Value.SelectAll),
		}
	case "set_selection":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				start, _ := tender.ToInt(args[0])
				end, _ := tender.ToInt(args[1])
				f.Value.SetSelection(start, end)
				return tender.NullValue, nil
			},
		}
	case "set_on_tab_focus":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				f.Value.SetOnTabFocus(func(){
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_value_change":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				f.Value.SetOnValueChange(func(value float64){
					tender.WrapFuncCall(vm, args[0], &tender.Float{Value: value})
				})
				return tender.NullValue, nil
			},
		}
	case "on_tab_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_value_change":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
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
				f.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	// Position and size methods
	case "set_pos":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				x, _ := tender.ToInt(args[0])
				y, _ := tender.ToInt(args[1])
				f.Value.SetPos(x, y)
				return tender.NullValue, nil
			},
		}
	case "set_size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				width, _ := tender.ToInt(args[0])
				height, _ := tender.ToInt(args[1])
				f.Value.SetSize(width, height)
				return tender.NullValue, nil
			},
		}
	case "set_x":
		res = &tender.NativeFunction{
			Value: FuncAIR(f.Value.SetX),
		}
	case "set_y":
		res = &tender.NativeFunction{
			Value: FuncAIR(f.Value.SetY),
		}
	case "set_width":
		res = &tender.NativeFunction{
			Value: FuncAIR(f.Value.SetWidth),
		}
	case "set_height":
		res = &tender.NativeFunction{
			Value: FuncAIR(f.Value.SetHeight),
		}
	}
	return
}

// WUIIntUpDown updates - similar to FloatUpDown but with int values
func (i *WUIIntUpDown) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "set_min":
		res = &tender.NativeFunction{
			Value: FuncAIR(i.Value.SetMin),
		}
	case "set_max":
		res = &tender.NativeFunction{
			Value: FuncAIR(i.Value.SetMax),
		}
	case "set_value":
		res = &tender.NativeFunction{
			Value: FuncAIR(i.Value.SetValue),
		}
	case "value":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(i.Value.Value())}, nil
			},
		}
	case "min":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(i.Value.Min())}, nil
			},
		}
	case "max":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(i.Value.Max())}, nil
			},
		}
	case "min_max":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				min, max := i.Value.MinMax()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(min)},
					&tender.Int{Value: int64(max)},
				}}, nil
			},
		}
	case "set_min_max":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				min, _ := tender.ToInt(args[0])
				max, _ := tender.ToInt(args[1])
				i.Value.SetMinMax(min, max)
				return tender.NullValue, nil
			},
		}
	case "cursor_position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				start, end := i.Value.CursorPosition()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(start)},
					&tender.Int{Value: int64(end)},
				}}, nil
			},
		}
	case "set_cursor_position":
		res = &tender.NativeFunction{
			Value: FuncAIR(i.Value.SetCursorPosition),
		}
	case "select_all":
		res = &tender.NativeFunction{
			Value: FuncAR(i.Value.SelectAll),
		}
	case "set_selection":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				start, _ := tender.ToInt(args[0])
				end, _ := tender.ToInt(args[1])
				i.Value.SetSelection(start, end)
				return tender.NullValue, nil
			},
		}
	case "set_on_tab_focus":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				i.Value.SetOnTabFocus(func(){
					tender.WrapFuncCall(vm, args[0])
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_value_change":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				i.Value.SetOnValueChange(func(value int){
					tender.WrapFuncCall(vm, args[0], &tender.Int{Value: int64(value)})
				})
				return tender.NullValue, nil
			},
		}
	case "on_tab_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return tender.NullValue, nil
			},
		}
	case "on_value_change":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
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
				i.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	// Position and size methods
	case "set_pos":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				x, _ := tender.ToInt(args[0])
				y, _ := tender.ToInt(args[1])
				i.Value.SetPos(x, y)
				return tender.NullValue, nil
			},
		}
	case "set_size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				width, _ := tender.ToInt(args[0])
				height, _ := tender.ToInt(args[1])
				i.Value.SetSize(width, height)
				return tender.NullValue, nil
			},
		}
	case "set_x":
		res = &tender.NativeFunction{
			Value: FuncAIR(i.Value.SetX),
		}
	case "set_y":
		res = &tender.NativeFunction{
			Value: FuncAIR(i.Value.SetY),
		}
	case "set_width":
		res = &tender.NativeFunction{
			Value: FuncAIR(i.Value.SetWidth),
		}
	case "set_height":
		res = &tender.NativeFunction{
			Value: FuncAIR(i.Value.SetHeight),
		}
	case "set_visible":
		res = &tender.NativeFunction{
			Value: FuncABR(i.Value.SetVisible),
		}
	case "visible":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if i.Value.Visible() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	}
	return
}

// WUILabel updates
func (l *WUILabel) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	case "set_text":
		res = &tender.NativeFunction{
			Value: FuncASR(l.Value.SetText),
		}
	case "text":
		res = &tender.NativeFunction{
			Value: FuncARS(l.Value.Text),
		}
	case "set_alignment":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				alignment, _ := tender.ToInt(args[0])
				l.Value.SetAlignment(wui.TextAlignment(alignment))
				return tender.NullValue, nil
			},
		}
	case "alignment":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(l.Value.Alignment())}, nil
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
				l.Value.SetFont(font)
				return tender.NullValue, nil
			},
		}
	case "font":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				font := l.Value.Font()
				return &WUIFont{Value: font}, nil
			},
		}
	case "focus":
		res = &tender.NativeFunction{
			Value: FuncAR(l.Value.Focus),
		}
	case "has_focus":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if l.Value.HasFocus() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
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
				l.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	}
	return
}

// WUIPaintBox updates - needs WUIFont and WUICanvas wrapper types
func (p *WUIPaintBox) IndexGet(index tender.Object) (res tender.Object, err error) {
	strIdx, ok := index.(*tender.String)
	if !ok {
		return nil, tender.ErrInvalidIndexType
	}

	switch strIdx.Value {
	// Setters
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
				p.Value.SetBounds(x, y, width, height)
				return tender.NullValue, nil
			},
		}
	case "set_anchors":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				horizontal, _ := tender.ToInt(args[0])
				vertical, _ := tender.ToInt(args[1])
				p.Value.SetAnchors(wui.Anchor(horizontal), wui.Anchor(vertical))
				return tender.NullValue, nil
			},
		}
	case "set_enabled":
		res = &tender.NativeFunction{
			Value: FuncABR(p.Value.SetEnabled),
		}
	case "set_height":
		res = &tender.NativeFunction{
			Value: FuncAIR(p.Value.SetHeight),
		}
	case "set_horizontal_anchor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				anchor, _ := tender.ToInt(args[0])
				p.Value.SetHorizontalAnchor(wui.Anchor(anchor))
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
				p.Value.SetOnMouseMove(func(x, y int) {
					tender.WrapFuncCall(vm, args[0],
						&tender.Int{Value: int64(x)},
						&tender.Int{Value: int64(y)})
				})
				return tender.NullValue, nil
			},
		}
	case "set_on_paint":
		res = &tender.NativeFunction{
			NeedVMObj: true,
			Value: func(args ...tender.Object) (tender.Object, error) {
				vm := args[0].(*tender.VMObj).Value
				args = args[1:]
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				p.Value.SetOnPaint(func(canvas *wui.Canvas) {
					wrapper := &WUICanvas{Value: canvas}
					tender.WrapFuncCall(vm, args[0], wrapper)
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
				p.Value.SetOnResize(func() {
					tender.WrapFuncCall(vm, args[0])
				})
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
				p.Value.SetPosition(x, y)
				return tender.NullValue, nil
			},
		}
	case "set_size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 2 {
					return nil, tender.ErrWrongNumArguments
				}
				width, _ := tender.ToInt(args[0])
				height, _ := tender.ToInt(args[1])
				p.Value.SetSize(width, height)
				return tender.NullValue, nil
			},
		}
	case "set_vertical_anchor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 1 {
					return nil, tender.ErrWrongNumArguments
				}
				anchor, _ := tender.ToInt(args[0])
				p.Value.SetVerticalAnchor(wui.Anchor(anchor))
				return tender.NullValue, nil
			},
		}
	case "set_visible":
		res = &tender.NativeFunction{
			Value: FuncABR(p.Value.SetVisible),
		}
	case "set_width":
		res = &tender.NativeFunction{
			Value: FuncAIR(p.Value.SetWidth),
		}
	case "set_x":
		res = &tender.NativeFunction{
			Value: FuncAIR(p.Value.SetX),
		}
	case "set_y":
		res = &tender.NativeFunction{
			Value: FuncAIR(p.Value.SetY),
		}

	// Getters
	case "anchors":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				horizontal, vertical := p.Value.Anchors()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(horizontal)},
					&tender.Int{Value: int64(vertical)},
				}}, nil
			},
		}
	case "bounds":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				x, y, width, height := p.Value.Bounds()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(x)},
					&tender.Int{Value: int64(y)},
					&tender.Int{Value: int64(width)},
					&tender.Int{Value: int64(height)},
				}}, nil
			},
		}
	case "enabled":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if p.Value.Enabled() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "handle":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.Handle())}, nil
			},
		}
	case "height":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.Height())}, nil
			},
		}
	case "horizontal_anchor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.HorizontalAnchor())}, nil
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
	case "on_paint":
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
	case "parent":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				parent := p.Value.Parent()
				if parent == nil {
					return tender.NullValue, nil
				}
				return tender.NullValue, nil
			},
		}
	case "position":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				x, y := p.Value.Position()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(x)},
					&tender.Int{Value: int64(y)},
				}}, nil
			},
		}
	case "size":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				width, height := p.Value.Size()
				return &tender.Array{Value: []tender.Object{
					&tender.Int{Value: int64(width)},
					&tender.Int{Value: int64(height)},
				}}, nil
			},
		}
	case "vertical_anchor":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.VerticalAnchor())}, nil
			},
		}
	case "visible":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				if p.Value.Visible() {
					return tender.TrueValue, nil
				}
				return tender.FalseValue, nil
			},
		}
	case "width":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.Width())}, nil
			},
		}
	case "x":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.X())}, nil
			},
		}
	case "y":
		res = &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				if len(args) != 0 {
					return nil, tender.ErrWrongNumArguments
				}
				return &tender.Int{Value: int64(p.Value.Y())}, nil
			},
		}

	// Methods
	case "paint":
		res = &tender.NativeFunction{
			Value: FuncAR(p.Value.Paint),
		}
	}
	return
}


type WUICanvas struct {
	tender.ObjectImpl
	Value *wui.Canvas
}

func (c *WUICanvas) TypeName() string { return "canvas" }
func (c *WUICanvas) String() string   { return "<canvas>" }
func (c *WUICanvas) Copy() tender.Object {
	return &WUICanvas{Value: c.Value}
}

func (c *WUICanvas) IndexGet(index tender.Object) (res tender.Object, err error) {
    strIdx, ok := index.(*tender.String)
    if !ok {
        return nil, tender.ErrInvalidIndexType
    }

    switch strIdx.Value {
    case "arc":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 7 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                width, _ := tender.ToInt(args[2])
                height, _ := tender.ToInt(args[3])
                fromClockAngle, _ := tender.ToFloat64(args[4])
                dAngle, _ := tender.ToFloat64(args[5])
                color, ok := extractColor(args[6])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[6].TypeName(),
                    }
                }
                c.Value.Arc(x, y, width, height, fromClockAngle, dAngle, color)
                return tender.NullValue, nil
            },
        }
    case "clear_draw_regions":
        res = &tender.NativeFunction{
            Value: FuncAR(c.Value.ClearDrawRegions),
        }
    case "draw_ellipse":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 5 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                width, _ := tender.ToInt(args[2])
                height, _ := tender.ToInt(args[3])
                color, ok := extractColor(args[4])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[4].TypeName(),
                    }
                }
                c.Value.DrawEllipse(x, y, width, height, color)
                return tender.NullValue, nil
            },
        }
    case "draw_image":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 4 {
                    return nil, tender.ErrWrongNumArguments
                }
                img, ok := args[0].(*WUIImage)
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "image",
                        Expected: "image",
                        Found:    args[0].TypeName(),
                    }
                }
                rect, ok := args[1].(*WUIRectangle)
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "rectangle",
                        Expected: "rectangle",
                        Found:    args[1].TypeName(),
                    }
                }
                destX, _ := tender.ToInt(args[2])
                destY, _ := tender.ToInt(args[3])
                c.Value.DrawImage(img.Value, rect.Value, destX, destY)
                return tender.NullValue, nil
            },
        }
    case "draw_pie":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 7 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                width, _ := tender.ToInt(args[2])
                height, _ := tender.ToInt(args[3])
                fromClockAngle, _ := tender.ToFloat64(args[4])
                dAngle, _ := tender.ToFloat64(args[5])
                color, ok := extractColor(args[6])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[6].TypeName(),
                    }
                }
                c.Value.DrawPie(x, y, width, height, fromClockAngle, dAngle, color)
                return tender.NullValue, nil
            },
        }
    case "draw_rect":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 5 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                width, _ := tender.ToInt(args[2])
                height, _ := tender.ToInt(args[3])
                color, ok := extractColor(args[4])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[4].TypeName(),
                    }
                }
                c.Value.DrawRect(x, y, width, height, color)
                return tender.NullValue, nil
            },
        }
    case "fill_ellipse":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 5 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                width, _ := tender.ToInt(args[2])
                height, _ := tender.ToInt(args[3])
                color, ok := extractColor(args[4])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[4].TypeName(),
                    }
                }
                c.Value.FillEllipse(x, y, width, height, color)
                return tender.NullValue, nil
            },
        }
    case "fill_pie":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 7 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                width, _ := tender.ToInt(args[2])
                height, _ := tender.ToInt(args[3])
                fromClockAngle, _ := tender.ToFloat64(args[4])
                dAngle, _ := tender.ToFloat64(args[5])
                color, ok := extractColor(args[6])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[6].TypeName(),
                    }
                }
                c.Value.FillPie(x, y, width, height, fromClockAngle, dAngle, color)
                return tender.NullValue, nil
            },
        }
    case "fill_rect":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 5 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                width, _ := tender.ToInt(args[2])
                height, _ := tender.ToInt(args[3])
                color, ok := extractColor(args[4])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[4].TypeName(),
                    }
                }
                c.Value.FillRect(x, y, width, height, color)
                return tender.NullValue, nil
            },
        }
    case "handle":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(c.Value.Handle())}, nil
            },
        }
    case "height":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(c.Value.Height())}, nil
            },
        }
    case "line":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 5 {
                    return nil, tender.ErrWrongNumArguments
                }
                x1, _ := tender.ToInt(args[0])
                y1, _ := tender.ToInt(args[1])
                x2, _ := tender.ToInt(args[2])
                y2, _ := tender.ToInt(args[3])
                color, ok := extractColor(args[4])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[4].TypeName(),
                    }
                }
                c.Value.Line(x1, y1, x2, y2, color)
                return tender.NullValue, nil
            },
        }
    case "polygon":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                arr, ok := args[0].(*tender.Array)
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "points",
                        Expected: "array",
                        Found:    args[0].TypeName(),
                    }
                }
                points := make([]wui.Point, len(arr.Value))
                for i, obj := range arr.Value {
                    pointArr, ok := obj.(*tender.Array)
                    if !ok || len(pointArr.Value) != 2 {
                        return nil, tender.ErrInvalidArgumentType{
                            Name:     fmt.Sprintf("point %d", i),
                            Expected: "array of 2 ints",
                            Found:    obj.TypeName(),
                        }
                    }
                    x, _ := tender.ToInt32(pointArr.Value[0])
                    y, _ := tender.ToInt32(pointArr.Value[1])
                    points[i] = wui.Point{X: x, Y: y}
                }
                color, ok := extractColor(args[1])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[1].TypeName(),
                    }
                }
                c.Value.Polygon(points, color)
                return tender.NullValue, nil
            },
        }
    case "polyline":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                arr, ok := args[0].(*tender.Array)
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "points",
                        Expected: "array",
                        Found:    args[0].TypeName(),
                    }
                }
                points := make([]wui.Point, len(arr.Value))
                for i, obj := range arr.Value {
                    pointArr, ok := obj.(*tender.Array)
                    if !ok || len(pointArr.Value) != 2 {
                        return nil, tender.ErrInvalidArgumentType{
                            Name:     fmt.Sprintf("point %d", i),
                            Expected: "array of 2 ints",
                            Found:    obj.TypeName(),
                        }
                    }
                    x, _ := tender.ToInt32(pointArr.Value[0])
                    y, _ := tender.ToInt32(pointArr.Value[1])
                    points[i] = wui.Point{X: x, Y: y}
                }
                color, ok := extractColor(args[1])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[1].TypeName(),
                    }
                }
                c.Value.Polyline(points, color)
                return tender.NullValue, nil
            },
        }
    case "pop_draw_region":
        res = &tender.NativeFunction{
            Value: FuncAR(c.Value.PopDrawRegion),
        }
    case "push_draw_region":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 4 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                width, _ := tender.ToInt(args[2])
                height, _ := tender.ToInt(args[3])
                c.Value.PushDrawRegion(x, y, width, height)
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
                c.Value.SetFont(font)
                return tender.NullValue, nil
            },
        }
    case "size":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                width, height := c.Value.Size()
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "text_extent":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 1 {
                    return nil, tender.ErrWrongNumArguments
                }
                s, ok := tender.ToString(args[0])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "text",
                        Expected: "string",
                        Found:    args[0].TypeName(),
                    }
                }
                width, height := c.Value.TextExtent(s)
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "text_out":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 4 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                s, ok := tender.ToString(args[2])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "text",
                        Expected: "string",
                        Found:    args[2].TypeName(),
                    }
                }
                color, ok := extractColor(args[3])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[3].TypeName(),
                    }
                }
                c.Value.TextOut(x, y, s, color)
                return tender.NullValue, nil
            },
        }
    case "text_rect":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 6 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                w, _ := tender.ToInt(args[2])
                h, _ := tender.ToInt(args[3])
                s, ok := tender.ToString(args[4])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "text",
                        Expected: "string",
                        Found:    args[4].TypeName(),
                    }
                }
                color, ok := extractColor(args[5])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[5].TypeName(),
                    }
                }
                c.Value.TextRect(x, y, w, h, s, color)
                return tender.NullValue, nil
            },
        }
    case "text_rect_extent":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 2 {
                    return nil, tender.ErrWrongNumArguments
                }
                s, ok := tender.ToString(args[0])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "text",
                        Expected: "string",
                        Found:    args[0].TypeName(),
                    }
                }
                givenWidth, _ := tender.ToInt(args[1])
                width, height := c.Value.TextRectExtent(s, givenWidth)
                return &tender.Array{Value: []tender.Object{
                    &tender.Int{Value: int64(width)},
                    &tender.Int{Value: int64(height)},
                }}, nil
            },
        }
    case "text_rect_format":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 7 {
                    return nil, tender.ErrWrongNumArguments
                }
                x, _ := tender.ToInt(args[0])
                y, _ := tender.ToInt(args[1])
                w, _ := tender.ToInt(args[2])
                h, _ := tender.ToInt(args[3])
                s, ok := tender.ToString(args[4])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "text",
                        Expected: "string",
                        Found:    args[4].TypeName(),
                    }
                }
                format, ok := extractFormat(args[5])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "format",
                        Expected: "format",
                        Found:    args[5].TypeName(),
                    }
                }
                color, ok := extractColor(args[6])
                if !ok {
                    return nil, tender.ErrInvalidArgumentType{
                        Name:     "color",
                        Expected: "color",
                        Found:    args[6].TypeName(),
                    }
                }
                c.Value.TextRectFormat(x, y, w, h, s, format, color)
                return tender.NullValue, nil
            },
        }
    case "width":
        res = &tender.NativeFunction{
            Value: func(args ...tender.Object) (tender.Object, error) {
                if len(args) != 0 {
                    return nil, tender.ErrWrongNumArguments
                }
                return &tender.Int{Value: int64(c.Value.Width())}, nil
            },
        }
    }
    return
}


// Image wrapper
type WUIImage struct {
	tender.ObjectImpl
	Value *wui.Image
}

func (i *WUIImage) TypeName() string { return "image" }
func (i *WUIImage) String() string   { return "<image>" }
func (i *WUIImage) Copy() tender.Object {
	return &WUIImage{Value: i.Value}
}

// Rectangle wrapper
type WUIRectangle struct {
	tender.ObjectImpl
	Value wui.Rectangle
}

func (r *WUIRectangle) TypeName() string { return "rectangle" }
func (r *WUIRectangle) String() string   { 
	x, y, w, h := r.Value.X, r.Value.Y, r.Value.Width, r.Value.Height
	return fmt.Sprintf("<rectangle x:%d y:%d w:%d h:%d>", x, y, w, h)
}
func (r *WUIRectangle) Copy() tender.Object {
	return &WUIRectangle{Value: r.Value}
}

// Format wrapper
type WUIFormat struct {
	tender.ObjectImpl
	Value wui.Format
}

func (f *WUIFormat) TypeName() string { return "format" }
func (f *WUIFormat) String() string   { return "<format>" }
func (f *WUIFormat) Copy() tender.Object {
	return &WUIFormat{Value: f.Value}
}

// Helper functions for extraction

func extractColor(obj tender.Object) (wui.Color, bool) {
	switch c := obj.(type) {
	case *WUIColor:
		return c.Value, true
	default:
		// Try to create color from integers or other representations
		// For now, return a default color and false
		return wui.Color(0), false
	}
}

func extractFormat(obj tender.Object) (wui.Format, bool) {
	switch f := obj.(type) {
	case *WUIFormat:
		return f.Value, true
	default:
		return wui.Format(0), false
	}
}

// Function to create RGB color
func wuiRGB(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 3 {
		err = tender.ErrWrongNumArguments
		return
	}

	r, ok1 := tender.ToInt(args[0])
	g, ok2 := tender.ToInt(args[1])
	b, ok3 := tender.ToInt(args[2])
	
	if !ok1 || !ok2 || !ok3 {
		err = tender.ErrInvalidArgumentType{
			Name:     "color components",
			Expected: "integers",
			Found:    fmt.Sprintf("%s, %s, %s", args[0].TypeName(), args[1].TypeName(), args[2].TypeName()),
		}
		return
	}

	color := wui.RGB(uint8(r), uint8(g), uint8(b))
	return &WUIColor{Value: color}, nil
}

// Function to create rectangle
func wuiRect(args ...tender.Object) (ret tender.Object, err error) {
	if len(args) != 4 {
		err = tender.ErrWrongNumArguments
		return
	}

	x, ok1 := tender.ToInt(args[0])
	y, ok2 := tender.ToInt(args[1])
	width, ok3 := tender.ToInt(args[2])
	height, ok4 := tender.ToInt(args[3])
	
	if !ok1 || !ok2 || !ok3 || !ok4 {
		err = tender.ErrInvalidArgumentType{
			Name:     "rectangle components",
			Expected: "integers",
			Found:    fmt.Sprintf("%s, %s, %s, %s", args[0].TypeName(), args[1].TypeName(), args[2].TypeName(), args[3].TypeName()),
		}
		return
	}

	rect := wui.Rect(x, y, width, height)
	return &WUIRectangle{Value: rect}, nil
}