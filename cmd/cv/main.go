package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/berquerant/cv"
)

func main() {
	var (
		delimiter = flag.String("d", ",", "delimiter")
	)
	flag.Usage = Usage
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fail("input format and output format are required")
	}

	var (
		src = cv.TypeFromString(args[0])
		dst = cv.TypeFromString(args[1])
	)

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("%v\n", err)
	}

	c := cv.New(src, dst, []rune(*delimiter)[0])
	result, err := c.Translate(input)
	if err != nil {
		fail("%v\n", err)
	}
	fmt.Printf("%s", result)
}

func fail(format string, v ...any) {
	fmt.Fprintf(os.Stderr, format, v...)
	os.Exit(1)
}

func Usage() {
	fmt.Fprintf(os.Stderr, "%s\n", usage)
	flag.PrintDefaults()
}

const usage = `cv - Translate data

Usage:
  cv INPUT_FORMAT OUTPUT_FORMAT [flags]

Read data from standard input and convert it to the specified format, then write it to standard output.

Valid values for each format are json, yaml, toml, and csv.
For csv format, the delimiter can be specified with the -d option. If it is not specified, it will be , (comma).
`
