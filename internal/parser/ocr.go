package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"git.home/c6bank-transactions/internal/image"
	"git.home/c6bank-transactions/internal/parser/ocr"
)

// 0     1            2             3          4          5        6               7                8
// Data; Nome Cartão; Final Cartão; Categoria; Descrição; Parcela; Valor (em US$); Cotação (em R$); Valor (em R$)

const (
	space            = " "
	lfRune           = '\n'
	processingText   = "Em processamento"
	installmentsText = "Parcela"
	cardText         = "Cartao final"
	firstInstallment = "1/"
	tesseractBin     = "/usr/bin/tesseract"
)

var (
	empty Transaction

	regexDate         = regexp.MustCompile(`^(\d{2})\/(\d{2})\s*`)
	regexCard         = regexp.MustCompile(`Cartao final\s*(\d+)`)
	regexValue        = regexp.MustCompile(`R\$\s*(-?[0-9.]+\,\d+)`)
	regexInstallments = regexp.MustCompile(`Parcela\s*(\d+)\s*de\s*(\d+)`)
)

func scanImage(file io.ReadSeeker, ref string) ([]Transaction, error) {
	cropped, err := image.Crop(file)
	if err != nil {
		return nil, err
	}

	text, err := ocr.Parse(cropped)
	if err != nil {
		return nil, err
	}

	return ScanImageLines(Time{}, text, ref)
}

func ScanImageLines(ct CurrentTime, text io.Reader, ref string) ([]Transaction, error) {
	var (
		transactions []Transaction
		current      string
	)

	reader := bufio.NewReader(text)

	for {
		line, err := reader.ReadString(lfRune)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		if line == lf { // blank lines
			continue
		}

		if len(current) > 0 && regexDate.MatchString(line) {
			if transaction := parseTransaction(ct, current, ref); transaction != empty {
				transactions = append(transactions, transaction)
			}

			current = current[:0]
		}

		current += line
	}

	if transaction := parseTransaction(ct, current, ref); transaction != empty {
		transactions = append(transactions, transaction)
	}

	installments, err := parseInstallments(transactions, ref)
	if err != nil {
		return nil, err
	}

	return append(transactions, installments...), nil
}

func parseInstallments(ts []Transaction, referenceText string) ([]Transaction, error) {
	var (
		installments []Transaction
		current      int
		total        int
		memo         string
		err          error
	)

	reference, err := time.Parse(refFormat, referenceText)
	if err != nil {
		return nil, err
	}

	for _, t := range ts {
		if !t.Installment {
			continue
		}

		current, total, memo, err = parseInstannmentText(t.Memo)
		if err != nil {
			return nil, err
		}

		for i := 1; current+i <= total; i++ {
			installmentDate := t.Date.AddDate(0, i, 0)
			referenceDate := reference.AddDate(0, i, 0)

			installments = append(installments, Transaction{
				Date:        installmentDate,
				Payee:       t.Payee,
				Memo:        fmt.Sprintf("%d/%d %s %s", current+i, total, memo, referenceDate.Format(refFormat)),
				Amount:      t.Amount,
				Installment: true,
				Future:      true,
			})
		}
	}

	return installments, nil
}

func parseInstannmentText(text string) (int, int, string, error) {
	parts := strings.SplitN(text, " ", 3)
	installmentsText := strings.SplitN(parts[0], "/", 2)

	fmt.Printf("parts: %#v\n", parts)

	current, err := strconv.Atoi(installmentsText[0])
	if err != nil {
		return 0, 0, "", err
	}

	total, err := strconv.Atoi(installmentsText[1])
	if err != nil {
		return 0, 0, "", err
	}

	return current, total, parts[1], nil
}

func parseTransaction(ct CurrentTime, line, ref string) Transaction {
	if strings.Contains(line, processingText) {
		return empty
	}

	var transaction Transaction

	// installments

	if strings.Contains(line, installmentsText) {
		transaction.Memo += parseRegex(line, regexInstallments) + space
		transaction.Installment = true

		if transaction.Memo[:2] != firstInstallment {
			return empty
		}
	}
	line = regexInstallments.ReplaceAllString(line, "")

	// date

	if err := transaction.ParseDate(ct, line); err != nil {
		return empty
	}
	line = line[5:]

	// card

	if strings.Contains(line, cardText) {
		transaction.Memo += parseRegex(line, regexCard) + space
	}
	line = regexCard.ReplaceAllString(line, "")

	// value

	transaction.Amount = parseRegex(line, regexValue)
	line = regexValue.ReplaceAllString(line, "")

	// payee

	transaction.Payee = strings.TrimSpace(strings.ReplaceAll(line, lf, ""))

	// memo

	transaction.Memo += ref

	return transaction
}

func ParseDate(ct CurrentTime, date string) time.Time {
	now := ct.Now()
	year := now.Year()

	month, err := strconv.Atoi(date[3:5])
	if err != nil {
		return time.Time{}
	}

	day, err := strconv.Atoi(date[:2])
	if err != nil {
		return time.Time{}
	}

	if month > int(now.Month()) {
		year = now.AddDate(-1, 0, 0).Year()
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func parseRegex(t string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(t)
	switch len(matches) {
	case 2:
		return matches[1]
	case 3:
		return fmt.Sprintf("%s/%s", matches[1], matches[2])
	default:
		return ""
	}
}

func __parseDate(t string) string {
	now := time.Now()
	year := now.Format("2006")

	date := t[:5]
	month, err := strconv.Atoi(t[3:5])
	if err != nil {
		return ""
	}

	if month > int(now.Month()) {
		year = now.AddDate(-1, 0, 0).Format("2006")
	}

	return fmt.Sprintf("%s/%s", date, year)
}
