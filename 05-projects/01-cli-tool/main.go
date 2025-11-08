package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/alyxpink/go-training/jq/query"
)

var (
	version = "1.0.0"
	compact = flag.Bool("c", false, "compact output")
	table   = flag.Bool("t", false, "table output")
	raw     = flag.Bool("r", false, "raw output (no quotes)")
	showVer = flag.Bool("v", false, "show version")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if *showVer {
		fmt.Printf("jq version %s\n", version)
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	queryStr := args[0]
	files := args[1:]

	q, err := query.Parse(queryStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query error: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		if err := processInput(os.Stdin, q, "stdin"); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	} else {
		for _, filename := range files {
			f, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening %s: %v\n", filename, err)
				os.Exit(1)
			}
			if err := processInput(f, q, filename); err != nil {
				f.Close()
				fmt.Fprintf(os.Stderr, "error processing %s: %v\n", filename, err)
				os.Exit(1)
			}
			f.Close()
		}
	}
}

func processInput(r io.Reader, q *query.Query, filename string) error {
	panic("not implemented")
}

func outputResult(data interface{}) error {
	panic("not implemented")
}

func usage() {
	fmt.Fprintf(os.Stderr, `jq - JSON query tool

Usage:
  jq [flags] query [files...]

Flags:
  -c    compact output
  -t    table output
  -r    raw output (no quotes)
  -v    show version
  -h    show help

Examples:
  jq '.name' data.json
  echo '{"name": "Alice"}' | jq '.name'
  jq -t '.users[]' data.json
`)
}
