package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Date        time.Time
	Payee       string
	Memo        string
	Amount      string
	Installment bool
	Future      bool
}

func (t *Transaction) ParseDate(ct CurrentTime, date string) error {
	now := ct.Now()
	year := now.Year()

	month, err := strconv.Atoi(date[3:5])
	if err != nil {
		return err
	}

	day, err := strconv.Atoi(date[:2])
	if err != nil {
		return err
	}

	if t.Installment {
		parts := strings.SplitN(t.Memo, "/", 2)

		if parts[0] != "1" {
			installment, err := strconv.Atoi(parts[0])
			if err != nil {
				return err
			}

			installmentDate, err := time.Parse(dateFormat, fmt.Sprintf("%s/%d", date, year))
			if err != nil {
				return err
			}

			t.Date = installmentDate.AddDate(0, installment, 0)

			return nil
		}
	}

	if month > int(now.Month()) {
		year = now.AddDate(-1, 0, 0).Year()
	}

	t.Date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	return nil
}

func (ts Transaction) CSVLine() []string {
	return []string{ts.Date.Format(dateFormat), ts.Payee, ts.Memo, ts.Amount}
}
