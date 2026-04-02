package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testdata = func() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "..", "..", "internal", "parser", "testdata")
}()

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantErr    string
		wantOutHas string
		wantCount  int    // expected number of data rows (excluding header)
		wantExact  string // must appear exactly once in output
	}{
		{
			name:     "no args shows usage",
			args:     nil,
			wantCode: 1,
			wantErr:  "Usage:",
		},
		{
			name:     "file not found",
			args:     []string{"nonexistent.csv"},
			wantCode: 1,
			wantErr:  "open file",
		},
		{
			name:     "unsupported format",
			args:     []string{filepath.Join(testdata, "dummy.xlsx")},
			wantCode: 1,
			wantErr:  "unsupported file format",
		},
		{
			name:       "single CSV file",
			args:       []string{filepath.Join(testdata, "Fatura_2026-01-15.csv")},
			wantCode:   0,
			wantOutHas: "MERCADO EXTRA",
		},
		{
			name:     "multiple CSV files deduplicates",
			args:     []string{
				filepath.Join(testdata, "Fatura_2026-01-15.csv"),
				filepath.Join(testdata, "Fatura_2026-01-15.csv"),
			},
			wantCode:   0,
			wantOutHas: "MERCADO EXTRA",
			wantCount:  4, // same file twice, deduplicated back to 4 transactions
			wantExact:  "MERCADO EXTRA", // must appear exactly once
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			code := run(tt.args, &stdout, &stderr)

			assert.Equal(t, tt.wantCode, code)

			if tt.wantErr != "" {
				assert.Contains(t, stderr.String(), tt.wantErr)
			}

			if tt.wantOutHas != "" {
				assert.Contains(t, stdout.String(), tt.wantOutHas)
			}

			if tt.wantCount > 0 {
				lines := strings.Count(stdout.String(), "\n") - 1 // exclude header
				assert.Equal(t, tt.wantCount, lines, "expected %d data rows", tt.wantCount)
			}

			if tt.wantExact != "" {
				assert.Equal(t, 1, strings.Count(stdout.String(), tt.wantExact),
					"%q should appear exactly once (deduplication check)", tt.wantExact)
			}
		})
	}
}

func TestRun_OutputToFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "output.csv")

	var stdout, stderr bytes.Buffer
	code := run([]string{"-o", path, filepath.Join(testdata, "Fatura_2026-01-15.csv")}, &stdout, &stderr)

	require.Equal(t, 0, code, "stderr: %s", stderr.String())

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "MERCADO EXTRA")
}
