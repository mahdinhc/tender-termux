import { tcell } from "console"

// initialize the raw tcell screen
var screen = tcell.new_screen()
screen.init()

// application state
var running = true
var player_x = 10
var player_y = 5

// define styles using the supported map properties
var player_style = {
    "fg": "cyan",
    "bold": true,
    "blink": true
}

var border_style = {
    "fg": "#00ff00",
    "bg": "black"
}

// helper function to draw a horizontal line
fn draw_hline(y, width, ch) {
    for x in range(0, width) {
        screen.set_content(x, y, ch, border_style)
    }
}

// main event loop
for running {
    // clear previous frame
    screen.clear()
    
    // fetch dynamic terminal size
    var size = screen.size()
    var w = size.width
    var h = size.height
    
    // render boundaries
    draw_hline(0, w, '-')
    draw_hline(h - 1, w, '-')
    
    // render player
    screen.set_content(player_x, player_y, '+', player_style)
    
    // present the buffer to the terminal
    screen.show()
    
    // poll for events (this blocks until input or resize occurs)
    var ev = screen.poll_event()
    
    if ev != null {
        if ev.kind == "key" {
            var key_name = ev.name
            var ch = ev.rune
            
            // handle exit
            if key_name == "esc" || key_name == "ctrl+c" || ch == 'q' {
                running = false
            }
            
            // handle movement (with collision bounds)
            if key_name == "up" || ch == 'w' {
                if player_y > 1 { 
                    player_y -= 1 
                }
            } 
            else if key_name == "down" || ch == 's' {
                if player_y < h - 2 { 
                    player_y += 1 
                }
            } 
            else if key_name == "left" || ch == 'a' {
                if player_x > 0 { 
                    player_x -= 1 
                }
            } 
            else if key_name == "right" || ch == 'd' {
                if player_x < w - 1 { 
                    player_x += 1 
                }
            }
            
        } 
        else if ev.kind == "resize" {
            // force terminal state synchronization when the window changes size
            screen.sync()
        }
    }
}

// teardown
// it is critical to call fini() to restore the terminal's original state (cursor visibility, cooked mode)
screen.fini()