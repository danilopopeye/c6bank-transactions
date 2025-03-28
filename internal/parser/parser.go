package parser

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"git.home/c6bank-transactions/internal/qif"
)

// Line is: date, payee, memo, value
type Line [4]string

func Parse(name string, file multipart.File, size int64, password string, invoiceRef string) (io.Reader, string, error) {
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
	case ".csv":
		qtype = qif.CreditCardType
		var reference time.Time

		if strings.HasPrefix(name, "Fatura_") {
			reference, err = time.Parse(time.DateOnly, name[7:17])
			if err != nil {
				return nil, "", fmt.Errorf("the filename should be in the form Fatura_YYYY-MM-DD.csv")
			}
		}

		lines, err = scanCSVRows(reference, file)
	case ".jpg", ".png":
		return parseImages(outputname, file, invoiceRef)
	default:
		return nil, "", fmt.Errorf("invalid file %s", name)
	}

	if err != nil {
		return nil, "", err
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

func parseImages(name string, file io.ReadSeeker, invoiceRef string) (io.Reader, string, error) {
	transactions, err := ScanImage(file, invoiceRef)
	if err != nil {
		return nil, "", err
	}

	output, err := TransactionsToCSV(transactions)

	return output, name + "-parsed.csv", err
}
