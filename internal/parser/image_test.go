package parser_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"git.home/c6bank-transactions/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tCurr = "COMPRA A VISTA"
	tNew  = "NOVO PARCELAMENTO"
	tOld  = "PARCELAMENTO ANTIGO"
)

var (
	phCur = []byte("__CURRENT__")
	phNew = []byte("__NEW__")
	phOld = []byte("__OLD__")

	mockTime MockTime
)

type MockTime struct{}

func (MockTime) Now() time.Time {
	return time.Date(1985, time.August, 5, 0, 0, 0, 0, time.UTC)
}

func TestScanImageLines(t *testing.T) {
	t.Parallel()

	text := transactionText(t)

	lines, err := parser.ScanImageLines(mockTime, text, "1985-09-01")
	require.NoError(t, err)

	transactions := [][]string{
		{"1985-08-01", "COMPRA A VISTA", "4321 09/1985", "167,91", "false", "false"},
		{"1985-08-02", "NOVO PARCELAMENTO", "1/4 4321 09/1985", "70,75", "true", "false"},
		{"1985-09-02", "NOVO PARCELAMENTO", "2/4 4321 10/1985", "70,75", "true", "true"},
		{"1985-10-02", "NOVO PARCELAMENTO", "3/4 4321 11/1985", "70,75", "true", "true"},
		{"1985-11-02", "NOVO PARCELAMENTO", "4/4 4321 12/1985", "70,75", "true", "true"},
	}

	assert.Len(t, lines, len(transactions))

	content := fmt.Sprintf("%v", lines)

	assert.Equal(t, 1, strings.Count(content, tCurr))
	assert.Equal(t, 4, strings.Count(content, tNew))
	// assert.Equal(t, 5, strings.Count(content, tOld))

	for i, line := range lines {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			transaction := transactions[i]

			assert.Equal(t, transaction[0], line.Date.UTC().Format(time.DateOnly))
			assert.Equal(t, transaction[1], line.Payee)
			assert.Equal(t, transaction[2], line.Memo)
			assert.Equal(t, transaction[3], line.Amount)
			assert.Equal(t, transaction[4], fmt.Sprint(line.Installment))
			assert.Equal(t, transaction[5], fmt.Sprint(line.Future))
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	month := mockTime.Now().Month()
	year := mockTime.Now().Year()

	test := map[string]struct {
		input string
		want  time.Time
	}{
		"invalid": {
			input: "inv/alid",
			want:  time.Time{},
		},
		"current": {
			input: fmt.Sprintf("01/%02d", month),
			want:  time.Date(year, month, 1, 0, 0, 0, 0, time.Local),
		},
		"next": {
			input: fmt.Sprintf("01/%02d", month+1),
			want:  time.Date(year-1, month+1, 1, 0, 0, 0, 0, time.Local),
		},
	}

	for name, tt := range test {
		t.Run(name+tt.input, func(t *testing.T) {
			t.Parallel()

			got := parser.ParseDate(mockTime, tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func transactionText(t *testing.T) io.Reader {
	t.Helper()

	month := mockTime.Now().Month()
	oldMonth := mockTime.Now().AddDate(0, -3, 0).Month()

	curr := []byte(fmt.Sprintf("01/%02d", month))
	inst := []byte(fmt.Sprintf("02/%02d", month))
	old := []byte(fmt.Sprintf("03/%02d", oldMonth))

	fixture, err := os.Open("../../test/fixtures/transactions.txt")
	require.NoError(t, err)

	text, err := io.ReadAll(fixture)
	require.NoError(t, err)

	replaced := bytes.ReplaceAll(text, phOld, old)
	replaced = bytes.ReplaceAll(replaced, phCur, curr)
	replaced = bytes.ReplaceAll(replaced, phNew, inst)

	return bytes.NewBuffer(replaced)
}
