# Tender

**Tender** is a general-purpose programming language specially designed for image processing, 2D graphics, scripting, and more!

## Overview

Tender compiles into bytecode and executes on a stack-based virtual machine (VM) written in native Golang.

## Features
- **Simple and highly readable syntax**  
- **Compiles to bytecode**  
- **Supports rich [built-in functions](docs/builtins.md)**  
- **Includes an extensive [standard library](docs/stdlib.md)**  
- **Optimized for 2D graphics**  

### Supported Standard Library

- [math](docs/stdlib-math.md): Mathematical constants and functions  
- [os](docs/stdlib-os.md): Platform-independent interface to OS functionality  
- [strings](docs/stdlib-strings.md): String conversion, manipulation, and regular expressions  
- [times](docs/stdlib-times.md): Time-related functions  
- [rand](docs/stdlib-rand.md): Random number generation  
- [fmt](docs/stdlib-fmt.md): Formatting functions  
- [json](docs/stdlib-json.md): JSON handling functions  
- [base64](docs/stdlib-base64.md): Base64 encoding and decoding  
- [hex](docs/stdlib-hex.md): Hexadecimal encoding and decoding  
- [colors](docs/stdlib-colors.md): Functions to print colored text to the terminal  
- [gzip](docs/stdlib-gzip.md): Gzip compression and decompression  
- [zip](docs/stdlib-zip.md): ZIP archive manipulation  
- [tar](docs/stdlib-tar.md): TAR archive creation and reading  
- [bufio](docs/stdlib-bufio.md): Buffered I/O functions  
- [crypto](docs/stdlib-crypto.md): Cryptographic functions  
- [path](docs/stdlib-path.md): File path manipulation  
- [image](docs/stdlib-image.md): Image manipulation  
- [canvas](docs/stdlib-canvas.md): Drawing functions for canvases  
- [dll](docs/stdlib-dll.md): Dynamic link library interactions  
- [io](docs/stdlib-io.md): Input and output functions  
- [audio](docs/stdlib-audio.md): Audio processing  
- [net](docs/stdlib-net.md): Networking functions  
- [http](docs/stdlib-http.md): HTTP client and server utilities  
- [websocket](docs/stdlib-websocket.md): WebSocket communication utilities  
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

- **[Runtime Types](docs/runtime-types.md)**  
- **[Built-in Functions](docs/builtins.md)**  
- **[Operators](docs/operators.md)**  
- **[Standard Library](docs/stdlib.md)**  

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