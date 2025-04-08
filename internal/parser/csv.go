package parser

import (
	"bytes"
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
	refFormat  = "01/2006"
	minus      = '-'
	fee        = "ANUIDADE DIFERENCIADA"
	unique     = "Única"
)

func scanCSVRows(reference time.Time, file io.Reader) ([]Line, error) {
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

		fixValue(&record[8])

		if record[5] == unique {
			record[5] = parseMemo(reference, record[2], 0, 0)
		} else {
			if err = handleInstallments(reference, record, &lines); err != nil {
				return nil, err
			}

			continue
		}

		date, err := time.Parse(dateFormat, record[0])
		if err != nil || date.IsZero() {
			continue // Skip appending if the date is empty or invalid
		}

		lines = append(lines, Line{record[0], record[4], record[5], record[8]})
	}

	return lines, nil
}

func handleInstallments(reference time.Time, record []string, lines *[]Line) error {
	purchase, card, payee, installment, value := record[0], record[2], record[4], record[5], record[8]

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

	if current > 1 {
		dateFixed := date.AddDate(0, current-1, 0)
		*lines = append(*lines, Line{
			dateFixed.Format(dateFormat), payee, parseMemo(reference, card, current, total),
			value,
		})

		return nil
	}

	for ; current <= total; current++ {
		ref := reference.AddDate(0, current-1, 0)
		dateFixed := date.AddDate(0, current-1, 0)
		memo := fmt.Sprintf("%d/%d %s %02d/%04d", current, total, card, ref.Month(), ref.Year())

		*lines = append(*lines, Line{dateFixed.Format(dateFormat), payee, memo, value})
	}

	return nil
}

func fixValue(value *string) {
	runes := []rune(*value)

	if runes[0] == minus {
		*value = string(runes[1:])
	} else {
		*value = "-" + *value
	}
}

func parseMemo(ref time.Time, card string, current, total int) string {
	var buf strings.Builder

	if current > 0 && total > 0 {
		buf.WriteString(fmt.Sprintf("%d/%d ", current, total))
	}

	buf.WriteString(card)

	if !ref.IsZero() {
		buf.WriteString(" ")
		buf.WriteString(ref.Format(refFormat))
	}

	return buf.String()
}

func TransactionsToCSV(transactions []Transaction) (io.Reader, error) {
	buf := new(bytes.Buffer)

	writer := csv.NewWriter(buf)
	// writer.Comma = ';'

	if err := writer.Write([]string{"Date", "Payee", "Memo", "Value"}); err != nil {
		return nil, err
	}

	for _, ts := range transactions {
		if err := writer.Write(ts.CSVLine()); err != nil {
			return nil, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf, nil
}
