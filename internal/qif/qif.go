package qif

import (
	"bytes"
	"io"
	"text/template"

	"github.com/segmentio/fasthash/fnv1a"
)

type Transaction struct {
	ID       uint64
	Date     string
	Amount   string
	Payee    string
	Memo     string
	Category string
}

type QIFType string

const (
	BankType       QIFType = "Bank"  // !Type:Bank  | Cash Flow: Checking & Savings Account
	CreditCardType QIFType = "CCard" // !Type:CCard | Cash Flow: Credit Card Account
	dateFormat             = "02/01/2006"

	// txFmt intentionally starts with a newline
	txFmt = `
N{{.ID}}
D{{.Date}}
P{{.Payee}}
T{{.Amount}}
{{- if .Memo}}
M{{.Memo}}
{{- end}}
^`
)

var txTemplate = template.Must(template.New("txFmt").Parse(txFmt))

func Parse(qtype QIFType, transactions []Transaction) (io.Reader, error) {
	buff := new(bytes.Buffer)
	buff.WriteString("!Type:")
	buff.WriteString(string(qtype))

	for _, tx := range transactions {
		tx.ID = fnv1a.HashString64(
			tx.Date + tx.Payee + tx.Amount,
		)

		if err := txTemplate.Execute(buff, tx); err != nil {
			return nil, err
		}
	}

	return buff, nil
}

// Field  Indicator Explanation
// D      Date
// T      Amount
// C      Cleared status
// N      Num (check or reference number)
// P      Payee
// M      Memo
// A      Address (up to five lines; the sixth line is an optional message)
// L      Category (Category/Subcategory/Transfer/Class)
// S      Category in split (Category/Transfer/Class)
// E      Memo in split
// $      Dollar amount of split
// ^      End of the entry
