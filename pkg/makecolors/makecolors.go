// Package makecolors provides colored terminal text output with support for
// ANSI escape codes, rich console markup, hex colors, and cross-platform
// compatibility (Windows 10+, Linux, macOS).
//
// # Calling styles
//
//	MakeColors("text")                                       // no color, plain
//	MakeColors("text", Options{Foreground: "red"})           // named color
//	MakeColors("text", Options{Foreground: "#FF0000"})       // hex foreground
//	MakeColors("text", Options{Foreground: "#F00", Background: "#00FFFF"}) // hex fg+bg
//	MakeColors("[bold red]text[/]")                          // rich markup, no Options needed
//
// Options is always optional — omit it entirely when using rich markup or
// when you just want plain text pass-through.
//
// Author: Hadi Cahyadi <cumulus13@gmail.com>
// License: MIT
package makecolors

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/cumulus13/go-make_colors/pkg/hexansi"
)

// ─── ANSI code tables ──────────────────────────────────────────────────────────

var fgCodes = map[string]string{
	"black":        "30",
	"red":          "31",
	"green":        "32",
	"yellow":       "33",
	"blue":         "34",
	"magenta":      "35",
	"cyan":         "36",
	"white":        "37",
	"lightblack":   "90",
	"lightgrey":    "90",
	"lightred":     "91",
	"lightgreen":   "92",
	"lightyellow":  "93",
	"lightblue":    "94",
	"lightmagenta": "95",
	"lightcyan":    "96",
	"lightwhite":   "97",
}

var bgCodes = map[string]string{
	"black":           "40",
	"red":             "41",
	"green":           "42",
	"yellow":          "43",
	"blue":            "44",
	"magenta":         "45",
	"cyan":            "46",
	"white":           "47",
	"on_black":        "40",
	"on_red":          "41",
	"on_green":        "42",
	"on_yellow":       "43",
	"on_blue":         "44",
	"on_magenta":      "45",
	"on_cyan":         "46",
	"on_white":        "47",
	"lightblack":      "100",
	"lightgrey":       "100",
	"lightred":        "101",
	"lightgreen":      "102",
	"lightyellow":     "103",
	"lightblue":       "104",
	"lightmagenta":    "105",
	"lightcyan":       "106",
	"lightwhite":      "107",
	"on_lightblack":   "100",
	"on_lightgrey":    "100",
	"on_lightred":     "101",
	"on_lightgreen":   "102",
	"on_lightyellow":  "103",
	"on_lightblue":    "104",
	"on_lightmagenta": "105",
	"on_lightcyan":    "106",
	"on_lightwhite":   "107",
}

var attrCodes = map[string]string{
	"bold":          "1",
	"dim":           "2",
	"italic":        "3",
	"underline":     "4",
	"blink":         "5",
	"reverse":       "7",
	"strikethrough": "9",
	"strike":        "9",
}

var abbreviations = map[string]string{
	"b":  "black",
	"bk": "black",
	"bl": "blue",
	"r":  "red",
	"rd": "red",
	"re": "red",
	"g":  "green",
	"gr": "green",
	"ge": "green",
	"y":  "yellow",
	"ye": "yellow",
	"yl": "yellow",
	"m":  "magenta",
	"mg": "magenta",
	"ma": "magenta",
	"c":  "cyan",
	"cy": "cyan",
	"cn": "cyan",
	"w":  "white",
	"wh": "white",
	"wi": "white",
	"wt": "white",
	"lb": "lightblue",
	"lr": "lightred",
	"lg": "lightgreen",
	"ly": "lightyellow",
	"lm": "lightmagenta",
	"lc": "lightcyan",
	"lw": "lightwhite",
	"lk": "lightblack",
}

var knownAttrs = []string{
	"bold", "dim", "italic", "underline", "blink", "reverse", "strikethrough", "strike",
}

// Reset is the ANSI reset sequence.
const Reset = "\x1b[0m"

// ─── Environment helpers ───────────────────────────────────────────────────────

