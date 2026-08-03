// tender@v0.0.4
canvas := import("canvas")
canvas.load_font("testfont", 12, "IBMPlexSans-Bold.ttf")

ctx := canvas.new_context(100, 100)

ctx.set_font("testfont")
ctx.hex("#00ff00")
ctx.text("Wooh", 10, 10)

ctx.save_png("out.png")