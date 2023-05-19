package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	file              string
	pass              string
	csvHeaders        = []string{"Date", "Payee", "Memo", "Outflow", "Inflow"}
	transactionRegexp = regexp.MustCompile(`(?P<date>[0-9/]{10}) (?P<payee>[A-Z0-9., ]+)-?(?P<memo>.*)\s+(?P<doc>[0-9]{12})\s+(?P<value>[0-9.]+,[0-9]{2})\s+(?P<type>[CD])`)
)

func main() {
	log.SetFlags(0)

	args := len(os.Args)

	log.Printf("%v", os.Args)

	if args > 1 {
		file = os.Args[1]
	}

	if args > 2 {
		pass = os.Args[2]
	}

	if file == "" {
		log.Printf("USAGE: %s path/to/file.(pdf|csv) [password]", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	csvFile := csv.NewWriter(os.Stdout)
	if err := csvFile.Write(csvHeaders); err != nil {
		panic(err)
	}

	if strings.HasSuffix(file, "pdf") {
		scanPdfRows(file, pass, csvFile)
	} else if strings.HasSuffix(file, "csv") {
		scanCSVRows(file, csvFile)
	}

	csvFile.Flush()

	if err := csvFile.Error(); err != nil {
		panic(err)
	}
}

func scanCSVRows(filename string, csvFile *csv.Writer) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	lines, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	// 0              1              2               3         4         5       6              7               8
	// Data de Compra;Nome no Cartão;Final do Cartão;Categoria;Descrição;Parcela;Valor (em US$);Cotação (em R$);Valor (em R$)
	for i, line := range lines {
		if i == 0 {
			continue
		}

		date, category, payee, installment, value := line[0], line[3], line[4], line[5], line[8]
		amount, installments := parseAmount(line[8], line[5])

		if category == "-" {
			fmt.Fprintf(os.Stderr, "ignored: %#v\n", line)
			continue
		}
		if amount.IsZero() {
			fmt.Fprintf(os.Stderr, "installment %s: %#v\n", installment, line)
			continue
		}

		var memo string

		if installments != "" {
			memo = fmt.Sprintf("%sx de R$ %s", installments, value)
		}

		if err := csvFile.Write([]string{date, payee, memo, amount.StringFixedBank(2), ""}); err != nil {
			panic(err)
		}
	}

	println(len(lines))
}

func scanPdfRows(file, pass string, csvFile *csv.Writer) {
	content := readPdf(file, pass)
	scanner := bufio.NewScanner(bytes.NewBufferString(content))

	for scanner.Scan() {
		line := transactionRegexp.FindStringSubmatch(scanner.Text())
		if len(line) != 7 {
			continue
		}

		var outflow decimal.Decimal
		var inflow decimal.Decimal

		date, payee, memo, cd := line[1], strings.TrimSpace(line[2]), strings.TrimSpace(line[3]), line[6]
		amount, _ := parseAmount(line[5], "")

		switch cd {
		case "C":
			inflow = amount
		case "D":
			outflow = amount
		default:
			println(line[0])
			panic("wrong type of CD")
		}

		if err := csvFile.Write([]string{date, payee, memo, outflow.StringFixedBank(2), inflow.StringFixedBank(2)}); err != nil {
			panic(err)
		}

		println("purchase:", line[0])
	}
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

func readPdf(path, pass string) string {
	bytes, err := exec.Command("pdftotext", "-layout", "-upw", pass, path, "-").CombinedOutput()
	if err != nil {
		log.Fatalf("error: %s\n%s", err, bytes)
	}

	return string(bytes)
}
