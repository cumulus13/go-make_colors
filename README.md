# go-make_colors

[![Go Reference](https://pkg.go.dev/badge/github.com/cumulus13/go-make_colors.svg)](https://pkg.go.dev/github.com/cumulus13/go-make_colors)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive, production-ready Go library for colored terminal text output with ANSI escape codes.  
Port of the Python [make_colors](https://github.com/cumulus13/make_colors) library.

---

## Features

| Feature | Description |
|---|---|
| 16 standard colors | `black red green yellow blue magenta cyan white` + light variants |
| Color abbreviations | `r g bl y m c w b lb lr lg ly lm lc lw lk` |
| Combined format strings | `"red-yellow"` `"bold-red-black"` `"lb_r"` `"italic,blue,white"` |
| Rich markup | `[bold red on yellow]text[/]` |
| Hex colors in markup | `[#FF6347 on #000000]text[/]` |
| Text attributes | `bold dim italic underline blink reverse strikethrough` |
| Attribute detection | Embedded in format strings: `"bold-red"` `"underline_green"` |
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

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/cumulus13/go-make_colors/pkg/makecolors"
)

func main() {
    // ── Functional helpers ──────────────────────────────────────────────────
    fmt.Println(makecolors.Red("Error!"))
    fmt.Println(makecolors.Green("Success"))
    fmt.Println(makecolors.Bold("Important"))

    // ── Sprint: fg + bg + attrs ─────────────────────────────────────────────
    fmt.Println(makecolors.Sprint("Warning", "yellow", "black", "bold"))

    // ── MakeColors: full options ────────────────────────────────────────────
    fmt.Println(makecolors.MakeColors("Info", makecolors.Options{
        Foreground: "lightblue",
        Background: "black",
        Attrs:      []string{"bold"},
    }))

    // ── Combined format strings ─────────────────────────────────────────────
    fmt.Println(makecolors.MakeColors("Hello", makecolors.Options{Foreground: "bold-red-yellow"}))
    fmt.Println(makecolors.MakeColors("World", makecolors.Options{Foreground: "italic_blue_black"}))

    // ── Color abbreviations ─────────────────────────────────────────────────
    fmt.Println(makecolors.MakeColors("short", makecolors.Options{Foreground: "lb", Background: "r"}))
    //                                                                          ^^lightblue  ^red

    // ── Rich markup ─────────────────────────────────────────────────────────
    fmt.Println(makecolors.MakeColors("[bold red]ERROR[/] [white]Database unavailable[/]",
        makecolors.Options{}))

    fmt.Println(makecolors.MakeColors("[white on blue]INFO[/]", makecolors.Options{}))

    // ── Hex colors ──────────────────────────────────────────────────────────
    fmt.Println(makecolors.MakeColors("[#FF6347 on #000000]Tomato on black[/]",
        makecolors.Options{}))

    // ── Pre-built Color objects ─────────────────────────────────────────────
    errColor := makecolors.NewColor("white", "red", "bold")
    fmt.Println(errColor.Format("FATAL"))
    fmt.Println(errColor.Sprintf("Exit code: %d", 1))

    // ── Console helper ──────────────────────────────────────────────────────
    c := makecolors.NewConsole(nil) // nil → os.Stdout
    c.Error("Something went wrong")
    c.Warn("Disk almost full")
    c.Info("Server started on :8080")
    c.Success("Deployment complete")
    c.Rich("[bold cyan]go-make_colors[/] is ready!")
    c.Status("CRITICAL", "System overload", "white", "red")
}
```

---

## Color Reference

### Standard colors
`black` `red` `green` `yellow` `blue` `magenta` `cyan` `white`

### Light variants
`lightblack` `lightred` `lightgreen` `lightyellow` `lightblue` `lightmagenta` `lightcyan` `lightwhite`

### Abbreviations

| Abbr | Full name    | Abbr | Full name     |
|------|-------------|------|--------------|
| `r`  | red         | `lb` | lightblue    |
| `g`  | green       | `lr` | lightred     |
| `bl` | blue        | `lg` | lightgreen   |
| `y`  | yellow      | `ly` | lightyellow  |
| `m`  | magenta     | `lm` | lightmagenta |
| `c`  | cyan        | `lc` | lightcyan    |
| `w`  | white       | `lw` | lightwhite   |
| `b`  | black       | `lk` | lightblack   |

### Text attributes
`bold` `dim` `italic` `underline` `blink` `reverse` `strikethrough`

---

## Combined Format Strings

Attributes and colors can be combined into a single `Foreground` string using `-`, `_`, or `,` as delimiters.  
Order of tokens does not matter.

```go
makecolors.MakeColors("text", makecolors.Options{Foreground: "bold-red-yellow"})
makecolors.MakeColors("text", makecolors.Options{Foreground: "italic_blue_white"})
makecolors.MakeColors("text", makecolors.Options{Foreground: "underline,green"})
makecolors.MakeColors("text", makecolors.Options{Foreground: "bold-underline-white-red"})
```

---

## Rich Markup

```go
// Single color
makecolors.MakeColors("[red]error[/]", opts)

// With background
makecolors.MakeColors("[white on red]ERROR[/]", opts)

// With style
makecolors.MakeColors("[bold italic green]success[/]", opts)

// Hex colors
makecolors.MakeColors("[#FF6347]tomato[/]", opts)
makecolors.MakeColors("[bold #FF69B4 on #000000]pink on black[/]", opts)

// Multiple segments
makecolors.MakeColors("[bold red][ERROR][/] [white]message here[/]", opts)

// Escaped brackets in content
makecolors.MakeColors(`[green]\[OK\] done[/]`, opts)
```

---

## Hex ↔ ANSI Conversion

```go
import "github.com/cumulus13/go-make_colors/pkg/hexansi"

// Hex → ANSI
result, err := hexansi.Convert("#FF0000", hexansi.ModeTrueColor)
fmt.Printf("%sRed text%s\n", result.FG, hexansi.Reset)

// Color name → hex
hex, err := hexansi.NameToHex("tomato")  // "#FF6347"

// Color name → ANSI
result, err = hexansi.NameToANSI("skyblue", hexansi.Mode256)

// Hex → nearest color name
name, err := hexansi.HexToColorName("#FF6347")  // "Tomato"
```

---

## Syntax Highlighting

```go
import "github.com/cumulus13/go-make_colors/pkg/syntax"

code := `def hello():
    return "Hello, World!"`

// Quick print
syntax.Print(code, syntax.Options{
    Lexer:       "python",
    Theme:       "monokai",
    LineNumbers: syntax.LineNumberAbsolute,
})

// Or get string
highlighted, err := syntax.Sprint(code, syntax.DefaultOptions())
fmt.Print(highlighted)
```

---

## Environment Variables

| Variable | Values | Effect |
|---|---|---|
| `MAKE_COLORS` | `0` | Disable all color output |
| `MAKE_COLORS_FORCE` | `1` or `True` | Force color even in non-TTY (file redirect, CI) |
| `MAKE_COLORS_DEBUG` | `1`, `true`, `True` | Print parsing debug info to stderr |

---

## CLI Tool

```bash
go install github.com/cumulus13/go-make_colors/cmd/make_colors@latest
```

```
Usage:
  make_colors [flags]

Flags:
  -v                   Print version and exit
  -t                   Run built-in test suite
  -c <fg> <bg>         Show ANSI codes for a color pair
  -m <markup>          Print Rich markup text
  -s <code>            Syntax-highlight a code snippet
  -l <lexer>           Language lexer for -s (default: auto)
  -theme <name>        Chroma theme for -s (default: monokai)
  -n                   Show line numbers for -s
```

**Examples:**
```bash
make_colors -m "[bold red]ERROR[/] [white]Database unavailable[/]"
make_colors -c lightblue red
make_colors -s "print('hello')" -l python -theme monokai -n
```

---

## Package Structure

```
go-make_colors/
├── cmd/
│   └── make_colors/       # CLI tool
│       └── main.go
├── pkg/
│   ├── hexansi/           # Hex↔ANSI conversion
│   │   ├── hexansi.go
│   │   └── hexansi_test.go
│   ├── makecolors/        # Core coloring library
│   │   ├── makecolors.go  # MakeColors, GetSort, Colorize, rich markup
│   │   ├── color.go       # Color/Colors types
│   │   ├── console.go     # Console helper
│   │   └── makecolors_test.go
│   └── syntax/            # Syntax highlighting
│       └── syntax.go
├── go.mod
├── go.sum
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
