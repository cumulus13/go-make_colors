package makecolors_test

import (
	"os"
	"strings"
	"testing"

	"github.com/cumulus13/go-make_colors/pkg/makecolors"
)

// ─── Helpers ──────────────────────────────────────────────────────────────────

func forceColor(t *testing.T) func() {
	t.Helper()
	old := os.Getenv("MAKE_COLORS_FORCE")
	os.Setenv("MAKE_COLORS_FORCE", "1")
	return func() { os.Setenv("MAKE_COLORS_FORCE", old) }
}

func stripped(s string) string { return makecolors.StripANSI(s) }

// ─── ColorMap / abbreviations ─────────────────────────────────────────────────

func TestColorMap(t *testing.T) {
	cases := map[string]string{
		"r":   "red",
		"g":   "green",
		"bl":  "blue",
		"y":   "yellow",
		"m":   "magenta",
		"c":   "cyan",
		"w":   "white",
		"b":   "black",
		"lb":  "lightblue",
		"lr":  "lightred",
		"lg":  "lightgreen",
		"ly":  "lightyellow",
		"lm":  "lightmagenta",
		"lc":  "lightcyan",
		"lw":  "lightwhite",
		"lk":  "lightblack",
		"bk":  "black",
		"rd":  "red",
		"wh":  "white",
		// full names pass through unchanged
		"red":   "red",
		"green": "green",
	}
	for in, want := range cases {
		got := makecolors.ColorMap(in)
		if got != want {
			t.Errorf("ColorMap(%q) = %q, want %q", in, got, want)
		}
	}
}

// ─── GetSort ──────────────────────────────────────────────────────────────────

func TestGetSort_DataOnly(t *testing.T) {
	cases := []struct {
		data       string
		wantFG     string
		wantBG     string
		wantAttrs  []string
	}{
		{"red-yellow", "red", "yellow", nil},
		{"red_yellow", "red", "yellow", nil},
		{"red,yellow", "red", "yellow", nil},
		{"r-b", "red", "black", nil},
		{"lb_r", "lightblue", "red", nil},
		{"bold-red", "red", "", []string{"bold"}},
		{"bold-red-yellow", "red", "yellow", []string{"bold"}},
		{"italic_blue_yellow", "blue", "yellow", []string{"italic"}},
		{"red", "red", "", nil},
		{"", "white", "", nil},
	}
	for _, tc := range cases {
		spec := makecolors.GetSort(tc.data, "", "", nil)
		if spec.Foreground != tc.wantFG {
			t.Errorf("GetSort(%q).FG = %q, want %q", tc.data, spec.Foreground, tc.wantFG)
		}
		if spec.Background != tc.wantBG {
			t.Errorf("GetSort(%q).BG = %q, want %q", tc.data, spec.Background, tc.wantBG)
		}
		for _, a := range tc.wantAttrs {
			found := false
			for _, ga := range spec.Attrs {
				if ga == a {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("GetSort(%q).Attrs missing %q, got %v", tc.data, a, spec.Attrs)
			}
		}
	}
}

func TestGetSort_ExplicitFGBG(t *testing.T) {
	spec := makecolors.GetSort("", "lightblue", "red", nil)
	if spec.Foreground != "lightblue" {
		t.Errorf("fg = %q, want lightblue", spec.Foreground)
	}
	if spec.Background != "red" {
		t.Errorf("bg = %q, want red", spec.Background)
	}
}

func TestGetSort_Attrs(t *testing.T) {
	spec := makecolors.GetSort("", "red", "", []string{"bold", "underline"})
	has := func(a string) bool {
		for _, x := range spec.Attrs {
			if x == a {
				return true
			}
		}
		return false
	}
	if !has("bold") || !has("underline") {
		t.Errorf("expected bold+underline in %v", spec.Attrs)
	}
}

// ─── MakeColors – plain text ──────────────────────────────────────────────────

func TestMakeColors_Plain(t *testing.T) {
	defer forceColor(t)()

	text := "hello"
	result := makecolors.MakeColors(text, makecolors.Options{Foreground: "red", Force: true})
	if stripped(result) != text {
		t.Errorf("stripped result %q ≠ %q", stripped(result), text)
	}
	if !strings.HasPrefix(result, "\x1b[") {
		t.Errorf("expected ANSI prefix, got %q", result)
	}
	if !strings.HasSuffix(result, makecolors.Reset) {
		t.Errorf("expected ANSI reset suffix, got %q", result)
	}
}

func TestMakeColors_Empty(t *testing.T) {
	result := makecolors.MakeColors("", makecolors.Options{Foreground: "red", Force: true})
	if result != "" {
		t.Errorf("empty input should return empty, got %q", result)
	}
}

func TestMakeColors_DisabledByEnv(t *testing.T) {
	os.Setenv("MAKE_COLORS", "0")
	defer os.Unsetenv("MAKE_COLORS")

	result := makecolors.MakeColors("hello", makecolors.Options{Foreground: "red"})
	if result != "hello" {
		t.Errorf("MAKE_COLORS=0 should return plain text, got %q", result)
	}
}

func TestMakeColors_ForcedOverridesDisabled(t *testing.T) {
	os.Setenv("MAKE_COLORS", "0")
	defer os.Unsetenv("MAKE_COLORS")

	result := makecolors.MakeColors("hello", makecolors.Options{Foreground: "red", Force: true})
	if !strings.Contains(result, "\x1b[") {
		t.Errorf("Force=true should emit ANSI even with MAKE_COLORS=0, got %q", result)
	}
}

// ─── Rich markup ──────────────────────────────────────────────────────────────

func TestMakeColors_RichMarkup_Simple(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("[red]hello[/]", makecolors.Options{Force: true})
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want %q", stripped(result), "hello")
	}
	if !strings.Contains(result, "\x1b[") {
		t.Error("expected ANSI codes in rich markup output")
	}
}

