package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"aisim/internal/model"
)

// REPL holds the mutable state for an interactive session.
type REPL struct {
	params model.Params
	in     *bufio.Scanner
	out    io.Writer
}

// New creates a REPL with the given initial parameters.
func New(p model.Params) *REPL {
	return &REPL{
		params: p,
		in:     bufio.NewScanner(os.Stdin),
		out:    os.Stdout,
	}
}

// NewWithIO creates a REPL with custom I/O (useful for testing).
func NewWithIO(p model.Params, in io.Reader, out io.Writer) *REPL {
	return &REPL{
		params: p,
		in:     bufio.NewScanner(in),
		out:    out,
	}
}

// Run starts the REPL loop.
func (r *REPL) Run() {
	fmt.Fprintln(r.out, "AI Layoff Trap Simulator  (type 'help' for commands, 'quit' to exit)")
	fmt.Fprintf(r.out, "Loaded: paper base parameters (N=%.4g, w=%.2f, c=%.2f, k=%.2f)\n\n",
		r.params.N, r.params.W, r.params.C, r.params.K)

	for {
		fmt.Fprint(r.out, "aisim> ")
		if !r.in.Scan() {
			break
		}
		line := strings.TrimSpace(r.in.Text())
		if !r.dispatch(line) {
			break
		}
	}
}
