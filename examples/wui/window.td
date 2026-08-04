// tender@v0.0.9
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