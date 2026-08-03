import "fmt"
import "console"
import "os"

fmt.fprint(console.stdout(), "Hello".red, "World".green, "\n")

fmt.fprintln(console.stderr(), "Hello".red, "World".green)

fmt.fprintln(os.stdout(), "Hello".red, "World".green, "\n")