// Package hexansi provides conversion between hex colors and ANSI escape codes.
// It supports truecolor (24-bit), 256-color, and 16-color terminal modes.
//
// Author: Hadi Cahyadi <cumulus13@gmail.com>
// License: MIT
package hexansi

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Mode represents the ANSI color mode.
type Mode int

const (
	// ModeTrueColor uses 24-bit RGB ANSI escape codes.
	ModeTrueColor Mode = iota
	// Mode256 uses 256-color ANSI escape codes.
	Mode256
	// Mode16 uses basic 16-color ANSI escape codes.
	Mode16
)

// Reset is the ANSI reset sequence.
const Reset = "\x1b[0m"

// Result holds the ANSI escape sequences for a converted color.
type Result struct {
	FG  string
	BG  string
	RGB [3]uint8
}

// ColorDatabase maps common color names to their RGB values.
var colorDatabase = map[string][3]uint8{
	"black":               {0, 0, 0},
	"white":               {255, 255, 255},
	"gray":                {128, 128, 128},
	"grey":                {128, 128, 128},
	"silver":              {192, 192, 192},
	"red":                 {255, 0, 0},
	"lime":                {0, 255, 0},
	"blue":                {0, 0, 255},
	"yellow":              {255, 255, 0},
	"magenta":             {255, 0, 255},
	"cyan":                {0, 255, 255},
	"maroon":              {128, 0, 0},
	"green":               {0, 128, 0},
	"navy":                {0, 0, 128},
	"olive":               {128, 128, 0},
	"purple":              {128, 0, 128},
	"teal":                {0, 128, 128},
	"orange":              {255, 165, 0},
	"pink":                {255, 192, 203},
	"brown":               {165, 42, 42},
	"gold":                {255, 215, 0},
	"tan":                 {210, 180, 140},
	"tomato":              {255, 99, 71},
	"coral":               {255, 127, 80},
	"crimson":             {220, 20, 60},
	"darkred":             {139, 0, 0},
	"deeppink":            {255, 20, 147},
	"hotpink":             {255, 105, 180},
	"blueviolet":          {138, 43, 226},
	"indigo":              {75, 0, 130},
	"slateblue":           {106, 90, 205},
	"mediumpurple":        {147, 112, 219},
	"deepskyblue":         {0, 191, 255},
	"skyblue":             {135, 206, 235},
	"lightblue":           {173, 216, 230},
	"powderblue":          {176, 224, 230},
	"turquoise":           {64, 224, 208},
	"aquamarine":          {127, 255, 212},
	"springgreen":         {0, 255, 127},
	"lightgreen":          {144, 238, 144},
	"palegreen":           {152, 251, 152},
	"forestgreen":         {34, 139, 34},
	"yellowgreen":         {154, 205, 50},
	"greenyellow":         {173, 255, 47},
	"lightyellow":         {255, 255, 224},
	"lemonchiffon":        {255, 250, 205},
	"sandybrown":          {244, 164, 96},
	"chocolate":           {210, 105, 30},
	"saddlebrown":         {139, 69, 19},
	"sienna":              {160, 82, 45},
	"khaki":               {240, 230, 140},
	"lavender":            {230, 230, 250},
	"plum":                {221, 160, 221},
	"violet":              {238, 130, 238},
	"orchid":              {218, 112, 214},
	"fuchsia":             {255, 0, 255},
	"steelblue":           {70, 130, 180},
	"cornflowerblue":      {100, 149, 237},
	"dodgerblue":          {30, 144, 255},
	"mediumblue":          {0, 0, 205},
	"midnightblue":        {25, 25, 112},
	"cadetblue":           {95, 158, 160},
	"darkturquoise":       {0, 206, 209},
	"mediumturquoise":     {72, 209, 204},
	"paleturquoise":       {175, 238, 238},
	"aqua":                {0, 255, 255},
	"darkcyan":            {0, 139, 139},
	"lightseagreen":       {32, 178, 170},
	"seagreen":            {46, 139, 87},
	"mediumseagreen":      {60, 179, 113},
	"darkseagreen":        {143, 188, 143},
	"darkolivegreen":      {85, 107, 47},
	"lawngreen":           {124, 252, 0},
	"chartreuse":          {127, 255, 0},
	"limegreen":           {50, 205, 50},
	"darkgreen":           {0, 100, 0},
	"wheat":               {245, 222, 179},
	"burlywood":           {222, 184, 135},
	"rosybrown":           {188, 143, 143},
}

// HexToRGB converts a hex color string to an [3]uint8 RGB tuple.
// Accepts formats: "#RRGGBB", "RRGGBB", "#RGB", "RGB".
func HexToRGB(hex string) ([3]uint8, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return [3]uint8{}, fmt.Errorf("invalid hex color: #%s", hex)
	}
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return [3]uint8{}, fmt.Errorf("invalid hex color: #%s", hex)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return [3]uint8{}, fmt.Errorf("invalid hex color: #%s", hex)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return [3]uint8{}, fmt.Errorf("invalid hex color: #%s", hex)
	}
	return [3]uint8{uint8(r), uint8(g), uint8(b)}, nil
}

