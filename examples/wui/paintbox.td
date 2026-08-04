// tender@v0.0.9
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