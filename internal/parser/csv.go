package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// 0     1            2             3          4          5        6               7                8
// Data; Nome Cartão; Final Cartão; Categoria; Descrição; Parcela; Valor (em US$); Cotação (em R$); Valor (em R$)

const (
	dateFormat = "02/01/2006"
	minus      = '-'
	fee        = "ANUIDADE DIFERENCIADA"
)

func scanCSVRows(file io.Reader, invoiceRef string, installmentH string) ([]Line, error) {
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	_, err := csvReader.Read() // header
	if err != nil {
		return nil, err
	}

	var lines []Line

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if record[8][0] == minus {
			record[8] = record[8][1:]
		} else {
			record[8] = fmt.Sprintf("-%s", record[8])
		}

		if record[5] == "Única" {
			record[5] = ""
		} else if strings.ToUpper(record[4]) == fee {
			// INFO do not handle installments
		} else {
			if err := handleInstallments(record, &lines, invoiceRef, installmentH); err != nil {
				return nil, err
			}

			continue
		}

		addInvoiceReference(&record[5], invoiceRef)

		lines = append(lines, Line{record[0], record[4], record[5], record[8]})
	}

	return lines, nil
}

func handleInstallments(record []string, lines *[]Line, invoiceRef string, installmentH string) error {
	purchase, payee, installment, value := record[0], record[4], record[5], record[8]

	parts := strings.SplitN(installment, "/", 2)

	current, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	total, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	date, err := time.Parse(dateFormat, purchase)
	if err != nil {
		return err
	}

	if installmentH == "current_mont" {
		if current > 1 {
			addInvoiceReference(&installment, invoiceRef)

			dateFixed := date.AddDate(0, current-1, 0)
			*lines = append(*lines, Line{dateFixed.Format(dateFormat), payee, installment, value})

			return nil
		}
	}

	var future string

	addInvoiceReference(&future, invoiceRef)

	for ; current <= total; current++ {
		dateFixed := date.AddDate(0, current-1, 0)
		memo := fmt.Sprintf("%d/%d - %s", current, total, future)
		*lines = append(*lines, Line{dateFixed.Format(dateFormat), payee, memo, value})
		// future = " futuro"
	}

	return nil
}

func addInvoiceReference(origin *string, invoiceRef string) {

	if len(invoiceRef) > 0 {
		if len(*origin) > 0 {
			*origin = fmt.Sprintf("%s - %s", *origin, invoiceRef)
		} else {
			*origin = invoiceRef
		}
	}

}
