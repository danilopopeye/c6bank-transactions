package parser_test

import (
	"testing"
	"time"

	"git.home/c6bank-transactions/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestTime_Now(t *testing.T) {
	pt := parser.Time{}

	assert.Equal(t, time.Now().Unix(), pt.Now().Unix())
}
