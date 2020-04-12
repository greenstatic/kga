package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateTypeString(t *testing.T) {
	assert_ := assert.New(t)

	type test struct {
		input string
		err   error
	}

	tests := []test{
		{"basic", nil},
		{"manifest", nil},
		{"helm", nil},
		{"", InvalidTypeStringError("")},
		{" ", InvalidTypeStringError("")},
		{"test", InvalidTypeStringError("")},
		{".", InvalidTypeStringError("")},
	}

	for i, tst := range tests {
		err := ValidateTypeString(tst.input)

		assert_.IsType(tst.err, err, fmt.Sprintf("Failed test: %d", i))
	}
}
