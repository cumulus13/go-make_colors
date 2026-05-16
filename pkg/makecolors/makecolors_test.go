package makecolors_test

import (
	"os"
	"strings"
	"testing"

	"github.com/cumulus13/go-make_colors/pkg/makecolors"
)

// ─── helpers ──────────────────────────────────────────────────────────────────

func forceColor(t *testing.T) func() {
	t.Helper()
	old := os.Getenv("MAKE_COLORS_FORCE")
	os.Setenv("MAKE_COLORS_FORCE", "1")
	return func() { os.Setenv("MAKE_COLORS_FORCE", old) }
}

func stripped(s string) string { return makecolors.StripANSI(s) }

// ─── Options is fully optional ────────────────────────────────────────────────

func TestMakeColors_NoOptions(t *testing.T) {
	// No Options at all → plain passthrough (no TTY in test runner)
	result := makecolors.MakeColors("hello")
	if result != "hello" {
		// If a TTY happened to be detected the string will have ANSI — strip and check.
		if stripped(result) != "hello" {
			t.Errorf("stripped = %q, want hello", stripped(result))
		}
	}
}

func TestMakeColors_EmptyOptions(t *testing.T) {
	defer forceColor(t)()
	// Empty Options{} with Force env → just white text (default fg)
	result := makecolors.MakeColors("hello", makecolors.Options{})
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
}

// ─── Hex colors in Options fields ─────────────────────────────────────────────

func TestMakeColors_HexForeground(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("hello", makecolors.Options{Foreground: "#00FFFF", Force: true})
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	// Truecolor sequence: ESC[38;2;0;255;255m
	if !strings.Contains(result, "38;2;0;255;255") {
		t.Errorf("expected truecolor fg sequence, got %q", result)
	}
}

func TestMakeColors_HexBackground(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("hello", makecolors.Options{
		Foreground: "white",
		Background: "#FF0000",
		Force:      true,
	})
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	// Truecolor bg: ESC[48;2;255;0;0m
	if !strings.Contains(result, "48;2;255;0;0") {
		t.Errorf("expected truecolor bg sequence, got %q", result)
	}
}

func TestMakeColors_HexBothFgAndBg(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("hello", makecolors.Options{
		Foreground: "#00FFFF",
		Background: "#FF0000",
		Force:      true,
	})
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	if !strings.Contains(result, "38;2;0;255;255") {
		t.Errorf("missing truecolor fg in %q", result)
	}
	if !strings.Contains(result, "48;2;255;0;0") {
		t.Errorf("missing truecolor bg in %q", result)
	}
}

func TestMakeColors_ShortHex(t *testing.T) {
	defer forceColor(t)()

	// #F00 should expand to #FF0000
	result := makecolors.MakeColors("hello", makecolors.Options{Foreground: "#F00", Force: true})
	if !strings.Contains(result, "38;2;255;0;0") {
		t.Errorf("expected #F00 → rgb(255,0,0) truecolor, got %q", result)
	}
}

func TestMakeColors_HexNoHash(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("hello", makecolors.Options{Foreground: "00FFFF", Force: true})
	if !strings.Contains(result, "38;2;0;255;255") {
		t.Errorf("expected hex without # to work, got %q", result)
	}
}

// ─── Hex + attrs ──────────────────────────────────────────────────────────────

func TestMakeColors_HexWithAttrs(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("hello", makecolors.Options{
		Foreground: "#FF6347",
		Attrs:      []string{"bold"},
		Force:      true,
	})
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	if !strings.Contains(result, "38;2;255;99;71") {
		t.Errorf("missing truecolor fg in %q", result)
	}
}

// ─── Sprint with hex ──────────────────────────────────────────────────────────

func TestSprint_HexFG(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.Sprint("hello", "#00FF00", "")
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	if !strings.Contains(result, "38;2;0;255;0") {
		t.Errorf("expected green truecolor, got %q", result)
	}
}

func TestSprint_HexBG(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.Sprint("hello", "white", "#0000FF")
	if !strings.Contains(result, "48;2;0;0;255") {
		t.Errorf("expected blue truecolor bg, got %q", result)
	}
}

// ─── Named colors still work ──────────────────────────────────────────────────

