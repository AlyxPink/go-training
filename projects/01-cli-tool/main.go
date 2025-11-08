package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/alyxpink/go-training/jq/formatter"
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

	// TODO: Parse the query string into a Query object
	// Hint: q, err := query.Parse(queryStr)
	q, err := query.Parse(queryStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query error: %v\n", err)
		os.Exit(1)
	}

	// TODO: Process each input (files or stdin)
	// Hint: Call processInput for each file, or stdin if no files
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
	// TODO: Decode JSON from reader
	var data interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// TODO: Execute query on data
	result, err := q.Execute(data)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	// TODO: Format and output result
	// Hint: Use formatter package based on flags
	return outputResult(result)
}

func outputResult(data interface{}) error {
	var fmt formatter.Formatter

	// TODO: Select formatter based on flags
	if *table {
		fmt = &formatter.TableFormatter{}
	} else if *compact {
		fmt = &formatter.JSONFormatter{Compact: true}
	} else if *raw {
		fmt = &formatter.RawFormatter{}
	} else {
		fmt = &formatter.JSONFormatter{Compact: false}
	}

	output, err := fmt.Format(data)
	if err != nil {
		return err
	}

	fmt.Print(output)
	return nil
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
