package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/segmentio/fasthash/fnv1a"
)

// ParseFile opens the file at path, detects its format by extension,
// and returns the parsed transactions.
func ParseFile(path string) ([]Transaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", path, err)
	}
	defer f.Close()

	name := filepath.Base(path)
	ext := strings.ToLower(filepath.Ext(name))

	switch ext {
	case ".csv":
		if len(name) != 21 || !strings.HasPrefix(name, "Fatura_") {
			return nil, fmt.Errorf("%w: %s", ErrWrongCSVFilename, name)
		}

		reference, err := time.Parse(time.DateOnly, name[7:17])
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrWrongCSVFilename, name)
		}

		lines, err := scanCSVRows(reference, f)
		if err != nil {
			return nil, fmt.Errorf("parse CSV %s: %w", path, err)
		}

		return linesToTypedTransactions(lines), nil

	case ".jpg", ".png":
		transactions, err := ScanImage(f, false)
		if err != nil {
			return nil, fmt.Errorf("parse image %s: %w", path, err)
		}

		return transactions, nil

	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// Deduplicate removes duplicate transactions based on Date+Payee+Amount+Memo.
func Deduplicate(transactions []Transaction) []Transaction {
	seen := make(map[uint64]struct{}, len(transactions))
	result := make([]Transaction, 0, len(transactions))

	for _, t := range transactions {
		h := fnv1a.Init64
		h = fnv1a.AddString64(h, t.Date.Format(dateFormat))
		h = fnv1a.AddString64(h, t.Payee)
		h = fnv1a.AddString64(h, t.Amount)
		h = fnv1a.AddString64(h, t.Memo)

		if _, exists := seen[h]; exists {
			continue
		}

		seen[h] = struct{}{}
		result = append(result, t)
	}

	return result
}

// linesToTypedTransactions converts []Line (string dates) to []Transaction (typed dates).
func linesToTypedTransactions(lines []Line) []Transaction {
	transactions := make([]Transaction, 0, len(lines))

	for _, l := range lines {
		date, err := time.Parse(dateFormat, l[0])
		if err != nil || date.IsZero() {
			continue
		}

		transactions = append(transactions, Transaction{
			Date:   date,
			Payee:  l[1],
			Memo:   l[2],
			Amount: l[3],
		})
	}

	return transactions
}