func envIs(key string, vals ...string) bool {
	v := os.Getenv(key)
	for _, want := range vals {
		if v == want {
			return true
		}
	}
	return false
}

func debugEnabled() bool  { return envIs("MAKE_COLORS_DEBUG", "1", "true", "True") }
func colorDisabled() bool { return envIs("MAKE_COLORS", "0") }
func colorForced() bool   { return envIs("MAKE_COLORS_FORCE", "1", "True") }

// ─── Terminal color-support detection ─────────────────────────────────────────

// SupportsColor reports whether the current terminal supports ANSI color output.
func SupportsColor() bool {
	if runtime.GOOS == "windows" {
		if os.Getenv("ANSICON") != "" || os.Getenv("ConEmuANSI") == "ON" {
			return true
		}
		if os.Getenv("WT_SESSION") != "" {
			return true
		}
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// ─── isHexColor ───────────────────────────────────────────────────────────────

// isHexColor reports whether s looks like a hex color (#RGB, #RRGGBB, RGB, RRGGBB).
func isHexColor(s string) bool {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 3 && len(s) != 6 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// ─── Hex → raw ANSI ───────────────────────────────────────────────────────────

// hexToFG converts a hex color string to a raw truecolor foreground ANSI sequence.
// Returns ("", false) if the input is not a valid hex color.
func hexToFG(color string) (string, bool) {
	if !isHexColor(color) {
		return "", false
	}
	res, err := hexansi.Convert(color, hexansi.ModeTrueColor)
	if err != nil {
		return "", false
	}
	return res.FG, true
}

// hexToBG converts a hex color string to a raw truecolor background ANSI sequence.
func hexToBG(color string) (string, bool) {
	if !isHexColor(color) {
		return "", false
	}
	res, err := hexansi.Convert(color, hexansi.ModeTrueColor)
	if err != nil {
		return "", false
	}
	return res.BG, true
}

// ─── Color map / abbreviation expansion ───────────────────────────────────────

// ColorMap expands short color codes to their full names.
func ColorMap(color string) string {
	if color == "" {
		return color
	}
	if full, ok := abbreviations[color]; ok {
		return full
	}
	return color
}

// ─── Attribute extraction ──────────────────────────────────────────────────────

var delimRe = regexp.MustCompile(`[-_,]+`)

func extractAttrs(text string) (string, []string) {
	if text == "" {
		return text, nil
	}
	var found []string
	cleaned := text
	for _, attr := range knownAttrs {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(attr) + `\b`)
		if re.MatchString(cleaned) {
			a := attr
			if a == "strike" {
				a = "strikethrough"
			}
			found = append(found, a)
			cleaned = re.ReplaceAllString(cleaned, "")
		}
	}
	cleaned = delimRe.ReplaceAllString(cleaned, "-")
	cleaned = strings.Trim(cleaned, "-_,")
	return strings.TrimSpace(cleaned), found
}

func uniqueAttrs(in []string) []string {
	seen := make(map[string]bool, len(in))
	var out []string
	for _, a := range in {
		if !seen[a] {
			seen[a] = true
			out = append(out, a)
		}
	}
	return out
}

// ─── ColorSpec & GetSort ───────────────────────────────────────────────────────

// ColorSpec is the fully resolved color specification.
// Foreground and Background may be either a named color ("red") or a raw ANSI
// escape sequence (when a hex color was supplied).
type ColorSpec struct {
	Foreground   string
	Background   string
	Attrs        []string
	fgIsRaw      bool // true = Foreground is already a complete ANSI sequence
	bgIsRaw      bool // true = Background is already a complete ANSI sequence
}

// GetSort parses a combined color specification string and returns a ColorSpec.
// Each of foreground and background may be a named color, an abbreviation,
// a combined "fg-bg" string, or a hex color ("#FF0000").
func GetSort(data, foreground, background string, attrs []string) ColorSpec {
	detected := append([]string(nil), attrs...)

	if data != "" {
		var dataAttrs []string
		data, dataAttrs = extractAttrs(data)
		detected = append(detected, dataAttrs...)

		if strings.ContainsAny(data, "-_,") {
			parts := nonEmpty(delimRe.Split(data, -1))
			switch len(parts) {
			case 0:
			case 1:
				foreground = parts[0]
			default:
				foreground = parts[0]
				background = parts[1]
			}
		} else {
			foreground = data
		}
	}

	if foreground != "" {
		var fgAttrs []string
		foreground, fgAttrs = extractAttrs(foreground)
		detected = append(detected, fgAttrs...)
	}
	if background != "" {
		var bgAttrs []string
		background, bgAttrs = extractAttrs(background)
		detected = append(detected, bgAttrs...)
	}

	if foreground != "" && len(foreground) > 2 && strings.ContainsAny(foreground, "-_,") {
		parts := nonEmpty(delimRe.Split(foreground, -1))
		if len(parts) >= 2 {
			foreground, background = parts[0], parts[1]
		} else if len(parts) == 1 {
			foreground = parts[0]
		}
	} else if background != "" && len(background) > 2 && strings.ContainsAny(background, "-_,") {
		parts := nonEmpty(delimRe.Split(background, -1))
		if len(parts) >= 2 {
			foreground, background = parts[0], parts[1]
		} else if len(parts) == 1 {
			background = parts[0]
		}
	}

	if foreground == "" {
		foreground = "white"
	}

	// Expand abbreviations (only for non-hex values)
	if !isHexColor(foreground) && len(foreground) <= 2 {
		foreground = ColorMap(foreground)
	}
	if background != "" && !isHexColor(background) && len(background) <= 2 {
		background = ColorMap(background)
	}

	spec := ColorSpec{
		Foreground: strings.TrimSpace(foreground),
		Background: strings.TrimSpace(background),
		Attrs:      uniqueAttrs(detected),
	}

	// Resolve hex colors to raw ANSI sequences immediately
	if isHexColor(spec.Foreground) {
		if raw, ok := hexToFG(spec.Foreground); ok {
			spec.Foreground = raw
			spec.fgIsRaw = true
		}
	}
	if spec.Background != "" && isHexColor(spec.Background) {
		if raw, ok := hexToBG(spec.Background); ok {
			spec.Background = raw
			spec.bgIsRaw = true
		}
	}

	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "[DEBUG] GetSort → fg=%q(raw=%v) bg=%q(raw=%v) attrs=%v\n",
			spec.Foreground, spec.fgIsRaw, spec.Background, spec.bgIsRaw, spec.Attrs)
	}

	return spec
}

