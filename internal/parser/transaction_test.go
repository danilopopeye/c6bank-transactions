package parser_test

import (
	"testing"
	"time"

	"git.home/c6bank-transactions/internal/parser"
	"github.com/stretchr/testify/assert"
)

// mockTime: 1985-08-05

func TestTransaction_ParseDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		date        string
		expected    string
		transaction parser.Transaction
	}{
		{name: "current year, before now", date: "21/03", expected: "1985-03-21"},
		{name: "later same year", date: "21/09", expected: "1984-09-21"},
		{
			name: "first installment", date: "21/08", expected: "1985-08-21",
			transaction: parser.Transaction{Installment: true, Memo: "1/4 4321 01/1970"},
		},
		{
			name: "old installment", date: "21/05", expected: "1985-08-21",
			transaction: parser.Transaction{Installment: true, Memo: "3/4 4321 01/1970"},
		},
		{
			name: "new year", date: "21/07", expected: "1986-03-21",
			transaction: parser.Transaction{Installment: true, Memo: "8/12 4321 01/1970"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			tx := test.transaction

			assert.NoError(t, tx.ParseDate(mockTime, test.date))
			assert.Equal(t, test.expected, tx.Date.Format(time.DateOnly))
		})
	}
}

func TestTransaction_CSVLine(t *testing.T) {
	t.Parallel()

	ts := parser.Transaction{
		Date:   time.Date(1985, time.December, 26, 0, 0, 0, 0, time.UTC),
		Payee:  "Payee",
		Memo:   "Memo",
		Amount: "123.45",
	}

	csv := ts.CSVLine()
	assert.Equal(t, []string{"26/12/1985", "Payee", "Memo", "123.45"}, csv)
}
