graphics := import("graphics")

win := graphics.new_window(400, 400, "Tender Event System Demo")

px := 400
py := 300

win.on_keydown(fn(key, x, y) {
    println("Key Down:", key)
})

win.on_keyup(fn(key, x, y) {
    println("Key Up:", key)
})

win.on_keyhold(fn(key) {
    if key == "w" || key == "up" {
        py -= 3
    }
    if key == "s" || key == "down" {
        py += 3
    }
    if key == "a" || key == "left" {
        px -= 3
    }
    if key == "d" || key == "right" {
        px += 3
    }
})

var X, Y, col = 10, 20, "#0f0"

win.on_mousedown(fn(btn, x, y) {
    col = "#f00"
    println("Mouse Down:", btn, "at", x, y)
})

win.on_mouseup(fn(btn, x, y) {
    col = "#ff0"
    println("Mouse Up:", btn, "at", x, y)
})


win.on_mousemove(fn(x, y) {
    X, Y = x, y
})

win.on_mouseenter(fn() {
    col = "#0f0"
    println("Mouse Entered Window")
})

win.on_mouseleave(fn() {
    println("Mouse Left Window")
})

win.on_resize(fn(w, h) {
    println("Window Resized to:", w, h)
})

win.on_draw(fn() {
    win.clear(0.12, 0.12, 0.15)
    win.rgb(0.2, 0.7, 1.0)
    win.rect(px - 20, py - 20, 40, 40)
    win.fill()
    win.hex(col)
    win.circle(X, Y, 5)
    win.fill()
})

win.run()
