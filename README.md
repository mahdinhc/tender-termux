# Tender

**Tender** is a general-purpose programming language specially designed for image processing, 2D graphics, scripting, and more!

## Overview

Tender is compiled and executed as bytecode on a stack-based virtual machine (VM) written in native Golang.

## Features
- **Simple and highly readable syntax**  
- **Compiles into bytecode**  
- **Supports rich built-in modules**  
- **Optimized for 2D graphics**  

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
ctx.hex("#0f0")        // Set color to green
ctx.dash(4, 2)         // Define dashed stroke
ctx.rect(25, 25, 50, 50)  // Draw a rectangle
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

1. **Download the repository.**  
2. **Run the `install.sh` script** to install Tender on your system.

### Windows (Manual Installation)
1. **Download the `tender.exe` binary** for your system from the [bin](bin) directory, along with the `pkg` folder from [pkg](pkg).  
2. **Copy the files to your desired location** with the following structure:

```bash
├───bin
│   └───tender.exe
└───pkg
    │   ansi.td
    │   cinf.td
    │   console.td
    │   enum.td
    │   fs.td
    │   matrix.td
    │   messagebox.td
    │   utf8.td
    │   vec2.td
    │   xml.td
    └───helper
```

3. **Add the path to the `bin` folder** to your system's environment variables.

---

## Documentation

- **[Runtime Types](https://github.com/2dprototype/tender/blob/master/docs/runtime-types.md)**  
- **[Built-ins](https://github.com/2dprototype/tender/blob/master/docs/builtins.md)**  
- **[Operators](https://github.com/2dprototype/tender/blob/master/docs/operators.md)**  
- **[Standard Library](https://github.com/2dprototype/tender/blob/master/docs/stdlib.md)**  

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

Syntax highlighting is currently available only for **Notepad++**. Download the configuration file from [here](misc/syntax/npp_tender.xml).

---

## License

Tender is distributed under the [MIT License](LICENSE). Additional licenses for third-party dependencies are provided in [LICENSE_GOLANG](LICENSE_GOLANG) and [LICENSE_TENGO](LICENSE_TENGO).

---

## Acknowledgments

Tender is written in Go and inspired by **Tengo**. We extend our gratitude to the contributors of Tengo for their valuable work.
