package parser_test

import (
	"testing"

	"git.home/c6bank-transactions/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTransactions(t *testing.T) {
	t.Parallel()

	transactionsLines := []parser.Line{
		{"26/12/0000", "BKIN ", "", "55,76"},
		{"26/12/0000", "PANIFICIO ", "", "41,00"},
		{"26/12/0000", "PADARIAFERPAO ", "", "11,55"},
		{"24/12/0000", "SUPERMERCADO JVA II ", "", "331,40"},
		{"24/12/0000", "MERCADOLIVRE*FORCATOT ", "", "422,96"},
		{"24/12/0000", "PADARIAFERPAO ", "", "21,91"},
		{"24/12/0000", "PADARIA ROMERA ", "", "108,60"},
		{"23/12/0000", "MINUTO PA-5858 ", "", "15,99"},
	}

	installmentTransactionsLines := []parser.Line{
		{"23/12/0000", "COMPRA PARCELADA", "1/2 - jan", "123,45"},
		{"23/01/0001", "COMPRA PARCELADA", "2/2 - jan", "123,45"},
	}

	ocrText := `26/12
BKIN R$ 55,76
Cartao final 6269

26/12
PANIFICIO R$ 41,00
Cartao final 6269

26/12
PADARIAFERPAO R$ 11,55
Cartao final 6269

24/12
SUPERMERCADO JVA II R$ 331,40
Cartao final 6269

24/12
MERCADOLIVRE*FORCATOT R$ 422,96
Cartao final 6269

24/12
PADARIAFERPAO R$ 21,91
Cartao final 6269

24/12
PADARIA ROMERA R$ 108,60
Cartao final 6269

23/12
MINUTO PA-5858 R$ 15,99
Cartao final 6269

23/12
PORTO SEGURO SEGUROS
Cartao final 6269

R$ 887,85
Parcela 1 de 2`

	ocrTextInstallment := `23/12
MINUTO PA-5858 R$ 15,99
Cartao final 6269

23/12
COMPRA PARCELADA
Cartao final 6269

R$ 123,45
Parcela 1 de 2`

	t.Run("simpleTransactionMatch", func(t *testing.T) {
		t.Parallel()

		transactions := parser.GetTransactions(ocrText)
		assert.EqualValues(t, transactions, transactionsLines)
	})

	t.Run("installmentTransactionMatch", func(t *testing.T) {
		t.Parallel()

		transactions, err := parser.GetInstallmentTransactions(ocrTextInstallment, "jan", "")
		require.NoError(t, err)

		assert.EqualValues(t, transactions, installmentTransactionsLines)
	})
}
