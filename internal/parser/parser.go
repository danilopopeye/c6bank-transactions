package parser

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"git.home/c6bank-transactions/internal/qif"
)

// Line is: date, payee, memo, value
type Line [4]string

func Parse(name string, file multipart.File, size int64, password string, invoiceRef string, installmentH string) (io.Reader, error) {
	var (
		err   error
		qtype qif.QIFType
		lines []Line
	)

	switch strings.ToLower(filepath.Ext(name)) {
	case ".pdf":
		qtype = qif.BankType
		lines, err = scanPDFRows(file, password, size)
	case ".csv":
		qtype = qif.CreditCardType
		lines, err = scanCSVRows(file, invoiceRef, installmentH)
	case ".jpg":
		qtype = qif.CreditCardType
		lines, err = scanImageRows(file, invoiceRef, installmentH)
	default:
		return nil, fmt.Errorf("invalid file %s", name)
	}

	if err != nil {
		return nil, err
	}

	return qif.Parse(qtype, linesToTransactions(lines))
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
