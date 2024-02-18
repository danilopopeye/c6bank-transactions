package parser

import (
	"fmt"
	"path/filepath"
	"strings"
)

func IsValid(fname, ftype string) error {
	if !isPDF(fname, ftype) && !isCSV(fname, ftype) {
		return fmt.Errorf("format not allowed, only a PDF or CSV")
	}
	return nil
}

func isCSV(fname, ftype string) bool {
	return strings.HasPrefix(ftype, "text/plain") && filepath.Ext(fname) == ".csv"
}

func isPDF(fname, ftype string) bool {
	return ftype == "application/pdf" && filepath.Ext(fname) == ".pdf"
}