func TestMakeColors_RichMarkup_WithBackground(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("[white on red]hello[/]", makecolors.Options{Force: true})
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want %q", stripped(result), "hello")
	}
}

func TestMakeColors_RichMarkup_MultipleSegments(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("[red]Error:[/] [white]message[/]", makecolors.Options{Force: true})
	plain := stripped(result)
	if !strings.Contains(plain, "Error:") || !strings.Contains(plain, "message") {
		t.Errorf("unexpected plain text: %q", plain)
	}
}

func TestMakeColors_RichMarkup_WithAttrs(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("[bold red]important[/]", makecolors.Options{Force: true})
	if stripped(result) != "important" {
		t.Errorf("stripped = %q, want important", stripped(result))
	}
	// bold code is "1"
	if !strings.Contains(result, "1") {
		t.Errorf("expected bold code in %q", result)
	}
}

func TestMakeColors_RichMarkup_HexColor(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.MakeColors("[#FF0000]red hex[/]", makecolors.Options{Force: true})
	if stripped(result) != "red hex" {
		t.Errorf("stripped = %q, want 'red hex'", stripped(result))
	}
	if !strings.Contains(result, "38;2;255;0;0") {
		t.Errorf("expected truecolor sequence in %q", result)
	}
}

// ─── Sprint / functional helpers ──────────────────────────────────────────────

func TestSprint(t *testing.T) {
	defer forceColor(t)()

	result := makecolors.Sprint("test", "green", "")
	if stripped(result) != "test" {
		t.Errorf("stripped = %q", stripped(result))
	}
}

func TestPredefinedColors(t *testing.T) {
	defer forceColor(t)()

	funcs := []func(string) string{
		makecolors.Red, makecolors.Green, makecolors.Blue, makecolors.Yellow,
		makecolors.Magenta, makecolors.Cyan, makecolors.White, makecolors.LightRed,
	}
	for _, fn := range funcs {
		result := fn("hello")
		if stripped(result) != "hello" {
			t.Errorf("stripped = %q, want hello", stripped(result))
		}
	}
}

// ─── Color type ───────────────────────────────────────────────────────────────

func TestColorFormat(t *testing.T) {
	defer forceColor(t)()

	c := makecolors.NewColor("red", "black")
	result := c.Format("hello")
	if stripped(result) != "hello" {
		t.Errorf("stripped = %q, want hello", stripped(result))
	}
	if !strings.HasPrefix(result, "\x1b[") {
		t.Errorf("expected ANSI prefix in %q", result)
	}
}

func TestColorSprintf(t *testing.T) {
	defer forceColor(t)()

	c := makecolors.NewColor("blue", "")
	result := c.Sprintf("count=%d", 42)
	if stripped(result) != "count=42" {
		t.Errorf("stripped = %q, want count=42", stripped(result))
	}
}

// ─── StripANSI ────────────────────────────────────────────────────────────────

func TestStripANSI(t *testing.T) {
	input := "\x1b[1;31mhello\x1b[0m world \x1b[38;2;255;0;0mred\x1b[0m"
	want := "hello world red"
	got := makecolors.StripANSI(input)
	if got != want {
		t.Errorf("StripANSI = %q, want %q", got, want)
	}
}

// ─── SupportsColor ────────────────────────────────────────────────────────────

func TestSupportsColor_ReturnsBool(t *testing.T) {
	// Just check it doesn't panic; actual value depends on the test environment.
	_ = makecolors.SupportsColor()
}
