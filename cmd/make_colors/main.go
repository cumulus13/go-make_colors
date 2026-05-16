// Command make_colors is a CLI tool for colorizing terminal text using
// the go-make_colors library.
//
// Usage:
//
//	make_colors [flags]
//	make_colors -c <foreground> <background>   # get ANSI code for a color pair
//	make_colors -t                              # run built-in test suite
//	make_colors -m "[bold red]Hello[/]"        # print rich markup
//	make_colors -s "code here" -l python       # syntax highlight a snippet
//
// Author: Hadi Cahyadi <cumulus13@gmail.com>
// License: MIT
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cumulus13/go-make_colors/pkg/makecolors"
	"github.com/cumulus13/go-make_colors/pkg/syntax"
)

const version = "1.0.0"

func main() {
	var (
		flagConvert  = flag.Bool("c", false, "Get ANSI code for a foreground/background pair (requires two positional args)")
		flagTest     = flag.Bool("t", false, "Run built-in test suite")
		flagMarkup   = flag.String("m", "", "Print Rich markup text")
		flagSnippet  = flag.String("s", "", "Syntax-highlight a code snippet")
		flagLexer    = flag.String("l", "auto", "Language for syntax highlighting (default: auto)")
		flagTheme    = flag.String("theme", "monokai", "Chroma theme for syntax highlighting")
		flagLines    = flag.Bool("n", false, "Show line numbers in syntax output")
		flagVersion  = flag.Bool("v", false, "Print version and exit")
	)

	flag.Usage = usage
	flag.Parse()

	if *flagVersion {
		fmt.Printf("go-make_colors v%s\n", version)
		fmt.Println("Author: Hadi Cahyadi <cumulus13@gmail.com>")
		fmt.Println("Home:   https://github.com/cumulus13/go-make_colors")
		return
	}

	if *flagMarkup != "" {
		fmt.Println(makecolors.MakeColors(*flagMarkup, makecolors.Options{Force: true}))
		return
	}

	if *flagSnippet != "" {
		opts := syntax.DefaultOptions()
		opts.Lexer = *flagLexer
		opts.Theme = *flagTheme
		if *flagLines {
			opts.LineNumbers = syntax.LineNumberAbsolute
		}
		if err := syntax.Print(*flagSnippet, opts); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *flagConvert {
		args := flag.Args()
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "error: -c requires two arguments: <foreground> <background>")
			os.Exit(1)
		}
		fg, bg := args[0], args[1]
		spec := makecolors.GetSort("", fg, bg, nil)
		sample := makecolors.MakeColors("████ Sample Text ████", makecolors.Options{
			Foreground: spec.Foreground,
			Background: spec.Background,
			Attrs:      spec.Attrs,
			Force:      true,
		})
		fmt.Printf("foreground : %s\n", spec.Foreground)
		fmt.Printf("background : %s\n", spec.Background)
		fmt.Printf("preview    : %s\n", sample)
		return
	}

	if *flagTest {
		runTests()
		return
	}

	// No flags → run tests and show help
	runTests()
	fmt.Println()
	usage()
}

// ─── Usage ────────────────────────────────────────────────────────────────────

func usage() {
	c := makecolors.NewConsole(os.Stdout)
	c.Rich("[bold lightcyan]go-make_colors[/] [dim]v" + version + "[/]")
	c.Rich("[lightblue]Author:[/] Hadi Cahyadi <cumulus13@gmail.com>")
	c.Rich("[lightblue]Home  :[/] https://github.com/cumulus13/go-make_colors")
	fmt.Println()

	c.Rich("[bold lightmagenta]Usage:[/]")
	fmt.Println("  make_colors [flags]")
	fmt.Println()

	c.Rich("[bold lightmagenta]Flags:[/]")
	rows := [][2]string{
		{"-v", "Print version and exit"},
		{"-t", "Run built-in test suite"},
		{"-c <fg> <bg>", "Show ANSI codes for a color pair"},
		{"-m <markup>", `Print Rich markup (e.g. "[bold red]hi[/]")`},
		{"-s <code>", "Syntax-highlight a code snippet"},
		{"-l <lexer>", "Language lexer for -s (default: auto)"},
		{"-theme <name>", "Chroma theme for -s (default: monokai)"},
		{"-n", "Show line numbers for -s"},
	}
	for _, row := range rows {
		flag_ := makecolors.Sprint(row[0], "lightyellow", "", "bold")
		fmt.Printf("  %-32s %s\n", flag_, row[1])
	}
	fmt.Println()

	c.Rich("[bold lightmagenta]Environment variables:[/]")
	envRows := [][2]string{
		{"MAKE_COLORS=0", "Disable all color output"},
		{"MAKE_COLORS_FORCE=1", "Force color even in non-TTY output"},
		{"MAKE_COLORS_DEBUG=1", "Print parsing debug information"},
	}
	for _, row := range envRows {
		k := makecolors.Sprint(row[0], "lightgreen", "")
		fmt.Printf("  %-35s %s\n", k, row[1])
	}
}

