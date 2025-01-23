package parser

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 0     1            2             3          4          5        6               7                8
// Data; Nome Cartão; Final Cartão; Categoria; Descrição; Parcela; Valor (em US$); Cotação (em R$); Valor (em R$)

const (
	tesseractBin = "/usr/bin/tesseract"
)

var simpleTransactionRegex = regexp.MustCompile(`(\d{2}\/\d{2})\s+(.+)\s*R\$\s*([0-9.]+\,\d+)\s+.+(\d{4})`)
var simpleProcessingTransactionRegex = regexp.MustCompile(`(?i)(\d{2}\/\d{2})\s+(.+)\s*Cartao final\s*(\d+)\s*Em processamento\s*R\$\s*([0-9.]+\,\d+)`)
var installmentTransactionRegex = regexp.MustCompile(`(\d{2}\/\d{2})\s+(.+)\s*.+(\d{4})\s*R\$\s*([0-9.]+\,\d+)\s+\w+\s+(\d+)\s+\w+\s+(\d+)`)

// var currentDate time.Time
var currentDate time.Time

func runOCR(file io.Reader) (string, error) {
	var ocrOutput strings.Builder

	cmd := exec.Command(tesseractBin, "stdin", "stdout", "--psm", "4")
	cmd.Stdin = file
	cmd.Stdout = &ocrOutput

	if err := cmd.Run(); err != nil {
		return "", err
	}

	currentDate = time.Now()

	return ocrOutput.String(), nil
}

func GetTransactions(t string) []line {
	var transactions []line

	simpleMatches := simpleTransactionRegex.FindAllStringSubmatch(t, -1)

	for _, m := range simpleMatches {
		transactions = append(transactions, line{fixYear(m[1]), m[2], "", m[3]})
	}

	pMatches := simpleProcessingTransactionRegex.FindAllStringSubmatch(t, -1)
	for _, m := range pMatches {
		transactions = append(transactions, line{fixYear(m[1]), m[2], "", m[4]})
	}

	return transactions
}

func GetInstallmentTransactions(t string, invoiceRef string) ([]line, error) {
	var transactions []line

	installmentMatches := installmentTransactionRegex.FindAllStringSubmatch(t, -1)

	for _, m := range installmentMatches {

		m[1] = fixYear(m[1])

		if err := handleInstallments([]string{m[1], "", "", "", m[2], fmt.Sprintf("%s/%s", m[5], m[6]), "", "", m[4]}, &transactions, invoiceRef); err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func scanImageRows(file io.Reader, invoiceRef string) ([]line, error) {

	text, err := runOCR(file)
	if err != nil {
		return nil, err
	}

	installments, err := GetInstallmentTransactions(text, invoiceRef)
	if err != nil {
		return nil, err
	}

	return append(installments, GetTransactions(text)...), nil
}

func fixYear(date string) string {

	d := strings.Split(date, "/")

	transactionMonth, err := strconv.Atoi(d[1])
	if err != nil {
		// ... handle error
		return date
	}

	if int(currentDate.Month()) < transactionMonth {
		return fmt.Sprintf("%s/%s", date, currentDate.AddDate(-1, 0, 0).Format("2006"))

	} else {
		return fmt.Sprintf("%s/%s", date, currentDate.Format("2006"))
	}
}
