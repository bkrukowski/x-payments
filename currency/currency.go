package currency

var (
	AED = Currency{
		Code:          "AED",
		DecimalDigits: 2,
	}
	USD = Currency{
		Code:          "USD",
		DecimalDigits: 2,
	}
)

type Currency struct {
	Code          string // Code represents ISO427
	DecimalDigits uint
}

func (c Currency) Is(c2 Currency) bool {
	return c.Code == c2.Code
}

// integerDivider represents the number of fractional units that convert to a single integer unit,
// e.g. 100 cents represents 1 USD
func (c Currency) integerDivider() uint {
	max := uint(1)

	for i := uint(0); i < c.DecimalDigits; i++ {
		max = max * 10
	}

	return max
}

// maxFractional return maximum allowed fractional value, e.g. 99 for USD
func (c Currency) maxFractional() uint {
	return c.integerDivider() - 1
}
