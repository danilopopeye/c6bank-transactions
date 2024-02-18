package parser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	csvHeaders        = []string{"Date", "Payee", "Memo", "Outflow", "Inflow"}
	transactionRegexp = regexp.MustCompile(`(?P<date>[0-9/]{10}) (?P<payee>[A-Z0-9., ]+)-?(?P<memo>.*)\s+(?P<doc>[0-9]{12})\s+(?P<value>[0-9.]+,[0-9]{2})\s+(?P<type>[CD])`)
)

func Parse(name string, file io.Reader, password string) (io.Reader, error) {
	buff := &bytes.Buffer{}

	csvFile := csv.NewWriter(buff)
	csvFile.Comma = ';'

	if err := csvFile.Write(csvHeaders); err != nil {
		return nil, err
	}

	var err error
	switch filepath.Ext(name) {
	case ".pdf":
		err = scanPDFRows(file, password, csvFile)
	case ".csv":
		err = scanCSVRows(file, password, csvFile)
	default:
		return nil, fmt.Errorf("invalid file %s", name)
	}

	if err != nil {
		return nil, err
	}

	csvFile.Flush()

	if err := csvFile.Error(); err != nil {
		return nil, err
	}

	return buff, nil
}

func parseAmount(value, installment string) (decimal.Decimal, string) {
	value = strings.ReplaceAll(value, ".", "")
	value = strings.ReplaceAll(value, ",", ".")
	amount := decimal.RequireFromString(value)

	if installment == "Única" || installment == "" {
		return amount, ""
	}

	parts := strings.SplitN(installment, "/", 2)

	if parts[0] == "1" {
		return amount.Mul(decimal.RequireFromString(parts[1])), parts[1]
	}

	return decimal.Zero, ""
}
