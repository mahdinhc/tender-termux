// tender@v0.0.9
import "wui"
import "image"

fn get_img_data() {
	img := image.new(255, 255)
	for i := 0; i < 255; i++ {
		for j := 0; j < 255; j++ {
			color := [i, 0, j, 255]
			img.set(i, j, color)
		}
	}	
	return img.encode("png")
}

var w, h = 500, 400

win := wui.new_window()
win.set_title("Drawing Example")
win.set_size(w, h)

pb := wui.new_paintbox()
pb.set_bounds(0, 0, w, h)

img := wui.new_image(get_img_data())
pb.set_on_paint(fn(canvas) {
    canvas.fill_rect(0, 0, w/2, h/2, wui.rgb_color(255, 0, 0))
    canvas.fill_ellipse(300, 50, 100, 100, wui.rgb_color(0, 0, 255))
    canvas.set_font(wui.new_font({name: "Arial", height: 24}))
    canvas.text_out(w/2, h/2, "Hello Canvas!", wui.rgb_color(255,255,255))
    canvas.line(50, 300, 400, 350, wui.rgb_color(0, 255, 0))
	canvas.draw_image(img, wui.rect(10,10,200,100), 10,10)
})

win.set_on_resize(fn() {
	w, h = win.height(), win.width()
	pb.paint()
})

win.add(pb)
win.show()