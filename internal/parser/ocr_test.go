package parser_test

import (
	"testing"

	"git.home/c6bank-transactions/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTransactions(t *testing.T) {

	transactionsRes := [][]string{
		{"26/12", "BKIN ", "55,76"},
		{"26/12", "PANIFICIO ", "", "41,00"},
		{"26/12", "PADARIAFERPAO ", "", "11,55"},
		{"24/12", "SUPERMERCADO JVA II ", "", "331,40"},
		{"24/12", "MERCADOLIVRE*FORCATOT ", "", "422,96"},
		{"24/12", "PADARIAFERPAO ", "", "21,91"},
		{"24/12", "PADARIA ROMERA ", "", "108,60"},
		{"23/12", "MINUTO PA-5858 ", "", "15,99"},
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
Parcela 1 de 10`

	t.Run("simpleTransactionMatch", func(t *testing.T) {
		transactions := parser.GetTransactions(ocrText)
		assert.Equal(t, transactions, transactionsRes)
		for i, tr := range transactions {
			t.Logf("%d -> Data(%s) Payee(%s) Parcela(%s) Valor(%s)", i, tr[0], tr[1], tr[2], tr[3])
		}

	})

	t.Run("installmentTransactionMatch", func(t *testing.T) {
		transactions, err := parser.GetInstallmentTransactions(ocrText, "jan")

		require.NoError(t, err)

		assert.Equal(t, transactions, transactionsRes)

		for i, tr := range transactions {
			t.Logf("%d -> Data(%s) Payee(%s) Parcela(%s) Valor(%s)", i, tr[0], tr[1], tr[2], tr[3])
		}

	})
}
