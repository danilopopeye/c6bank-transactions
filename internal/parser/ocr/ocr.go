package ocr

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

const TesseractBin = "tesseract"

var ErrOCRParse = errors.New("ocr parse error")

func Parse(file io.Reader) (io.Reader, error) {
	var (
		ocrOutput bytes.Buffer
		ocrError  bytes.Buffer
	)

	cmd := exec.Command(TesseractBin, "stdin", "stdout", "--psm", "4", "-l", "por", "-l", "eng")
	cmd.Stdin = file
	cmd.Stdout = &ocrOutput
	cmd.Stderr = &ocrError

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s - %s", ErrOCRParse, err, ocrError.String())
	}

	return &ocrOutput, nil
}