func nonEmpty(ss []string) []string {
	var out []string
	for _, s := range ss {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// ─── ANSI sequence builder ────────────────────────────────────────────────────

// buildANSI constructs the opening ANSI escape for the given ColorSpec.
// Handles both named colors (looked up in tables) and raw sequences (hex colors).
func buildANSI(spec ColorSpec) string {
	// Fast path: both fg and bg are raw sequences — combine them with attr codes.
	if spec.fgIsRaw || spec.bgIsRaw {
		var sb strings.Builder
		for _, a := range spec.Attrs {
			if code, ok := attrCodes[a]; ok {
				sb.WriteString("\x1b[" + code + "m")
			}
		}
		if spec.bgIsRaw && spec.Background != "" {
			sb.WriteString(spec.Background)
		} else if !spec.bgIsRaw && spec.Background != "" {
			bgKey := strings.TrimPrefix(spec.Background, "on_")
			if code, ok := bgCodes[bgKey]; ok {
				sb.WriteString("\x1b[" + code + "m")
			} else if code, ok := bgCodes[spec.Background]; ok {
				sb.WriteString("\x1b[" + code + "m")
			}
		}
		if spec.fgIsRaw && spec.Foreground != "" {
			sb.WriteString(spec.Foreground)
		} else if !spec.fgIsRaw && spec.Foreground != "" {
			if code, ok := fgCodes[spec.Foreground]; ok {
				sb.WriteString("\x1b[" + code + "m")
			}
		}
		if sb.Len() == 0 {
			return ""
		}
		return sb.String() // note: no trailing reset here; Colorize adds it
	}

	// Normal path: both are named colors.
	var codes []string
	for _, a := range spec.Attrs {
		if code, ok := attrCodes[a]; ok {
			codes = appendUniq(codes, code)
		}
	}
	if spec.Background != "" {
		bgKey := strings.TrimPrefix(spec.Background, "on_")
		if code, ok := bgCodes[bgKey]; ok {
			codes = appendUniq(codes, code)
		} else if code, ok := bgCodes[spec.Background]; ok {
			codes = appendUniq(codes, code)
		}
	}
	if spec.Foreground != "" {
		if code, ok := fgCodes[spec.Foreground]; ok {
			codes = appendUniq(codes, code)
		}
	}
	if len(codes) == 0 {
		return ""
	}
	return "\x1b[" + strings.Join(codes, ";") + "m"
}

func appendUniq(sl []string, v string) []string {
	for _, s := range sl {
		if s == v {
			return sl
		}
	}
	return append(sl, v)
}

// Colorize wraps text in ANSI codes described by spec.
// Works for both named-color specs and hex-color specs.
func Colorize(text string, spec ColorSpec) string {
	open := buildANSI(spec)
	if open == "" {
		return text
	}
	return open + text + Reset
}

// ─── Options ──────────────────────────────────────────────────────────────────

// Options configures a MakeColors call.
// Every field is optional — the zero value means "no color / use defaults".
//
// Foreground and Background accept:
//   - Named colors:   "red", "lightblue", "cyan", …
//   - Abbreviations:  "r", "lb", "c", …
//   - Hex colors:     "#FF0000", "#F00", "FF0000"
//   - Combined:       "bold-red-yellow", "lb_r"  (Foreground only)
type Options struct {
	Foreground string   // fg color: name, abbreviation, hex, or combined string
	Background string   // bg color: name, abbreviation, or hex
	Attrs      []string // ["bold", "italic", "underline", …]
	Force      bool     // force ANSI output even when stdout is not a TTY
}

// ─── MakeColors ───────────────────────────────────────────────────────────────

// MakeColors applies color formatting to text.
//
// Options is variadic — you may omit it entirely:
//
//	MakeColors("hello")                                    // plain passthrough
//	MakeColors("hello", Options{Foreground: "red"})        // named color
//	MakeColors("hello", Options{Foreground: "#00FFFF"})    // hex color
//	MakeColors("[bold red]hello[/]")                       // rich markup, no Options
//	MakeColors("[bold red]hello[/]", Options{Force: true}) // markup + force
//
// Environment variables:
//
//	MAKE_COLORS=0         → disable all output
//	MAKE_COLORS_FORCE=1   → force output regardless of TTY
//	MAKE_COLORS_DEBUG=1   → print parsing info to stderr
func MakeColors(text string, opts ...Options) string {
	if text == "" {
		return ""
	}
	o := mergeOpts(opts)

	if strings.Contains(text, "[") && strings.Contains(text, "[/]") {
		return applyRichMarkup(text, o)
	}
	return applyPlain(text, o)
}

// mergeOpts returns the first Options element, or a zero-value Options if none
// was provided. Multiple elements are intentionally not merged — only the first
// is used; the variadic signature purely makes the argument optional.
func mergeOpts(opts []Options) Options {
	if len(opts) == 0 {
		return Options{}
	}
	return opts[0]
}

// applyPlain applies color to non-markup text.
func applyPlain(text string, o Options) string {
	// A completely empty Options means no-op (plain text).
	if o.Foreground == "" && o.Background == "" && len(o.Attrs) == 0 && !o.Force && !colorForced() {
		if colorDisabled() || !SupportsColor() {
			return text
		}
	}

	fg := o.Foreground
	if fg == "" && (o.Background != "" || len(o.Attrs) > 0) {
		fg = "white" // sensible default when only bg/attrs are set
	}

	spec := GetSort("", fg, o.Background, o.Attrs)

	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "[DEBUG] applyPlain fg=%q bg=%q attrs=%v\n",
			spec.Foreground, spec.Background, spec.Attrs)
	}

	if !shouldColor(o.Force) {
		return text
	}
	return Colorize(text, spec)
}

