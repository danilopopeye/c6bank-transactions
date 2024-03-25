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

func scanCSVRows(file io.Reader) ([]line, error) {
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	_, err := csvReader.Read() // header
	if err != nil {
		return nil, err
	}

	var lines []line

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
			if err := handleInstallments(record, &lines); err != nil {
				return nil, err
			}

			continue
		}

		lines = append(lines, line{record[0], record[4], record[5], record[8]})
	}

	return lines, nil
}

func handleInstallments(record []string, lines *[]line) error {
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

	if current > 1 {
		dateFixed := date.AddDate(0, current-1, 0)
		*lines = append(*lines, line{dateFixed.Format(dateFormat), payee, installment, value})
		return nil
	}

	var future string

	for ; current <= total; current++ {
		dateFixed := date.AddDate(0, current-1, 0)
		memo := fmt.Sprintf("%d/%d%s", current, total, future)
		*lines = append(*lines, line{dateFixed.Format(dateFormat), payee, memo, value})
		future = " futuro"
	}

	return nil
}
