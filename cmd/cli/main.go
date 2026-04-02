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
	output := flag.String("o", "", "output CSV file (defaults to stdout)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <file1> [file2 ...]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Parse C6 Bank transaction files into a single CSV.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Supported formats: CSV, PNG, JPG/JPEG")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var all []parser.Transaction

	for _, path := range flag.Args() {
		transactions, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "error generating CSV: %v\n", err)
		os.Exit(1)
	}

	var w io.Writer = os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	}

	if _, err := io.Copy(w, r); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}
}
