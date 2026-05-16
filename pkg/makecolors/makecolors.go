// Package makecolors provides colored terminal text output with support for
// ANSI escape codes, rich console markup, hex colors, and cross-platform
// compatibility (Windows 10+, Linux, macOS).
//
// Features:
//   - Standard and light foreground/background colors
//   - Color abbreviations (r, g, bl, lb, ...)
//   - Combined format strings ("red-yellow", "bold-red", "lb_r")
//   - Rich markup: "[bold red on yellow]text[/]"
//   - Hex color support in markup: "[#FF0000 on #00FF00]text[/]"
//   - Text attributes: bold, dim, italic, underline, blink, reverse, strikethrough
//   - Environment variable controls: MAKE_COLORS, MAKE_COLORS_FORCE, MAKE_COLORS_DEBUG
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

// fgCodes maps color names to their ANSI foreground codes.
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

// bgCodes maps color names to their ANSI background codes.
var bgCodes = map[string]string{
	"black":          "40",
	"red":            "41",
	"green":          "42",
	"yellow":         "43",
	"blue":           "44",
	"magenta":        "45",
	"cyan":           "46",
	"white":          "47",
	"on_black":       "40",
	"on_red":         "41",
	"on_green":       "42",
	"on_yellow":      "43",
	"on_blue":        "44",
	"on_magenta":     "45",
	"on_cyan":        "46",
	"on_white":       "47",
	"lightblack":     "100",
	"lightgrey":      "100",
	"lightred":       "101",
	"lightgreen":     "102",
	"lightyellow":    "103",
	"lightblue":      "104",
	"lightmagenta":   "105",
	"lightcyan":      "106",
	"lightwhite":     "107",
	"on_lightblack":  "100",
	"on_lightgrey":   "100",
	"on_lightred":    "101",
	"on_lightgreen":  "102",
	"on_lightyellow": "103",
	"on_lightblue":   "104",
	"on_lightmagenta": "105",
	"on_lightcyan":   "106",
	"on_lightwhite":  "107",
}

// attrCodes maps attribute names to their ANSI codes.
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

// abbreviations maps short color codes to full color names.
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

// knownAttrs is the ordered list of recognizable text attributes.
var knownAttrs = []string{
	"bold", "dim", "italic", "underline", "blink", "reverse", "strikethrough", "strike",
}

// Reset is the ANSI reset sequence.
const Reset = "\x1b[0m"

// ─── Environment helpers ────────────────────────────────────────────────────────

func envIs(key string, vals ...string) bool {
	v := os.Getenv(key)
	for _, want := range vals {
		if v == want {
			return true
		}
	}
	return false
}

func debugEnabled() bool { return envIs("MAKE_COLORS_DEBUG", "1", "true", "True") }
func colorDisabled() bool { return envIs("MAKE_COLORS", "0") }
func colorForced() bool   { return envIs("MAKE_COLORS_FORCE", "1", "True") }

// ─── Terminal color-support detection ─────────────────────────────────────────

// SupportsColor reports whether the current terminal supports ANSI color output.
func SupportsColor() bool {
	if runtime.GOOS == "windows" {
		// On Windows 10 1511+ VT processing is supported; we trust the env.
		if os.Getenv("ANSICON") != "" || os.Getenv("ConEmuANSI") == "ON" {
			return true
		}
		// Check TERM or WT_SESSION (Windows Terminal)
		if os.Getenv("WT_SESSION") != "" {
			return true
		}
	}
	// Standard Unix TTY check via file stat
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// ─── Color map / abbreviation expansion ───────────────────────────────────────

// ColorMap expands short color codes to their full names.
// Returns the input unchanged if it is already a full name or unknown short code.
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

// extractAttrs strips recognized attribute tokens from a color string and
// returns the cleaned string together with the found attribute names.
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
	// Collapse leftover delimiters
	cleaned = delimRe.ReplaceAllString(cleaned, "-")
	cleaned = strings.Trim(cleaned, "-_,")
	return strings.TrimSpace(cleaned), found
}

// uniqueAttrs returns a deduplicated slice preserving order.
func uniqueAttrs(in []string) []string {
	seen := make(map[string]bool, len(in))
	out := in[:0:0]
	for _, a := range in {
		if !seen[a] {
			seen[a] = true
			out = append(out, a)
		}
	}
	return out
}

// ─── GetSort ──────────────────────────────────────────────────────────────────

// ColorSpec is the result of parsing a color specification.
type ColorSpec struct {
	Foreground string
	Background string // empty string means "none"
	Attrs      []string
}

