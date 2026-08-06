package stdlib

import (
	"github.com/2dprototype/tender"
	"github.com/gdamore/tcell/v2"
)

var tcellModule = &tender.ImmutableMap{
	Value: map[string]tender.Object{
		"new_screen": &tender.NativeFunction{
			Value: func(args ...tender.Object) (tender.Object, error) {
				s, err := tcell.NewScreen()
				if err != nil {
					return wrapError(err), nil
				}

				methods := map[string]tender.Object{
					"init": &tender.NativeFunction{
						Value: FuncARE(s.Init),
					},
					"fini": &tender.NativeFunction{
						Value: FuncAR(s.Fini),
					},
					"clear": &tender.NativeFunction{
						Value: FuncAR(s.Clear),
					},
					"show": &tender.NativeFunction{
						Value: FuncAR(s.Show),
					},
					"sync": &tender.NativeFunction{
						Value: FuncAR(s.Sync),
					},
					"set_content": &tender.NativeFunction{
						Value: func(args ...tender.Object) (tender.Object, error) {
							if len(args) < 4 { return nil, tender.ErrWrongNumArguments }
							x, _ := tender.ToInt(args[0])
							y, _ := tender.ToInt(args[1])
							ch, _ := tender.ToRune(args[2])
							
							style := tcell.StyleDefault
							if styleMap, ok := args[3].(*tender.Map); ok {
								style = buildTCellStyle(styleMap)
							}
							s.SetContent(x, y, ch, nil, style)
							return tender.NullValue, nil
						},
					},
					"size": &tender.NativeFunction{
						Value: func(args ...tender.Object) (tender.Object, error) {
							w, h := s.Size()
							return &tender.Map{
								Value: map[string]tender.Object{
									"width":  &tender.Int{Value: int64(w)},
									"height": &tender.Int{Value: int64(h)},
								},
							}, nil
						},
					},
					"poll_event": &tender.NativeFunction{
						Value: func(args ...tender.Object) (tender.Object, error) {
							ev := s.PollEvent()
							switch ev := ev.(type) {
							case *tcell.EventKey:
								return &tender.Map{
									Value: map[string]tender.Object{
										"kind": &tender.String{Value: "key"},
										"key":  &tender.Int{Value: int64(ev.Key())},
										"rune": &tender.Char{Value: ev.Rune()},
										"name": &tender.String{Value: ev.Name()},
									},
								}, nil
							case *tcell.EventResize:
								w, h := ev.Size()
								return &tender.Map{
									Value: map[string]tender.Object{
										"kind":   &tender.String{Value: "resize"},
										"width":  &tender.Int{Value: int64(w)},
										"height": &tender.Int{Value: int64(h)},
									},
								}, nil
							}
							return tender.NullValue, nil
						},
					},
				}
				return &tender.ImmutableMap{Value: methods}, nil
			},
		},
	},
}

// Helper function to dynamically map tender maps to tcell styling
func buildTCellStyle(m *tender.Map) tcell.Style {
	st := tcell.StyleDefault
	for k, v := range m.Value {
		switch k {
		case "color", "fg", "foreground":
			if s, ok := tender.ToString(v); ok {
				st = st.Foreground(tcell.GetColor(s))
			}
		case "bg", "background":
			if s, ok := tender.ToString(v); ok {
				st = st.Background(tcell.GetColor(s))
			}
		case "bold":
			if b, _ := tender.ToBool(v); b {
				st = st.Bold(true)
			}
		case "underline":
			if b, _ := tender.ToBool(v); b {
				st = st.Underline(true)
			}
		case "blink":
			if b, _ := tender.ToBool(v); b {
				st = st.Blink(true)
			}
		case "reverse":
			if b, _ := tender.ToBool(v); b {
				st = st.Reverse(true)
			}
		case "italic":
			if b, _ := tender.ToBool(v); b {
				st = st.Italic(true)
			}
		case "strike":
			if b, _ := tender.ToBool(v); b {
				st = st.StrikeThrough(true)
			}
		}
	}
	return st
}