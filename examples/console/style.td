// tender@v0.0.8
import "console"

// Create red text
red_text := console.style("Warning!", {color: "#f00", bold: true})

// Create bordered panel
panel := console.style("Content", {
  border: "rounded",
  padding: 2,
  color: "#00FF22",
  background: "#222222"
})

// Print to colorable stdout
println(red_text)
println(panel)