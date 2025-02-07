package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/ledongthuc/pdf"
)

const lf = "\n"

// 0    1    2     3    4     5
// line date payee memo value type
var transactionRegexp = regexp.MustCompile(`(?P<date>[0-9/]{10})\s+(?P<payee>[A-Z0-9., ]+)\s*-?\s*(?P<memo>.*)\s+[0-9]{12}\s+(?P<value>[0-9.]+,[0-9]{2})\s+(?P<type>[CD])`)

func scanPDFRows(file io.ReaderAt, pass string, size int64) ([]Line, error) {
	content, err := readPDF(file, size, pass)
	if err != nil {
		return nil, err
	}

	var lines []Line
	scanner := bufio.NewScanner(content)

	for scanner.Scan() {
		record := transactionRegexp.FindStringSubmatch(scanner.Text())
		if len(record) != 6 {
			continue
		}

		if record[5] == "D" {
			record[4] = fmt.Sprintf("-%s", record[4])
		}

		lines = append(lines, Line{record[1], record[2], record[3], record[4]})
	}

	return lines, nil
}

func readPDF(file io.ReaderAt, size int64, pass string) (io.Reader, error) {
	reader, err := pdf.NewReaderEncrypted(file, size, func() string { return pass })
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)
	totalPage := reader.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := reader.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, err := p.GetTextByRow()
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			for _, word := range row.Content {
				buff.WriteString(word.S)
				buff.WriteString(" ")
			}

			buff.WriteString(lf)
		}
	}

	return buff, nil
}
