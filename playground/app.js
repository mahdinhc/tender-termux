// Tender WASM Playground Application Controller

document.addEventListener("DOMContentLoaded", () => {
    const codeEditor = document.getElementById("codeEditor");
    const lineNumbers = document.getElementById("lineNumbers");
    const consoleOutput = document.getElementById("consoleOutput");
    const astOutput = document.getElementById("astOutput");
    const statusDot = document.getElementById("statusDot");
    const statusText = document.getElementById("statusText");
    const runBtn = document.getElementById("runBtn");
    const parseBtn = document.getElementById("parseBtn");
    const clearBtn = document.getElementById("clearBtn");
    const exampleSelect = document.getElementById("exampleSelect");
    const charCount = document.getElementById("charCount");
    const modulesGrid = document.getElementById("modulesGrid");

    let capturedLogs = "";

    // Sample Scripts
    const SAMPLES = {
        hello: `fmt := import("fmt")

println("Welcome to the WASM Tender Playground!")
println("======================================")

name := "Developer"
message := fmt.sprintf("Hello, %s! Tender is running natively inside WebAssembly.", name)
println(message)`,

        fibonacci: `fib := fn(n) {
    if n <= 1 { return n }
    return fib(n - 1) + fib(n - 2)
}

println("Calculating Fibonacci sequence up to N=12:")
for i := 0; i <= 12; i++ {
    println("fib(" + string(i) + ") =", fib(i))
}`,

        matrix: `// Tender Matrix Operations
m1 := matrix([
    [1, 2, 3],
    [4, 5, 6]
])

println("Matrix m1:")
println(m1)

m2 := matrix([
    [7, 8],
    [9, 1],
    [2, 3]
])

println("Matrix m2:")
println(m2)

println("Matrix Multiplication (m1 * m2):")
println(m1 * m2)`,

        math_complex: `math := import("math")
cmplx := import("cmplx")

println("Math constants:")
println("Pi =", math.pi)
println("E =", math.e)

println("Sqrt(144) =", math.sqrt(144))
println("2^10 =", math.pow(2, 10))

z1 := complex(3, 4)
println("Complex number z1 =", z1)
println("Abs(z1) =", cmplx.abs(z1))`,

        strings_regex: `strings := import("strings")

text := "  Tender Programming Language WASM Playground  "
trimmed := strings.trim_space(text)
println("Trimmed:", trimmed)
println("Upper:", strings.to_upper(trimmed))
println("Contains 'WASM':", strings.contains(trimmed, "WASM"))

words := strings.split(trimmed, " ")
println("Word count:", len(words))
println("Words array:", words)`,

        json_base64: `json := import("json")
base64 := import("base64")

data := {
    "title": "WASM Tender Playground",
    "version": "1.0",
    "supported_modules": ["math", "json", "canvas", "crypto", "strings"]
}

encodedJSON := json.encode(data)
println("Encoded JSON:", string(encodedJSON))

encodedB64 := base64.encode(encodedJSON)
println("Base64 Encoded:", encodedB64)

decodedB64 := base64.decode(encodedB64)
println("Base64 Decoded:", string(decodedB64))`,

        structs: `User := struct {
    name string
    role string
    score int
}

user1 := User{name: "Alice", role: "Admin", score: 95}
println("User Struct:", user1)
println("User Name:", user1.name)
println("User Score:", user1.score)`,

        concurrency: `// Goroutines & Channels in Tender
ch := chan(5)

go fn() {
    for i := 1; i <= 3; i++ {
        ch <- i * 10
    }
}()

println("Reading values from channel:")
println("Channel Value 1:", <-ch)
println("Channel Value 2:", <-ch)
println("Channel Value 3:", <-ch)`,

        canvas: `canvas := import("canvas")

ctx := canvas.new_context(200, 200)
println("Created Canvas:", ctx)`
    };

    // Intercept standard output from WebAssembly
    if (globalThis.fs) {
        const origWriteSync = globalThis.fs.writeSync;
        const textDecoder = new TextDecoder("utf-8");
        globalThis.fs.writeSync = function(fd, buf) {
            if (fd === 1 || fd === 2) {
                const text = textDecoder.decode(buf);
                capturedLogs += text;
            }
            return origWriteSync ? origWriteSync.apply(this, arguments) : buf.length;
        };
    }

    // Code Editor Line Numbers & Character Count Update
    function updateEditor() {
        const text = codeEditor.value;
        const lines = text.split("\n").length;
        lineNumbers.innerHTML = Array.from({
            length: lines
        }, (_, i) => i + 1).join("<br>");
        charCount.textContent = `${text.length} chars`;
    }

    codeEditor.addEventListener("input", updateEditor);
    codeEditor.addEventListener("scroll", () => {
        lineNumbers.scrollTop = codeEditor.scrollTop;
    });

    // Example Selector
    exampleSelect.addEventListener("change", (e) => {
        const key = e.target.value;
        if (SAMPLES[key]) {
            codeEditor.value = SAMPLES[key];
            updateEditor();
        }
    });

    // Load Initial Code
    codeEditor.value = SAMPLES.hello;
    updateEditor();

    // Tabs logic
    const tabBtns = document.querySelectorAll(".tab-btn");
    const tabContents = document.querySelectorAll(".tab-content");

    tabBtns.forEach((btn) => {
        btn.addEventListener("click", () => {
            tabBtns.forEach((b) => b.classList.remove("active"));
            tabContents.forEach((c) => c.classList.remove("active"));
            btn.classList.add("active");
            const target = document.getElementById(btn.dataset.tab);
            if (target) target.classList.add("active");
        });
    });

    // Populate Modules Tab
    function renderModules() {
        if (typeof globalThis.tenderGetModules !== "function") return;
        try {
            let modulesRaw = globalThis.tenderGetModules();
            let modules = typeof modulesRaw === "string" ? JSON.parse(modulesRaw) : modulesRaw;
            modulesGrid.innerHTML = "";
            if (Array.isArray(modules)) {
                modules.forEach((mod) => {
                    const card = document.createElement("div");
                    card.className = "module-card";
                    card.innerHTML = `
            <div class="module-name">${mod}</div>
            <div class="module-desc">import("${mod}")</div>
          `;
                    modulesGrid.appendChild(card);
                });
            }
        } catch (e) {
            console.warn("Failed to render modules:", e);
        }
    }

    // Execute Tender Code
    function runCode() {
        if (typeof globalThis.tenderRun !== "function") {
            consoleOutput.innerHTML = `<span class="output-error">WASM Engine not ready yet.</span>`;
            return;
        }

        capturedLogs = "";
        consoleOutput.innerHTML = `<span class="output-log">Running Tender script...</span>\n`;

        setTimeout(() => {
            let resRaw = globalThis.tenderRun(codeEditor.value);
            let res = typeof resRaw === "string" ? JSON.parse(resRaw) : resRaw;

            let html = "";
            if (capturedLogs) {
                html += `<div class="output-log">${escapeHtml(capturedLogs)}</div>`;
            }

            if (res && res.error) {
                html += `<div class="output-error">❌ ${escapeHtml(res.error)}</div>`;
            } else if (res && res.variables && Object.keys(res.variables).length > 0) {
                html += `<div class="output-meta"><strong>Scope Variables:</strong>\n`;
                for (const [k, v] of Object.entries(res.variables)) {
                    html += `  ${escapeHtml(k)} = ${escapeHtml(v)}\n`;
                }
                html += `</div>`;
            }

            if (res && res.durationMs !== undefined) {
                html += `<div class="output-meta">Execution completed in ${res.durationMs}ms</div>`;
            }

            consoleOutput.innerHTML = html || "<span class=\"output-log\">(Program produced no output)</span>";
        }, 10);
    }

    // Inspect AST
    function parseCode() {
        if (typeof globalThis.tenderParse !== "function") {
            astOutput.textContent = "WASM Engine not ready.";
            return;
        }

        let resRaw = globalThis.tenderParse(codeEditor.value);
        let res = typeof resRaw === "string" ? JSON.parse(resRaw) : resRaw;

        if (res && res.error) {
            astOutput.textContent = `Parse Error: ${res.error}`;
        } else if (res && res.ast) {
            astOutput.textContent = res.ast;
        }
    }

    runBtn.addEventListener("click", runCode);
    parseBtn.addEventListener("click", () => {
        parseCode();
        // Switch to AST tab
        document.querySelector('[data-tab="astTab"]').click();
    });

    clearBtn.addEventListener("click", () => {
        consoleOutput.innerHTML = `<span class="output-log">Console cleared.</span>`;
        capturedLogs = "";
    });

    // Keyboard Shortcuts (Ctrl+Enter or Cmd+Enter to run)
    document.addEventListener("keydown", (e) => {
        if (e.shiftKey && e.key === "Enter") {
            e.preventDefault();
            runCode();
        }
    });

    function escapeHtml(str) {
        return str
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    // WebAssembly Initialization
    async function initWasm() {
        try {
            const go = new Go();
            let wasmPath = "tender.wasm";

            statusText.textContent = "LOADING";
            const response = await fetch(wasmPath);
            if (!response.ok) {
                throw new Error(`Failed to load ${wasmPath} (HTTP ${response.status}). Ensure tender.wasm is built.`);
            }

            let result;
            if (WebAssembly.instantiateStreaming) {
                try {
                    result = await WebAssembly.instantiateStreaming(response.clone(), go.importObject);
                } catch (e) {
                    const bytes = await response.arrayBuffer();
                    result = await WebAssembly.instantiate(bytes, go.importObject);
                }
            } else {
                const bytes = await response.arrayBuffer();
                result = await WebAssembly.instantiate(bytes, go.importObject);
            }

            go.run(result.instance);

            statusDot.classList.add("ready");
            statusText.textContent = "Ready";
            runBtn.disabled = false;
            consoleOutput.innerHTML = `<span class="output-log">Tender WASM Engine loaded successfully. Click 'Run Code' or press Ctrl+Enter to execute.</span>`;

            renderModules();
        } catch (err) {
            console.warn("WASM Init Notice:", err);
            statusText.textContent = "WASM Missing";
            consoleOutput.innerHTML = `
<div class="output-error">
<strong>Notice: tender.wasm build step required.</strong><br>
To generate the WASM binary for this playground, run the following build command in the root/playground directory:<br><br>
<code style="color: #38bdf8; background: rgba(0,0,0,0.4); padding: 4px 8px; border-radius: 4px;">
cd playground && set GOOS=js&& set GOARCH=wasm&& go build -o tender.wasm main.go
</code>
</div>`;
        }
    }

    initWasm();
});