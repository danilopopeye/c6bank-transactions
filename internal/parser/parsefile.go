package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	case ".pdf":
		stat, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("stat file %s: %w", path, err)
		}

		lines, err := scanPDFRows(f, "", stat.Size())
		if err != nil {
			return nil, fmt.Errorf("parse PDF %s: %w", path, err)
		}

		return linesToTypedTransactions(lines), nil

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
