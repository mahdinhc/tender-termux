# colors Module

The `colors` module provides IOWriter for printing colored text to the terminal.

## Constants
- `colors.stdout`: `&stdlib.IOWriter{}`
- `colors.stderr`: `&stdlib.IOWriter{}`

### Example Usage

```go
import "fmt"
import "colors"

fmt.fprint(colors.stdout, "Hello".red, "World".green, "\n")
fmt.fprintln(colors.stderr, "Hello".red, "World".green)
```