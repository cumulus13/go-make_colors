package makecolors

import (
	"fmt"
	"io"
	"os"
)

// Console provides a stateful, writer-bound interface for printing colored text.
// The zero value writes to os.Stdout with color auto-detected.
//
// Example:
//
//	c := makecolors.NewConsole(os.Stdout)
//	c.Print("hello", makecolors.Options{Foreground: "green"})
type Console struct {
	w     io.Writer
	Force bool // force color even when terminal doesn't support it
}

// NewConsole returns a Console that writes to w.
func NewConsole(w io.Writer) *Console {
	if w == nil {
		w = os.Stdout
	}
	return &Console{w: w}
}

// Stdout is a package-level Console bound to os.Stdout.
var Stdout = NewConsole(os.Stdout)

// Print writes colored text + newline to the console's writer.
func (c *Console) Print(text string, opts Options) {
	opts.Force = opts.Force || c.Force
	fmt.Fprintln(c.w, MakeColors(text, opts))
}

// Printf formats and prints colored text to the console's writer.
func (c *Console) Printf(opts Options, format string, args ...interface{}) {
	opts.Force = opts.Force || c.Force
	fmt.Fprintln(c.w, MakeColors(fmt.Sprintf(format, args...), opts))
}

// Error prints text in red (convenience for error messages).
func (c *Console) Error(text string) {
	fmt.Fprintln(c.w, MakeColors(text, Options{Foreground: "lightred", Force: c.Force}))
}

// Warn prints text in yellow (convenience for warnings).
func (c *Console) Warn(text string) {
	fmt.Fprintln(c.w, MakeColors(text, Options{Foreground: "yellow", Force: c.Force}))
}

// Info prints text in cyan (convenience for informational messages).
func (c *Console) Info(text string) {
	fmt.Fprintln(c.w, MakeColors(text, Options{Foreground: "cyan", Force: c.Force}))
}

// Success prints text in green (convenience for success messages).
func (c *Console) Success(text string) {
	fmt.Fprintln(c.w, MakeColors(text, Options{Foreground: "lightgreen", Force: c.Force}))
}

// Debug prints text in dim white (convenience for debug output).
func (c *Console) Debug(text string) {
	fmt.Fprintln(c.w, MakeColors(text, Options{Foreground: "lightblack", Force: c.Force}))
}

// Status prints a labeled status line, e.g. "[ERROR] message".
//
//	c.Status("ERROR", "Database unavailable", "white", "red")
func (c *Console) Status(label, message, labelFG, labelBG string) {
	tag := MakeColors("["+label+"]", Options{Foreground: labelFG, Background: labelBG, Force: c.Force})
	fmt.Fprintln(c.w, tag+" "+message)
}

// Rich prints text using Rich-style markup.
//
//	c.Rich("[bold red]Error:[/] something went wrong")
func (c *Console) Rich(text string) {
	fmt.Fprintln(c.w, MakeColors(text, Options{Force: c.Force}))
}