// GetSort parses a combined color specification string and returns a ColorSpec.
//
// The data parameter may use any of:
//   - "red-yellow"         → fg=red, bg=yellow
//   - "bold-red-black"     → fg=red, bg=black, attrs=[bold]
//   - "r_b"                → fg=red, bg=black (abbreviations)
//   - "lightblue"          → fg=lightblue
//   - "[bold red on blue]" → handled by the rich parser upstream
func GetSort(data, foreground, background string, attrs []string) ColorSpec {
	detected := append([]string(nil), attrs...)

	if data != "" {
		var dataAttrs []string
		data, dataAttrs = extractAttrs(data)
		detected = append(detected, dataAttrs...)

		if strings.ContainsAny(data, "-_,") {
			parts := delimRe.Split(data, -1)
			parts = nonEmpty(parts)
			switch len(parts) {
			case 0:
				// nothing
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

	// Strip attrs from explicit fg/bg strings
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

	// Handle nested delimiters in fg or bg after attr extraction
	if foreground != "" && len(foreground) > 2 && strings.ContainsAny(foreground, "-_,") {
		parts := nonEmpty(delimRe.Split(foreground, -1))
		if len(parts) >= 2 {
			foreground = parts[0]
			background = parts[1]
		} else if len(parts) == 1 {
			foreground = parts[0]
		}
	} else if background != "" && len(background) > 2 && strings.ContainsAny(background, "-_,") {
		parts := nonEmpty(delimRe.Split(background, -1))
		if len(parts) >= 2 {
			foreground = parts[0]
			background = parts[1]
		} else if len(parts) == 1 {
			background = parts[0]
		}
	}

	if foreground == "" {
		foreground = "white"
	}

	// Expand abbreviations
	if len(foreground) <= 2 {
		foreground = ColorMap(foreground)
	}
	if background != "" && len(background) <= 2 {
		background = ColorMap(background)
	}

	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "[DEBUG] GetSort → fg=%q bg=%q attrs=%v\n", foreground, background, detected)
	}

	return ColorSpec{
		Foreground: strings.TrimSpace(foreground),
		Background: strings.TrimSpace(background),
		Attrs:      uniqueAttrs(detected),
	}
}

func nonEmpty(ss []string) []string {
	out := ss[:0:0]
	for _, s := range ss {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// ─── ANSI sequence builder ────────────────────────────────────────────────────

// buildANSI constructs the opening ANSI escape sequence for the given color
// specification. Returns an empty string if no codes apply.
func buildANSI(fg, bg string, attrs []string) string {
	var codes []string

	for _, a := range attrs {
		if code, ok := attrCodes[a]; ok {
			codes = appendUniq(codes, code)
		}
	}
	if bg != "" {
		// strip leading "on_" for lookup
		bgKey := strings.TrimPrefix(bg, "on_")
		if code, ok := bgCodes[bgKey]; ok {
			codes = appendUniq(codes, code)
		} else if code, ok := bgCodes[bg]; ok {
			codes = appendUniq(codes, code)
		}
	}
	if fg != "" {
		if code, ok := fgCodes[fg]; ok {
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

// Colorize wraps text in ANSI escape codes described by the ColorSpec.
// Returns the plain text if the spec produces no codes.
func Colorize(text string, spec ColorSpec) string {
	open := buildANSI(spec.Foreground, spec.Background, spec.Attrs)
	if open == "" {
		return text
	}
	return open + text + Reset
}

// ─── Rich markup parser ───────────────────────────────────────────────────────

// markupResult holds a parsed markup section.
type markupResult struct {
	content string
	fg      string // may be ANSI raw if hex was used
	bg      string
	attrs   []string
	isFGRaw bool // true if fg is already a raw ANSI sequence
	isBGRaw bool
}

var (
	richPattern   = regexp.MustCompile(`\[([^\[\]]+?)\](.*?)\[/\]`)
	ansiStripRe   = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

const (
	escapedLeft  = "\x00ESCAPED_LEFT\x00"
	escapedRight = "\x00ESCAPED_RIGHT\x00"
)

// parseRichMarkup parses Rich-style markup from text and returns a slice of
// markup sections. Untagged portions get empty fg/bg.
func parseRichMarkup(text string) []markupResult {
	// Pre-process escaped brackets
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
		// match[0],match[1] = full match
		// match[2],match[3] = group 1 (markup)
		// match[4],match[5] = group 2 (content)
		markup := strings.ToLower(strings.TrimSpace(proc[match[2]:match[3]]))
		content := proc[match[4]:match[5]]

		// Prepend text before first match
		if i == 0 && match[0] > 0 {
			content = proc[:match[0]] + content
		}

		// Append text after [/] until next match or end
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

		// Restore escaped brackets in content
		content = strings.ReplaceAll(content, escapedLeft, "[")
		content = strings.ReplaceAll(content, escapedRight, "]")

		res := parseMarkupTag(markup)
		res.content = content
		results = append(results, res)
	}

	return results
}

// parseMarkupTag parses a markup tag string like "bold red on blue" or
// "#FF0000 on #00FF00 italic".
func parseMarkupTag(markup string) markupResult {
	var res markupResult
	parts := strings.Fields(markup)

	var colorParts []string
	for _, p := range parts {
		found := false
		for _, attr := range knownAttrs {
			if p == attr {
				a := p
				if a == "strike" {
					a = "strikethrough"
				}
				res.attrs = append(res.attrs, a)
				found = true
				break
			}
		}
		if !found {
			colorParts = append(colorParts, p)
		}
	}

	// colorParts may be: ["red"], ["red","on","blue"], ["#FF0000","on","#00FF00"]
	switch {
	case len(colorParts) >= 3 && colorParts[1] == "on":
		res.fg = resolveMarkupColor(colorParts[0], false, &res.isFGRaw)
		res.bg = resolveMarkupColor(colorParts[2], true, &res.isBGRaw)
	case len(colorParts) == 1:
		res.fg = resolveMarkupColor(colorParts[0], false, &res.isFGRaw)
	}

	return res
}

// resolveMarkupColor resolves a color token (name or hex) to an ANSI-ready
// string.  For hex colors it stores a raw ANSI sequence; for named colors it
// stores the name for normal lookup.
func resolveMarkupColor(token string, isBG bool, rawFlag *bool) string {
	if strings.HasPrefix(token, "#") || hexansi.IsHex(token) {
		result, err := hexansi.Convert(token, hexansi.ModeTrueColor)
		if err == nil {
			*rawFlag = true
			if isBG {
				return result.BG
			}
			return result.FG
		}
	}
	return token
}

// ─── MakeColors ────────────────────────────────────────────────────────────────

// Options configures a single MakeColors call.
type Options struct {
	Foreground string
	Background string
	Attrs      []string
	Force      bool // force color even when terminal doesn't support it
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Foreground: "white"}
}

// MakeColors applies color formatting to text.
//
// It supports:
//   - Plain text with foreground/background options
//   - Rich markup: "[bold red on yellow]text[/]"
//   - Combined format: "bold-red-yellow" as Foreground
//   - Hex colors inside markup tags
//
// The MAKE_COLORS=0 env var disables output; MAKE_COLORS_FORCE=1 forces it.
func MakeColors(text string, opts Options) string {
	if text == "" {
		return ""
	}

	// Rich markup detection
	if strings.Contains(text, "[") && strings.Contains(text, "[/]") {
		return applyRichMarkup(text, opts)
	}

	return applyPlain(text, opts)
}

// applyPlain applies color to plain (non-markup) text.
func applyPlain(text string, opts Options) string {
	spec := resolveOptions(opts)

	if debugEnabled() {
		fmt.Fprintf(os.Stderr, "[DEBUG] applyPlain fg=%q bg=%q attrs=%v\n",
			spec.Foreground, spec.Background, spec.Attrs)
	}

	if !shouldColor(opts.Force) {
		return text
	}
	return Colorize(text, spec)
}

// resolveOptions converts an Options into a ColorSpec via GetSort.
func resolveOptions(opts Options) ColorSpec {
	fg := opts.Foreground
	if fg == "" {
		fg = "white"
	}
	return GetSort("", fg, opts.Background, opts.Attrs)
}

// applyRichMarkup parses and applies rich markup formatting.
func applyRichMarkup(text string, opts Options) string {
	sections := parseRichMarkup(text)

	var sb strings.Builder
	for _, sec := range sections {
		if sec.content == "" {
			continue
		}

		fg := sec.fg
		if fg == "" {
			fg = opts.Foreground
			if fg == "" {
				fg = "white"
			}
		}
		bg := sec.bg
		if bg == "" {
			bg = opts.Background
		}

		var colored string
		if sec.isFGRaw || sec.isBGRaw {
			// Raw ANSI already computed (hex colors)
			colored = applyRawANSI(sec.content, sec.fg, sec.isFGRaw, sec.bg, sec.isBGRaw, sec.attrs)
		} else {
			spec := GetSort("", fg, bg, sec.attrs)
			colored = Colorize(sec.content, spec)
		}
		sb.WriteString(colored)
	}

	output := sb.String()
	if !shouldColor(opts.Force) {
		return ansiStripRe.ReplaceAllString(output, "")
	}
	return output
}

// applyRawANSI combines raw ANSI fg/bg sequences with attribute codes.
func applyRawANSI(text, fgRaw string, hasFG bool, bgRaw string, hasBG bool, attrs []string) string {
	var sb strings.Builder

	// Attribute codes
	for _, a := range attrs {
		if code, ok := attrCodes[a]; ok {
			// Inject into stream
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

// shouldColor decides whether to emit ANSI codes for this call.
func shouldColor(force bool) bool {
	if force || colorForced() {
		return true
	}
	if colorDisabled() {
		return false
	}
	return SupportsColor()
}

// ─── Convenience helpers ───────────────────────────────────────────────────────

// Make is the shortest alias for MakeColors.
func Make(text string, opts Options) string { return MakeColors(text, opts) }

// Sprint returns text colored with the given foreground (and optional background).
func Sprint(text, foreground, background string, attrs ...string) string {
	return MakeColors(text, Options{
		Foreground: foreground,
		Background: background,
		Attrs:      attrs,
	})
}

// Sprintf formats and colors text. The color spec is applied to the entire result.
func Sprintf(foreground, background string, attrs []string, format string, a ...interface{}) string {
	return Sprint(fmt.Sprintf(format, a...), foreground, background, attrs...)
}

// Fprintln writes a colored line to the given writer.
func Fprintln(w interface{ WriteString(string) (int, error) }, text, foreground, background string, attrs ...string) {
	s := Sprint(text, foreground, background, attrs...)
	w.WriteString(s + "\n") //nolint:errcheck
}

// Println prints colored text followed by a newline to stdout.
func Println(text, foreground, background string, attrs ...string) {
	fmt.Println(Sprint(text, foreground, background, attrs...))
}

// StripANSI removes all ANSI escape sequences from a string.
func StripANSI(s string) string {
	return ansiStripRe.ReplaceAllString(s, "")
}

// ─── Functional color helpers ─────────────────────────────────────────────────
// These mirror the dynamically-generated Python functions.

// Red returns text in red.
func Red(text string) string { return Sprint(text, "red", "") }

// Green returns text in green.
func Green(text string) string { return Sprint(text, "green", "") }

// Blue returns text in blue.
func Blue(text string) string { return Sprint(text, "blue", "") }

// Yellow returns text in yellow.
func Yellow(text string) string { return Sprint(text, "yellow", "") }

// Magenta returns text in magenta.
func Magenta(text string) string { return Sprint(text, "magenta", "") }

// Cyan returns text in cyan.
func Cyan(text string) string { return Sprint(text, "cyan", "") }

// White returns text in white.
func White(text string) string { return Sprint(text, "white", "") }

// Black returns text in black (visible on light backgrounds).
func Black(text string) string { return Sprint(text, "black", "") }

// LightRed returns text in light red.
func LightRed(text string) string { return Sprint(text, "lightred", "") }

// LightGreen returns text in light green.
func LightGreen(text string) string { return Sprint(text, "lightgreen", "") }

// LightBlue returns text in light blue.
func LightBlue(text string) string { return Sprint(text, "lightblue", "") }

// LightYellow returns text in light yellow.
func LightYellow(text string) string { return Sprint(text, "lightyellow", "") }

// LightMagenta returns text in light magenta.
func LightMagenta(text string) string { return Sprint(text, "lightmagenta", "") }

// LightCyan returns text in light cyan.
func LightCyan(text string) string { return Sprint(text, "lightcyan", "") }

// LightWhite returns text in light white.
func LightWhite(text string) string { return Sprint(text, "lightwhite", "") }

// Bold returns text with bold attribute.
func Bold(text string) string { return Sprint(text, "white", "", "bold") }

// Dim returns text with dim attribute.
func Dim(text string) string { return Sprint(text, "white", "", "dim") }

// Italic returns text with italic attribute.
func Italic(text string) string { return Sprint(text, "white", "", "italic") }

// Underline returns text with underline attribute.
func Underline(text string) string { return Sprint(text, "white", "", "underline") }

// Strikethrough returns text with strikethrough attribute.
func Strikethrough(text string) string { return Sprint(text, "white", "", "strikethrough") }
