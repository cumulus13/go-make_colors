//go:build ignore

// Run with: go run examples/demo.go
package main

import (
	"fmt"
	"os"

	"github.com/cumulus13/go-make_colors/pkg/hexansi"
	. "github.com/cumulus13/go-make_colors/pkg/makecolors" // dot-import = Python-style bare names
	"github.com/cumulus13/go-make_colors/pkg/syntax"
)

func main() {
	os.Setenv("MAKE_COLORS_FORCE", "1")

	sep := func(title string) {
		fmt.Println()
		line := "════════════════════════════════════════════════"
		fmt.Println(Sprint(line, "lightcyan", "", "bold"))
		fmt.Println(Sprint("  "+title, "lightcyan", "", "bold"))
		fmt.Println(Sprint(line, "lightcyan", "", "bold"))
	}

	// ── Options is always optional ────────────────────────────────────────────
	sep("Options is always optional")

	fmt.Println(MakeColors("No options at all — plain passthrough"))
	fmt.Println(MakeColors("Empty Options{} — default white", Options{}))
	fmt.Println(MakeColors("Only Force=true", Options{Force: true}))

	// ── Named colors ──────────────────────────────────────────────────────────
	sep("Named colors in Options")

	fmt.Println(MakeColors("Named fg only",    Options{Foreground: "cyan"}))
	fmt.Println(MakeColors("Named fg+bg",       Options{Foreground: "white", Background: "blue"}))
	fmt.Println(MakeColors("Named fg+attrs",    Options{Foreground: "red", Attrs: []string{"bold"}}))

	// ── Hex colors in Options ─────────────────────────────────────────────────
	sep("Hex colors in Options{Foreground / Background}")

	fmt.Println(MakeColors("Hello in Cyan",          Options{Foreground: "#00FFFF"}))
	fmt.Println(MakeColors("Red hex fg",              Options{Foreground: "#FF0000"}))
	fmt.Println(MakeColors("Short hex #F00",          Options{Foreground: "#F00"}))
	fmt.Println(MakeColors("No-hash hex FF6347",      Options{Foreground: "FF6347"}))
	fmt.Println(MakeColors("Hex fg + named bg",       Options{Foreground: "#FF69B4", Background: "black"}))
	fmt.Println(MakeColors("Named fg + hex bg",       Options{Foreground: "white",   Background: "#8B0000"}))
	fmt.Println(MakeColors("Both hex",                Options{Foreground: "#00FFFF", Background: "#FF0000"}))
	fmt.Println(MakeColors("Hex + bold attr",         Options{Foreground: "#FFD700", Attrs: []string{"bold"}}))
	fmt.Println(MakeColors("Hex fg + bg + italic",    Options{
		Foreground: "#00FF7F",
		Background: "#000080",
		Attrs:      []string{"italic"},
	}))

	// ── Sprint with hex ───────────────────────────────────────────────────────
	sep("Sprint(text, fg, bg, attrs...) with hex")

	fmt.Println(Sprint("Sprint hex fg",      "#00FFFF", ""))
	fmt.Println(Sprint("Sprint hex bg",      "white",   "#8B008B"))
	fmt.Println(Sprint("Sprint both hex",    "#FFFF00", "#000080"))
	fmt.Println(Sprint("Sprint hex + bold",  "#FF6347", "",       "bold"))

	// ── Combined format strings ───────────────────────────────────────────────
	sep("Combined format strings")

	fmt.Println(MakeColors("bold-red",           Options{Foreground: "bold-red"}))
	fmt.Println(MakeColors("italic-blue-yellow",  Options{Foreground: "italic_blue_yellow"}))
	fmt.Println(MakeColors("lb_r abbreviations",  Options{Foreground: "lb_r"}))

	// ── Rich markup — no Options needed ──────────────────────────────────────
	sep("Rich markup (Options optional)")

	fmt.Println(MakeColors("[red]No Options needed[/]"))
	fmt.Println(MakeColors("[bold green on black]Bold green[/]"))
	fmt.Println(MakeColors("[bold white on red][ERROR][/] [lightred]something failed[/]"))
	fmt.Println(MakeColors("[#00FFFF]Hex in markup[/]"))
	fmt.Println(MakeColors("[bold #FF6347 on #000000]Hex markup bold[/]"))

	// ── Dot-import: bare function names (Python-style) ────────────────────────
	sep("Dot-import bare names (. \"…/makecolors\")")

	fmt.Println(Red("Red()"))
	fmt.Println(Green("Green()"))
	fmt.Println(Cyan("Cyan()"))
	fmt.Println(LightBlue("LightBlue()"))
	fmt.Println(Bold("Bold()"))
	fmt.Println(Underline("Underline()"))

	// ── NewColor with hex ─────────────────────────────────────────────────────
	sep("NewColor with hex")

	tomato  := NewColor("#FF6347", "#000000")
	skyblue := NewColor("#87CEEB", "", "bold")
	fmt.Println(tomato.Format("  Tomato on black"))
	fmt.Println(skyblue.Format("  Bold sky blue"))
	fmt.Println(tomato.Sprintf("  Exit code: %d", 1))

	// ── Console helper ────────────────────────────────────────────────────────
	sep("Console helper")

	c := NewConsole(os.Stdout)
	c.Force = true
	c.Error("  c.Error()")
	c.Warn("  c.Warn()")
	c.Info("  c.Info()")
	c.Success("  c.Success()")
	c.Debug("  c.Debug()")
	c.Status("ERROR", "Database unavailable", "white", "red")
	c.Rich("[bold #00FFFF]Hex color[/] in [italic yellow]Rich()[/]")

	// ── hexansi package ───────────────────────────────────────────────────────
	sep("hexansi package")

	for _, h := range []string{"#FF0000", "#FFA500", "#00FFFF", "#008000", "#0000FF", "#FF69B4"} {
		res, _ := hexansi.Convert(h, hexansi.ModeTrueColor)
		name, _ := hexansi.HexToColorName(h)
		fmt.Printf("  %s██%s  %-10s  %s\n", res.FG, hexansi.Reset, h, name)
	}

	// ── Syntax highlight ──────────────────────────────────────────────────────
	sep("Syntax highlighting")

	code := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`
	syntax.Print(code, syntax.Options{
		Lexer:       "go",
		Theme:       "monokai",
		LineNumbers: syntax.LineNumberAbsolute,
		TrueColor:   true,
	})

	sep("Done ✓")
}
