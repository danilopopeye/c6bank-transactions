package ocr_test

import (
	"io"
	"os"
	"os/exec"
	"testing"

	"git.home/c6bank-transactions/internal/parser/ocr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const parsedText = `30/01
ESQUINA LISBOA Em processamento

LANCHON SAO PAU R$ 64,24
Cartão final 6137

30/01
E-GR COMERCI*EGR Em processamento

R$ 48,03
Comer SAO PAU Parcela 1 de 2
Cartão final 4432

29/01

MP *ALIEXPRESS R$ 167,91
Cartão final 4432

29/01

VENUTO R$ 92,21
Cartão final 6137

29/01

DROGASIL1164 R$ 96,53
Cartão final 8240

29/01

APPLE.COM/BILL R$ 14,90
Cartão final 4432
`

func TestParse(t *testing.T) {
	t.Parallel()

	if _, err := exec.LookPath(ocr.TesseractBin); err != nil {
		t.Skip("missing `tesseract` binary")
	}

	fixture, err := os.Open("../../../test/fixtures/cropped.png")
	require.NoError(t, err)

	reader, err := ocr.Parse(fixture)
	require.NoError(t, err)

	parsedBytes, err := io.ReadAll(reader)
	require.NoError(t, err)

	assert.Equal(t, parsedText, string(parsedBytes))
}
