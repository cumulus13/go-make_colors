package makecolors

import "fmt"

// Color holds a pre-computed ANSI sequence for a fg+bg+attrs combination.
// Foreground and Background accept named colors, abbreviations, or hex strings.
//
//	c := NewColor("red", "black")
//	c := NewColor("#FF0000", "#000000")
//	c := NewColor("white", "red", "bold")
type Color struct {
	spec ColorSpec
	open string // pre-built opening ANSI sequence
}

// NewColor creates a Color for the given foreground, optional background, and
// optional attribute list. fg and bg may be named colors or hex strings.
func NewColor(foreground, background string, attrs ...string) *Color {
	spec := GetSort("", foreground, background, attrs)
	return &Color{spec: spec, open: buildANSI(spec)}
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

// String implements fmt.Stringer, returning the opening sequence so a Color
// can be used directly in format strings.
func (c *Color) String() string { return c.open }

// Colors is an alias kept for naming compatibility.
type Colors = Color

// NewColors is an alias for NewColor.
var NewColors = NewColor

// ─── Pre-built singletons ─────────────────────────────────────────────────────

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
