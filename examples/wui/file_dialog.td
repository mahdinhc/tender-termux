// tender@v0.0.9
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