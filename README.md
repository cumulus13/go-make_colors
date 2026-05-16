# go-make_colors

[![Go Reference](https://pkg.go.dev/badge/github.com/cumulus13/go-make_colors.svg)](https://pkg.go.dev/github.com/cumulus13/go-make_colors)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive, production-ready Go library for colored terminal text output.  
Port of the Python [make_colors](https://github.com/cumulus13/make_colors) library.

---

## Features

| Feature | Description |
|---|---|
| Named colors | `black red green yellow blue magenta cyan white` + `light*` variants |
| Color abbreviations | `r g bl y m c w b lb lr lg ly lm lc lw lk` |
| **Hex colors** | `#FF0000` `#F00` `FF0000` — in `Options` fields **and** markup |
| Combined format strings | `"bold-red-yellow"` `"lb_r"` `"italic,blue,white"` |
| Rich markup | `[bold red on yellow]text[/]` |
| Text attributes | `bold dim italic underline blink reverse strikethrough` |
| **Options is optional** | `MakeColors("text")` — omit `Options{}` entirely |
| Hex ↔ ANSI conversion | Truecolor (24-bit), 256-color, 16-color modes |
| Syntax highlighting | Via [Chroma](https://github.com/alecthomas/chroma) |
| Environment controls | `MAKE_COLORS` `MAKE_COLORS_FORCE` `MAKE_COLORS_DEBUG` |
| Cross-platform | Windows 10+, Linux, macOS |

---

## Installation

```bash
go get github.com/cumulus13/go-make_colors
```

---

## Import styles

Go requires a package prefix by default. Choose whichever style suits your code:

```go
// 1. Standard — explicit prefix (recommended for libraries)
import "github.com/cumulus13/go-make_colors/pkg/makecolors"

makecolors.MakeColors("hello", makecolors.Options{Foreground: "red"})
makecolors.Red("error")

// 2. Short alias — compact but still clear
import mc "github.com/cumulus13/go-make_colors/pkg/makecolors"

mc.MakeColors("hello", mc.Options{Foreground: "red"})
mc.Red("error")

// 3. Dot-import — Python-style bare names (from make_colors import *)
import . "github.com/cumulus13/go-make_colors/pkg/makecolors"

MakeColors("hello", Options{Foreground: "red"})
Red("error")
```

> The dot-import (`. "package"`) is the Go equivalent of Python's
> `from make_colors import make_colors, MakeColors, ...` — all exported names
> land directly in your file's scope with no prefix.

---

## Quick start

```go
package main

import (
    "fmt"
    . "github.com/cumulus13/go-make_colors/pkg/makecolors"
)

func main() {
    // Options is always optional
    fmt.Println(MakeColors("Plain text, no color"))
    fmt.Println(MakeColors("Plain text, empty options", Options{}))

    // Named colors
    fmt.Println(MakeColors("Red text",   Options{Foreground: "red"}))
    fmt.Println(MakeColors("On blue bg", Options{Foreground: "white", Background: "blue"}))

    // Hex colors — fg, bg, or both
    fmt.Println(MakeColors("Cyan hex",        Options{Foreground: "#00FFFF"}))
    fmt.Println(MakeColors("Short hex",       Options{Foreground: "#F00"}))
    fmt.Println(MakeColors("No-hash hex",     Options{Foreground: "FF6347"}))
    fmt.Println(MakeColors("Hex bg",          Options{Foreground: "white",   Background: "#8B0000"}))
    fmt.Println(MakeColors("Both hex",        Options{Foreground: "#00FFFF", Background: "#FF0000"}))
    fmt.Println(MakeColors("Hex + bold",      Options{Foreground: "#FFD700", Attrs: []string{"bold"}}))

    // Combined format string — fg, bg, and attrs in one string
    fmt.Println(MakeColors("Bold red on yellow", Options{Foreground: "bold-red-yellow"}))
    fmt.Println(MakeColors("Italic blue on white", Options{Foreground: "italic_blue_white"}))

    // Rich markup — Options not needed at all
    fmt.Println(MakeColors("[red]Error[/]"))
    fmt.Println(MakeColors("[bold green on black]Success[/]"))
    fmt.Println(MakeColors("[bold white on red][ERROR][/] [lightred]something failed[/]"))
    fmt.Println(MakeColors("[#00FFFF]Hex color in markup[/]"))
    fmt.Println(MakeColors("[bold #FF6347 on #000000]Tomato bold[/]"))

    // Sprint — positional args, no Options struct
    fmt.Println(Sprint("hello", "red", ""))
    fmt.Println(Sprint("hello", "#FF0000", "#000000"))
    fmt.Println(Sprint("hello", "red", "", "bold", "underline"))

    // Single-color helpers
    fmt.Println(Red("error"))
    fmt.Println(Green("ok"))
    fmt.Println(LightBlue("info"))
    fmt.Println(Bold("important"))
}
```

---

## `MakeColors` signature

```go
func MakeColors(text string, opts ...Options) string
```

`Options` is **variadic** — you may pass zero or one value:

| Call | Effect |
|---|---|
| `MakeColors("text")` | Plain passthrough |
| `MakeColors("text", Options{})` | Same — zero value |
| `MakeColors("text", Options{Foreground: "red"})` | Named color |
| `MakeColors("text", Options{Foreground: "#FF0000"})` | Hex fg |
| `MakeColors("text", Options{Background: "#00FFFF"})` | Hex bg, default fg |
| `MakeColors("[bold red]text[/]")` | Rich markup, no Options |
| `MakeColors("[bold red]text[/]", Options{Force: true})` | Markup + force output |

### `Options` fields

```go
type Options struct {
    Foreground string   // named color, abbreviation, hex, or combined string
    Background string   // named color, abbreviation, or hex
    Attrs      []string // ["bold", "italic", "underline", …]
    Force      bool     // force ANSI output even when stdout is not a TTY
}
```

---

## Color reference

### Named colors
`black` `red` `green` `yellow` `blue` `magenta` `cyan` `white`  
`lightblack` `lightred` `lightgreen` `lightyellow` `lightblue` `lightmagenta` `lightcyan` `lightwhite`

### Abbreviations

| Abbr | Full         | Abbr | Full          |
|------|-------------|------|--------------|
| `r`  | red         | `lb` | lightblue    |
| `g`  | green       | `lr` | lightred     |
| `bl` | blue        | `lg` | lightgreen   |
| `y`  | yellow      | `ly` | lightyellow  |
| `m`  | magenta     | `lm` | lightmagenta |
| `c`  | cyan        | `lc` | lightcyan    |
| `w`  | white       | `lw` | lightwhite   |
| `b`  | black       | `lk` | lightblack   |

### Hex colors

Accepted anywhere a color name is accepted — `Options.Foreground`, `Options.Background`,
`Sprint()` args, `NewColor()` args, and inside rich markup tags.

```go
"#FF0000"   // standard 6-digit with hash
"#F00"      // 3-digit shorthand
"FF0000"    // 6-digit without hash
```

All hex colors use **truecolor (24-bit)** ANSI sequences (`ESC[38;2;R;G;Bm`).

### Text attributes
`bold` `dim` `italic` `underline` `blink` `reverse` `strikethrough`

---

## Combined format strings

Embed fg, bg, and attributes into a single `Foreground` string using `-`, `_`, or `,`.
Token order does not matter.

```go
Options{Foreground: "bold-red-yellow"}         // bold, fg=red, bg=yellow
Options{Foreground: "italic_blue_white"}        // italic, fg=blue, bg=white
Options{Foreground: "underline,green,black"}    // underline, fg=green, bg=black
Options{Foreground: "bold-underline-white-red"} // bold+underline, fg=white, bg=red
Options{Foreground: "lb_r"}                     // fg=lightblue, bg=red
```

---

## Rich markup

```go
// Single color
MakeColors("[red]error[/]")

// With background
MakeColors("[white on red]ERROR[/]")

// With attributes
MakeColors("[bold italic green]success[/]")

// Hex colors
MakeColors("[#FF6347]tomato[/]")
MakeColors("[bold #FF69B4 on #000000]pink on black[/]")

// Multiple segments — Options still optional
MakeColors("[bold red][ERROR][/] [white]message here[/]")

// Escaped brackets in content
MakeColors(`[green]\[OK\] done[/]`)
```

---

## `Sprint` — positional args, no struct

```go
// Sprint(text, fg, bg, attrs...)
Sprint("hello", "red", "")
Sprint("hello", "red", "black", "bold")
Sprint("hello", "#FF0000", "")
Sprint("hello", "#FF0000", "#000000")
Sprint("hello", "#FFD700", "", "bold", "underline")
Sprint("hello", "lb", "r")           // abbreviations
Sprint("hello", "bold-red-yellow", "") // combined format
```

---

## `NewColor` — reusable formatter

```go
// Named colors
err  := NewColor("white", "red", "bold")
warn := NewColor("black", "yellow")
info := NewColor("white", "blue")

// Hex colors
tomato  := NewColor("#FF6347", "#000000")
skyblue := NewColor("#87CEEB", "", "bold")

fmt.Println(err.Format("FATAL"))
fmt.Println(tomato.Sprintf("Exit code: %d", 1))

// Use as a format verb (implements fmt.Stringer)
fmt.Printf("%scolored%s\n", tomato, Reset)
```

Pre-built singletons are also available:

```go
ColorRed.Format("error")
ColorLightGreen.Format("ok")
ColorBold.Format("important")
```

---

## `Console` helper

```go
c := NewConsole(nil)  // nil → os.Stdout
c.Force = true        // force ANSI even outside TTY

c.Error("something went wrong")    // lightred
c.Warn("disk almost full")         // yellow
c.Info("listening on :8080")       // cyan
c.Success("deployment complete")   // lightgreen
c.Debug("x=42")                    // lightblack/dim

c.Status("ERROR", "DB unavailable", "white", "red")
c.Rich("[bold cyan]Rich[/] markup in [italic yellow]Console.Rich()[/]")
c.Print("custom", Options{Foreground: "#FF6347", Attrs: []string{"bold"}})
```

---

## Hex ↔ ANSI conversion (`hexansi` package)

```go
import "github.com/cumulus13/go-make_colors/pkg/hexansi"

// Hex → ANSI (truecolor, 256, or 16)
res, err := hexansi.Convert("#FF0000", hexansi.ModeTrueColor)
fmt.Printf("%sRed text%s\n", res.FG, hexansi.Reset)

res, _ = hexansi.Convert("#FF0000", hexansi.Mode256)
res, _ = hexansi.Convert("#FF0000", hexansi.Mode16)

// Color name → hex
hex, err := hexansi.NameToHex("tomato")     // "#FF6347"
hex, err  = hexansi.NameToHex("skyblue")    // "#87CEEB"

// Color name → ANSI directly
res, err = hexansi.NameToANSI("tomato", hexansi.ModeTrueColor)

// Hex → nearest human name
name, err := hexansi.HexToColorName("#FF6347") // "Tomato"

// Auto-detect hex or name
res, err = hexansi.ToANSI("#FF0000", hexansi.ModeTrueColor)
res, err = hexansi.ToANSI("skyblue", hexansi.ModeTrueColor)
```

---

## Syntax highlighting (`syntax` package)

```go
import "github.com/cumulus13/go-make_colors/pkg/syntax"

code := `def hello():
    return "Hello, World!"`

// Quick print to stdout
syntax.Print(code, syntax.Options{
    Lexer:       "python",
    Theme:       "monokai",
    LineNumbers: syntax.LineNumberAbsolute,
    TrueColor:   true,
})

// Get as string
highlighted, err := syntax.Sprint(code, syntax.DefaultOptions())
fmt.Print(highlighted)

// Object style
s := syntax.New(code, syntax.Options{Lexer: "python", Theme: "solarized-dark"})
fmt.Print(s)  // implements fmt.Stringer

// Available themes / lexers
themes := syntax.AvailableThemes()
lexers := syntax.AvailableLexers()
```

### `syntax.Options` fields

```go
type Options struct {
    Lexer          string         // "python", "go", "auto" (default)
    Theme          string         // Chroma theme, default "monokai"
    LineNumbers    LineNumberMode // LineNumberNone / LineNumberAbsolute / LineNumberRelative
    StartLine      int            // first line number (default 1)
    TabSize        int            // spaces per tab (default 4)
    CodeWidth      int            // wrap at column (0 = disabled)
    WordWrap       bool
    HighlightLines []int          // 1-based line numbers to emphasise
    TrueColor      bool           // true = 24-bit, false = 256-color
}
```

---

## Environment variables

| Variable | Values | Effect |
|---|---|---|
| `MAKE_COLORS` | `0` | Disable all color output — functions return plain text |
| `MAKE_COLORS_FORCE` | `1` or `True` | Force ANSI output even when stdout is not a TTY |
| `MAKE_COLORS_DEBUG` | `1`, `true`, `True` | Print color-parsing debug info to stderr |

```bash
MAKE_COLORS=0          go run main.go   # plain output
MAKE_COLORS_FORCE=1    go run main.go   # colors even when piped
MAKE_COLORS_DEBUG=1    go run main.go   # show parsing trace
```

---

## CLI tool

```bash
go install github.com/cumulus13/go-make_colors/cmd/make_colors@latest
```

```
Usage:
  make_colors [flags]

Flags:
  -v                   Print version and exit
  -t                   Run built-in test suite
  -c <fg> <bg>         Show ANSI codes and preview for a color pair
  -m <markup>          Print Rich markup text
  -s <code>            Syntax-highlight a code snippet
  -l <lexer>           Language lexer for -s  (default: auto)
  -theme <name>        Chroma theme for -s    (default: monokai)
  -n                   Show line numbers for -s
```

**Examples:**
```bash
make_colors -m "[bold red]ERROR[/] [white]Database unavailable[/]"
make_colors -c "#FF6347" black
make_colors -s "print('hello')" -l python -theme monokai -n
```

---

## First-time setup

```bash
# After unzipping, run once to download dependencies and build the CLI:

# Windows
setup.bat

# Linux / macOS
./setup.sh
```

Or manually:

```bash
go mod tidy
go test ./...
go run examples/demo.go
go build -o make_colors ./cmd/make_colors
```

---

## Package structure

```
go-make_colors/
├── cmd/
│   └── make_colors/        CLI tool
│       └── main.go
├── pkg/
│   ├── hexansi/            Hex ↔ ANSI conversion
│   │   ├── hexansi.go
│   │   └── hexansi_test.go
│   ├── makecolors/         Core library
│   │   ├── makecolors.go   MakeColors · GetSort · Colorize · rich parser
│   │   ├── color.go        Color / Colors types + singletons
│   │   ├── console.go      Console helper
│   │   └── makecolors_test.go
│   └── syntax/             Syntax highlighting (Chroma)
│       └── syntax.go
├── examples/
│   └── demo.go             Full runnable demo
├── go.mod
├── setup.bat / setup.sh
└── README.md
```

---

## License

MIT © [Hadi Cahyadi](cumulus13@gmail.com)

## 👤 Author
        
[Hadi Cahyadi](mailto:cumulus13@gmail.com)
    

[![Buy Me a Coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/cumulus13)

[![Donate via Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/cumulus13)
 
[Support me on Patreon](https://www.patreon.com/cumulus13)