func TestMakeColors_NamedColor(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("hello", makecolors.Options{Foreground: "red", Force: true})
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	if !strings.Contains(result, "31") { // red = ESC[31m
		t.Errorf("expected red code 31 in %q", result)
	}
}

// ─── Rich markup (Options still optional) ────────────────────────────────────

func TestMakeColors_RichMarkup_NoOptions(t *testing.T) {
	defer forceColor(t)()

	// Rich markup works with zero Options
	result := makecolors.MakeColors("[red]hello[/]")
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
}

func TestMakeColors_RichMarkup_HexInTag(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("[#00FFFF]hello[/]")
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	if !strings.Contains(result, "38;2;0;255;255") {
		t.Errorf("expected truecolor in markup, got %q", result)
	}
}

func TestMakeColors_RichMarkup_ForceViaOptions(t *testing.T) {
	result := makecolors.MakeColors("[red]hello[/]", makecolors.Options{Force: true})
	if !strings.Contains(result, "\x1b[") {
		t.Errorf("expected ANSI with Force=true, got %q", result)
	}
}

// ─── MAKE_COLORS=0 disables ───────────────────────────────────────────────────

func TestMakeColors_Disabled(t *testing.T) {
	os.Setenv("MAKE_COLORS", "0")
	defer os.Unsetenv("MAKE_COLORS")

	result := makecolors.MakeColors("hello", makecolors.Options{Foreground: "#FF0000"})
	if result != "hello" {
		t.Errorf("MAKE_COLORS=0: expected plain text, got %q", result)
	}
}

func TestMakeColors_Force_OverridesDisabled(t *testing.T) {
	os.Setenv("MAKE_COLORS", "0")
	defer os.Unsetenv("MAKE_COLORS")

	result := makecolors.MakeColors("hello", makecolors.Options{Foreground: "#FF0000", Force: true})
	if !strings.Contains(result, "38;2;255;0;0") {
		t.Errorf("Force=true should emit ANSI even with MAKE_COLORS=0, got %q", result)
	}
}

// ─── GetSort with hex ─────────────────────────────────────────────────────────

func TestGetSort_HexForeground(t *testing.T) {
	spec := makecolors.GetSort("", "#FF0000", "", nil)
	// After GetSort, Foreground should be a raw ANSI sequence (starts with ESC)
	if !strings.HasPrefix(spec.Foreground, "\x1b[") {
		t.Errorf("expected raw ANSI for hex fg, got %q", spec.Foreground)
	}
}

func TestGetSort_HexBackground(t *testing.T) {
	spec := makecolors.GetSort("", "white", "#00FF00", nil)
	if !strings.HasPrefix(spec.Background, "\x1b[") {
		t.Errorf("expected raw ANSI for hex bg, got %q", spec.Background)
	}
}

// ─── ColorMap ─────────────────────────────────────────────────────────────────

func TestColorMap(t *testing.T) {
	cases := map[string]string{
		"r": "red", "g": "green", "bl": "blue", "y": "yellow",
		"m": "magenta", "c": "cyan", "w": "white", "b": "black",
		"lb": "lightblue", "lr": "lightred", "lg": "lightgreen",
		"ly": "lightyellow", "lm": "lightmagenta", "lc": "lightcyan",
		"lw": "lightwhite", "lk": "lightblack",
		"red": "red", "green": "green", // full names pass through
	}
	for in, want := range cases {
		if got := makecolors.ColorMap(in); got != want {
			t.Errorf("ColorMap(%q) = %q, want %q", in, got, want)
		}
	}
}

// ─── StripANSI ────────────────────────────────────────────────────────────────

func TestStripANSI(t *testing.T) {
	input := "\x1b[1;31mhello\x1b[0m world \x1b[38;2;255;0;0mred\x1b[0m"
	want := "hello world red"
	if got := makecolors.StripANSI(input); got != want {
		t.Errorf("StripANSI = %q, want %q", got, want)
	}
}

// ─── Color object ─────────────────────────────────────────────────────────────

func TestNewColor_Hex(t *testing.T) {
	defer forceColor(t)()

	c := makecolors.NewColor("#FF0000", "#000000")
	result := c.Format("hello")
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	if !strings.Contains(result, "38;2;255;0;0") {
		t.Errorf("expected truecolor fg in Color.Format, got %q", result)
	}
}
