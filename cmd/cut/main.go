package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	icut "cut/internal/cut"
)

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("cut", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var fieldsFlag string
	var delimiterFlag string
	var separatedOnly bool

	fs.StringVar(&fieldsFlag, "f", "", "fields to select (e.g., '1,3-5,7-')")
	fs.StringVar(&delimiterFlag, "d", "\t", "field delimiter (single character or string)")
	fs.BoolVar(&separatedOnly, "s", false, "only print lines with the delimiter present")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if fieldsFlag == "" {
		_, _ = fmt.Fprintln(stderr, "missing -f fields specification")
		fs.Usage()
		return 2
	}

	selector, err := icut.ParseFields(fieldsFlag)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "invalid -f:", err)
		return 2
	}

	proc := icut.Processor{
		Delimiter:      []byte(delimiterFlag),
		SeparatedOnly:  separatedOnly,
		FieldSelection: selector,
	}

	if err := proc.Process(stdin, stdout); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
