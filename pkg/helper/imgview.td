import "os"
import "fs"
import "path"

TEMP := os.getenv("temp")
SCRIPT := path.join(TEMP, "temp_quick_img_view.td")
IMG_SRC := path.join(TEMP, "temp_quick_img.png")

script := `
import "canvas"
import "image"
import "math"
import "os"

// Load and decode the image once
img_bytes := embed("temp_quick_img.png")
img_obj := image.decode(img_bytes)
if is_error(img_obj) {
    println(img_obj)
    os.exit(0)
} 
img_size := img_obj.bounds().size
img_w := img_size.width
img_h := img_size.height

w := 600
h := 400

fn draw(ctx, win_w, win_h) {
    ctx.hex("#000")
    ctx.clear()

    // Compute scale to fit inside the window (no stretching)
    scale_x := float(win_w) / float(img_w)
    scale_y := float(win_h) / float(img_h)
    scale := math.min(scale_x, scale_y)

    draw_w := float(img_w) * scale
    draw_h := float(img_h) * scale
    x := (float(win_w) - draw_w) / 2.0
    y := (float(win_h) - draw_h) / 2.0

    // Draw the scaled image centered
    ctx.push()
    ctx.translate(x, y)
    ctx.scale(scale, scale)
    ctx.drawimage(img_bytes, 0, 0)
    ctx.pop()
}

canvas.new_window(w, h, "Quick Image View", fn(window){
    ctx := window.new_context(w, h)
    draw(ctx, w, h)
    window.update(w, h)

    for {
        e := window.next_event()
        if e.kind == "size" {
            w = e.width_px
            h = e.height_px
            ctx = window.new_context(w, h)
            draw(ctx, w, h)
            window.update(w, h)
        }
        else if e.kind == "lifecycle" && e.prev == 3 && e.next == 0 {
            break
        }
    }
})

`

export fn(img) {
	fs.writefile(SCRIPT, script)
	fs.writefile(IMG_SRC, img)
	cmd := os.exec("tender", SCRIPT)
	cmd.output() |> string |> println
	fs.remove(SCRIPT)
	fs.remove(IMG_SRC)
}