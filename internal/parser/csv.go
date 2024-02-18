package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "02/01/2006"

func scanCSVRows(file io.Reader, _ string, csvFile *csv.Writer) error {
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	lines, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	// 0     1            2             3          4          5        6               7                8
	// Data; Nome Cartão; Final Cartão; Categoria; Descrição; Parcela; Valor (em US$); Cotação (em R$); Valor (em R$)
	for i, line := range lines {
		if i == 0 {
			continue
		}

		date, card, payee, memo, outflow, inflow := line[0], line[2], line[4], line[5], line[8], ""

		if outflow[0] == '-' {
			outflow, inflow = "", outflow
		}

		if memo != "Única" {
			if err := handleInstallments(line, csvFile); err != nil {
				return err
			}

			continue
		}

		// log.Printf("%v\n", []string{date, payee, fmt.Sprintf("%s; %s", card, memo), outflow, inflow})
		if err := csvFile.Write([]string{date, payee, fmt.Sprintf("%s; %s", card, memo), outflow, inflow}); err != nil {
			return err
		}
	}

	return nil
}

func handleInstallments(line []string, csvFile *csv.Writer) error {
	purchase, card, payee, installment, value := line[0], line[2], line[4], line[5], line[8]

	parts := strings.SplitN(installment, "/", 2)
	// if parts[0] == parts[1] {
	// 	return nil
	// }

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
		log.Printf("%s %s :: %s -> %s", payee, installment, date.Format(dateFormat), dateFixed.Format(dateFormat))

		if err := csvFile.Write([]string{dateFixed.Format(dateFormat), payee, fmt.Sprintf("%s; %d/%d", card, current, total), value, ""}); err != nil {
			return err
		}

		return nil
	}

	// if payee != "HS MORUMBI" {
	// 	return nil
	// }

	for ; current <= total; current++ {
		dateFixed := date.AddDate(0, current-1, 0)

		if err := csvFile.Write([]string{dateFixed.Format(dateFormat), payee, fmt.Sprintf("%s; %d/%d", card, current, total), value, ""}); err != nil {
			return err
		}
	}

	return nil
}
