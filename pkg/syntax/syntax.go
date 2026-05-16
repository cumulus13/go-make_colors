// Package syntax provides terminal syntax highlighting via the Chroma library,
// mirroring the Python make_colors Syntax class.
//
// Author: Hadi Cahyadi <cumulus13@gmail.com>
// License: MIT
package syntax

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// LineNumberMode controls how line numbers are displayed.
type LineNumberMode int

const (
	LineNumberNone     LineNumberMode = iota // no line numbers
	LineNumberAbsolute                        // 1, 2, 3 …
	LineNumberRelative                        // relative to the first visible line
)

// Options configures a Syntax highlighting operation.
type Options struct {
	Lexer          string         // language name or "auto"
	Theme          string         // Chroma style name (default: "monokai")
	LineNumbers    LineNumberMode
	StartLine      int   // first line number when LineNumberAbsolute (default: 1)
	TabSize        int   // spaces per tab (default: 4)
	CodeWidth      int   // max column width for wrapping (0 = disabled)
	WordWrap       bool
	HighlightLines []int // 1-based line numbers to highlight
	TrueColor      bool  // use truecolor formatter (vs 256-color)
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Lexer:     "auto",
		Theme:     "monokai",
		StartLine: 1,
		TabSize:   4,
	}
}

// Syntax is a code highlighter.
type Syntax struct {
	code string
	opts Options
}

// New creates a Syntax instance for the given code.
func New(code string, opts Options) *Syntax {
	if opts.TabSize <= 0 {
		opts.TabSize = 4
	}
	if opts.Theme == "" {
		opts.Theme = "monokai"
	}
	if opts.StartLine <= 0 {
		opts.StartLine = 1
	}
	if opts.TabSize != 4 {
		code = strings.ReplaceAll(code, "\t", strings.Repeat(" ", opts.TabSize))
	}
	return &Syntax{code: code, opts: opts}
}

// Highlight returns the highlighted code as a string.
func (s *Syntax) Highlight() (string, error) {
	lexer := resolveLexer(s.opts.Lexer, s.code)
	style := resolveStyle(s.opts.Theme)
	formatter := resolveFormatter(s.opts.TrueColor)

	// Chroma v2: lexer.Tokenise() returns an Iterator, not []Token.
	iter, err := lexer.Tokenise(nil, s.code)
	if err != nil {
		return "", fmt.Errorf("syntax: tokenise: %w", err)
	}

	var sb strings.Builder
	if err := formatter.Format(&sb, style, iter); err != nil {
		return "", fmt.Errorf("syntax: format: %w", err)
	}

	highlighted := sb.String()

	if s.opts.LineNumbers != LineNumberNone {
		highlighted = addLineNumbers(highlighted, s.opts)
	}

	return highlighted, nil
}

// Print writes the highlighted code to w (or os.Stdout if w is nil).
func (s *Syntax) Print(w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	out, err := s.Highlight()
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, out)
	return err
}

// PrintStdout is a convenience wrapper around Print(os.Stdout).
func (s *Syntax) PrintStdout() error { return s.Print(os.Stdout) }

// String implements fmt.Stringer. Errors are swallowed; use Highlight for
// production code that needs error handling.
func (s *Syntax) String() string {
	out, _ := s.Highlight()
	return out
}

// ─── Package-level helpers ────────────────────────────────────────────────────

// Print is a convenience function that highlights code and writes it to stdout.
func Print(code string, opts Options) error {
	return New(code, opts).Print(os.Stdout)
}

// Sprint highlights code and returns it as a string.
func Sprint(code string, opts Options) (string, error) {
	return New(code, opts).Highlight()
}

// AvailableThemes returns a sorted list of all Chroma style names.
func AvailableThemes() []string {
	return styles.Names()
}

// AvailableLexers returns a list of all Chroma lexer names.
func AvailableLexers() []string {
	var names []string
	for _, l := range lexers.GlobalLexerRegistry.Lexers {
		names = append(names, l.Config().Name)
	}
	return names
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

func resolveLexer(name, code string) chroma.Lexer {
	if name == "" || strings.EqualFold(name, "auto") {
		l := lexers.Analyse(code)
		if l != nil {
			return chroma.Coalesce(l)
		}
		return chroma.Coalesce(lexers.Fallback)
	}
	l := lexers.Get(name)
	if l == nil {
		return chroma.Coalesce(lexers.Fallback)
	}
	return chroma.Coalesce(l)
}

func resolveStyle(name string) *chroma.Style {
	s := styles.Get(name)
	if s == nil {
		return styles.Fallback
	}
	return s
}

func resolveFormatter(trueColor bool) chroma.Formatter {
	name := "terminal256"
	if trueColor {
		name = "terminal16m"
	}
	f := formatters.Get(name)
	if f == nil {
		return formatters.Fallback
	}
	return f
}

func addLineNumbers(highlighted string, opts Options) string {
	lines := strings.Split(highlighted, "\n")
	// Chroma typically appends a trailing newline; trim it for clean counting.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	digits := len(fmt.Sprintf("%d", opts.StartLine+len(lines)-1))
	fmtStr := fmt.Sprintf("%%%dd: ", digits)

	var sb strings.Builder
	for i, line := range lines {
		lineNum := opts.StartLine + i
		prefix := fmt.Sprintf(fmtStr, lineNum)

		if contains(opts.HighlightLines, lineNum) {
			prefix = "\x1b[1;41m" + prefix + "\x1b[0m"
		}

		sb.WriteString(prefix)
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func contains(sl []int, v int) bool {
	for _, x := range sl {
		if x == v {
			return true
		}
	}
	return false
}