// ─── Test suite ───────────────────────────────────────────────────────────────

func runTests() {
	c := makecolors.NewConsole(os.Stdout)
	force := makecolors.Options{Force: true}

	sep := func(title string) {
		line := strings.Repeat("═", 60)
		c.Rich(fmt.Sprintf("[bold lightcyan]%s[/]", line))
		c.Rich(fmt.Sprintf("[bold lightcyan]  %s[/]", title))
		c.Rich(fmt.Sprintf("[bold lightcyan]%s[/]", line))
	}

	// ── Basic colors ──────────────────────────────────────────────────────────
	sep("Standard Foreground Colors")
	for _, name := range []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"} {
		s := makecolors.MakeColors(fmt.Sprintf("  ● %-14s Sample text", name),
			makecolors.Options{Foreground: name, Background: "black", Force: true})
		fmt.Println(s)
	}

	sep("Light Foreground Colors")
	for _, name := range []string{"lightblack", "lightred", "lightgreen", "lightyellow",
		"lightblue", "lightmagenta", "lightcyan", "lightwhite"} {
		s := makecolors.MakeColors(fmt.Sprintf("  ● %-15s Sample text", name),
			makecolors.Options{Foreground: name, Background: "black", Force: true})
		fmt.Println(s)
	}

	// ── Background colors ─────────────────────────────────────────────────────
	sep("Background Colors")
	for _, bg := range []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"} {
		fg := "white"
		if bg == "white" || bg == "cyan" || bg == "yellow" {
			fg = "black"
		}
		s := makecolors.MakeColors(fmt.Sprintf("  on_%-10s  Sample", bg),
			makecolors.Options{Foreground: fg, Background: bg, Force: true})
		fmt.Println(s)
	}

	// ── Abbreviations ─────────────────────────────────────────────────────────
	sep("Color Abbreviations")
	abbrevPairs := [][2]string{
		{"r", "red"}, {"g", "green"}, {"bl", "blue"}, {"y", "yellow"},
		{"m", "magenta"}, {"c", "cyan"}, {"w", "white"}, {"b", "black"},
		{"lb", "lightblue"}, {"lr", "lightred"}, {"lg", "lightgreen"},
		{"ly", "lightyellow"}, {"lm", "lightmagenta"}, {"lc", "lightcyan"},
		{"lw", "lightwhite"},
	}
	for _, p := range abbrevPairs {
		abbr := makecolors.Sprint("'"+p[0]+"'", p[1], "black", "bold")
		full := makecolors.Sprint("'"+p[1]+"'", p[1], "black")
		fmt.Printf("  %-22s → %s\n", abbr, full)
	}

	// ── Combined format ───────────────────────────────────────────────────────
	sep("Combined Format Strings")
	combos := []string{"red-yellow", "blue_white", "g-b", "lb_r", "w,m", "bold-red", "italic-blue-yellow"}
	for _, combo := range combos {
		s := makecolors.MakeColors(fmt.Sprintf("  %-22s", combo),
			makecolors.Options{Foreground: combo, Force: true})
		fmt.Printf("%s → %s\n", s,
			makecolors.MakeColors("Hello World", makecolors.Options{Foreground: combo, Force: true}))
	}

	// ── Attributes ────────────────────────────────────────────────────────────
	sep("Text Attributes")
	for _, attr := range []string{"bold", "dim", "italic", "underline", "blink", "reverse", "strikethrough"} {
		s := makecolors.MakeColors(fmt.Sprintf("  %-16s Sample text", attr),
			makecolors.Options{Foreground: "white", Attrs: []string{attr}, Force: true})
		fmt.Println(s)
	}

	// ── Rich markup ───────────────────────────────────────────────────────────
	sep("Rich Markup Format")
	richExamples := []string{
		"[red]This is red text[/]",
		"[bold green]Bold green text[/]",
		"[white on blue]White on blue[/]",
		"[italic yellow on black]Italic yellow on black[/]",
		"[bold white on red][ERROR][/] [lightred]Something failed[/]",
		"[bold blue][INFO][/] [white]Server started on :8080[/]",
		"[bold yellow][WARNING][/] [lightyellow]High memory usage[/]",
		`[bold cyan]Hex color:[/] [#FF6347]Tomato hex color[/]`,
	}
	for _, ex := range richExamples {
		fmt.Println(makecolors.MakeColors(ex, force))
	}

	// ── Hex colors ────────────────────────────────────────────────────────────
	sep("Hex Colors via Markup")
	hexExamples := []string{
		"[#FF0000]Red via hex[/]",
		"[#00FF00]Green via hex[/]",
		"[#0000FF]Blue via hex[/]",
		"[bold #FF69B4 on #000000]Hot pink on black[/]",
	}
	for _, ex := range hexExamples {
		fmt.Println(makecolors.MakeColors(ex, force))
	}

	// ── Color singletons ──────────────────────────────────────────────────────
	sep("Pre-built Color Formatters")
	fmt.Println(makecolors.ColorRed.Format("  ColorRed.Format(\"...\")"))
	fmt.Println(makecolors.ColorGreen.Format("  ColorGreen.Format(\"...\")"))
	fmt.Println(makecolors.ColorLightBlue.Format("  ColorLightBlue.Format(\"...\")"))
	fmt.Println(makecolors.ColorBold.Format("  ColorBold.Format(\"...\")"))
	fmt.Println(makecolors.ColorUnderline.Format("  ColorUnderline.Format(\"...\")"))

	// ── Console helper ────────────────────────────────────────────────────────
	sep("Console Helper Methods")
	c.Error("  Console.Error: something went wrong")
	c.Warn("  Console.Warn: disk is almost full")
	c.Info("  Console.Info: listening on :8080")
	c.Success("  Console.Success: deployment complete")
	c.Debug("  Console.Debug: x=42, y=13")
	c.Status("CRITICAL", "System shutdown required", "white", "red")

	// ── Log-level simulation ──────────────────────────────────────────────────
	sep("Log Level Simulation")
	logLines := []string{
		"[bold white on black][DEBUG][/] [cyan]Database connection established[/]",
		"[bold blue on black][INFO] [/] [white]User authenticated successfully[/]",
		"[bold yellow on black][WARN] [/] [lightyellow]Memory usage above 80%[/]",
		"[bold white on red][ERROR][/] [lightred]Failed to connect to database[/]",
		"[bold white on red][FATAL][/] [white on red]System shutdown required[/]",
	}
	for _, line := range logLines {
		fmt.Println(makecolors.MakeColors(line, force))
	}

	// ── StripANSI ─────────────────────────────────────────────────────────────
	sep("StripANSI")
	colored := makecolors.MakeColors("[bold red]Colored text[/]", force)
	stripped := makecolors.StripANSI(colored)
	fmt.Printf("  Colored : %s\n", colored)
	fmt.Printf("  Stripped: %q\n", stripped)

	// ── Color palette ─────────────────────────────────────────────────────────
	sep("Full Color Palette")
	allColors := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"lightblack", "lightred", "lightgreen", "lightyellow",
		"lightblue", "lightmagenta", "lightcyan", "lightwhite",
	}
	fmt.Print("  ")
	for i, name := range allColors {
		block := makecolors.MakeColors("██", makecolors.Options{Foreground: name, Force: true})
		fmt.Print(block)
		if i == 7 {
			fmt.Println()
			fmt.Print("  ")
		}
	}
	fmt.Println()

	sep("All Tests Complete ✓")
	c.Rich("[bold lightgreen]go-make_colors is working correctly![/]")
}
