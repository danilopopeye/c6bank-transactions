package parser_test

import (
	"testing"

	"git.home/c6bank-transactions/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fname   string
		ftype   string
		wantErr error
	}{
		{"valid PDF", "example.pdf", "application/pdf", nil},
		{"valid CSV", "example.csv", "text/plain", nil},
		{"valid JPG", "example.jpg", "image/jpeg", nil},
		{"valid JPEG", "example.jpeg", "image/jpeg", nil},
		{"valid PNG", "example.png", "image/png", nil},
		{"invalid file type", "example.txt", "text/plain", parser.ErrInvalidFormat},
		{"invalid file extension", "example.pdfx", "application/pdf", parser.ErrInvalidFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.IsValid(tt.fname, tt.ftype)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
