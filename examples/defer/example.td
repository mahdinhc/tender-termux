// tender@v0.0.8
import "fmt"

fn main() {
    fmt.println("starting function...")

    defer fmt.println("first defer written (runs last)")
    defer fmt.println("second defer written (runs second)")
    defer fmt.println("third defer written (runs first)")

    fmt.println("finished main logic.")
}

main()