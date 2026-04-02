package parser_test

import (
	"testing"
	"time"

	"git.home/c6bank-transactions/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		wantErr     bool
		errContains string
		wantLen     int
	}{
		{
			name:    "valid CSV",
			path:    "testdata/Fatura_2026-01-15.csv",
			wantErr: false,
			wantLen: 4, // 1 unique + 3 installments (1/3)
		},
		{
			name:        "CSV with wrong filename",
			path:        "testdata/wrong-name.csv",
			wantErr:     true,
			errContains: "Fatura_",
		},
		{
			name:        "unsupported format",
			path:        "testdata/dummy.xlsx",
			wantErr:     true,
			errContains: "unsupported file format",
		},
		{
			name:        "file not found",
			path:        "testdata/nonexistent.csv",
			wantErr:     true,
			errContains: "open file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transactions, err := parser.ParseFile(tt.path)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			assert.Len(t, transactions, tt.wantLen)
		})
	}
}

func TestParseFile_CSVContent(t *testing.T) {
	t.Parallel()

	transactions, err := parser.ParseFile("testdata/Fatura_2026-01-15.csv")
	require.NoError(t, err)
	require.Len(t, transactions, 4)

	assert.Equal(t, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), transactions[0].Date)
	assert.Equal(t, "MERCADO EXTRA", transactions[0].Payee)
	// Amazon BR 1/3 generates 3 installments: Jan, Feb, Mar
	assert.Equal(t, time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC), transactions[1].Date)
	assert.Equal(t, "AMAZON BR", transactions[1].Payee)
	assert.Equal(t, time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC), transactions[2].Date)
	assert.Equal(t, time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC), transactions[3].Date)
}

func TestDeduplicate(t *testing.T) {
	t.Parallel()

	date := time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name         string
		transactions []parser.Transaction
		wantLen      int
	}{
		{
			name:         "empty input",
			transactions: nil,
			wantLen:      0,
		},
		{
			name: "no duplicates",
			transactions: []parser.Transaction{
				{Date: date, Payee: "A", Memo: "m1", Amount: "10.00"},
				{Date: date, Payee: "B", Memo: "m2", Amount: "20.00"},
			},
			wantLen: 2,
		},
		{
			name: "removes exact duplicates",
			transactions: []parser.Transaction{
				{Date: date, Payee: "A", Memo: "m1", Amount: "10.00"},
				{Date: date, Payee: "A", Memo: "m1", Amount: "10.00"},
			},
			wantLen: 1,
		},
		{
			name: "same date payee amount but different memo",
			transactions: []parser.Transaction{
				{Date: date, Payee: "A", Memo: "m1", Amount: "10.00"},
				{Date: date, Payee: "A", Memo: "m2", Amount: "10.00"},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := parser.Deduplicate(tt.transactions)
			assert.Len(t, result, tt.wantLen)
		})
	}
}
