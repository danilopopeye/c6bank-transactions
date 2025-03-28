package parser_test

import (
	"encoding/csv"
	"testing"
	"time"

	"git.home/c6bank-transactions/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsToCSV(t *testing.T) {
	t.Parallel()

	transaction := parser.Transaction{Date: time.Now(), Payee: "Payee", Memo: "Memo", Amount: "123.45"}
	transactions := []parser.Transaction{transaction}

	reader, err := parser.TransactionsToCSV(transactions)
	require.NoError(t, err)

	csvReader := csv.NewReader(reader)

	csvLines, err := csvReader.ReadAll()
	require.NoError(t, err)

	header := csvLines[0]
	assert.Equal(t, "Date", header[0])
	assert.Equal(t, "Payee", header[1])
	assert.Equal(t, "Memo", header[2])
	assert.Equal(t, "Value", header[3])

	line := csvLines[1]
	assert.Equal(t, transaction.Date.Format("02/01/2006"), line[0])
	assert.Equal(t, transaction.Payee, line[1])
	assert.Equal(t, transaction.Memo, line[2])
	assert.Equal(t, transaction.Amount, line[3])
}
