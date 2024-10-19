# Tender

**Tender** is a general-purpose programming language specially designed for image processing, 2D graphics, scripting, and more! Here is a quick [tutorial](pages/tutorial.md).

## Overview

Tender compiles into bytecode and executes on a stack-based virtual machine (VM) written in native Golang.

## Features
- **Simple and highly readable syntax**  
- **Compiles to bytecode**  
- **Supports rich [built-in functions](pages/builtins.md)**  
- **Includes an extensive [standard library](pages/stdlib.md)**  
- **Optimized for 2D graphics**  

### Supported Standard Library

- [math](pages/stdlib-math.md): Mathematical constants and functions  
- [os](pages/stdlib-os.md): Platform-independent interface to OS functionality  
- [strings](pages/stdlib-strings.md): String conversion, manipulation, and regular expressions  
- [times](pages/stdlib-times.md): Time-related functions  
- [rand](pages/stdlib-rand.md): Random number generation  
- [fmt](pages/stdlib-fmt.md): Formatting functions  
- [json](pages/stdlib-json.md): JSON handling functions  
- [base64](pages/stdlib-base64.md): Base64 encoding and decoding  
- [hex](pages/stdlib-hex.md): Hexadecimal encoding and decoding  
- [colors](pages/stdlib-colors.md): Functions to print colored text to the terminal  
- [gzip](pages/stdlib-gzip.md): Gzip compression and decompression  
- [zip](pages/stdlib-zip.md): ZIP archive manipulation  
- [tar](pages/stdlib-tar.md): TAR archive creation and reading  
- [bufio](pages/stdlib-bufio.md): Buffered I/O functions  
- [crypto](pages/stdlib-crypto.md): Cryptographic functions  
- [path](pages/stdlib-path.md): File path manipulation  
- [image](pages/stdlib-image.md): Image manipulation  
- [canvas](pages/stdlib-canvas.md): Drawing functions for canvases  
- [dll](pages/stdlib-dll.md): Dynamic link library interactions  
- [io](pages/stdlib-io.md): Input and output functions  
- [audio](pages/stdlib-audio.md): Audio processing  
- [net](pages/stdlib-net.md): Networking functions  
- [http](pages/stdlib-http.md): HTTP client and server utilities  
- [websocket](pages/stdlib-websocket.md): WebSocket communication utilities  
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

- **[Runtime Types](pages/runtime-types.md)**  
- **[Built-in Functions](pages/builtins.md)**  
- **[Operators](pages/operators.md)**  
- **[Standard Library](pages/stdlib.md)**  

## Examples

Explore various examples demonstrating Tender’s features in the [examples](https://github.com/2dprototype/tender/blob/main/examples) directory.

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

Syntax highlighting is currently available only for **Notepad++**. Download the configuration file [here](https://github.com/2dprototype/tender/blob/main/misc/syntax/npp_tender.xml).

---

## License

Tender is distributed under the [MIT License](https://github.com/2dprototype/tender/blob/main/LICENSE), with additional licenses provided for third-party dependencies. See [LICENSE_GOLANG](https://github.com/2dprototype/tender/blob/main/LICENSE_GOLANG) and [LICENSE_TENGO](https://github.com/2dprototype/tender/blob/main/LICENSE_TENGO) for more information.

---

## Acknowledgments

Tender is written in Go, based on Tengo. We extend our gratitude to the contributors of Tengo for their valuable work.