package currency

import (
	"fmt"
	"strconv"
)

type InvalidFractional struct {
	Currency Currency
	Given    uint
	Max      uint
}

func newInvalidFractional(currency Currency, given uint, max uint) *InvalidFractional {
	return &InvalidFractional{Currency: currency, Given: given, Max: max}
}

func (i *InvalidFractional) Error() string {
	return fmt.Sprintf("invalid fractional value for %s, max %d, given %d", i.Currency.Code, i.Max, i.Given)
}

type Amount struct {
	Currency   Currency
	Integer    uint
	Fractional uint
}

func NewAmount(currency Currency, integer uint, fractional uint) (Amount, error) {
	if max := currency.maxFractional(); fractional > max {
		return Amount{}, newInvalidFractional(currency, fractional, max)
	}

	return Amount{
		Currency:   currency,
		Integer:    integer,
		Fractional: fractional,
	}, nil
}

// MustNewAmount panic on the invalid input.
// Personally, I would panic in [NewAmount], mention that in the doc comment, and did not write [MustNewAmount],
// IMO it's proper idiomatic solution in GO - it's an invalid input that never should reach the constructor,
// it's not an unexpected behaviour that we want to gracefully handle.
// I am happy to elaborate about that in person.
//
// Example from stdlib:
//
//	reflect.TypeOf(nil).Kind()
//
// Output:
//
//	panic: reflect: call of reflect.Value.Len on int Value
func MustNewAmount(currency Currency, integer uint, fractional uint) Amount {
	a, err := NewAmount(currency, integer, fractional)
	if err != nil {
		panic(fmt.Sprintf("currency.MustNewAmount: %s", err.Error()))
	}

	return a
}

// NewAmountFromFractions convert 999 cents to 9.99 USD.
// Designed to store integer instead of float in the DB - see precision error.
func NewAmountFromFractions(currency Currency, fractions uint) Amount {
	return MustNewAmount(currency, fractions/currency.integerDivider(), fractions%currency.integerDivider())
}

func (s Amount) String() string {
	return fmt.Sprintf(
		"%d.%0"+strconv.Itoa(int(s.Currency.DecimalDigits))+"d %s",
		s.Integer,
		s.Fractional,
		s.Currency.Code,
	)
}

func (s Amount) ToFractional() uint {
	return s.Integer*s.Currency.maxFractional() + s.Fractional
}