// rgbToANSI16 converts RGB to a basic 16-color ANSI code.
func rgbToANSI16(r, g, b uint8) int {
	const threshold = 128
	if int(r) < threshold && int(g) < threshold && int(b) < threshold {
		if r == g && g == b {
			if r < 64 {
				return 30
			}
			return 90
		}
	}
	var bits int
	if int(r) > threshold {
		bits |= 1
	}
	if int(g) > threshold {
		bits |= 2
	}
	if int(b) > threshold {
		bits |= 4
	}
	if bits == 0 {
		return 30
	}
	base := bits + 29
	if int(r) > 200 || int(g) > 200 || int(b) > 200 {
		base += 60
	}
	return base
}

// rgbToANSI256 converts RGB to the nearest ANSI 256-color code.
func rgbToANSI256(r, g, b uint8) int {
	if r == g && g == b {
		if r < 8 {
			return 16
		}
		if r > 248 {
			return 231
		}
		return int(math.Round(float64(r-8)/247.0*24)) + 232
	}
	to6cube := func(v uint8) int {
		if int(v) < 48 {
			return 0
		}
		if int(v) < 115 {
			return 1
		}
		return (int(v) - 35) / 40
	}
	ir := to6cube(r)
	ig := to6cube(g)
	ib := to6cube(b)
	return 16 + 36*ir + 6*ig + ib
}

// Convert converts a hex color string to ANSI escape sequences using the given mode.
func Convert(hexColor string, mode Mode) (Result, error) {
	rgb, err := HexToRGB(hexColor)
	if err != nil {
		return Result{}, err
	}
	r, g, b := rgb[0], rgb[1], rgb[2]

	var fg, bg string
	switch mode {
	case ModeTrueColor:
		fg = fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
		bg = fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
	case Mode256:
		code := rgbToANSI256(r, g, b)
		fg = fmt.Sprintf("\x1b[38;5;%dm", code)
		bg = fmt.Sprintf("\x1b[48;5;%dm", code)
	case Mode16:
		code := rgbToANSI16(r, g, b)
		fg = fmt.Sprintf("\x1b[%dm", code)
		bg = fmt.Sprintf("\x1b[%dm", code+10)
	default:
		return Result{}, fmt.Errorf("invalid mode: %d", mode)
	}
	return Result{FG: fg, BG: bg, RGB: rgb}, nil
}

// NameToHex converts a color name to a hex string "#RRGGBB".
// Returns an error if the name is not found.
func NameToHex(name string) (string, error) {
	norm := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(name, " ", "")))
	if rgb, ok := colorDatabase[norm]; ok {
		return fmt.Sprintf("#%02X%02X%02X", rgb[0], rgb[1], rgb[2]), nil
	}
	// partial match
	for k, rgb := range colorDatabase {
		if strings.Contains(norm, k) || strings.Contains(k, norm) {
			return fmt.Sprintf("#%02X%02X%02X", rgb[0], rgb[1], rgb[2]), nil
		}
	}
	return "", fmt.Errorf("color name %q not found", name)
}

// NameToANSI converts a color name to ANSI escape sequences.
func NameToANSI(name string, mode Mode) (Result, error) {
	hex, err := NameToHex(name)
	if err != nil {
		return Result{}, err
	}
	return Convert(hex, mode)
}

// ToANSI accepts either a hex color or color name and converts to ANSI sequences.
func ToANSI(input string, mode Mode) (Result, error) {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "#") || isHex(input) {
		return Convert(input, mode)
	}
	return NameToANSI(input, mode)
}

// IsHex reports whether s looks like a 3- or 6-digit hex string (no leading #).
func IsHex(s string) bool { return isHex(s) }

// isHex reports whether s looks like a 3- or 6-digit hex string.
func isHex(s string) bool {
	if len(s) != 3 && len(s) != 6 {
		return false
	}
	_, err := strconv.ParseUint(s, 16, 32)
	return err == nil
}

// HexToColorName returns the nearest human-readable color name for a hex value.
func HexToColorName(hexColor string) (string, error) {
	rgb, err := HexToRGB(hexColor)
	if err != nil {
		return "", err
	}
	r, g, b := float64(rgb[0]), float64(rgb[1]), float64(rgb[2])

	minDist := math.MaxFloat64
	closest := "unknown"
	for name, crgb := range colorDatabase {
		dr := r - float64(crgb[0])
		dg := g - float64(crgb[1])
		db := b - float64(crgb[2])
		dist := math.Sqrt(dr*dr + dg*dg + db*db)
		if dist < minDist {
			minDist = dist
			closest = name
		}
	}

	brightness := (r + g + b) / 3
	if minDist < 30 {
		return strings.Title(closest), nil
	} else if brightness < 50 {
		return "Dark " + strings.Title(closest), nil
	} else if brightness > 200 {
		return "Light " + strings.Title(closest), nil
	}
	return strings.Title(closest), nil
}
