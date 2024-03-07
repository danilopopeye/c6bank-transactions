package qif

// https://github.com/rockstardevs/csvtoqfx

import (
	"bytes"
	"encoding/csv"
	"io"
	"text/template"
)

// !Type:Bank
// D12/19/18
// PCOUNTY WASTE 12/18 PURCHASE 804-8439288 VA
// T-25.00
// C*
// ^

const (
	dateFormat = "01/02/2006"
	headerFmt  = `!Account
N{{.Name}}
T{{.Type}}
^
!Type:{{.Type}}`

	// txFmt intentionally starts with a newline
	txFmt = `
D{{.Date}}
P{{.Payee}}
T{{.Amount}}
{{- if .Memo}}
M{{.Memo}}
{{- end}}
N{{.ID}}
^`
)

var (
	headerTemplate = template.Must(template.New("headerFmt").Parse(headerFmt))
	txTemplate     = template.Must(template.New("txFmt").Parse(txFmt))
)

type header struct {
	Name string
	Type string
}

type transaction struct {
	ID     string // - N{{.ID}}
	Date   string
	Payee  string
	Memo   string
	Amount string
}

// !Type:Bank  | Cash Flow: Checking & Savings Account
// !Type:CCard | Cash Flow: Credit Card Account

func ParseCSV(file io.Reader) (io.Reader, error) {
	buff := new(bytes.Buffer)
	buff.WriteString("!Type:CCard")

	// hdr := header{Name: "cartão de crédito", Type: "CCard"}
	// if err := headerTemplate.Execute(buff, hdr); err != nil {
	// 	return nil, err
	// }

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	_, err := csvReader.Read() // header
	if err != nil {
		return nil, err
	}

	for {
		// 0     1            2             3          4          5        6               7                8
		// Data; Nome Cartão; Final Cartão; Categoria; Descrição; Parcela; Valor (em US$); Cotação (em R$); Valor (em R$)
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		tx := transaction{
			ID:     record[0] + record[2] + record[4] + record[5] + record[8],
			Date:   record[0],
			Payee:  record[4],
			Amount: record[8],
		}

		if record[5] != "Única" {
			tx.Memo = record[5]
		}

		if err := txTemplate.Execute(buff, tx); err != nil {
			return nil, err
		}
	}

	return buff, nil
}
