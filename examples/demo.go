//go:build ignore

// examples/demo.go – run with: go run examples/demo.go
package main

import (
	"fmt"
	"os"

	"github.com/cumulus13/go-make_colors/pkg/hexansi"
	"github.com/cumulus13/go-make_colors/pkg/makecolors"
	"github.com/cumulus13/go-make_colors/pkg/syntax"
)

func main() {
	// Force colors so this demo renders even when stdout is redirected.
	os.Setenv("MAKE_COLORS_FORCE", "1")

	sep := func(title string) {
		fmt.Println()
		fmt.Println(makecolors.Sprint("════════════════════════════════════════", "lightcyan", "", "bold"))
		fmt.Println(makecolors.Sprint("  "+title, "lightcyan", "", "bold"))
		fmt.Println(makecolors.Sprint("════════════════════════════════════════", "lightcyan", "", "bold"))
	}

	// ── Functional helpers ────────────────────────────────────────────────────
	sep("Functional Helpers")
	fmt.Println(makecolors.Red("Red()"))
	fmt.Println(makecolors.Green("Green()"))
	fmt.Println(makecolors.Blue("Blue()"))
	fmt.Println(makecolors.Yellow("Yellow()"))
	fmt.Println(makecolors.Magenta("Magenta()"))
	fmt.Println(makecolors.Cyan("Cyan()"))
	fmt.Println(makecolors.LightRed("LightRed()"))
	fmt.Println(makecolors.LightGreen("LightGreen()"))
	fmt.Println(makecolors.LightBlue("LightBlue()"))
	fmt.Println(makecolors.Bold("Bold()"))
	fmt.Println(makecolors.Italic("Italic()"))
	fmt.Println(makecolors.Underline("Underline()"))
	fmt.Println(makecolors.Strikethrough("Strikethrough()"))

	// ── Sprint ────────────────────────────────────────────────────────────────
	sep("Sprint(text, fg, bg, attrs...)")
	fmt.Println(makecolors.Sprint("White on red, bold", "white", "red", "bold"))
	fmt.Println(makecolors.Sprint("Lightblue on black", "lightblue", "black"))
	fmt.Println(makecolors.Sprint("Bold + underline", "green", "", "bold", "underline"))

	// ── Combined format strings ───────────────────────────────────────────────
	sep("Combined Format Strings (Foreground field)")
	combos := []string{
		"red-yellow",
		"bold-red-black",
		"italic_blue_white",
		"underline,green,black",
		"bold-underline-white-red",
		"lb_r",   // lightblue on red
		"dim_w",  // dim white
	}
	for _, combo := range combos {
		s := makecolors.MakeColors("Sample text", makecolors.Options{Foreground: combo, Force: true})
		fmt.Printf("  %-30s → %s\n", combo, s)
	}

	// ── Rich markup ───────────────────────────────────────────────────────────
	sep("Rich Markup")
	markupExamples := []string{
		"[red]Simple red[/]",
		"[bold green]Bold green[/]",
		"[white on blue]White on blue[/]",
		"[bold italic yellow on black]Bold italic yellow on black[/]",
		"[underline lightcyan]Underlined light cyan[/]",
		"[bold white on red][FATAL][/] [lightred]System crash[/]",
		`[green]\[OK\] Escaped brackets in content[/]`,
	}
	for _, ex := range markupExamples {
		fmt.Println(makecolors.MakeColors(ex, makecolors.Options{Force: true}))
	}

	// ── Hex colors ────────────────────────────────────────────────────────────
	sep("Hex Colors")
	hexExamples := []string{
		"[#FF0000]Pure red hex[/]",
		"[#00FF00]Pure green hex[/]",
		"[#0000FF]Pure blue hex[/]",
		"[#FF6347]Tomato[/]",
		"[bold #FF69B4 on #000000]Hot pink on black[/]",
		"[#FFFFFF on #800080]White on purple[/]",
	}
	for _, ex := range hexExamples {
		fmt.Println(makecolors.MakeColors(ex, makecolors.Options{Force: true}))
	}

	// ── Color type ────────────────────────────────────────────────────────────
	sep("Color / Colors Objects")
	errColor := makecolors.NewColor("white", "red", "bold")
	warnColor := makecolors.NewColor("black", "yellow")
	infoColor := makecolors.NewColor("white", "blue")

	fmt.Println(errColor.Format("  This is an error"))
	fmt.Println(warnColor.Format("  This is a warning"))
	fmt.Println(infoColor.Format("  This is informational"))
	fmt.Println(errColor.Sprintf("  Exit code: %d", 127))

	// Use as format verb
	fmt.Printf("  %s%s%s\n", errColor, "Raw ANSI usage", makecolors.Reset)

	// Pre-built singletons
	fmt.Println(makecolors.ColorRed.Format("  ColorRed singleton"))
	fmt.Println(makecolors.ColorLightGreen.Format("  ColorLightGreen singleton"))
	fmt.Println(makecolors.ColorBold.Format("  ColorBold singleton"))

	// ── Console helper ────────────────────────────────────────────────────────
	sep("Console Helper")
	c := makecolors.NewConsole(os.Stdout)
	c.Force = true
	c.Error("  c.Error()")
	c.Warn("  c.Warn()")
	c.Info("  c.Info()")
	c.Success("  c.Success()")
	c.Debug("  c.Debug()")
	c.Status("ERROR", "Database unavailable", "white", "red")
	c.Status("INFO", "Server started on :8080", "white", "blue")
	c.Rich("[bold cyan]  c.Rich()[/] with [italic yellow]inline markup[/]")

	// ── GetSort / ColorSpec ───────────────────────────────────────────────────
	sep("GetSort / ColorSpec")
	specs := []struct{ data, fg, bg string }{
		{"red-yellow", "", ""},
		{"bold-red-black", "", ""},
		{"", "lb", "r"},
		{"italic-blue-white", "", ""},
		{"", "white", "on_red"},
	}
	for _, s := range specs {
		spec := makecolors.GetSort(s.data, s.fg, s.bg, nil)
		desc := s.data
		if desc == "" {
			desc = fmt.Sprintf("fg=%q bg=%q", s.fg, s.bg)
		}
		colored := makecolors.Colorize("████", spec)
		fmt.Printf("  %-28s → fg=%-14s bg=%-14s attrs=%v %s\n",
			desc, spec.Foreground, spec.Background, spec.Attrs, colored)
	}

	// ── StripANSI ─────────────────────────────────────────────────────────────
	sep("StripANSI")
	colored := makecolors.MakeColors("[bold red]Colorful text[/]", makecolors.Options{Force: true})
	plain := makecolors.StripANSI(colored)
	fmt.Printf("  Colored  : %s\n", colored)
	fmt.Printf("  Plain    : %q\n", plain)

	// ── hexansi package ───────────────────────────────────────────────────────
	sep("hexansi Package")

	for _, hex := range []string{"#FF0000", "#FFA500", "#FFFF00", "#008000", "#0000FF", "#800080"} {
		res, _ := hexansi.Convert(hex, hexansi.ModeTrueColor)
		name, _ := hexansi.HexToColorName(hex)
		fmt.Printf("  %s██%s  %-10s  %s\n", res.FG, hexansi.Reset, hex, name)
	}

	fmt.Println()
	for _, name := range []string{"red", "orange", "green", "blue", "purple", "pink", "gold"} {
		hex, err := hexansi.NameToHex(name)
		if err != nil {
			continue
		}
		res, _ := hexansi.Convert(hex, hexansi.ModeTrueColor)
		fmt.Printf("  %s██%s  %-12s → %s\n", res.FG, hexansi.Reset, name, hex)
	}

	// ── Syntax highlighting ───────────────────────────────────────────────────
	sep("Syntax Highlighting")
	code := `package main

import "fmt"

func main() {
    msg := "Hello, World!"
    fmt.Println(msg)
}`
	opts := syntax.Options{
		Lexer:       "go",
		Theme:       "monokai",
		LineNumbers: syntax.LineNumberAbsolute,
		StartLine:   1,
		TrueColor:   true,
	}
	if err := syntax.Print(code, opts); err != nil {
		fmt.Fprintf(os.Stderr, "syntax error: %v\n", err)
	}

	// ── Log-level simulation ──────────────────────────────────────────────────
	sep("Log Level Simulation")
	logLines := []string{
		"[bold white on black][DEBUG][/] [cyan]Cache miss for key user:1234[/]",
		"[bold blue on black][INFO] [/] [white]HTTP GET /api/users 200 OK[/]",
		"[bold yellow on black][WARN] [/] [lightyellow]Response time > 500ms[/]",
		"[bold white on red][ERROR][/] [lightred]Failed to connect: timeout[/]",
		"[bold white on red][FATAL][/] [white on red]Out of memory – shutting down[/]",
	}
	for _, line := range logLines {
		fmt.Println(makecolors.MakeColors(line, makecolors.Options{Force: true}))
	}

	// ── Progress bar ──────────────────────────────────────────────────────────
	sep("Progress Bar Simulation")
	for i := 0; i <= 20; i += 4 {
		pct := i * 5
		var color string
		switch {
		case pct < 30:
			color = "red"
		case pct < 70:
			color = "yellow"
		default:
			color = "green"
		}
		filled := string([]rune("█████████████████████")[:i])
		empty := string([]rune("░░░░░░░░░░░░░░░░░░░░░")[: 20-i])
		bar := makecolors.Sprint(fmt.Sprintf("[%s%s] %3d%%", filled, empty, pct), color, "", "bold")
		fmt.Println("  " + bar)
	}

	sep("Done ✓")
}
