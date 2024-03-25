package qif_test

import (
	"fmt"
	"io"
	"testing"

	"git.home/c6bank-transactions/internal/qif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	transactions := []qif.Transaction{
		{
			ID:     0,
			Date:   "01/01/1111",
			Amount: "123,45",
			Payee:  "with memo",
			Memo:   "memo",
		},
		{
			ID:     1,
			Date:   "02/02/2222",
			Amount: "987,65",
			Payee:  "without memo",
		},
	}

	renderedTx := `
N17620778198587367266
D01/01/1111
Pwith memo
T123,45
Mmemo
^
N930274905452525158
D02/02/2222
Pwithout memo
T987,65
^`

	qifTypes := []qif.QIFType{qif.BankType, qif.CreditCardType}

	for _, qtype := range qifTypes {
		t.Run(string(qtype), func(t *testing.T) {
			parsed, err := qif.Parse(qtype, transactions)
			assert.NoError(t, err)

			output, err := io.ReadAll(parsed)
			require.NoError(t, err)
			rendered := fmt.Sprintf("!Type:%s%s", qtype, renderedTx)

			assert.Equal(t, rendered, string(output))
		})
	}
}
