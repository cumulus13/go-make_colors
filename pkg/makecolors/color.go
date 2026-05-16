package makecolors

import "fmt"

// Color holds a pre-computed ANSI sequence for a foreground+background combination.
// It can be used as a formatter for multiple strings without re-parsing the spec.
//
// Example:
//
//	c := makecolors.NewColor("red", "black")
//	fmt.Println(c.Format("Error!"))
//	fmt.Println(c.Sprintf("Exit code: %d", 1))
type Color struct {
	spec  ColorSpec
	open  string // opening ANSI sequence
}

// NewColor creates a Color for the given foreground and optional background.
func NewColor(foreground, background string, attrs ...string) *Color {
	spec := GetSort("", foreground, background, attrs)
	return &Color{
		spec: spec,
		open: buildANSI(spec.Foreground, spec.Background, spec.Attrs),
	}
}

// Format wraps text in this Color's ANSI sequence.
func (c *Color) Format(text string) string {
	if c.open == "" {
		return text
	}
	return c.open + text + Reset
}

// Sprintf formats the string then wraps it in this Color's ANSI sequence.
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.Format(fmt.Sprintf(format, a...))
}

// ANSIOpen returns the raw opening ANSI escape sequence (without the reset).
func (c *Color) ANSIOpen() string { return c.open }

// Spec returns the resolved ColorSpec.
func (c *Color) Spec() ColorSpec { return c.spec }

// String implements fmt.Stringer, returning the opening sequence.
// This allows a Color to be used directly in format strings alongside text.
func (c *Color) String() string { return c.open }

// Colors is an alias for Color kept for naming compatibility with the Python API.
type Colors = Color

// NewColors is an alias for NewColor.
var NewColors = NewColor

// ─── Pre-built color singletons ───────────────────────────────────────────────
// These are lazily constructed on first call to avoid init-order issues.

var (
	ColorRed          = NewColor("red", "")
	ColorGreen        = NewColor("green", "")
	ColorBlue         = NewColor("blue", "")
	ColorYellow       = NewColor("yellow", "")
	ColorMagenta      = NewColor("magenta", "")
	ColorCyan         = NewColor("cyan", "")
	ColorWhite        = NewColor("white", "")
	ColorBlack        = NewColor("black", "")
	ColorLightRed     = NewColor("lightred", "")
	ColorLightGreen   = NewColor("lightgreen", "")
	ColorLightBlue    = NewColor("lightblue", "")
	ColorLightYellow  = NewColor("lightyellow", "")
	ColorLightMagenta = NewColor("lightmagenta", "")
	ColorLightCyan    = NewColor("lightcyan", "")
	ColorLightWhite   = NewColor("lightwhite", "")
	ColorBold         = NewColor("white", "", "bold")
	ColorDim          = NewColor("white", "", "dim")
	ColorItalic       = NewColor("white", "", "italic")
	ColorUnderline    = NewColor("white", "", "underline")
)
