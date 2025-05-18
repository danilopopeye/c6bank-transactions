package parser

import (
	"fmt"
	"path/filepath"
	"strings"
)

var ErrInvalidFormat = fmt.Errorf("format not allowed, only: PDF, JPEG/PNG or CSV")

func IsValid(fname, ftype string) error {
	if !isPDF(fname, ftype) && !isCSV(fname, ftype) && !isImage(fname, ftype) {
		return ErrInvalidFormat
	}

	return nil
}

func isCSV(fname, ftype string) bool {
	return strings.HasPrefix(ftype, "text/plain") && filepath.Ext(fname) == ".csv"
}

func isPDF(fname, ftype string) bool {
	return ftype == "application/pdf" && filepath.Ext(fname) == ".pdf"
}

func isImage(fname, ftype string) bool {
	switch ext := strings.ToLower(filepath.Ext(fname)); ext {
	case ".jpg", ".jpeg":
		return strings.HasPrefix(ftype, "image/jpeg")
	case ".png":
		return strings.HasPrefix(ftype, "image/png")
	default:
		return false
	}
}
