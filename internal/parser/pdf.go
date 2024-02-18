package parser

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

func scanPDFRows(file io.Reader, pass string, csvFile *csv.Writer) error {
	content, err := readPDF(file, pass)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(content))

	for scanner.Scan() {
		line := transactionRegexp.FindStringSubmatch(scanner.Text())
		if len(line) != 7 {
			continue
		}

		var inflow, outflow string
		date, payee, memo, cd, amount := line[1], strings.TrimSpace(line[2]), strings.TrimSpace(line[3]), line[6], line[5]

		switch cd {
		case "C":
			inflow = amount
		case "D":
			outflow = amount
		default:
			return fmt.Errorf("wrong type of transaction: %s -- %s", cd, line)
		}

		if err := csvFile.Write([]string{date, payee, memo, outflow, inflow}); err != nil {
			return err
		}
	}

	return nil
}

func readPDF(file io.Reader, pass string) ([]byte, error) {
	cmd := exec.Command("pdftotext", "-layout", "-upw", pass, "-", "/dev/stdout")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(stdin, file); err != nil {
		return nil, err
	}

	stdin.Close()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("coud not parse pdf: %s - %s", bytes.TrimSpace(out), err)
		return nil, fmt.Errorf("%s", bytes.TrimSpace(out))
	}

	return out, err
}
