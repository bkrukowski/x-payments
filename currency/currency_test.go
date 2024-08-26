package currency_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"payments/currency"
)

func TestCurrency_Is(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		a        currency.Currency
		b        currency.Currency
		expected bool
	}{
		{
			a:        currency.AED,
			b:        currency.AED,
			expected: true,
		},
		{
			a:        currency.AED,
			b:        currency.USD,
			expected: false,
		},
		{
			a:        currency.USD,
			b:        currency.AED,
			expected: false,
		},
	}

	for i, s := range scenarios {
		s := s

		t.Run(fmt.Sprintf("scenario %d", i), func(t *testing.T) {
			assert.Equal(t, s.expected, s.a.Is(s.b))
		})
	}
}
