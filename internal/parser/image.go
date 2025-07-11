package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"slices"
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
	coma             = ","
	lfRune           = '\n'
	processingText   = "Em processamento"
	installmentsText = "Parcela"
	cardText         = "Cartao final"
	cardTextAccent   = "Cartão final"
	firstInstallment = "1/"
	tesseractBin     = "/usr/bin/tesseract"
)

var (
	ErrInvalidReference = fmt.Errorf("could not parse reference")
	empty               Transaction

	regextMultiSpace  = regexp.MustCompile(`\s+`)
	regexDate         = regexp.MustCompile(`^(\d{2})\/(\d{2})\s*`)
	regexCard         = regexp.MustCompile(`Cart[aã]o final\s*(\d+)`)
	regexValue        = regexp.MustCompile(`R\$\s*(-?[0-9.]+[, ]\d+)`)
	regexInstallments = regexp.MustCompile(`Parcela.*(\d+)\s*de\s*(\d+)`)
	regexReference    = regexp.MustCompile(`(janeiro|fevereiro|março|abril|maio|junho|julho|agosto|setembro|outubro|novembro|dezembro)`)

	months = []string{
		"", "janeiro", "fevereiro", "março", "abril", "maio",
		"junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
	}
)

func ScanImage(file io.ReadSeeker, includeProcessing bool) ([]Transaction, error) {
	cropped, reference, err := image.Crop(file)
	if err != nil {
		return nil, err
	}

	text, err := ocr.Parse(cropped)
	if err != nil {
		return nil, err
	}

	refText, err := ocr.Parse(reference)
	if err != nil {
		return nil, err
	}

	month, err := ParseRef(refText)
	if err != nil {
		return nil, err
	}

	return ScanImageLines(Time{}, text, month, includeProcessing)
}

func ScanImageLines(ct CurrentTime, text io.Reader, ref time.Time, includeProcessing bool) ([]Transaction, error) {
	var (
		transactions []Transaction
		current      string
	)

	reader := bufio.NewReader(text)

	refText := ref.Format(refFormat)

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
			if transaction := parseTransaction(ct, current, refText, includeProcessing); transaction != empty {
				transactions = append(transactions, transaction)
			}

			current = current[:0]
		}

		current += line
	}

	if transaction := parseTransaction(ct, current, refText, includeProcessing); transaction != empty {
		transactions = append(transactions, transaction)
	}

	installments, err := parseInstallments(transactions, ref)
	if err != nil {
		return nil, err
	}

	return append(transactions, installments...), nil
}

func parseInstallments(ts []Transaction, refTime time.Time) ([]Transaction, error) {
	var (
		installments []Transaction
		current      int
		total        int
		memo         string
		err          error
	)

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
			referenceDate := refTime.AddDate(0, i, 0)

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

func parseTransaction(ct CurrentTime, line, ref string, includeProcessing bool) Transaction {
	if strings.Contains(line, processingText) {
		if !includeProcessing {
			return empty
		}

		line = strings.ReplaceAll(line, processingText, "")
	}

	// fmt.Println("line:", line)

	var transaction Transaction

	// installments

	if strings.Contains(line, installmentsText) {
		transaction.Memo += parseRegex(line, regexInstallments) + space
		transaction.Installment = true

		if transaction.Memo == space || transaction.Memo[:2] != firstInstallment {
			return empty
		}
	}
	line = regexInstallments.ReplaceAllString(line, "")

	// fmt.Println("installments:", line)

	// date

	if err := transaction.ParseDate(ct, line); err != nil {
		return empty
	}
	line = line[5:]

	// fmt.Println("date:", line)

	// card

	if strings.Contains(line, cardText) || strings.Contains(line, cardTextAccent) {
		transaction.Memo += parseRegex(line, regexCard) + space
	}
	line = regexCard.ReplaceAllString(line, "")

	// fmt.Println("card:", line)

	// value

	transaction.Amount = parseRegex(line, regexValue)
	line = regexValue.ReplaceAllString(line, "")

	// fmt.Println("value:", line)

	// payee

	transaction.Payee = strings.TrimSpace(
		regextMultiSpace.ReplaceAllString(
			strings.ReplaceAll(line, lf, " "), space,
		),
	)

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
		return strings.Replace(matches[1], space, coma, 1)
	case 3:
		return fmt.Sprintf("%s/%s", matches[1], matches[2])
	default:
		return ""
	}
}

func ParseRef(reference io.Reader) (time.Time, error) {
	btext, err := io.ReadAll(reference)
	if err != nil {
		return time.Time{}, err
	}
	text := string(btext)

	matches := regexReference.FindStringSubmatch(text)
	if matches == nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidReference, text)
	}

	index := slices.Index(months, matches[1])
	if index == -1 {
		return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidReference, text)
	}

	now := time.Now()

	// TODO should handle next year
	return time.Date(now.Year(), time.Month(index), 1, 0, 0, 0, 0, time.Local), nil
}