// applyRichMarkup parses and applies rich markup formatting.
func applyRichMarkup(text string, o Options) string {
	sections := parseRichMarkup(text)

	var sb strings.Builder
	for _, sec := range sections {
		if sec.content == "" {
			continue
		}

		fg := sec.fg
		if fg == "" {
			fg = o.Foreground
		}
		bg := sec.bg
		if bg == "" {
			bg = o.Background
		}

		var colored string
		if sec.isFGRaw || sec.isBGRaw {
			colored = applyRawANSI(sec.content, sec.fg, sec.isFGRaw, sec.bg, sec.isBGRaw, sec.attrs)
		} else {
			spec := GetSort("", fg, bg, sec.attrs)
			colored = Colorize(sec.content, spec)
		}
		sb.WriteString(colored)
	}

	output := sb.String()
	if !shouldColor(o.Force) {
		return ansiStripRe.ReplaceAllString(output, "")
	}
	return output
}

// applyRawANSI combines raw ANSI fg/bg sequences with attribute codes.
// Used when hex colors were resolved to raw sequences during markup parsing.
func applyRawANSI(text, fgRaw string, hasFG bool, bgRaw string, hasBG bool, attrs []string) string {
	var sb strings.Builder
	for _, a := range attrs {
		if code, ok := attrCodes[a]; ok {
			sb.WriteString("\x1b[" + code + "m")
		}
	}
	if hasBG && bgRaw != "" {
		sb.WriteString(bgRaw)
	}
	if hasFG && fgRaw != "" {
		sb.WriteString(fgRaw)
	}
	sb.WriteString(text)
	sb.WriteString(Reset)
	return sb.String()
}

