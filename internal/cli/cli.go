// Package cli handles command-line parsing and exit statuses.
package cli

import (
	"flag"
	"fmt"
	"github.com/rumpl/format-structs/internal/formatter"
	"io"
)

// Run executes the command with explicit arguments and output destinations.
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("format-structs", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", ".", "project directory")
	check := flags.Bool("check", false, "check formatting without writing")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: format-structs [flags] [file.go | directory | directory/...]...")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if err := formatter.Run(*root, *check, flags.Args(), stdout); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
