# Stdlib colors

The `colors` module provides IOWriter for printing colored text to the terminal.

## Functions
- `colors.stdout()`: returns `&stdlib.IOWriter{}`
- `colors.stderr()`: returns `&stdlib.IOWriter{}`

### Example Usage

```go
import "fmt"
import "colors"

fmt.fprint(colors.stdout(), "Hello".red, "World".green, "\n")
fmt.fprintln(colors.stderr(), "Hello".red, "World".green)
```