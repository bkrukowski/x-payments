package currency_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"payments/currency"
)

func TestAmount_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		amount currency.Amount
		want   string
	}{
		{
			name:   "150.15 USD",
			amount: currency.MustNewAmount(currency.USD, 150, 15),
			want:   "150.15 USD",
		},
		{
			name:   "2000.00 AED",
			amount: currency.MustNewAmount(currency.AED, 2000, 0),
			want:   "2000.00 AED",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.amount.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAmount(t *testing.T) {
	t.Parallel()

	t.Run("Incorrect input", func(t *testing.T) {
		t.Parallel()

		type input struct {
			currency.Currency
			integer    uint
			fractional uint
		}

		scenarios := []struct {
			name     string
			input    input
			expected string
		}{
			{
				name: "1.100 USD",
				input: input{
					Currency:   currency.USD,
					integer:    1,
					fractional: 100,
				},
				expected: "invalid fractional value for USD, max 99, given 100",
			},
		}

		for _, s := range scenarios {
			s := s

			t.Run(s.name, func(t *testing.T) {
				a, err := currency.NewAmount(s.input.Currency, s.input.integer, s.input.fractional)
				require.EqualError(t, err, s.expected)
				assert.Zero(t, a)
			})

			t.Run(fmt.Sprintf("Must %s", s.name), func(t *testing.T) {
				defer func() {
					r := recover()
					require.NotEmpty(t, r)
					assert.Equal(t, fmt.Sprintf("currency.MustNewAmount: %s", s.expected), r)
				}()

				assert.Zero(
					t,
					currency.MustNewAmount(s.input.Currency, s.input.integer, s.input.fractional),
				)
			})
		}
	})
}

func TestNewAmountFromFractions(t *testing.T) {
	t.Parallel()

	type args struct {
		currency  currency.Currency
		fractions uint
	}

	tests := []struct {
		name string
		args args
		want currency.Amount
	}{
		{
			name: "1005 cents",
			args: args{
				currency:  currency.USD,
				fractions: 1005,
			},
			want: currency.MustNewAmount(currency.USD, 10, 5),
		},
		{
			name: "100 cents",
			args: args{
				currency:  currency.USD,
				fractions: 100,
			},
			want: currency.MustNewAmount(currency.USD, 1, 0),
		},
		{
			name: "999 fils",
			args: args{
				currency:  currency.AED,
				fractions: 999,
			},
			want: currency.MustNewAmount(currency.AED, 9, 99),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equalf(
				t,
				tt.want,
				currency.NewAmountFromFractions(tt.args.currency, tt.args.fractions),
				"NewAmountFromFractions(%v, %v)",
				tt.args.currency,
				tt.args.fractions,
			)
		})
	}
}
