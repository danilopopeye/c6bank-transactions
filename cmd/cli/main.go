package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"git.home/c6bank-transactions/internal/parser"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)
	output := fs.String("o", "", "output CSV file (defaults to stdout)")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: %s [flags] <file1> [file2 ...]\n", "cli")
		fmt.Fprintln(stderr, "Parse C6 Bank transaction files into a single CSV.")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Supported formats: CSV, PNG, JPG/JPEG")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return 1
	}

	var all []parser.Transaction

	for _, path := range fs.Args() {
		transactions, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		all = append(all, transactions...)
	}

	all = parser.Deduplicate(all)

	slices.SortFunc(all, func(a, b parser.Transaction) int {
		if a.Date.Before(b.Date) {
			return -1
		}
		if a.Date.After(b.Date) {
			return 1
		}
		return strings.Compare(a.Payee, b.Payee)
	})

	r, err := parser.TransactionsToCSV(all)
	if err != nil {
		fmt.Fprintf(stderr, "error generating CSV: %v\n", err)
		return 1
	}

	var w io.Writer = stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(stderr, "error creating output file: %v\n", err)
			return 1
		}
		defer f.Close()
		w = f
	}

	if _, err := io.Copy(w, r); err != nil {
		fmt.Fprintf(stderr, "error writing output: %v\n", err)
		return 1
	}

	return 0
}
