// tender@v0.0.8
import "canvas"

w := 600
h := 400


fn draw(ctx) {
	ctx.hex(`#121212`)
	ctx.clear()
	ctx.hex("#0f0")
	ctx.dash(4)
	ctx.rect(w/2-50, h/2-50, 100, 100)
	// ctx.drawimage(embed("logo.png"), w/2-50, h/2-50)
	ctx.stroke()	
}

canvas.new_window(w, h, "Hello Tender", fn(window){
	ctx := window.new_context(w, h)
	draw(ctx)
	window.update(w, h)
	for {
		e := window.next_event()
		if e.kind == "size" {
			w = e.width_px 
			h = e.height_px
			ctx = window.new_context(e.width_px, e.height_px)
			draw(ctx)
			window.update(w, h)
		}
		else if e.kind == "lifecycle" && e.prev == 3 && e.next == 0 {
			break
		}
	}
})
	
