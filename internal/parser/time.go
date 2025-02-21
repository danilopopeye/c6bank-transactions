package parser

import "time"

type CurrentTime interface {
	Now() time.Time
}

var _ CurrentTime = (*Time)(nil)

type Time struct{}

func (Time) Now() time.Time {
	return time.Now()
}