func shouldColor(force bool) bool {
	if force || colorForced() {
		return true
	}
	if colorDisabled() {
		return false
	}
	return SupportsColor()
}

// ─── Rich markup parser ───────────────────────────────────────────────────────

type markupResult struct {
	content  string
	fg       string
	bg       string
	attrs    []string
	isFGRaw  bool
	isBGRaw  bool
}

var (
	richPattern = regexp.MustCompile(`\[([^\[\]]+?)\](.*?)\[/\]`)
	ansiStripRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

const (
	escapedLeft  = "\x00ESCAPED_LEFT\x00"
	escapedRight = "\x00ESCAPED_RIGHT\x00"
)

func parseRichMarkup(text string) []markupResult {
	proc := strings.ReplaceAll(text, `\[`, escapedLeft)
	proc = strings.ReplaceAll(proc, `\]`, escapedRight)

	matches := richPattern.FindAllStringSubmatchIndex(proc, -1)
	if len(matches) == 0 {
		final := strings.ReplaceAll(proc, escapedLeft, "[")
		final = strings.ReplaceAll(final, escapedRight, "]")
		return []markupResult{{content: final}}
	}

	var results []markupResult
	for i, match := range matches {
		markup := strings.ToLower(strings.TrimSpace(proc[match[2]:match[3]]))
		content := proc[match[4]:match[5]]

		if i == 0 && match[0] > 0 {
			content = proc[:match[0]] + content
		}

		afterClose := match[1]
		var afterText string
		if i < len(matches)-1 {
			afterText = proc[afterClose:matches[i+1][0]]
		} else {
			afterText = proc[afterClose:]
		}
		if afterText != "" {
			content += afterText
		}

		content = strings.ReplaceAll(content, escapedLeft, "[")
		content = strings.ReplaceAll(content, escapedRight, "]")

		res := parseMarkupTag(markup)
		res.content = content
		results = append(results, res)
	}
	return results
}

func parseMarkupTag(markup string) markupResult {
	var res markupResult
	parts := strings.Fields(markup)

	var colorParts []string
	for _, p := range parts {
		isAttr := false
		for _, attr := range knownAttrs {
			if p == attr {
				a := p
				if a == "strike" {
					a = "strikethrough"
				}
				res.attrs = append(res.attrs, a)
				isAttr = true
				break
			}
		}
		if !isAttr {
			colorParts = append(colorParts, p)
		}
	}

	switch {
	case len(colorParts) >= 3 && colorParts[1] == "on":
		res.fg = resolveMarkupColor(colorParts[0], false, &res.isFGRaw)
		res.bg = resolveMarkupColor(colorParts[2], true, &res.isBGRaw)
	case len(colorParts) == 1:
		res.fg = resolveMarkupColor(colorParts[0], false, &res.isFGRaw)
	}

	return res
}

func resolveMarkupColor(token string, isBG bool, rawFlag *bool) string {
	if isHexColor(token) {
		res, err := hexansi.Convert(token, hexansi.ModeTrueColor)
		if err == nil {
			*rawFlag = true
			if isBG {
				return res.BG
			}
			return res.FG
		}
	}
	return token
}

// ─── Convenience functions ────────────────────────────────────────────────────

// Make is the shortest alias for MakeColors.
func Make(text string, opts ...Options) string { return MakeColors(text, opts...) }

// Sprint returns colored text: Sprint(text, fg, bg, attrs...).
// fg and bg may be named colors, abbreviations, or hex strings.
//
//	Sprint("hello", "red", "")
//	Sprint("hello", "#FF0000", "#000000")
//	Sprint("hello", "bold-red", "")
//	Sprint("hello", "red", "", "bold", "underline")
func Sprint(text, foreground, background string, attrs ...string) string {
	return MakeColors(text, Options{
		Foreground: foreground,
		Background: background,
		Attrs:      attrs,
	})
}

// Sprintf formats then colorizes: Sprintf(fg, bg, attrs, format, args...).
func Sprintf(foreground, background string, attrs []string, format string, a ...interface{}) string {
	return Sprint(fmt.Sprintf(format, a...), foreground, background, attrs...)
}

// Println prints colored text + newline to stdout.
func Println(text, foreground, background string, attrs ...string) {
	fmt.Println(Sprint(text, foreground, background, attrs...))
}

// StripANSI removes all ANSI escape sequences from s.
func StripANSI(s string) string { return ansiStripRe.ReplaceAllString(s, "") }

// ─── Single-color helper functions ────────────────────────────────────────────

func Red(text string) string          { return Sprint(text, "red", "") }
func Green(text string) string        { return Sprint(text, "green", "") }
func Blue(text string) string         { return Sprint(text, "blue", "") }
func Yellow(text string) string       { return Sprint(text, "yellow", "") }
func Magenta(text string) string      { return Sprint(text, "magenta", "") }
func Cyan(text string) string         { return Sprint(text, "cyan", "") }
func White(text string) string        { return Sprint(text, "white", "") }
func Black(text string) string        { return Sprint(text, "black", "") }
func LightRed(text string) string     { return Sprint(text, "lightred", "") }
func LightGreen(text string) string   { return Sprint(text, "lightgreen", "") }
func LightBlue(text string) string    { return Sprint(text, "lightblue", "") }
func LightYellow(text string) string  { return Sprint(text, "lightyellow", "") }
func LightMagenta(text string) string { return Sprint(text, "lightmagenta", "") }
func LightCyan(text string) string    { return Sprint(text, "lightcyan", "") }
func LightWhite(text string) string   { return Sprint(text, "lightwhite", "") }
func Bold(text string) string         { return Sprint(text, "white", "", "bold") }
func Dim(text string) string          { return Sprint(text, "white", "", "dim") }
func Italic(text string) string       { return Sprint(text, "white", "", "italic") }
func Underline(text string) string    { return Sprint(text, "white", "", "underline") }
func Strikethrough(text string) string { return Sprint(text, "white", "", "strikethrough") }
