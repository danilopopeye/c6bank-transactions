package parser

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"git.home/c6bank-transactions/internal/qif"
)

var ErrWrongCSVFilename = errors.New("the filename should be in the form Fatura_YYYY-MM-DD.csv")

// Line is: date, payee, memo, value
type Line [4]string

func Parse(name string, file multipart.File, size int64, password string) (io.Reader, string, error) {
	var (
		err   error
		qtype qif.QIFType
		lines []Line
	)

	outputname := strings.TrimSuffix(name, filepath.Ext(name))

	switch strings.ToLower(filepath.Ext(name)) {
	case ".pdf":
		qtype = qif.BankType

		lines, err = scanPDFRows(file, password, size)
		if err != nil {
			return nil, "", err
		}
	case ".csv":
		qtype = qif.CreditCardType

		if len(name) != 21 || !strings.HasPrefix(name, "Fatura_") {
			return nil, "", ErrWrongCSVFilename
		}

		reference, err := time.Parse(time.DateOnly, name[7:17])
		if err != nil {
			fmt.Printf("ERROR error parsing CSV filename %q: %v", name, err)

			return nil, "", ErrWrongCSVFilename
		}

		lines, err = scanCSVRows(reference, file)
		if err != nil {
			return nil, "", err
		}
	case ".jpg", ".png":
		return parseImages(outputname, file)
	default:
		return nil, "", fmt.Errorf("invalid file %s", name)
	}

	output, err := qif.Parse(qtype, linesToTransactions(lines))

	return output, outputname + ".qif", err
}

func linesToTransactions(lines []Line) []qif.Transaction {
	qt := make([]qif.Transaction, 0, len(lines))

	for _, l := range lines {
		qt = append(qt, qif.Transaction{
			Date:   l[0],
			Payee:  l[1],
			Memo:   l[2],
			Amount: l[3],
		})
	}

	return qt
}

func parseImages(name string, file io.ReadSeeker) (io.Reader, string, error) {
	transactions, err := ScanImage(file)
	if err != nil {
		return nil, "", err
	}

	output, err := TransactionsToCSV(transactions)

	return output, name + "-parsed.csv", err
}
