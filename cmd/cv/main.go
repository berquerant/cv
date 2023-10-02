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
		src       = flag.String("i", "auto", "source format")
		dst       = flag.String("o", "json", "destination format")
	)
	flag.Usage = Usage
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	c := cv.New(cv.TypeFromString(*src), cv.TypeFromString(*dst), []rune(*delimiter)[0])
	result, err := c.Translate(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s", result)
}

func Usage() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

const usage = `cv - Translate data

Usage:
  cv [flags]

Read data from standard input and convert it to the specified format, then write it to standard output.

The input format is specified with the -i option. If it is not specified, it will be automatically determined.
The output format is specified with the -o option. If it is not specified, it will be json.
Valid values for each format are json, yaml, toml, and csv.
For csv format, the delimiter can be specified with the -d option. If it is not specified, it will be , (comma).
`
