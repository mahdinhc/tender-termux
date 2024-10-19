# Tender

**Tender** is a general-purpose programming language specially designed for image processing, 2D graphics, scripting, and more! Here is a quick [tutorial](docs/pages/tutorial.md). Also check the [docs](https://2dprototype.github.io/tender)!

## Overview

Tender compiles into bytecode and executes on a stack-based virtual machine (VM) written in native Golang.

## Features
- **Simple and highly readable syntax**  
- **Compiles to bytecode**  
- **Supports rich [built-in functions](docs/pages/builtins.md)**  
- **Includes an extensive [standard library](docs/pages/stdlib.md)**  
- **Optimized for 2D graphics**  

### Supported Standard Library

- [math](docs/pages/stdlib-math.md): Mathematical constants and functions  
- [os](docs/pages/stdlib-os.md): Platform-independent interface to OS functionality  
- [strings](docs/pages/stdlib-strings.md): String conversion, manipulation, and regular expressions  
- [times](docs/pages/stdlib-times.md): Time-related functions  
- [rand](docs/pages/stdlib-rand.md): Random number generation  
- [fmt](docs/pages/stdlib-fmt.md): Formatting functions  
- [json](docs/pages/stdlib-json.md): JSON handling functions  
- [base64](docs/pages/stdlib-base64.md): Base64 encoding and decoding  
- [hex](docs/pages/stdlib-hex.md): Hexadecimal encoding and decoding  
- [colors](docs/pages/stdlib-colors.md): Functions to print colored text to the terminal  
- [gzip](docs/pages/stdlib-gzip.md): Gzip compression and decompression  
- [zip](docs/pages/stdlib-zip.md): ZIP archive manipulation  
- [tar](docs/pages/stdlib-tar.md): TAR archive creation and reading  
- [bufio](docs/pages/stdlib-bufio.md): Buffered I/O functions  
- [crypto](docs/pages/stdlib-crypto.md): Cryptographic functions  
- [path](docs/pages/stdlib-path.md): File path manipulation  
- [image](docs/pages/stdlib-image.md): Image manipulation  
- [canvas](docs/pages/stdlib-canvas.md): Drawing functions for canvases  
- [dll](docs/pages/stdlib-dll.md): Dynamic link library interactions  
- [io](docs/pages/stdlib-io.md): Input and output functions  
- [audio](docs/pages/stdlib-audio.md): Audio processing  
- [net](docs/pages/stdlib-net.md): Networking functions  
- [http](docs/pages/stdlib-http.md): HTTP client and server utilities  
- [websocket](docs/pages/stdlib-websocket.md): WebSocket communication utilities  
- **gob**: Gob Encoding/Ddecoding
- **csv**: CSV Encoding/Ddecoding

## Quick Start

1. **Install Tender on your machine.**  
2. **Copy the sample code below:**

```go
// Basic example
str1 := "hello"
str2 := "world"

println(str1 + " " + str2)
```

```go
// Canvas drawing example (similar to JS Canvas)
import "canvas"
	
var ctx = canvas.new_context(100, 100)
ctx.hex("#0f0")          // Set color to green
ctx.dash(4, 2)           // Define dashed stroke
ctx.rect(25, 25, 50, 50) // Draw a rectangle
ctx.stroke()

ctx.save_png("out.png")  // Save output as PNG
```

3. **Save your code as `hello.td`** (use the `.td` extension).  
4. **Run your script using the following command:**

```bash
tender hello.td
```

---

## Installation

### Using Go

1. Install the latest version of Go.  
2. Run the following command to install:

```bash
go install github.com/2dprototype/tender/cmd/tender@latest
```

### Manual Installation (Windows)

Precompiled binaries are available. Download them from the release tags.

---

## Documentation
Check the [docs](https://2dprototype.github.io/tender)!

- **[Runtime Types](docs/pages/runtime-types.md)**  
- **[Built-in Functions](docs/pages/builtins.md)**  
- **[Operators](docs/pages/operators.md)**  
- **[Standard Library](docs/pages/stdlib.md)**  

## Examples

Explore various examples demonstrating Tender’s features in the [examples](examples) directory.

---

## Dependencies

Tender uses the following dependencies:

- [go-mp3](https://github.com/hajimehoshi/go-mp3)  
- [gorilla/websocket](https://github.com/gorilla/websocket)  
- [ebitengine/oto/v3](https://github.com/ebitengine/oto/v3)  
- [exp/shiny](https://pkg.go.dev/golang.org/x/exp/shiny)  
- [fogleman/gg](https://github.com/fogleman/gg)  

---

## Syntax Highlighting

Syntax highlighting is currently available only for **Notepad++**. Download the configuration file [here](misc/syntax/npp_tender.xml).

---

## License

Tender is distributed under the [MIT License](LICENSE), with additional licenses provided for third-party dependencies. See [LICENSE_GOLANG](LICENSE_GOLANG) and [LICENSE_TENGO](LICENSE_TENGO) for more information.

---

## Acknowledgments

Tender is written in Go, based on Tengo. We extend our gratitude to the contributors of Tengo for their valuable work.